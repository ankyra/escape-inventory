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

type DependenciesDAO interface {
	GetAllReleasesWithoutProcessedDependencies() ([]*Release, error)
	SetDependencies(*Release, []*Dependency) error
	GetDependencies(*Release) ([]*Dependency, error)
	GetDownstreamDependencies(*Release) ([]*Dependency, error)
	GetDownstreamDependenciesFilteredBy(*Release, *DownstreamDependenciesFilter) ([]*Dependency, error)
}

type Dependency struct {
	Project     string `json:"project"`
	Application string `json:"name"`
	Version     string `json:"version"`
	BuildScope  bool   `json:"build"`
	DeployScope bool   `json:"deploy"`
	IsExtension bool   `json:"is_extension"`
}

func NewDependency(project, name, version string) *Dependency {
	return &Dependency{
		Project:     project,
		Application: name,
		Version:     version,
	}
}

type DownstreamDependenciesFilter struct {
	Namespaces []string
}

func NewDownstreamDependenciesFilter(namespaces []string) *DownstreamDependenciesFilter {
	return &DownstreamDependenciesFilter{
		Namespaces: namespaces,
	}
}
