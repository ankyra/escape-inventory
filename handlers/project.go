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
	"fmt"
	"net/http"

	"github.com/ankyra/escape-registry/dao"
	"github.com/ankyra/escape-registry/dao/types"
	"github.com/ankyra/escape-registry/model"
	"github.com/gorilla/mux"
)

func GetProjectsHandler(w http.ResponseWriter, r *http.Request) {
	result, err := dao.GetProjects()
	if err != nil {
		HandleError(w, r, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(result)
}

func GetProjectHandler(w http.ResponseWriter, r *http.Request) {
	project := mux.Vars(r)["project"]
	getProjectHandler(w, r, project)
}

func getProjectHandler(w http.ResponseWriter, r *http.Request, project string) {
	proj, err := model.GetProject(project)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(proj)

}

func AddProjectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		HandleError(w, r, model.NewUserError(fmt.Errorf("Empty body")))
		return
	}
	result := types.Project{}
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		HandleError(w, r, model.NewUserError(fmt.Errorf("Invalid JSON")))
		return
	}
	if err := model.AddProject(&result); err != nil {
		HandleError(w, r, err)
		return
	}
	w.WriteHeader(200)
}

func UpdateProjectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		HandleError(w, r, model.NewUserError(fmt.Errorf("Empty body")))
		return
	}
	result := types.Project{}
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		HandleError(w, r, model.NewUserError(fmt.Errorf("Invalid JSON")))
		return
	}
	if err := model.UpdateProject(&result); err != nil {
		HandleError(w, r, err)
		return
	}
	w.WriteHeader(201)
}
