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
		namespace.Logo,
		namespace.IsPublic)
}

func (s *SQLHelper) UpdateNamespace(namespace *Project) error {
	return s.PrepareAndExecUpdate(s.UpdateProjectQuery,
		namespace.Name,
		namespace.Description,
		namespace.OrgURL,
		namespace.Logo,
		namespace.Name,
		namespace.IsPublic)
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

func (s *SQLHelper) GetNamespacesByNames(namespaces []string) (map[string]*Project, error) {
	insertMarks := []string{}
	for i, _ := range namespaces {
		if s.UseNumericInsertMarks {
			insertMarks = append(insertMarks, "$"+strconv.Itoa(i+1))
		} else {
			insertMarks = append(insertMarks, "?")
		}
	}
	query := s.GetNamespacesByNamesQuery
	if len(namespaces) == 0 {
		return map[string]*Project{}, nil
	} else if len(namespaces) == 1 {
		query += " = " + insertMarks[0]
	} else {
		query += " IN (" + strings.Join(insertMarks, ", ") + ")"
	}
	interfaceNamespaces := []interface{}{}
	for _, n := range namespaces {
		interfaceNamespaces = append(interfaceNamespaces, n)
	}
	rows, err := s.PrepareAndQuery(query, interfaceNamespaces...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := map[string]*Project{}
	for rows.Next() {
		var name, description, orgURL, logo string
		var isPublic bool
		if err := rows.Scan(&name, &description, &orgURL, &logo, &isPublic); err != nil {
			return nil, err
		}
		prj, ok := result[name]
		if !ok {
			prj = &Project{
				Name:        name,
				Description: description,
				OrgURL:      orgURL,
				Logo:        logo,
				Permission:  "admin", // default permission for open source
				IsPublic:    isPublic,
			}
		}
		result[prj.Name] = prj
	}
	return result, nil
	return nil, nil
}

func (s *SQLHelper) GetNamespacesForUser(namespaces []string) (map[string]*Project, error) {
	insertMarks := []string{}
	for i, _ := range namespaces {
		if s.UseNumericInsertMarks {
			insertMarks = append(insertMarks, "$"+strconv.Itoa(i+1))
		} else {
			insertMarks = append(insertMarks, "?")
		}
	}
	query := s.GetNamespacesForUserQuery
	if len(namespaces) == 0 {
		return map[string]*Project{}, nil
	} else if len(namespaces) == 1 {
		query += " = " + insertMarks[0]
	} else {
		query += " IN (" + strings.Join(insertMarks, ", ") + ")"
	}
	interfaceNamespaces := []interface{}{}
	for _, n := range namespaces {
		interfaceNamespaces = append(interfaceNamespaces, n)
	}
	rows, err := s.PrepareAndQuery(query, interfaceNamespaces...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := map[string]*Project{}
	for rows.Next() {
		var name, description, orgURL, logo string
		var isPublic bool
		if err := rows.Scan(&name, &description, &orgURL, &logo, &isPublic); err != nil {
			return nil, err
		}
		prj, ok := result[name]
		if !ok {
			prj = &Project{
				Name:        name,
				Description: description,
				OrgURL:      orgURL,
				Logo:        logo,
				Permission:  "admin", // default permission for open source
				IsPublic:    isPublic,
			}
		}
		result[prj.Name] = prj
	}
	return result, nil
	return nil, nil
}

func (s *SQLHelper) GetNamespacesFilteredBy(f *NamespacesFilter) (map[string]*Project, error) {
	insertMarks := []string{}
	for i, _ := range f.Namespaces {
		if s.UseNumericInsertMarks {
			insertMarks = append(insertMarks, "$"+strconv.Itoa(i+1))
		} else {
			insertMarks = append(insertMarks, "?")
		}
	}
	query := s.GetNamespacesByNamesQuery
	if len(f.Namespaces) == 0 {
		return map[string]*Project{}, nil
	} else if len(f.Namespaces) == 1 {
		query += " = " + insertMarks[0]
	} else {
		query += " IN (" + strings.Join(insertMarks, ", ") + ")"
	}
	interfaceNamespaces := []interface{}{}
	for _, n := range f.Namespaces {
		interfaceNamespaces = append(interfaceNamespaces, n)
	}
	rows, err := s.PrepareAndQuery(query, interfaceNamespaces...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := map[string]*Project{}
	for rows.Next() {
		var name, description, orgURL, logo string
		var isPublic bool
		if err := rows.Scan(&name, &description, &orgURL, &logo, &isPublic); err != nil {
			return nil, err
		}
		prj, ok := result[name]
		if !ok {
			prj = &Project{
				Name:        name,
				Description: description,
				OrgURL:      orgURL,
				Logo:        logo,
				IsPublic:    isPublic,
				Permission:  "admin", // default permission for open source
			}
		}
		result[prj.Name] = prj
	}
	return result, nil
}

func (s *SQLHelper) HardDeleteNamespace(namespace string) error {
	if err := s.PrepareAndExec(s.HardDeleteProjectUnitSubscriptions, namespace); err != nil {
		return err
	}
	if err := s.PrepareAndExec(s.HardDeleteProjectReleaseDependenciesQuery, namespace); err != nil {
		return err
	}
	if err := s.PrepareAndExec(s.HardDeleteProjectPackageURIsQuery, namespace); err != nil {
		return err
	}
	if err := s.PrepareAndExec(s.HardDeleteProjectReleasesQuery, namespace); err != nil {
		return err
	}
	if err := s.PrepareAndExec(s.HardDeleteProjectApplicationsQuery, namespace); err != nil {
		return err
	}
	if err := s.PrepareAndExec(s.HardDeleteProjectQuery, namespace); err != nil {
		return err
	}
	return nil
}

func (s *SQLHelper) scanNamespace(rows *sql.Rows) (*Project, error) {
	var name, description, orgURL, logo string
	var isPublic bool
	if err := rows.Scan(&name, &description, &orgURL, &logo, &isPublic); err != nil {
		return nil, err
	}
	return &Project{
		Name:        name,
		Description: description,
		OrgURL:      orgURL,
		Logo:        logo,
		Permission:  "admin", // default permission for open source
		IsPublic:    isPublic,
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
