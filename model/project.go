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
	"fmt"

	"github.com/ankyra/escape-registry/dao"
	"github.com/ankyra/escape-registry/dao/types"
)

type ProjectPayload struct {
	*types.Project
	Units []string `json:"units"`
}

func GetProject(project string) (*ProjectPayload, error) {
	prj, err := dao.GetProject(project)
	if err != nil {
		return nil, err
	}
	units, err := GetApplications(project)
	if err != nil && !dao.IsNotFound(err) {
		return nil, err
	}
	return &ProjectPayload{
		prj,
		units,
	}, nil
}

func AddProject(p *types.Project) error {
	if p.Name == "" {
		return NewUserError(fmt.Errorf("Missing name"))
	}
	return dao.AddProject(p)
}

func UpdateProject(p *types.Project) error {
	if p.Name == "" {
		return NewUserError(fmt.Errorf("Missing name"))
	}
	return dao.UpdateProject(p)
}
