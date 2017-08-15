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

package model

import (
	"github.com/ankyra/escape-registry/dao"
	"github.com/ankyra/escape-registry/dao/types"
)

type ApplicationPayload struct {
	*types.Application
	Versions []string `json:"versions",omitempty`
}

func GetApplication(project, name string) (*ApplicationPayload, error) {
	_, err := dao.GetProject(project)
	if err != nil {
		return nil, err
	}
	app, err := dao.GetApplication(project, name)
	if err != nil {
		return nil, err
	}
	versions, err := dao.FindAllVersions(app)
	if err != nil {
		return nil, err
	}
	return &ApplicationPayload{
		app,
		versions,
	}, nil
}

func GetApplications(project string) (map[string]*types.Application, error) {
	_, err := dao.GetProject(project)
	if err != nil {
		return nil, err
	}
	return dao.GetApplications(project)
}

func GetApplicationVersions(project, name string) ([]string, error) {
	app, err := dao.GetApplication(project, name)
	if err != nil {
		return nil, err
	}
	result, err := dao.FindAllVersions(app)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, types.NotFound
	}
	return result, nil
}
