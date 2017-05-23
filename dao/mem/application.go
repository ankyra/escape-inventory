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
	. "github.com/ankyra/escape-registry/dao/types"
)

type mem_application struct {
	dao      *mem_dao
	project  string
	name     string
	releases map[string]ReleaseDAO
}

func newApplication(project, name string, dao *mem_dao) *mem_application {
	return &mem_application{
		dao:      dao,
		project:  project,
		name:     name,
		releases: map[string]ReleaseDAO{},
	}
}

func (a *mem_application) GetName() string {
	return a.name
}

func (a *mem_application) FindAllVersions() ([]string, error) {
	versions := []string{}
	for _, r := range a.releases {
		if r.GetApplication() == a {
			versions = append(versions, r.GetVersion())
		}
	}
	return versions, nil
}
