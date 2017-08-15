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

	"github.com/ankyra/escape-registry/model"
	"github.com/gorilla/mux"
)

func GetApplicationsHandler(w http.ResponseWriter, r *http.Request) {
	project := mux.Vars(r)["project"]
	getApplicationsHandler(w, r, project)
}

func getApplicationsHandler(w http.ResponseWriter, r *http.Request, project string) {
	apps, err := model.GetApplications(project)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(apps)
}

func GetApplicationHandler(w http.ResponseWriter, r *http.Request) {
	project := mux.Vars(r)["project"]
	name := mux.Vars(r)["name"]
	getApplicationHandler(w, r, project, name)
}

func getApplicationHandler(w http.ResponseWriter, r *http.Request, project, name string) {
	app, err := model.GetApplication(project, name)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(app)
}
