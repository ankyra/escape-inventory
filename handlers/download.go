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

	"github.com/ankyra/escape-inventory/metrics"
	"github.com/ankyra/escape-inventory/model"
	"github.com/gorilla/mux"
)

type downloadHandlerProvider struct {
	GetDownloadReadSeeker func(namespace, name, version string) (io.Reader, error)
}

func newDownloadHandlerProvider() *downloadHandlerProvider {
	return &downloadHandlerProvider{
		GetDownloadReadSeeker: model.GetDownloadReadSeeker,
	}
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	newDownloadHandlerProvider().DownloadHandler(w, r)
}

func (h *downloadHandlerProvider) DownloadHandler(w http.ResponseWriter, r *http.Request) {
	namespace := mux.Vars(r)["namespace"]
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]
	filename := name + "-" + version + ".tgz"
	reader, err := h.GetDownloadReadSeeker(namespace, name, version)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	metrics.DownloadCounter.Inc()
	w.Header().Set("Content-Type", "application/gzip")
	w.Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)
	w.WriteHeader(200)
	io.Copy(w, reader)
}
