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
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/ankyra/escape-inventory/dao/types"
	"github.com/gorilla/mux"
	. "gopkg.in/check.v1"
)

const (
	UploadURL     = "/api/v1/registry/{project}/units/{name}/versions/{version}/upload"
	uploadTestURL = "/api/v1/registry/project/units/name/versions/v1.0.0/upload"
)

func (s *suite) uploadMux() *mux.Router {
	return s.uploadMuxWithProvider(newUploadHandlerProvider())
}

func (s *suite) uploadMuxWithProvider(provider *uploadHandlerProvider) *mux.Router {
	r := mux.NewRouter()
	postRouter := r.Methods("POST").Subrouter()
	postRouter.Handle(UploadURL, http.HandlerFunc(provider.UploadHandler))
	return r
}

/*
	UploadHandler
*/

func (s *suite) Test_UploadHandler(c *C) {
	content := "package content"
	file := "my-package.tgz"
	os.RemoveAll(file)
	err := ioutil.WriteFile(file, []byte(content), 0644)
	c.Assert(err, IsNil)

	provider := &uploadHandlerProvider{
		UploadPackage: func(project, releaseId string, pkg io.ReadSeeker) error {
			return nil
		},
	}

	resp := s.testPOST_file(c, s.uploadMuxWithProvider(provider), uploadTestURL, file)
	c.Assert(resp.StatusCode, Equals, 200)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "")

	os.RemoveAll(file)
}

func (s *suite) Test_UploadHandler_fails_if_upload_fails(c *C) {
	content := "package content"
	file := "my-package.tgz"
	os.RemoveAll(file)
	err := ioutil.WriteFile(file, []byte(content), 0644)
	c.Assert(err, IsNil)

	provider := &uploadHandlerProvider{
		UploadPackage: func(project, releaseId string, pkg io.ReadSeeker) error {
			return types.AlreadyExists
		},
	}

	resp := s.testPOST_file(c, s.uploadMuxWithProvider(provider), uploadTestURL, file)
	c.Assert(resp.StatusCode, Equals, 409)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "Resource already exists")

	os.RemoveAll(file)
}

func (s *suite) Test_UploadHandler_fails_if_project_does_not_exist(c *C) {
	content := "package content"
	file := "my-package.tgz"
	os.RemoveAll(file)
	err := ioutil.WriteFile(file, []byte(content), 0644)
	c.Assert(err, IsNil)

	resp := s.testPOST_file(c, s.uploadMux(), uploadTestURL, file)
	c.Assert(resp.StatusCode, Equals, 404)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "")

	os.RemoveAll(file)
}

func (s *suite) Test_UploadHandler_fails_if_not_multipart_request(c *C) {
	resp := s.testPOST(c, s.uploadMux(), uploadTestURL, nil)
	c.Assert(resp.StatusCode, Equals, 400)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "request Content-Type isn't multipart/form-data")
}

func (s *suite) Test_UploadHandler_fails_if_file_form_field_missing(c *C) {
	resp := s.testPOST_file(c, s.uploadMux(), uploadTestURL, "")
	c.Assert(resp.StatusCode, Equals, 400)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "http: no such file")
}
