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
	"net/http"

	"github.com/ankyra/escape-inventory/dao"
	"github.com/ankyra/escape-inventory/dao/types"
	"github.com/gorilla/mux"
)

type feedHandlerProvider struct {
	GetFeedPage        func(pageSize int) ([]*types.FeedEvent, error)
	GetProjectFeedPage func(project string, pageSize int) ([]*types.FeedEvent, error)
}

func newFeedHandlerProvider() *feedHandlerProvider {
	return &feedHandlerProvider{
		GetFeedPage:        dao.GetFeedPage,
		GetProjectFeedPage: dao.GetProjectFeedPage,
	}
}

func FeedHandler(w http.ResponseWriter, r *http.Request) {
	newFeedHandlerProvider().FeedHandler(w, r)
}
func ProjectFeedHandler(w http.ResponseWriter, r *http.Request) {
	newFeedHandlerProvider().ProjectFeedHandler(w, r)
}

func (h *feedHandlerProvider) FeedHandler(w http.ResponseWriter, r *http.Request) {
	feed, err := h.GetFeedPage(types.FeedPageSize)
	ErrorOrJsonSuccess(w, r, feed, err)
}

func (h *feedHandlerProvider) ProjectFeedHandler(w http.ResponseWriter, r *http.Request) {
	project := mux.Vars(r)["project"]
	feed, err := h.GetProjectFeedPage(project, types.FeedPageSize)
	ErrorOrJsonSuccess(w, r, feed, err)
}
