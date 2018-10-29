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

type namespaceHandlerProvider struct {
	GetNamespaces   func() (map[string]*types.Project, error)
	GetNamespace    func(namespace string) (*model.NamespacePayload, error)
	AddNamespace    func(namespace *types.Project, username string) error
	UpdateNamespace func(namespace *types.Project) error

	GetNamespaceHooks    func(namespace string) (types.Hooks, error)
	UpdateNamespaceHooks func(namespace string, hooks types.Hooks) error

	HardDeleteNamespace func(namespace string) error
}

func newNamespaceHandlerProvider() *namespaceHandlerProvider {
	return &namespaceHandlerProvider{
		GetNamespaces:        dao.GetNamespaces,
		GetNamespace:         model.GetNamespace,
		AddNamespace:         model.AddNamespace,
		UpdateNamespace:      model.UpdateNamespace,
		GetNamespaceHooks:    model.GetNamespaceHooks,
		UpdateNamespaceHooks: model.UpdateNamespaceHooks,
		HardDeleteNamespace:  dao.HardDeleteNamespace,
	}
}

func GetNamespacesHandler(w http.ResponseWriter, r *http.Request) {
	newNamespaceHandlerProvider().GetNamespacesHandler(w, r)
}
func GetNamespaceHandler(w http.ResponseWriter, r *http.Request) {
	newNamespaceHandlerProvider().GetNamespaceHandler(w, r)
}
func AddNamespaceHandler(w http.ResponseWriter, r *http.Request) {
	newNamespaceHandlerProvider().AddNamespaceHandler(w, r)
}
func UpdateNamespaceHandler(w http.ResponseWriter, r *http.Request) {
	newNamespaceHandlerProvider().UpdateNamespaceHandler(w, r)
}
func GetNamespaceHooksHandler(w http.ResponseWriter, r *http.Request) {
	newNamespaceHandlerProvider().GetNamespaceHooksHandler(w, r)
}
func UpdateNamespaceHooksHandler(w http.ResponseWriter, r *http.Request) {
	newNamespaceHandlerProvider().UpdateNamespaceHooksHandler(w, r)
}
func HardDeleteNamespaceHandler(w http.ResponseWriter, r *http.Request) {
	newNamespaceHandlerProvider().HardDeleteNamespaceHandler(w, r)
}

func (h *namespaceHandlerProvider) GetNamespacesHandler(w http.ResponseWriter, r *http.Request) {
	result, err := h.GetNamespaces()
	ErrorOrJsonSuccess(w, r, result, err)
}

func (h *namespaceHandlerProvider) GetNamespaceHandler(w http.ResponseWriter, r *http.Request) {
	namespace := mux.Vars(r)["namespace"]
	proj, err := h.GetNamespace(namespace)
	ErrorOrJsonSuccess(w, r, proj, err)
}

func (h *namespaceHandlerProvider) AddNamespaceHandler(w http.ResponseWriter, r *http.Request) {
	result := types.Project{}
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		HandleError(w, r, model.NewUserError(fmt.Errorf("Invalid JSON")))
		return
	}
	if err := h.AddNamespace(&result, ""); err != nil {
		HandleError(w, r, err)
		return
	}
	w.WriteHeader(200)
}

func (h *namespaceHandlerProvider) UpdateNamespaceHandler(w http.ResponseWriter, r *http.Request) {
	namespace := mux.Vars(r)["namespace"]
	result := types.Project{}
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		HandleError(w, r, model.NewUserError(fmt.Errorf("Invalid JSON")))
		return
	}
	if result.Name != namespace {
		HandleError(w, r, model.NewUserError(fmt.Errorf("Namespace 'name' field doesn't correspond with URL")))
		return
	}
	if err := h.UpdateNamespace(&result); err != nil {
		HandleError(w, r, err)
		return
	}
	w.WriteHeader(201)
}

func (h *namespaceHandlerProvider) GetNamespaceHooksHandler(w http.ResponseWriter, r *http.Request) {
	namespace := mux.Vars(r)["namespace"]
	hooks, err := h.GetNamespaceHooks(namespace)
	ErrorOrJsonSuccess(w, r, hooks, err)
}

func (h *namespaceHandlerProvider) UpdateNamespaceHooksHandler(w http.ResponseWriter, r *http.Request) {
	namespace := mux.Vars(r)["namespace"]
	result := types.Hooks{}
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		HandleError(w, r, model.NewUserError(fmt.Errorf("Invalid JSON")))
		return
	}
	if err := h.UpdateNamespaceHooks(namespace, result); err != nil {
		HandleError(w, r, err)
		return
	}
	w.WriteHeader(201)
}

func (h *namespaceHandlerProvider) HardDeleteNamespaceHandler(w http.ResponseWriter, r *http.Request) {
	namespace := mux.Vars(r)["namespace"]
	ErrorOrSuccess(w, r, h.HardDeleteNamespace(namespace))
}
