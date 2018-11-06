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
	"time"

	core "github.com/ankyra/escape-core"
)

type ReleasesDAO interface {
	GetRelease(namespace, name, releaseId string) (*Release, error)
	GetReleaseByTag(namespace, name, tag string) (*Release, error)
	TagRelease(release *Release, tag string) error
	AddRelease(*Release) error
	UpdateRelease(*Release) error
	GetAllReleases() ([]*Release, error)
	GetPackageURIs(release *Release) ([]string, error)
	AddPackageURI(release *Release, uri string) error
	GetProviders(providerName string) (map[string]*MinimalReleaseMetadata, error)
	GetProvidersFilteredBy(providerName string, q *ProvidersFilter) (map[string]*MinimalReleaseMetadata, error)
	RegisterProviders(release *core.ReleaseMetadata) error
}

type Release struct {
	Application           *Application
	ReleaseId             string
	Version               string
	Metadata              *core.ReleaseMetadata
	ProcessedDependencies bool
	Downloads             int
	UploadedBy            string
	UploadedAt            time.Time
}

func NewRelease(app *Application, metadata *core.ReleaseMetadata) *Release {
	return &Release{
		Application: app,
		ReleaseId:   metadata.GetReleaseId(),
		Version:     metadata.Version,
		Metadata:    metadata,
	}
}

// Used for view purposes
type MinimalReleaseMetadata struct {
	Application string `json:"application"`
	Project     string `json:"project"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

func (m *MinimalReleaseMetadata) GetReleaseId() string {
	return m.Project + "/" + m.Application + "-v" + m.Version
}

type ProvidersFilter struct {
	Namespaces []string
}

func NewProvidersFilter(namespaces []string) *ProvidersFilter {
	return &ProvidersFilter{
		Namespaces: namespaces,
	}
}
