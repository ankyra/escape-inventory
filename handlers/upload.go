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
	"github.com/ankyra/escape-registry/metrics"
	"github.com/ankyra/escape-registry/model"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	project := mux.Vars(r)["project"]
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]
	releaseId := name + "-" + version
	f, _, err := r.FormFile("file")
	if err != nil {
		HandleError(w, r, model.NewUserError(err))
		return
	}
	if err := model.UploadPackage(project, releaseId, f); err != nil {
		HandleError(w, r, err)
		return
	}
	metrics.UploadCounter.Inc()
	w.WriteHeader(200)
}

func RegisterAndUploadHandler(w http.ResponseWriter, r *http.Request) {
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
	release, err := model.AddRelease(project, string(metadata))
	if err != nil {
		HandleError(w, r, err)
		return
	}
	f, _, err := r.FormFile("file")
	if err != nil {
		HandleError(w, r, model.NewUserError(err))
		return
	}
	if err := model.UploadPackage(project, release.GetReleaseId(), f); err != nil {
		HandleError(w, r, err)
		return
	}
	metrics.UploadCounter.Inc()
	w.WriteHeader(200)
}
