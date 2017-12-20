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
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/ankyra/escape-inventory/dao/types"
	"github.com/ankyra/escape-inventory/model"
	"github.com/gorilla/mux"

	. "gopkg.in/check.v1"
)

const (
	GetApplicationsURL     = "/api/v1/registry/{project}/units/"
	getApplicationsTestURL = "/api/v1/registry/project/units/"
	GetApplicationURL      = "/api/v1/registry/{project}/units/"
	getApplicationTestURL  = "/api/v1/registry/project/units/"
)

/*
	GetApplicationsHandler

*/

func (s *suite) getApplicationsMuxWithProvider(provider *applicationHandlerProvider) *mux.Router {
	r := mux.NewRouter()
	router := r.Methods("GET").Subrouter()
	router.Handle(GetApplicationsURL, http.HandlerFunc(provider.GetApplicationsHandler))
	return r
}

func (s *suite) Test_GetApplicationsHandler_happy_path(c *C) {
	var capturedProject string
	provider := &applicationHandlerProvider{
		GetApplications: func(project string) (map[string]*types.Application, error) {
			capturedProject = project
			return map[string]*types.Application{
				"my-app": types.NewApplication("project", "my-app"),
			}, nil
		},
	}
	resp := s.testGET(c, s.getApplicationsMuxWithProvider(provider), getApplicationsTestURL)
	c.Assert(resp.StatusCode, Equals, 200)
	c.Assert(capturedProject, Equals, "project")
	c.Assert(resp.Header.Get("Content-Type"), Equals, "application/json")
	result := map[string]*types.Application{}
	c.Assert(json.NewDecoder(resp.Body).Decode(&result), IsNil)
	c.Assert(result["my-app"].Name, Equals, "my-app")
}

func (s *suite) Test_GetApplicationsHandler_fails_if_get_applications_fails(c *C) {
	var capturedProject string
	provider := &applicationHandlerProvider{
		GetApplications: func(project string) (map[string]*types.Application, error) {
			capturedProject = project
			return nil, types.NotFound
		},
	}
	resp := s.testGET(c, s.getApplicationsMuxWithProvider(provider), getApplicationsTestURL)
	c.Assert(resp.StatusCode, Equals, 404)
	c.Assert(capturedProject, Equals, "project")
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "")
}

/*
	GetApplicationHandler

*/

func (s *suite) getApplicationMuxWithProvider(provider *applicationHandlerProvider) *mux.Router {
	r := mux.NewRouter()
	router := r.Methods("GET").Subrouter()
	router.Handle(GetApplicationURL, http.HandlerFunc(provider.GetApplicationHandler))
	return r
}

func (s *suite) Test_GetApplicationHandler_happy_path(c *C) {
	var capturedProject string
	provider := &applicationHandlerProvider{
		GetApplication: func(project, name string) (*model.ApplicationPayload, error) {
			capturedProject = project
			result := model.ApplicationPayload{
				Application: types.NewApplication("project", "name"),
				Versions:    []string{"1.0", "1.1"},
			}
			return &result, nil
		},
	}
	resp := s.testGET(c, s.getApplicationMuxWithProvider(provider), getApplicationTestURL)
	c.Assert(resp.StatusCode, Equals, 200)
	c.Assert(capturedProject, Equals, "project")
	c.Assert(resp.Header.Get("Content-Type"), Equals, "application/json")
	result := model.ApplicationPayload{}
	c.Assert(json.NewDecoder(resp.Body).Decode(&result), IsNil)
	c.Assert(result.Name, Equals, "name")
	c.Assert(result.Versions, HasLen, 2)
}

func (s *suite) Test_GetApplicationHandler_fails_if_get_application_fails(c *C) {
	var capturedProject string
	provider := &applicationHandlerProvider{
		GetApplication: func(project, name string) (*model.ApplicationPayload, error) {
			capturedProject = project
			return nil, types.NotFound
		},
	}
	resp := s.testGET(c, s.getApplicationMuxWithProvider(provider), getApplicationTestURL)
	c.Assert(resp.StatusCode, Equals, 404)
	c.Assert(capturedProject, Equals, "project")
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "")
}
