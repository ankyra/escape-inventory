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
	"bytes"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/ankyra/escape-inventory/dao/types"
	"github.com/gorilla/mux"
	. "gopkg.in/check.v1"
)

const (
	DownloadURL     = "/api/v1/registry/{project}/units/{name}/versions/{version}/download"
	downloadTestURL = "/api/v1/registry/project/units/name/versions/v1.0.0/download"
)

func (s *suite) downloadMuxWithProvider(provider *downloadHandlerProvider) *mux.Router {
	r := mux.NewRouter()
	postRouter := r.Methods("GET").Subrouter()
	postRouter.Handle(DownloadURL, http.HandlerFunc(provider.DownloadHandler))
	return r
}

/*
	DownloadHandler
*/

func (s *suite) Test_DownloadHandler_happy_path(c *C) {

	provider := &downloadHandlerProvider{
		GetDownloadReadSeeker: func(project, releaseId string) (io.Reader, error) {
			return bytes.NewReader([]byte("package data")), nil
		},
	}
	resp := s.testGET(c, s.downloadMuxWithProvider(provider), downloadTestURL)
	c.Assert(resp.StatusCode, Equals, 200)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "package data")
	c.Assert(resp.Header.Get("Content-Type"), Equals, "application/gzip")
	c.Assert(resp.Header.Get("Content-Disposition"), Equals, `attachment; filename="name-v1.0.0.tgz"`)
}

func (s *suite) Test_DownloadHandler_fails_if_download_fails(c *C) {
	provider := &downloadHandlerProvider{
		GetDownloadReadSeeker: func(project, releaseId string) (io.Reader, error) {
			return nil, types.NotFound
		},
	}
	resp := s.testGET(c, s.downloadMuxWithProvider(provider), downloadTestURL)
	c.Assert(resp.StatusCode, Equals, 404)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "")
}
