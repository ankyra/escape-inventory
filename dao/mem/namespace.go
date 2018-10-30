package mem

import (
	. "github.com/ankyra/escape-inventory/dao/types"
)

func (a *dao) UpdateNamespace(project *Project) error {
	_, exists := a.namespaceMetadata[project.Name]
	if !exists {
		return NotFound
	}
	a.namespaceMetadata[project.Name] = project
	return nil
}

func (a *dao) GetNamespaces() (map[string]*Project, error) {
	return a.namespaceMetadata, nil
}

func (a *dao) GetNamespacesByGroups(readGroups []string) (map[string]*Project, error) {
	result := map[string]*Project{}
	for name, namespaceMetadata := range a.namespaceMetadata {
		allowedGroups, found := a.acls[name]
		if found {
			namespaceMetadata.MatchingGroups = []string{}
			matchedGroups := map[string]bool{}
			for _, g := range readGroups {
				_, found := allowedGroups[g]
				if found {
					matchedGroups[g] = true
				}
			}
			_, found := allowedGroups["*"]
			if found {
				matchedGroups["*"] = true
			}
			if len(matchedGroups) > 0 {
				for key, _ := range matchedGroups {
					namespaceMetadata.MatchingGroups = append(namespaceMetadata.MatchingGroups, key)
				}
				result[name] = namespaceMetadata
			}
		}
	}
	return result, nil
}

func (a *dao) GetNamespacesByNames(namespaces []string) (map[string]*Project, error) {
	namespacesFound := map[string]*Project{}
	for _, name := range namespaces {
		namespace, ok := a.namespaceMetadata[name]
		if ok {
			namespacesFound[name] = namespace
		}
	}
	return namespacesFound, nil
}

func (a *dao) GetNamespace(namespace string) (*Project, error) {
	prj, ok := a.namespaceMetadata[namespace]
	if !ok {
		return nil, NotFound
	}
	return prj, nil
}

func (a *dao) AddNamespace(namespace *Project) error {
	namespace.Permission = "admin"
	_, exists := a.namespaceMetadata[namespace.Name]
	if exists {
		return AlreadyExists
	}
	_, ok := a.namespaces[namespace.Name]
	if !ok {
		a.namespaces[namespace.Name] = map[string]*application{}
	}
	a.namespaceMetadata[namespace.Name] = namespace
	a.namespaceHooks[namespace] = NewHooks()
	return nil
}

func (a *dao) GetNamespaceHooks(namespace *Project) (Hooks, error) {
	namespace, ok := a.namespaceMetadata[namespace.Name]
	if !ok {
		return nil, NotFound
	}
	return a.namespaceHooks[namespace], nil
}

func (a *dao) SetNamespaceHooks(namespace *Project, hooks Hooks) error {
	namespace, ok := a.namespaceMetadata[namespace.Name]
	if !ok {
		return NotFound
	}
	a.namespaceHooks[namespace] = hooks
	return nil
}

func (a *dao) HardDeleteNamespace(namespace string) error {
	namespaceMetadata, exists := a.namespaceMetadata[namespace]
	if !exists {
		return NotFound
	}
	toDelete := []*Application{}
	for app, _ := range a.subscriptions {
		if app.Project == namespace {
			toDelete = append(toDelete, app)
		}
	}
	for _, i := range toDelete {
		delete(a.subscriptions, i)
	}
	for _, app := range a.namespaces[namespace] {
		delete(a.applicationHooks, app.App)
		for _, rel := range a.apps[app.App].Releases {
			delete(a.releases, rel.Release)
		}
		delete(a.apps, app.App)
	}
	delete(a.namespaceMetadata, namespace)
	delete(a.namespaceHooks, namespaceMetadata)
	delete(a.namespaces, namespace)
	delete(a.acls, namespace)

	return nil
}
