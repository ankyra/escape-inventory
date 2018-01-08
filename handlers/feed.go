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
	"strconv"

	"github.com/ankyra/escape-inventory/dao/types"
	"github.com/ankyra/escape-inventory/model"
	"github.com/gorilla/mux"
)

type feedHandlerProvider struct {
	GetFeedPage    func(page int) ([]*types.FeedEvent, error)
	GetProjectFeed func(project string, page int) ([]*types.FeedEvent, error)
}

func newFeedHandlerProvider() *feedHandlerProvider {
	return &feedHandlerProvider{
		GetFeedPage:    model.GetFeedPage,
		GetProjectFeed: model.GetProjectFeed,
	}
}

func FeedHandler(w http.ResponseWriter, r *http.Request) {
	newFeedHandlerProvider().FeedHandler(w, r)
}
func ProjectFeedHandler(w http.ResponseWriter, r *http.Request) {
	newFeedHandlerProvider().ProjectFeedHandler(w, r)
}

func (h *feedHandlerProvider) FeedHandler(w http.ResponseWriter, r *http.Request) {
	page := 0
	pageStr := r.URL.Query().Get("page")
	if pageStr != "" {
		var err error
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			page = 0
		}
	}
	feed, err := h.GetFeedPage(page)
	ErrorOrJsonSuccess(w, r, feed, err)
}

func (h *feedHandlerProvider) ProjectFeedHandler(w http.ResponseWriter, r *http.Request) {
	page := 0
	pageStr := r.URL.Query().Get("page")
	if pageStr != "" {
		var err error
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			page = 0
		}
	}
	project := mux.Vars(r)["project"]
	feed, err := h.GetProjectFeed(project, page)
	ErrorOrJsonSuccess(w, r, feed, err)
}
