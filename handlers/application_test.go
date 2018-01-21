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
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/ankyra/escape-inventory/dao/types"
	"github.com/ankyra/escape-inventory/model"
	"github.com/gorilla/mux"

	. "gopkg.in/check.v1"
)

const (
	GetApplicationsURL     = "/api/v1/inventory/{project}/units/"
	getApplicationsTestURL = "/api/v1/inventory/project/units/"

	GetApplicationURL     = "/api/v1/inventory/{project}/units/{name}/"
	getApplicationTestURL = "/api/v1/inventory/project/units/name/"

	GetApplicationVersionsURL     = "/api/v1/inventory/{project}/units/{name}/versions/"
	getApplicationVersionsTestURL = "/api/v1/inventory/project/units/name/versions/"

	GetApplicationHooksURL     = "/api/v1/inventory/{project}/units/{name}/hooks/"
	getApplicationHooksTestURL = "/api/v1/inventory/project/units/name/hooks/"

	UpdateApplicationHooksURL     = "/api/v1/inventory/{project}/units/{name}/hooks/"
	updateApplicationHooksTestURL = "/api/v1/inventory/project/units/name/hooks/"
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
	var capturedProject, capturedName string
	provider := &applicationHandlerProvider{
		GetApplication: func(project, name string) (*model.ApplicationPayload, error) {
			capturedProject = project
			capturedName = name
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
	c.Assert(capturedName, Equals, "name")
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

/*
	GetApplicationVersionsHandler

*/

func (s *suite) getApplicationVersionsMuxWithProvider(provider *applicationHandlerProvider) *mux.Router {
	r := mux.NewRouter()
	router := r.Methods("GET").Subrouter()
	router.Handle(GetApplicationVersionsURL, http.HandlerFunc(provider.GetApplicationVersionsHandler))
	return r
}

func (s *suite) Test_GetApplicationVersionsHandler_happy_path(c *C) {
	var capturedProject, capturedName string
	provider := &applicationHandlerProvider{
		GetApplicationVersions: func(project, name string) ([]string, error) {
			capturedProject = project
			capturedName = name
			return []string{"1.0", "1.1"}, nil
		},
	}
	resp := s.testGET(c, s.getApplicationVersionsMuxWithProvider(provider), getApplicationVersionsTestURL)
	c.Assert(resp.StatusCode, Equals, 200)
	c.Assert(capturedProject, Equals, "project")
	c.Assert(capturedName, Equals, "name")
	result := []string{}
	c.Assert(json.NewDecoder(resp.Body).Decode(&result), IsNil)
	c.Assert(result, DeepEquals, []string{"1.0", "1.1"})
}

func (s *suite) Test_GetApplicationVersionsHandler_fails_if_GetApplicationVersions_fails(c *C) {
	provider := &applicationHandlerProvider{
		GetApplicationVersions: func(project, name string) ([]string, error) {
			return nil, types.NotFound
		},
	}
	resp := s.testGET(c, s.getApplicationVersionsMuxWithProvider(provider), getApplicationVersionsTestURL)
	c.Assert(resp.StatusCode, Equals, 404)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "")
}

/*
	GetApplicationHooksHandler

*/

func (s *suite) getApplicationHooksMuxWithProvider(provider *applicationHandlerProvider) *mux.Router {
	r := mux.NewRouter()
	router := r.Methods("GET").Subrouter()
	router.Handle(GetApplicationHooksURL, http.HandlerFunc(provider.GetApplicationHooksHandler))
	return r
}

func (s *suite) Test_GetApplicationHooksHandler_happy_path(c *C) {
	var capturedProject, capturedName string
	provider := &applicationHandlerProvider{
		GetApplicationHooks: func(project, name string) (types.Hooks, error) {
			capturedProject = project
			capturedName = name
			hooks := types.NewHooks()
			hooks["test"] = map[string]string{
				"ok": "what is this",
			}
			return hooks, nil
		},
	}
	resp := s.testGET(c, s.getApplicationHooksMuxWithProvider(provider), getApplicationHooksTestURL)
	c.Assert(resp.StatusCode, Equals, 200)
	c.Assert(capturedProject, Equals, "project")
	c.Assert(capturedName, Equals, "name")
	c.Assert(resp.Header.Get("Content-Type"), Equals, "application/json")
	result := types.NewHooks()
	c.Assert(json.NewDecoder(resp.Body).Decode(&result), IsNil)
	c.Assert(result["test"]["ok"], Equals, "what is this")
}

func (s *suite) Test_GetApplicationHooksHandler_fails_if_GetApplicationHooks_fails(c *C) {
	var capturedProject, capturedName string
	provider := &applicationHandlerProvider{
		GetApplicationHooks: func(project, name string) (types.Hooks, error) {
			capturedProject = project
			capturedName = name
			return nil, types.NotFound
		},
	}
	resp := s.testGET(c, s.getApplicationHooksMuxWithProvider(provider), getApplicationHooksTestURL)
	c.Assert(resp.StatusCode, Equals, 404)
	c.Assert(capturedProject, Equals, "project")
	c.Assert(capturedName, Equals, "name")
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "")
}

/*
	UpdateApplicationHooksHandler

*/

func (s *suite) updateApplicationHooksMuxWithProvider(provider *applicationHandlerProvider) *mux.Router {
	r := mux.NewRouter()
	router := r.Methods("PUT").Subrouter()
	router.Handle(UpdateApplicationHooksURL, http.HandlerFunc(provider.UpdateApplicationHooksHandler))
	return r
}

func (s *suite) Test_UpdateApplicationHooksHandler_happy_path(c *C) {
	provider := &applicationHandlerProvider{
		UpdateApplicationHooks: func(project, name string, hooks types.Hooks) error {
			return nil
		},
	}
	data := map[string]map[string]string{
		"test": map[string]string{
			"test": "test",
		},
	}
	resp := s.testPUT(c, s.updateApplicationHooksMuxWithProvider(provider), updateApplicationHooksTestURL, data)
	c.Assert(resp.StatusCode, Equals, 201)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "")
}

func (s *suite) Test_UpdateApplicationHooksHandler_fails_if_invalid_json(c *C) {
	serviceCalled := false
	provider := &applicationHandlerProvider{
		UpdateApplicationHooks: func(project, name string, hooks types.Hooks) error {
			serviceCalled = true
			return types.NotFound
		},
	}
	var data interface{}
	resp := s.testPUT(c, s.updateApplicationHooksMuxWithProvider(provider), updateApplicationHooksTestURL, data)
	c.Assert(resp.StatusCode, Equals, 400)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "Invalid JSON")
}

func (s *suite) Test_UpdateApplicationHooksHandler_fails_if_UpdateApplicationHooks_fails(c *C) {
	provider := &applicationHandlerProvider{
		UpdateApplicationHooks: func(project, name string, hooks types.Hooks) error {
			return types.NotFound
		},
	}
	data := map[string]map[string]string{
		"test": map[string]string{
			"test": "test",
		},
	}
	resp := s.testPUT(c, s.updateApplicationHooksMuxWithProvider(provider), updateApplicationHooksTestURL, data)
	c.Assert(resp.StatusCode, Equals, 404)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "")
}
