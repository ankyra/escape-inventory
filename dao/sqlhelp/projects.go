package sqlhelp

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"

	. "github.com/ankyra/escape-inventory/dao/types"
)

func (s *SQLHelper) AddProject(project *Project) error {
	return s.PrepareAndExecInsert(s.AddProjectQuery,
		project.Name,
		project.Description,
		project.OrgURL,
		project.Logo)
}

func (s *SQLHelper) UpdateProject(project *Project) error {
	return s.PrepareAndExecUpdate(s.UpdateProjectQuery,
		project.Name,
		project.Description,
		project.OrgURL,
		project.Logo,
		project.Name)
}

func (s *SQLHelper) GetProjectHooks(project *Project) (Hooks, error) {
	rows, err := s.PrepareAndQuery(s.GetProjectHooksQuery, project.Name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		return s.scanHooks(rows)
	}
	return nil, NotFound
}

func (s *SQLHelper) SetProjectHooks(project *Project, hooks Hooks) error {
	bytes, err := json.Marshal(hooks)
	if err != nil {
		return err
	}
	return s.PrepareAndExecUpdate(s.SetProjectHooksQuery,
		string(bytes),
		project.Name)
}

func (s *SQLHelper) GetProject(project string) (*Project, error) {
	rows, err := s.PrepareAndQuery(s.GetProjectQuery, project)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		return s.scanProject(rows)
	}
	return nil, NotFound
}

func (s *SQLHelper) GetProjects() (map[string]*Project, error) {
	rows, err := s.PrepareAndQuery(s.GetProjectsQuery)
	if err != nil {
		return nil, err
	}
	return s.scanProjects(rows)
}

func (s *SQLHelper) GetProjectsByGroups(readGroups []string) (map[string]*Project, error) {
	starFound := false
	for _, g := range readGroups {
		if g == "*" {
			starFound = true
			break
		}
	}
	if !starFound {
		readGroups = append(readGroups, "*")
	}
	insertMarks := []string{}
	for i, _ := range readGroups {
		if s.UseNumericInsertMarks {
			insertMarks = append(insertMarks, "$"+strconv.Itoa(i+1))
		} else {
			insertMarks = append(insertMarks, "?")
		}
	}
	query := s.GetProjectsByGroupsQuery
	if len(readGroups) == 1 {
		query += " = " + insertMarks[0]
	} else {
		query += "IN (" + strings.Join(insertMarks, ", ") + ")"
	}
	interfaceGroups := []interface{}{}
	for _, g := range readGroups {
		interfaceGroups = append(interfaceGroups, g)
	}

	rows, err := s.PrepareAndQuery(query, interfaceGroups...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := map[string]*Project{}
	groups := map[string]map[string]bool{}
	for rows.Next() {
		var name, description, orgURL, logo, matchedGroup string
		if err := rows.Scan(&name, &description, &orgURL, &logo, &matchedGroup); err != nil {
			return nil, err
		}
		prj, ok := result[name]
		if !ok {
			prj = &Project{
				Name:           name,
				Description:    description,
				OrgURL:         orgURL,
				Logo:           logo,
				Permission:     "admin", // default permission for open source
				MatchingGroups: []string{matchedGroup},
			}
			groups[name] = map[string]bool{
				matchedGroup: true,
			}
		} else {
			_, added := groups[name][matchedGroup]
			if !added {
				prj.MatchingGroups = append(prj.MatchingGroups, matchedGroup)
				groups[name][matchedGroup] = true
			}
		}
		result[prj.Name] = prj
	}
	return result, nil
}

func (s *SQLHelper) HardDeleteProject(project string) error {
	if err := s.PrepareAndExec(s.HardDeleteProjectACLQuery, project); err != nil {
		return err
	}
	if err := s.PrepareAndExec(s.HardDeleteProjectUnitSubscriptions, project); err != nil {
		return err
	}
	if err := s.PrepareAndExec(s.HardDeleteProjectReleaseDependenciesQuery, project); err != nil {
		return err
	}
	if err := s.PrepareAndExec(s.HardDeleteProjectPackageURIsQuery, project); err != nil {
		return err
	}
	if err := s.PrepareAndExec(s.HardDeleteProjectReleasesQuery, project); err != nil {
		return err
	}
	if err := s.PrepareAndExec(s.HardDeleteProjectApplicationsQuery, project); err != nil {
		return err
	}
	if err := s.PrepareAndExec(s.HardDeleteProjectQuery, project); err != nil {
		return err
	}
	return nil
}

func (s *SQLHelper) scanProject(rows *sql.Rows) (*Project, error) {
	var name, description, orgURL, logo string
	if err := rows.Scan(&name, &description, &orgURL, &logo); err != nil {
		return nil, err
	}
	return &Project{
		Name:        name,
		Description: description,
		OrgURL:      orgURL,
		Logo:        logo,
		Permission:  "admin", // default permission for open source
	}, nil
}

func (s *SQLHelper) scanProjects(rows *sql.Rows) (map[string]*Project, error) {
	defer rows.Close()
	result := map[string]*Project{}
	for rows.Next() {
		prj, err := s.scanProject(rows)
		if err != nil {
			return nil, err
		}
		result[prj.Name] = prj
	}
	return result, nil
}

func (s *SQLHelper) scanHooks(rows *sql.Rows) (Hooks, error) {
	var hooksString string
	if err := rows.Scan(&hooksString); err != nil {
		return nil, err
	}
	result := NewHooks()
	if err := json.Unmarshal([]byte(hooksString), &result); err != nil {
		return nil, err
	}
	return result, nil
}
