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
	"encoding/json"
	"github.com/ankyra/escape-core"
	"github.com/ankyra/escape-registry/model"
	"github.com/gorilla/mux"
	"net/http"
)

func DiffHandler(w http.ResponseWriter, r *http.Request) {
	project := mux.Vars(r)["project"]
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]
	diffWith := mux.Vars(r)["diffWith"]
	metadata, err := model.GetReleaseMetadata(project, name+"-"+version)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	releaseId := name + "-" + diffWith
	if diffWith == "" {
		prev, err := model.GetPreviousVersion(project, name, metadata.Version)
		if err != nil {
			HandleError(w, r, err)
			return
		}
		releaseId = name + "-v" + prev
	}
	previousMetadata, err := model.GetReleaseMetadata(project, releaseId)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	changes := []string{}
	for _, change := range core.Diff(previousMetadata, metadata) {
		changes = append(changes, change.ToString())
	}
	bytes, err := json.Marshal(changes)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(bytes)
}
