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

	"github.com/ankyra/escape-inventory/dao/types"
	"github.com/ankyra/escape-inventory/model"
	"github.com/gorilla/mux"
)

type dependencyHandlerProvider struct {
	GetDownstreamDependencies func(namespace, name, version string) ([]*types.Dependency, error)
	GetDependencyGraph        func(namespace, name, version string, downstreamFunc model.DownstreamDependenciesResolver) (*model.DependencyGraph, error)
}

func newDependencyHandlerProvider() *dependencyHandlerProvider {
	return &dependencyHandlerProvider{
		GetDownstreamDependencies: model.GetDownstreamDependencies,
		GetDependencyGraph:        model.GetDependencyGraph,
	}
}

func DownstreamHandler(w http.ResponseWriter, r *http.Request) {
	newDependencyHandlerProvider().DownstreamHandler(w, r)
}
func DependencyGraphHandler(w http.ResponseWriter, r *http.Request) {
	newDependencyHandlerProvider().DependencyGraphHandler(w, r)
}

func (h *dependencyHandlerProvider) DownstreamHandler(w http.ResponseWriter, r *http.Request) {
	namespace := mux.Vars(r)["namespace"]
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]
	deps, err := h.GetDownstreamDependencies(namespace, name, version)
	ErrorOrJsonSuccess(w, r, deps, err)
}

func (h *dependencyHandlerProvider) DependencyGraphHandler(w http.ResponseWriter, r *http.Request) {
	namespace := mux.Vars(r)["namespace"]
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]
	graph, err := h.GetDependencyGraph(namespace, name, version, nil)
	ErrorOrJsonSuccess(w, r, graph, err)
}
