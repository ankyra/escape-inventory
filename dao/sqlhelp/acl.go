package sqlhelp

import (
	. "github.com/ankyra/escape-inventory/dao/types"
	"github.com/ankyra/escape-middleware/errors"
)

func (s *SQLHelper) SetACL(project, group string, perm Permission) error {
	err := s.PrepareAndExecInsert(s.InsertACLQuery,
		project,
		group,
		int(perm))
	if err == errors.AlreadyExists {
		return s.PrepareAndExecUpdate(s.UpdateACLQuery,
			int(perm), project, group)
	}
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
	return s.PrepareAndExec(s.DeleteACLQuery,
		project,
		group)
}

func (s *SQLHelper) GetPermittedGroups(project string, perm Permission) ([]string, error) {
	rows, err := s.PrepareAndQuery(s.GetPermittedGroupsQuery, project, int(perm))
	if err != nil {
		return nil, err
	}
	return s.ReadRowsIntoStringArray(rows)
}
