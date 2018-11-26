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

type NamespacesDAO interface {
	GetNamespace(namespace string) (*Project, error)
	AddNamespace(*Project) error
	HardDeleteNamespace(namespace string) error
	UpdateNamespace(*Project) error
	GetNamespaces() (map[string]*Project, error)
	GetNamespacesByNames(namespaces []string) (map[string]*Project, error)
	GetNamespacesForUser(namespaces []string) (map[string]*Project, error)
	GetNamespacesFilteredBy(*NamespacesFilter) (map[string]*Project, error)
	GetNamespaceHooks(*Project) (Hooks, error)
	SetNamespaceHooks(*Project, Hooks) error
}

type Project struct {
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	OrgURL         string   `json:"org_url"`
	Logo           string   `json:"logo"`
	Hooks          Hooks    `json:"hooks,omitempty"` // only used for view purposes
	Permission     string   `json:"permission"`      // only used for view purposes
	MatchingGroups []string `json:"-"`               // used to work out highest permission in model
	IsPublic       bool     `json:"is_public"`
}

func NewProject(project string) *Project {
	return &Project{
		Name: project,
	}
}

type NamespacesFilter struct {
	Namespaces []string
}

func NewNamespacesFilter(namespaces []string) *NamespacesFilter {
	return &NamespacesFilter{
		Namespaces: namespaces,
	}
}
