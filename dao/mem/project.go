package mem

import (
	. "github.com/ankyra/escape-registry/dao/types"
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
			for _, g := range readGroups {
				_, found := allowedGroups[g]
				if found {
					result[name] = project
					break
				}
			}
			_, found := allowedGroups["*"]
			if found {
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
	return nil
}
