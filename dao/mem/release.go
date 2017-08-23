package mem

import (
	core "github.com/ankyra/escape-core"
	. "github.com/ankyra/escape-registry/dao/types"
)

func (a *dao) GetRelease(project, name, releaseId string) (*Release, error) {
	prj, ok := a.projects[project]
	if !ok {
		return nil, NotFound
	}
	app, ok := prj[name]
	if !ok {
		return nil, NotFound
	}
	release, ok := app.Releases[releaseId]
	if !ok {
		return nil, NotFound
	}
	return release.Release, nil
}

func (a *dao) AddRelease(project string, rel *core.ReleaseMetadata) (*Release, error) {
	apps, ok := a.projects[project]
	if !ok {
		apps = map[string]*application{}
	}
	key := rel.GetReleaseId()
	app, ok := apps[rel.Name]
	if !ok {
		return nil, NotFound
	}
	_, alreadyExists := app.Releases[key]
	if alreadyExists {
		return nil, AlreadyExists
	}
	newRelease := NewRelease(app.App, rel)
	app.Releases[key] = &release{
		Release:  newRelease,
		Packages: []string{},
	}
	apps[rel.Name] = app
	a.projects[project] = apps
	a.releases[newRelease] = app.Releases[key]
	return app.Releases[key].Release, nil
}

func (a *dao) GetAllReleases() ([]*Release, error) {
	result := []*Release{}
	for _, rel := range a.releases {
		result = append(result, rel.Release)
	}
	return result, nil
}

func (a *dao) GetPackageURIs(release *Release) ([]string, error) {
	r := a.releases[release]
	return r.Packages, nil
}

func (a *dao) AddPackageURI(release *Release, uri string) error {
	r := a.releases[release]
	for _, u := range r.Packages {
		if u == uri {
			return AlreadyExists
		}
	}
	r.Packages = append(r.Packages, uri)
	return nil
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
