package sqlhelp

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"

	. "github.com/ankyra/escape-inventory/dao/types"
)

func (s *SQLHelper) AddNamespace(namespace *Project) error {
	return s.PrepareAndExecInsert(s.AddProjectQuery,
		namespace.Name,
		namespace.Description,
		namespace.OrgURL,
		namespace.Logo)
}

func (s *SQLHelper) UpdateNamespace(namespace *Project) error {
	return s.PrepareAndExecUpdate(s.UpdateProjectQuery,
		namespace.Name,
		namespace.Description,
		namespace.OrgURL,
		namespace.Logo,
		namespace.Name)
}

func (s *SQLHelper) GetNamespaceHooks(namespace *Project) (Hooks, error) {
	rows, err := s.PrepareAndQuery(s.GetProjectHooksQuery, namespace.Name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		return s.scanHooks(rows)
	}
	return nil, NotFound
}

func (s *SQLHelper) SetNamespaceHooks(namespace *Project, hooks Hooks) error {
	bytes, err := json.Marshal(hooks)
	if err != nil {
		return err
	}
	return s.PrepareAndExecUpdate(s.SetProjectHooksQuery,
		string(bytes),
		namespace.Name)
}

func (s *SQLHelper) GetNamespace(namespace string) (*Project, error) {
	rows, err := s.PrepareAndQuery(s.GetProjectQuery, namespace)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		return s.scanNamespace(rows)
	}
	return nil, NotFound
}

func (s *SQLHelper) GetNamespaces() (map[string]*Project, error) {
	rows, err := s.PrepareAndQuery(s.GetProjectsQuery)
	if err != nil {
		return nil, err
	}
	return s.scanNamespaces(rows)
}

func (s *SQLHelper) GetNamespacesByGroups(readGroups []string) (map[string]*Project, error) {
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

func (s *SQLHelper) HardDeleteNamespace(project string) error {
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

func (s *SQLHelper) scanNamespace(rows *sql.Rows) (*Project, error) {
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

func (s *SQLHelper) scanNamespaces(rows *sql.Rows) (map[string]*Project, error) {
	defer rows.Close()
	result := map[string]*Project{}
	for rows.Next() {
		prj, err := s.scanNamespace(rows)
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
