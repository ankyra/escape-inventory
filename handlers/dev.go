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
	"net/http"

	"github.com/ankyra/escape-inventory/dao"
)

type devHandlerProvider struct {
	WipeDatabaseFunc func() error
}

func NewDevHandlerProvider() *devHandlerProvider {
	return &devHandlerProvider{
		WipeDatabaseFunc: dao.GlobalDAO.WipeDatabase,
	}
}

func WipeDatabaseHandler(w http.ResponseWriter, r *http.Request) {
	NewDevHandlerProvider().wipeDatabase(w, r)
}

func (h devHandlerProvider) wipeDatabase(w http.ResponseWriter, r *http.Request) {
	h.WipeDatabaseFunc()
	w.WriteHeader(200)
}
