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

package types

import (
	"fmt"
)

type Hooks map[string]map[string]string

func NewHooks() Hooks {
	return map[string]map[string]string{}
}

type DAO interface {
	NamespacesDAO
	ApplicationsDAO
	ReleasesDAO
	DependenciesDAO
	MetricsDAO

	WipeDatabase() error
}

var NotFound = fmt.Errorf("Not found")
var AlreadyExists = fmt.Errorf("Already exists")
var LimitError = fmt.Errorf("Plan limit exceeded")
var Unauthorized = fmt.Errorf("Unauthorized")
