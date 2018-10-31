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
	core "github.com/ankyra/escape-core"
	. "github.com/ankyra/escape-inventory/dao/types"
)

func (a *dao) GetProviders(providerName string) (map[string]*MinimalReleaseMetadata, error) {
	providers, ok := a.providers[providerName]
	if !ok {
		return nil, NotFound
	}
	result := map[string]*MinimalReleaseMetadata{}
	for _, release := range providers {
		result[release.GetReleaseId()] = release
	}
	return result, nil
}

func (a *dao) RegisterProviders(release *core.ReleaseMetadata) error {
	for _, provider := range release.Provides {
		name := provider.Name
		store, ok := a.providers[name]
		if !ok {
			store = map[string]*MinimalReleaseMetadata{}
		}
		key := release.Project + "-" + release.Name
		minidata := &MinimalReleaseMetadata{
			Application: release.Name,
			Project:     release.Project,
			Version:     release.Version,
			Description: release.Description,
		}
		rel, ok := store[key]
		if !ok {
			store[key] = minidata
		} else {
			current := core.NewSemanticVersion(rel.Version)
			maybeNew := core.NewSemanticVersion(release.Version)
			if !maybeNew.LessOrEqual(current) {
				store[key] = minidata
			}
		}
		a.providers[name] = store
	}
	return nil
}
