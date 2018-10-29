package mem

import (
	. "github.com/ankyra/escape-inventory/dao/types"
)

func (a *dao) SetACL(namespace, group string, perm Permission) error {
	groups, ok := a.acls[namespace]
	if !ok {
		groups = map[string]Permission{}
	}
	groups[group] = perm
	a.acls[namespace] = groups
	return nil
}

func (a *dao) GetACL(namespace string) (map[string]Permission, error) {
	groups, ok := a.acls[namespace]
	if !ok {
		groups = map[string]Permission{}
	}
	return groups, nil
}

func (a *dao) DeleteACL(namespace, group string) error {
	groups, ok := a.acls[namespace]
	if !ok {
		return nil
	}
	delete(groups, group)
	return nil
}
func (a *dao) GetPermittedGroups(namespace string, perm Permission) ([]string, error) {
	result := []string{}
	groups, ok := a.acls[namespace]
	if !ok {
		return result, nil
	}
	for group, groupPerm := range groups {
		if perm <= groupPerm {
			result = append(result, group)
		}
	}
	return result, nil
}
