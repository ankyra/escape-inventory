/*
Copyright 2017, 2018 Ankyra

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

type ApplicationPayload struct {
	*types.Application
	Versions []string `json:"versions",omitempty`
}

func GetApplication(namespace, name string) (*ApplicationPayload, error) {
	_, err := dao.GetProject(namespace)
	if err != nil {
		return nil, err
	}
	app, err := dao.GetApplication(namespace, name)
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

func GetApplications(namespace string) (map[string]*types.Application, error) {
	_, err := dao.GetProject(namespace)
	if err != nil {
		return nil, err
	}
	return dao.GetApplications(namespace)
}

func GetApplicationVersions(namespace, name string) ([]string, error) {
	app, err := dao.GetApplication(namespace, name)
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

func GetApplicationHooks(namespace, name string) (types.Hooks, error) {
	app, err := dao.GetApplication(namespace, name)
	if err != nil {
		return nil, err
	}
	return dao.GetApplicationHooks(app)
}

func UpdateApplicationHooks(namespace, name string, hooks types.Hooks) error {
	app, err := dao.GetApplication(namespace, name)
	if err != nil {
		return err
	}
	currentHooks, err := dao.GetApplicationHooks(app)
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
		case "build":
			newValue, err := parseEventsHookConfig(value)
			if err != nil {
				return err
			}
			currentHooks[key] = newValue
		default:
			return NewUserError(fmt.Errorf("Unknown hook type '%s'", key))
		}
	}
	return dao.SetApplicationHooks(app, hooks)
}

func parseBuildHookConfig(values map[string]string) (map[string]string, error) {
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
