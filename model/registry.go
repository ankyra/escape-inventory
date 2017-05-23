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

func Registry(project, name string) ([]string, error) {
	if name == "" {
		result := []string{}
		apps, err := dao.GetApplications(project)
		if err != nil {
			return nil, err
		}
		for _, app := range apps {
			result = append(result, app.GetName())
		}
		if len(result) == 0 {
			return nil, types.NotFound
		}
		return result, nil
	}
	app, err := dao.GetApplication(project, name)
	if err != nil {
		return nil, err
	}
	result, err := app.FindAllVersions()
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, types.NotFound
	}
	return result, nil
}
