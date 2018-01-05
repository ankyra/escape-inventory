package mem

import (
	. "github.com/ankyra/escape-inventory/dao/types"
)

func (a *dao) UpdateProject(project *Project) error {
	_, exists := a.projectMetadata[project.Name]
	if !exists {
		return NotFound
	}
	a.projectMetadata[project.Name] = project
	return nil
}

func (a *dao) GetProjects() (map[string]*Project, error) {
	return a.projectMetadata, nil
}

func (a *dao) GetProjectsByGroups(readGroups []string) (map[string]*Project, error) {
	result := map[string]*Project{}
	for name, project := range a.projectMetadata {
		allowedGroups, found := a.acls[name]
		if found {
			project.MatchingGroups = []string{}
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
					project.MatchingGroups = append(project.MatchingGroups, key)
				}
				result[name] = project
			}
		}
	}
	return result, nil
}

func (a *dao) GetProject(project string) (*Project, error) {
	prj, ok := a.projectMetadata[project]
	if !ok {
		return nil, NotFound
	}
	return prj, nil
}

func (a *dao) AddProject(project *Project) error {
	_, exists := a.projectMetadata[project.Name]
	if exists {
		return AlreadyExists
	}
	_, ok := a.projects[project.Name]
	if !ok {
		a.projects[project.Name] = map[string]*application{}
	}
	a.projectMetadata[project.Name] = project
	a.projectHooks[project] = NewHooks()
	return nil
}

func (a *dao) GetProjectHooks(project *Project) (Hooks, error) {
	project, ok := a.projectMetadata[project.Name]
	if !ok {
		return nil, NotFound
	}
	return a.projectHooks[project], nil
}

func (a *dao) SetProjectHooks(project *Project, hooks Hooks) error {
	project, ok := a.projectMetadata[project.Name]
	if !ok {
		return NotFound
	}
	a.projectHooks[project] = hooks
	return nil
}
