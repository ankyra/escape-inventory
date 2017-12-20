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
	"errors"
	"io"
	"io/ioutil"
	"net/http"

	core "github.com/ankyra/escape-core"
	"github.com/ankyra/escape-inventory/dao/types"
	"github.com/gorilla/mux"
	. "gopkg.in/check.v1"
)

const (
	RegisterURL     = "/api/v1/registry/{project}/register"
	registerTestURL = "/api/v1/registry/project/register"
)

func (s *suite) registerMuxWithProvider(provider *registerHandlerProvider) *mux.Router {
	r := mux.NewRouter()
	postRouter := r.Methods("POST").Subrouter()
	postRouter.Handle(RegisterURL, http.HandlerFunc(provider.RegisterHandler))
	return r
}

func (s *suite) Test_RegisterHandler_happy_path(c *C) {
	var capturedProject, capturedMetadata, capturedUsername string
	provider := &registerHandlerProvider{
		AddReleaseByUser: func(project, metadata, username string) (*core.ReleaseMetadata, error) {
			capturedProject = project
			capturedMetadata = metadata
			capturedUsername = username
			return core.NewReleaseMetadata("name", "1.0"), nil
		},
		ReadRequestBody: func(body io.Reader) ([]byte, error) {
			return []byte("metadata"), nil
		},
	}
	resp := s.testPOST(c, s.registerMuxWithProvider(provider), registerTestURL, nil)
	c.Assert(resp.StatusCode, Equals, 200)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "")
	c.Assert(capturedProject, Equals, "project")
	c.Assert(capturedMetadata, Equals, "metadata")
	c.Assert(capturedUsername, Equals, "")
}

func (s *suite) Test_RegisterHandler_fails_if_add_release_fails(c *C) {
	provider := &registerHandlerProvider{
		AddReleaseByUser: func(project, metadata, username string) (*core.ReleaseMetadata, error) {
			return nil, types.AlreadyExists
		},
		ReadRequestBody: func(body io.Reader) ([]byte, error) {
			return []byte{}, nil
		},
	}
	resp := s.testPOST(c, s.registerMuxWithProvider(provider), registerTestURL, nil)
	c.Assert(resp.StatusCode, Equals, 409)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "Resource already exists")
}

func (s *suite) Test_RegisterHandler_fails_if_body_read_fails(c *C) {
	provider := &registerHandlerProvider{
		ReadRequestBody: func(body io.Reader) ([]byte, error) {
			return nil, errors.New("couldn't read body")
		},
	}
	resp := s.testPOST(c, s.registerMuxWithProvider(provider), registerTestURL, nil)
	c.Assert(resp.StatusCode, Equals, 500)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "")
}
