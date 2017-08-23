package mem

import (
	. "github.com/ankyra/escape-registry/dao/types"
)

func (a *dao) SetACL(project, group string, perm Permission) error {
	groups, ok := a.acls[project]
	if !ok {
		groups = map[string]Permission{}
	}
	groups[group] = perm
	a.acls[project] = groups
	return nil
}

func (a *dao) GetACL(project string) (map[string]Permission, error) {
	groups, ok := a.acls[project]
	if !ok {
		groups = map[string]Permission{}
	}
	return groups, nil
}

func (a *dao) DeleteACL(project, group string) error {
	groups, ok := a.acls[project]
	if !ok {
		return nil
	}
	delete(groups, group)
	return nil
}
func (a *dao) GetPermittedGroups(project string, perm Permission) ([]string, error) {
	result := []string{}
	groups, ok := a.acls[project]
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
