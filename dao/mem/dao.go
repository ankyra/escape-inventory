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
	applications []ApplicationDAO
	releases     map[string]ReleaseDAO
}

func NewInMemoryDAO() DAO {
	return &mem_dao{
		applications: []ApplicationDAO{},
		releases:     map[string]ReleaseDAO{},
	}
}

func (a *mem_dao) GetApplications(project string) ([]ApplicationDAO, error) {
	return a.applications, nil
}

func (a *mem_dao) GetApplication(project, name string) (ApplicationDAO, error) {
	for _, app := range a.applications {
		if app.GetName() == name {
			return app, nil
		}
	}
	return nil, NotFound
}

func (a *mem_dao) GetRelease(project, releaseId string) (ReleaseDAO, error) {
	release, ok := a.releases[releaseId]
	if !ok {
		return nil, NotFound
	}
	return release, nil
}

func (a *mem_dao) AddRelease(project string, release *core.ReleaseMetadata) (ReleaseDAO, error) {
	key := release.GetReleaseId()
	_, alreadyExists := a.releases[key]
	if alreadyExists {
		return nil, AlreadyExists
	}
	var application *mem_application
	name := release.GetName()
	for _, app := range a.applications {
		if app.GetName() == name {
			application = app.(*mem_application)
		}
	}
	if application == nil {
		application = newApplication(name, a)
		a.applications = append(a.applications, application)
	}
	a.releases[key] = newRelease(release, application)
	application.releases[key] = a.releases[key]
	return a.releases[key], nil
}

func (a *mem_dao) GetAllReleases() ([]ReleaseDAO, error) {
	result := []ReleaseDAO{}
	for _, rel := range a.releases {
		result = append(result, rel)
	}
	return result, nil
}
