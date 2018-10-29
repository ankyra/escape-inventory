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

package mem

import (
	. "github.com/ankyra/escape-inventory/dao/types"
)

type application struct {
	App      *Application
	Releases map[string]*release
}

type release struct {
	Release      *Release
	Packages     []string
	Dependencies []*Dependency
}

type dao struct {
	namespaceMetadata map[string]*Project
	namespaceHooks    map[*Project]Hooks
	namespaces        map[string]map[string]*application
	apps              map[*Application]*application
	applicationHooks  map[*Application]Hooks
	subscriptions     map[*Application][]*Application
	releases          map[*Release]*release
	acls              map[string]map[string]Permission
	metrics           map[string]*Metrics
	providers         map[string]map[string]*MinimalReleaseMetadata
}

func NewInMemoryDAO() DAO {
	return &dao{
		namespaceMetadata: map[string]*Project{},
		namespaceHooks:    map[*Project]Hooks{},
		namespaces:        map[string]map[string]*application{},
		apps:              map[*Application]*application{},
		applicationHooks:  map[*Application]Hooks{},
		subscriptions:     map[*Application][]*Application{},
		releases:          map[*Release]*release{},
		acls:              map[string]map[string]Permission{},
		metrics:           map[string]*Metrics{},
		providers:         map[string]map[string]*MinimalReleaseMetadata{},
	}
}

func (a *dao) WipeDatabase() error {
	a.namespaceMetadata = map[string]*Project{}
	a.namespaceHooks = map[*Project]Hooks{}
	a.namespaces = map[string]map[string]*application{}
	a.apps = map[*Application]*application{}
	a.applicationHooks = map[*Application]Hooks{}
	a.subscriptions = map[*Application][]*Application{}
	a.releases = map[*Release]*release{}
	a.acls = map[string]map[string]Permission{}

	return nil
}
