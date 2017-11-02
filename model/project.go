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

	"github.com/ankyra/escape-inventory/dao"
	"github.com/ankyra/escape-inventory/dao/types"
)

type ProjectPayload struct {
	*types.Project
	Units map[string]*types.Application `json:"units"`
}

func GetProject(project string) (*ProjectPayload, error) {
	prj, err := dao.GetProject(project)
	if err != nil {
		return nil, err
	}
	units, err := dao.GetApplications(project)
	if err != nil && !dao.IsNotFound(err) {
		return nil, err
	}
	return &ProjectPayload{
		prj,
		units,
	}, nil
}

func GetProjectHooks(project string) (types.Hooks, error) {
	prj, err := dao.GetProject(project)
	if err != nil {
		return nil, err
	}
	return dao.GetProjectHooks(prj)
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

func UpdateProjectHooks(project string, hooks types.Hooks) error {
	prj, err := dao.GetProject(project)
	if err != nil {
		return err
	}
	currentHooks, err := dao.GetProjectHooks(prj)
	if err != nil {
		return err
	}
	for key, value := range hooks {
		switch key {
		case "slack":
			newValue, err := parseSlackHookConfig(value)
			if err != nil {
				return err
			}
			currentHooks[key] = newValue
		case "events":
			newValue, err := parseEventsHookConfig(value)
			if err != nil {
				return err
			}
			currentHooks[key] = newValue
		default:
			return NewUserError(fmt.Errorf("Unknown hook type '%s'", key))
		}
	}
	return dao.SetProjectHooks(prj, hooks)
}

func parseSlackHookConfig(values map[string]string) (map[string]string, error) {
	result := map[string]string{}
	for key, value := range values {
		switch key {
		case "username", "icon_emoji", "url":
			result[key] = value
		default:
			return nil, NewUserError(fmt.Errorf("Unknown  hook configuration key '%s'", key))
		}
	}
	return result, nil
}

func parseEventsHookConfig(values map[string]string) (map[string]string, error) {
	result := map[string]string{}
	for key, value := range values {
		switch key {
		case "url":
			result[key] = value
		default:
			return nil, NewUserError(fmt.Errorf("Unknown hook configuration key '%s'", key))
		}
	}
	return result, nil
}
