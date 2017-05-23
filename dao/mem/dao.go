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

type mem_dao struct {
	projects map[string]map[string]ApplicationDAO
}

func NewInMemoryDAO() DAO {
	return &mem_dao{
		projects: map[string]map[string]ApplicationDAO{},
	}
}

func (a *mem_dao) GetApplications(project string) ([]ApplicationDAO, error) {
	result := []ApplicationDAO{}
	for _, app := range a.projects[project] {
		result = append(result, app)
	}
	return result, nil
}

func (a *mem_dao) GetApplication(project, name string) (ApplicationDAO, error) {
	prj, ok := a.projects[project]
	if !ok {
		return nil, NotFound
	}
	result, ok := prj[name]
	if !ok {
		return nil, NotFound
	}
	return result, nil
}

func (a *mem_dao) GetRelease(project, name, releaseId string) (ReleaseDAO, error) {
	prj, ok := a.projects[project]
	if !ok {
		return nil, NotFound
	}
	app, ok := prj[name]
	if !ok {
		return nil, NotFound
	}
	release, ok := app.(*mem_application).releases[releaseId]
	if !ok {
		return nil, NotFound
	}
	return release, nil
}

func (a *mem_dao) AddRelease(project string, release *core.ReleaseMetadata) (ReleaseDAO, error) {
	apps, ok := a.projects[project]
	if !ok {
		apps = map[string]ApplicationDAO{}
	}
	key := release.GetReleaseId()
	app, ok := apps[release.GetName()]
	if !ok {
		app = newApplication(project, release.GetName(), a)
	}
	application := app.(*mem_application)
	_, alreadyExists := application.releases[key]
	if alreadyExists {
		return nil, AlreadyExists
	}
	application.releases[key] = newRelease(release, application)
	apps[release.GetName()] = app
	a.projects[project] = apps
	return application.releases[key], nil
}

func (a *mem_dao) GetAllReleases() ([]ReleaseDAO, error) {
	result := []ReleaseDAO{}
	for _, prj := range a.projects {
		for _, app := range prj {
			for _, rel := range app.(*mem_application).releases {
				result = append(result, rel)
			}
		}
	}
	return result, nil
}
