package sqlhelp

import (
	. "github.com/ankyra/escape-registry/dao/types"
)

func (s *SQLHelper) SetACL(project, group string, perm Permission) error {
	stmt, err := s.DB.Prepare(s.InsertACLQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(project, group, int(perm))
	if err != nil {
		if s.IsUniqueConstraintError(err) {
			stmt, err := s.DB.Prepare(s.UpdateACLQuery)
			if err != nil {
				return err
			}
			_, err = stmt.Exec(int(perm), project, group)
			return err
		}
		return err
	}
	return nil
}

func (s *SQLHelper) GetACL(project string) (map[string]Permission, error) {
	rows, err := s.PrepareAndQuery(s.GetACLQuery, project)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := map[string]Permission{}
	for rows.Next() {
		var group_name string
		var permission Permission
		if err := rows.Scan(&group_name, &permission); err != nil {
			return nil, err
		}
		result[group_name] = permission
	}
	return result, nil
}

func (s *SQLHelper) DeleteACL(project, group string) error {
	stmt, err := s.DB.Prepare(s.DeleteACLQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(project, group)
	return err
}

func (s *SQLHelper) GetPermittedGroups(project string, perm Permission) ([]string, error) {
	rows, err := s.PrepareAndQuery(s.GetPermittedGroupsQuery, project, int(perm))
	if err != nil {
		return nil, err
	}
	return s.ReadRowsIntoStringArray(rows)
}
