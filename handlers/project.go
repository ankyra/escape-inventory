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

	"github.com/ankyra/escape-inventory/dao"
	"github.com/ankyra/escape-inventory/dao/types"
	"github.com/ankyra/escape-inventory/model"
	"github.com/gorilla/mux"
)

type projectHandlerProvider struct {
	GetProjects   func() (map[string]*types.Project, error)
	GetProject    func(project string) (*model.ProjectPayload, error)
	AddProject    func(project *types.Project, username string) error
	UpdateProject func(project *types.Project) error

	GetProjectHooks    func(project string) (types.Hooks, error)
	UpdateProjectHooks func(project string, hooks types.Hooks) error

	HardDeleteProject func(project string) error
}

func newProjectHandlerProvider() *projectHandlerProvider {
	return &projectHandlerProvider{
		GetProjects:        dao.GetProjects,
		GetProject:         model.GetProject,
		AddProject:         model.AddProject,
		UpdateProject:      model.UpdateProject,
		GetProjectHooks:    model.GetProjectHooks,
		UpdateProjectHooks: model.UpdateProjectHooks,
		HardDeleteProject:  dao.HardDeleteProject,
	}
}

func GetProjectsHandler(w http.ResponseWriter, r *http.Request) {
	newProjectHandlerProvider().GetProjectsHandler(w, r)
}
func GetProjectHandler(w http.ResponseWriter, r *http.Request) {
	newProjectHandlerProvider().GetProjectHandler(w, r)
}
func AddProjectHandler(w http.ResponseWriter, r *http.Request) {
	newProjectHandlerProvider().AddProjectHandler(w, r)
}
func UpdateProjectHandler(w http.ResponseWriter, r *http.Request) {
	newProjectHandlerProvider().UpdateProjectHandler(w, r)
}
func GetProjectHooksHandler(w http.ResponseWriter, r *http.Request) {
	newProjectHandlerProvider().GetProjectHooksHandler(w, r)
}
func UpdateProjectHooksHandler(w http.ResponseWriter, r *http.Request) {
	newProjectHandlerProvider().UpdateProjectHooksHandler(w, r)
}
func HardDeleteProjectHandler(w http.ResponseWriter, r *http.Request) {
	newProjectHandlerProvider().HardDeleteProjectHandler(w, r)
}

func (h *projectHandlerProvider) GetProjectsHandler(w http.ResponseWriter, r *http.Request) {
	result, err := h.GetProjects()
	ErrorOrJsonSuccess(w, r, result, err)
}

func (h *projectHandlerProvider) GetProjectHandler(w http.ResponseWriter, r *http.Request) {
	project := mux.Vars(r)["project"]
	proj, err := h.GetProject(project)
	ErrorOrJsonSuccess(w, r, proj, err)
}

func (h *projectHandlerProvider) AddProjectHandler(w http.ResponseWriter, r *http.Request) {
	result := types.Project{}
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		HandleError(w, r, model.NewUserError(fmt.Errorf("Invalid JSON")))
		return
	}
	if err := h.AddProject(&result, ""); err != nil {
		HandleError(w, r, err)
		return
	}
	w.WriteHeader(200)
}

func (h *projectHandlerProvider) UpdateProjectHandler(w http.ResponseWriter, r *http.Request) {
	project := mux.Vars(r)["project"]
	result := types.Project{}
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		HandleError(w, r, model.NewUserError(fmt.Errorf("Invalid JSON")))
		return
	}
	if result.Name != project {
		HandleError(w, r, model.NewUserError(fmt.Errorf("Project 'name' field doesn't correspond with URL")))
		return
	}
	if err := h.UpdateProject(&result); err != nil {
		HandleError(w, r, err)
		return
	}
	w.WriteHeader(201)
}

func (h *projectHandlerProvider) GetProjectHooksHandler(w http.ResponseWriter, r *http.Request) {
	project := mux.Vars(r)["project"]
	hooks, err := h.GetProjectHooks(project)
	ErrorOrJsonSuccess(w, r, hooks, err)
}

func (h *projectHandlerProvider) UpdateProjectHooksHandler(w http.ResponseWriter, r *http.Request) {
	project := mux.Vars(r)["project"]
	result := types.Hooks{}
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		HandleError(w, r, model.NewUserError(fmt.Errorf("Invalid JSON")))
		return
	}
	if err := h.UpdateProjectHooks(project, result); err != nil {
		HandleError(w, r, err)
		return
	}
	w.WriteHeader(201)
}

func (h *projectHandlerProvider) HardDeleteProjectHandler(w http.ResponseWriter, r *http.Request) {
	project := mux.Vars(r)["project"]
	ErrorOrSuccess(w, r, h.HardDeleteProject(project))
}
