package mem

import (
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

func (a *dao) AddRelease(rel *Release) error {
	apps, ok := a.projects[rel.Application.Project]
	if !ok {
		apps = map[string]*application{}
	}
	key := rel.Metadata.GetReleaseId()
	app, ok := apps[rel.Application.Name]
	if !ok {
		return NotFound
	}
	_, alreadyExists := app.Releases[key]
	if alreadyExists {
		return AlreadyExists
	}
	app.Releases[key] = &release{
		Release:  rel,
		Packages: []string{},
	}
	apps[rel.Application.Name] = app
	a.projects[rel.Application.Project] = apps
	a.releases[rel] = app.Releases[key]
	return nil
}

func (a *dao) UpdateRelease(r *Release) error {
	return nil
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
