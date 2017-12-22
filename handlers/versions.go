/*
Copyright 2017 Ankyra

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

	core "github.com/ankyra/escape-core"
	"github.com/ankyra/escape-inventory/model"
	"github.com/gorilla/mux"
)

type versionHandlerProvider struct {
	GetReleaseMetadata func(project, name, version string) (*core.ReleaseMetadata, error)
	GetRelease         func(project, name, version string) (*model.ReleasePayload, error)
	GetNextVersion     func(project, name, prefix string) (string, error)
	GetPreviousVersion func(project, name, version string) (*core.ReleaseMetadata, error)
}

func newVersionHandlerProvider() *versionHandlerProvider {
	return &versionHandlerProvider{
		GetReleaseMetadata: model.GetReleaseMetadata,
		GetRelease:         model.GetRelease,
		GetNextVersion:     model.GetNextVersion,
		GetPreviousVersion: model.GetPreviousReleaseMetadata,
	}
}

func GetVersionHandler(w http.ResponseWriter, r *http.Request) {
	newVersionHandlerProvider().GetVersionHandler(w, r)
}
func NextVersionHandler(w http.ResponseWriter, r *http.Request) {
	newVersionHandlerProvider().NextVersionHandler(w, r)
}
func PreviousVersionHandler(w http.ResponseWriter, r *http.Request) {
	newVersionHandlerProvider().PreviousVersionHandler(w, r)
}

func (h *versionHandlerProvider) GetVersionHandler(w http.ResponseWriter, r *http.Request) {
	project := mux.Vars(r)["project"]
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]
	full := r.URL.Query().Get("full")
	var result interface{}
	var err error
	if full != "" {
		result, err = h.GetRelease(project, name, version)
	} else {
		result, err = h.GetReleaseMetadata(project, name, version)
	}
	ErrorOrJsonSuccess(w, r, result, err)
}

func (h *versionHandlerProvider) NextVersionHandler(w http.ResponseWriter, r *http.Request) {
	project := mux.Vars(r)["project"]
	name := mux.Vars(r)["name"]
	prefix := r.URL.Query().Get("prefix")
	version, err := h.GetNextVersion(project, name, prefix)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	w.Write([]byte(version))
}

func (h *versionHandlerProvider) PreviousVersionHandler(w http.ResponseWriter, r *http.Request) {
	project := mux.Vars(r)["project"]
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]
	metadata, err := h.GetPreviousVersion(project, name, version)
	ErrorOrJsonSuccess(w, r, metadata, err)
}
