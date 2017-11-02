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
	"net/http"

	"github.com/ankyra/escape-inventory/model"
	"github.com/gorilla/mux"
)

func RegistryHandler(w http.ResponseWriter, r *http.Request) {
	project := mux.Vars(r)["project"]
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]
	full := r.URL.Query().Get("full")
	releaseId := name + "-" + version
	var result interface{}
	if version == "" {
		result_, err := model.GetApplicationVersions(project, name)
		if err != nil {
			HandleError(w, r, err)
			return
		}
		result = result_
	} else if full == "" {
		result_, err := model.GetReleaseMetadata(project, releaseId)
		if err != nil {
			HandleError(w, r, err)
			return
		}
		result = result_
	} else {
		result_, err := model.GetRelease(project, releaseId)
		if err != nil {
			HandleError(w, r, err)
			return
		}
		result = result_
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(result)
}

func GetReleaseHandler(w http.ResponseWriter, r *http.Request) {
	project := mux.Vars(r)["project"]
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]
	releaseId := name + "-" + version
	result, err := model.GetRelease(project, releaseId)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(result)
}
