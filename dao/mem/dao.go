/*
Copyright 2017 Ankyra

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package mem

import (
	"github.com/ankyra/escape-core"
	. "github.com/ankyra/escape-registry/dao/types"
)

type application struct {
	App      *Application
	Releases map[string]*release
}

type release struct {
	Release  *Release
	Packages []string
}

type dao struct {
	projectMetadata map[string]*Project
	projects        map[string]map[string]*application
	apps            map[*Application]*application
	releases        map[*Release]*release
	acls            map[string]map[string]Permission
}

func NewInMemoryDAO() DAO {
	return &dao{
		projectMetadata: map[string]*Project{},
		projects:        map[string]map[string]*application{},
		apps:            map[*Application]*application{},
		releases:        map[*Release]*release{},
		acls:            map[string]map[string]Permission{},
	}
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
	a.projectMetadata[project.Name] = project
	return nil
}

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

func (a *dao) GetApplications(project string) ([]*Application, error) {
	result := []*Application{}
	for _, app := range a.projects[project] {
		result = append(result, app.App)
	}
	return result, nil
}

func (a *dao) GetApplication(project, name string) (*Application, error) {
	prj, ok := a.projects[project]
	if !ok {
		return nil, NotFound
	}
	result, ok := prj[name]
	if !ok {
		return nil, NotFound
	}
	return result.App, nil
}

func (a *dao) FindAllVersions(app *Application) ([]string, error) {
	application := a.apps[app]
	versions := []string{}
	if application == nil {
		return versions, nil
	}
	for _, r := range application.Releases {
		versions = append(versions, r.Release.Version)
	}
	return versions, nil
}

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
		unwrapped, err := a.GetApplication(project, rel.Name)
		if err != nil {
			unwrapped = NewApplication(project, rel.Name)
		}
		app = &application{
			App:      unwrapped,
			Releases: map[string]*release{},
		}
		a.apps[unwrapped] = app
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
