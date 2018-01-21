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
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ankyra/escape-inventory/dao/types"
	"github.com/ankyra/escape-inventory/model"
	"github.com/gorilla/mux"
)

type applicationHandlerProvider struct {
	GetApplications        func(project string) (map[string]*types.Application, error)
	GetApplication         func(project, name string) (*model.ApplicationPayload, error)
	GetApplicationVersions func(project, name string) ([]string, error)
	GetApplicationHooks    func(project, name string) (types.Hooks, error)
	UpdateApplicationHooks func(project, name string, hooks types.Hooks) error
}

func newApplicationHandlerProvider() *applicationHandlerProvider {
	return &applicationHandlerProvider{
		GetApplications:        model.GetApplications,
		GetApplication:         model.GetApplication,
		GetApplicationHooks:    model.GetApplicationHooks,
		GetApplicationVersions: model.GetApplicationVersions,
		UpdateApplicationHooks: model.UpdateApplicationHooks,
	}
}

func GetApplicationsHandler(w http.ResponseWriter, r *http.Request) {
	newApplicationHandlerProvider().GetApplicationsHandler(w, r)
}
func GetApplicationHandler(w http.ResponseWriter, r *http.Request) {
	newApplicationHandlerProvider().GetApplicationHandler(w, r)
}
func GetApplicationVersionsHandler(w http.ResponseWriter, r *http.Request) {
	newApplicationHandlerProvider().GetApplicationVersionsHandler(w, r)
}
func GetApplicationHooksHandler(w http.ResponseWriter, r *http.Request) {
	newApplicationHandlerProvider().GetApplicationHooksHandler(w, r)
}
func UpdateApplicationHooksHandler(w http.ResponseWriter, r *http.Request) {
	newApplicationHandlerProvider().UpdateApplicationHooksHandler(w, r)
}

func (h *applicationHandlerProvider) GetApplicationsHandler(w http.ResponseWriter, r *http.Request) {
	project := mux.Vars(r)["project"]
	apps, err := h.GetApplications(project)
	ErrorOrJsonSuccess(w, r, apps, err)
}

func (h *applicationHandlerProvider) GetApplicationHandler(w http.ResponseWriter, r *http.Request) {
	project := mux.Vars(r)["project"]
	name := mux.Vars(r)["name"]
	app, err := h.GetApplication(project, name)
	ErrorOrJsonSuccess(w, r, app, err)
}

func (h *applicationHandlerProvider) GetApplicationVersionsHandler(w http.ResponseWriter, r *http.Request) {
	project := mux.Vars(r)["project"]
	name := mux.Vars(r)["name"]
	versions, err := h.GetApplicationVersions(project, name)
	ErrorOrJsonSuccess(w, r, versions, err)
}

func (h *applicationHandlerProvider) GetApplicationHooksHandler(w http.ResponseWriter, r *http.Request) {
	project := mux.Vars(r)["project"]
	name := mux.Vars(r)["name"]
	hooks, err := h.GetApplicationHooks(project, name)
	ErrorOrJsonSuccess(w, r, hooks, err)
}

func (h *applicationHandlerProvider) UpdateApplicationHooksHandler(w http.ResponseWriter, r *http.Request) {
	project := mux.Vars(r)["project"]
	name := mux.Vars(r)["name"]
	result := types.Hooks{}
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		HandleError(w, r, model.NewUserError(fmt.Errorf("Invalid JSON")))
		return
	}
	if err := h.UpdateApplicationHooks(project, name, result); err != nil {
		HandleError(w, r, err)
		return
	}
	w.WriteHeader(201)
}
