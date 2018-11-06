package mem

import (
	. "github.com/ankyra/escape-inventory/dao/types"
)

func (a *dao) GetAllReleasesWithoutProcessedDependencies() ([]*Release, error) {
	result := []*Release{}
	for _, rel := range a.releases {
		if !rel.Release.ProcessedDependencies {
			result = append(result, rel.Release)
		}
	}
	return result, nil
}

func (a *dao) SetDependencies(release *Release, depends []*Dependency) error {
	r := a.releases[release]
	r.Dependencies = depends
	return nil
}

func (a *dao) GetDependencies(release *Release) ([]*Dependency, error) {
	r, ok := a.releases[release]
	if !ok {
		return []*Dependency{}, nil
	}
	return r.Dependencies, nil
}

func (a *dao) SetDependencyTree(release *Release, depends []*DependencyTree) error {
	return nil
}
func (a *dao) GetDependencyTree(release *Release) ([]*DependencyTree, error) {
	return nil, nil
}

func (a *dao) GetDownstreamDependencies(release *Release) ([]*Dependency, error) {
	project := release.Application.Project
	app := release.Application.Name
	version := release.Version
	result := []*Dependency{}
	for r, rel := range a.releases {
		found := false
		buildScope := false
		deployScope := false
		for _, dep := range rel.Dependencies {
			if dep.Project == project && dep.Application == app && dep.Version == version {
				found = true
				buildScope = buildScope || dep.BuildScope
				deployScope = deployScope || dep.DeployScope
			}
		}
		if found {
			d := Dependency{
				Project:     r.Application.Project,
				Application: r.Application.Name,
				Version:     r.Version,
				BuildScope:  buildScope,
				DeployScope: deployScope,
			}
			result = append(result, &d)
		}
	}
	return result, nil
}

func (a *dao) GetDownstreamDependenciesFilteredBy(release *Release, query *DownstreamDependenciesFilter) ([]*Dependency, error) {
	deps, err := a.GetDownstreamDependencies(release)
	if err != nil {
		return nil, err
	}
	result := []*Dependency{}
	for _, dep := range deps {
		for _, namespace := range query.Namespaces {
			if namespace == dep.Project {
				result = append(result, dep)
				break
			}
		}
	}
	return result, nil
}
