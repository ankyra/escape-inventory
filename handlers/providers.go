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

package handlers

import (
	"net/http"

	"github.com/ankyra/escape-inventory/dao"
	"github.com/ankyra/escape-inventory/dao/types"
)

type providerHandler struct {
	GetProviders func(providerName string) (map[string]*types.MinimalReleaseMetadata, error)
}

func newProviderHandler() *providerHandler {
	return &providerHandler{
		GetProviders: dao.GetProviders,
	}
}

func ProviderHandler(w http.ResponseWriter, r *http.Request) {
	newProviderHandler().providerHandler(w, r)
}

func (p *providerHandler) providerHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	providers, err := p.GetProviders(name)
	ErrorOrJsonSuccess(w, r, providers, err)
}
