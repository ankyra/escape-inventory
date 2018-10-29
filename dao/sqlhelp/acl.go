package sqlhelp

import (
	. "github.com/ankyra/escape-inventory/dao/types"
)

func (s *SQLHelper) SetACL(project, group string, perm Permission) error {
	err := s.PrepareAndExecInsert(s.InsertACLQuery,
		namespace,
		group,
		int(perm))
	if err == AlreadyExists {
		return s.PrepareAndExecUpdate(s.UpdateACLQuery,
			int(perm), namespace, group)
	}
	return err
}

func (s *SQLHelper) GetACL(namespace string) (map[string]Permission, error) {
	rows, err := s.PrepareAndQuery(s.GetACLQuery, namespace)
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

func (s *SQLHelper) DeleteACL(namespace, group string) error {
	return s.PrepareAndExec(s.DeleteACLQuery,
		namespace,
		group)
}

func (s *SQLHelper) GetPermittedGroups(namespace string, perm Permission) ([]string, error) {
	rows, err := s.PrepareAndQuery(s.GetPermittedGroupsQuery, namespace, int(perm))
	if err != nil {
		return nil, err
	}
	return s.ReadRowsIntoStringArray(rows)
}
