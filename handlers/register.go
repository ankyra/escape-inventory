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
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/ankyra/escape-inventory/model"
	"github.com/gorilla/mux"
)

type Membership struct {
	Group string `json:"name"`
}

type User struct {
	Name   string        `json:"username"`
	Groups []*Membership `json:"groups"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	project := mux.Vars(r)["project"]
	if r.Body == nil {
		HandleError(w, r, model.NewUserError(fmt.Errorf("Empty body")))
		return
	}
	metadata, err := ioutil.ReadAll(r.Body)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	username := ReadUsernameFromContext(r)
	if _, err := model.AddReleaseByUser(project, string(metadata), username); err != nil {
		HandleError(w, r, err)
		return
	}
	w.WriteHeader(200)
}

func ReadUsernameFromContext(r *http.Request) string {
	user := r.Context().Value("user")
	if user != nil {
		value := reflect.Indirect(reflect.ValueOf(user))
		return value.FieldByName("Name").String()
	}
	return ""
}
