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
	"io"
	"net/http"

	"github.com/ankyra/escape-inventory/cmd"
	"github.com/ankyra/escape-inventory/metrics"
	"github.com/ankyra/escape-inventory/model"
	"github.com/gorilla/mux"
)

type uploadHandlerProvider struct {
	UploadPackage func(namespace, releaseId string, pkg io.ReadSeeker) error
}

func newUploadHandlerProvider() *uploadHandlerProvider {
	return &uploadHandlerProvider{
		UploadPackage: model.UploadPackage,
	}
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	newUploadHandlerProvider().UploadHandler(w, r)
}

func (h *uploadHandlerProvider) UploadHandler(w http.ResponseWriter, r *http.Request) {
	namespace := mux.Vars(r)["namespace"]
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]
	releaseId := name + "-" + version
	f, _, err := r.FormFile("file")
	if err != nil {
		HandleError(w, r, model.NewUserError(err))
		return
	}
	if err := h.UploadPackage(namespace, releaseId, f); err != nil {
		HandleError(w, r, err)
		return
	}
	metrics.UploadCounter.Inc()
	username := ReadUsernameFromContext(r)
	var url string
	if cmd.Config != nil && cmd.Config.WebHook != "" {
		url = cmd.Config.WebHook
	}
	go model.CallWebHook(namespace, name, version, releaseId, username, url)
	w.WriteHeader(200)
}
