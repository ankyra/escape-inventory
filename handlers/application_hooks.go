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

	"github.com/ankyra/escape-inventory/dao/types"
	"github.com/ankyra/escape-inventory/model"
	"github.com/gorilla/mux"
)

func GetApplicationHooksHandler(w http.ResponseWriter, r *http.Request) {
	project := mux.Vars(r)["project"]
	name := mux.Vars(r)["name"]
	getApplicationHooksHandler(w, r, project, name)
}

func getApplicationHooksHandler(w http.ResponseWriter, r *http.Request, project, name string) {
	hooks, err := model.GetApplicationHooks(project, name)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(hooks)
}

func UpdateApplicationHooksHandler(w http.ResponseWriter, r *http.Request) {
	project := mux.Vars(r)["project"]
	name := mux.Vars(r)["name"]
	if r.Body == nil {
		HandleError(w, r, model.NewUserError(fmt.Errorf("Empty body")))
		return
	}
	result := types.Hooks{}
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		HandleError(w, r, model.NewUserError(fmt.Errorf("Invalid JSON")))
		return
	}
	if err := model.UpdateApplicationHooks(project, name, result); err != nil {
		HandleError(w, r, err)
		return
	}
	w.WriteHeader(201)
}
