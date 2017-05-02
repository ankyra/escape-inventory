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
	. "github.com/ankyra/escape-registry/dao/types"
)

func Registry(typ, name string) ([]string, error) {
	types, err := dao.GetReleaseTypes()
	if err != nil {
		return nil, err
	}
	if typ == "" {
		return types, nil
	}
	typeFound := false
	for _, t := range types {
		if t == typ {
			typeFound = true
		}
	}
	if !typeFound {
		return nil, NotFound
	}
	if name == "" {
		return dao.GetApplicationsByType(typ)
	}
	app, err := dao.GetApplication(typ, name)
	if err != nil {
		return nil, err
	}
	return app.FindAllVersions()
}
