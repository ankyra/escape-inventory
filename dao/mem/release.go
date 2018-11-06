package mem

import (
	. "github.com/ankyra/escape-inventory/dao/types"
)

func (a *dao) GetRelease(namespace, name, releaseId string) (*Release, error) {
	prj, ok := a.namespaces[namespace]
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

func (a *dao) GetReleaseByTag(namespaces, name, tag string) (*Release, error) {
	prj, ok := a.namespaces[namespace]
	if !ok {
		return nil, NotFound
	}
	app, ok := prj[name]
	if !ok {
		return nil, NotFound
	}
	release, ok := app.Tags[tag]
	if !ok {
		return nil, NotFound
	}
	return release.Release, nil
}

func (a *dao) TagRelease(rel *Release, tag string) error {
	prj, ok := a.namespaces[rel.Application.Project]
	if !ok {
		return NotFound
	}
	app, ok := prj[rel.Application.Name]
	if !ok {
		return NotFound
	}
	release, ok := app.Releases[rel.ReleaseId]
	if !ok {
		return NotFound
	}
	app.Tags[tag] = release
	return nil
}

func (a *dao) AddRelease(rel *Release) error {
	apps, ok := a.namespaces[rel.Application.Project]
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
	a.namespaces[rel.Application.Project] = apps
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
	r, ok := a.releases[release]
	if !ok {
		return []string{}, nil
	}
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
