package mem

import (
	. "github.com/ankyra/escape-registry/dao/types"
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
	r := a.releases[release]
	return r.Dependencies, nil
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

func (a *dao) GetDownstreamDependenciesByGroups(release *Release, readGroups []string) ([]*Dependency, error) {
	project := release.Application.Project
	app := release.Application.Name
	version := release.Version
	result := []*Dependency{}
	for r, rel := range a.releases {
		allowedGroups := a.acls[r.Application.Project]
		found := false
		for _, g := range readGroups {
			_, found = allowedGroups[g]
			if found {
				break
			}
		}
		if !found {
			_, found = allowedGroups["*"]
		}
		if !found {
			continue
		}
		found = false
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
