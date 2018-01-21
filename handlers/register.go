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
	"io/ioutil"
	"net/http"

	core "github.com/ankyra/escape-core"
	"github.com/ankyra/escape-inventory/model"
	"github.com/gorilla/mux"
)

type registerHandlerProvider struct {
	AddReleaseByUser func(project, metadata, username string) (*core.ReleaseMetadata, error)
	ReadRequestBody  func(body io.Reader) ([]byte, error)
}

func newRegisterHandlerProvider() *registerHandlerProvider {
	return &registerHandlerProvider{
		AddReleaseByUser: model.AddReleaseByUser,
		ReadRequestBody:  ioutil.ReadAll,
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	newRegisterHandlerProvider().RegisterHandler(w, r)
}

func (h *registerHandlerProvider) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	project := mux.Vars(r)["project"]
	metadata, err := h.ReadRequestBody(r.Body)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	username := ReadUsernameFromContext(r)
	if _, err := h.AddReleaseByUser(project, string(metadata), username); err != nil {
		HandleError(w, r, err)
		return
	}
	w.WriteHeader(200)
}
