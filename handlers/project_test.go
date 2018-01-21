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
	GetProjectsURL = "/api/v1/inventory/"

	GetProjectURL     = "/api/v1/inventory/{project}/"
	getProjectTestURL = "/api/v1/inventory/project/"

	AddProjectURL     = "/api/v1/inventory/{project}/add-project"
	addProjectTestURL = "/api/v1/inventory/project/add-project"

	UpdateProjectURL     = "/api/v1/inventory/{project}/"
	updateProjectTestURL = "/api/v1/inventory/project/"

	GetProjectHooksURL     = "/api/v1/inventory/{project}/hooks/"
	getProjectHooksTestURL = "/api/v1/inventory/project/hooks/"

	UpdateProjectHooksURL     = "/api/v1/inventory/{project}/hooks/"
	updateProjectHooksTestURL = "/api/v1/inventory/project/hooks/"

	HardDeleteProjectHooksURL     = "/api/v1/inventory/{project}/hard-delete"
	hardDeleteProjectHooksTestURL = "/api/v1/inventory/project/hard-delete"
)

/*
	GetProjectsHandler

*/

func (s *suite) getProjectsMuxWithProvider(provider *projectHandlerProvider) *mux.Router {
	r := mux.NewRouter()
	router := r.Methods("GET").Subrouter()
	router.Handle(GetProjectsURL, http.HandlerFunc(provider.GetProjectsHandler))
	return r
}

func (s *suite) Test_GetProjectsHandler_happy_path(c *C) {
	provider := &projectHandlerProvider{
		GetProjects: func() (map[string]*types.Project, error) {
			return map[string]*types.Project{
				"test": types.NewProject("test"),
			}, nil
		},
	}
	resp := s.testGET(c, s.getProjectsMuxWithProvider(provider), GetProjectsURL)
	c.Assert(resp.StatusCode, Equals, 200)
	c.Assert(resp.Header.Get("Content-Type"), Equals, "application/json")
	result := map[string]*types.Project{}
	c.Assert(json.NewDecoder(resp.Body).Decode(&result), IsNil)
	c.Assert(result, HasLen, 1)
	c.Assert(result["test"].Name, Equals, "test")
}

func (s *suite) Test_GetProjectsHandler_fails_if_GetProjects_fails(c *C) {
	provider := &projectHandlerProvider{
		GetProjects: func() (map[string]*types.Project, error) {
			return nil, types.NotFound
		},
	}
	resp := s.testGET(c, s.getProjectsMuxWithProvider(provider), GetProjectsURL)
	c.Assert(resp.StatusCode, Equals, 404)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "")
}

/*
	GetProjectHandler

*/

func (s *suite) getProjectMuxWithProvider(provider *projectHandlerProvider) *mux.Router {
	r := mux.NewRouter()
	router := r.Methods("GET").Subrouter()
	router.Handle(GetProjectURL, http.HandlerFunc(provider.GetProjectHandler))
	return r
}

func (s *suite) Test_GetProjectHandler_happy_path(c *C) {
	var capturedProject string
	provider := &projectHandlerProvider{
		GetProject: func(project string) (*model.ProjectPayload, error) {
			capturedProject = project
			return &model.ProjectPayload{
				Project: types.NewProject("test"),
			}, nil
		},
	}
	resp := s.testGET(c, s.getProjectMuxWithProvider(provider), getProjectTestURL)
	c.Assert(resp.StatusCode, Equals, 200)
	c.Assert(capturedProject, Equals, "project")
	c.Assert(resp.Header.Get("Content-Type"), Equals, "application/json")
	result := model.ProjectPayload{}
	c.Assert(json.NewDecoder(resp.Body).Decode(&result), IsNil)
	c.Assert(result.Project.Name, Equals, "test")
}

func (s *suite) Test_GetProjectHandler_fails_if_GetProject_fails(c *C) {
	provider := &projectHandlerProvider{
		GetProject: func(project string) (*model.ProjectPayload, error) {
			return nil, types.NotFound
		},
	}
	resp := s.testGET(c, s.getProjectMuxWithProvider(provider), getProjectTestURL)
	c.Assert(resp.StatusCode, Equals, 404)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "")
}

/*
	GetProjectHooksHandler

*/

func (s *suite) getProjectHooksMuxWithProvider(provider *projectHandlerProvider) *mux.Router {
	r := mux.NewRouter()
	router := r.Methods("GET").Subrouter()
	router.Handle(GetProjectHooksURL, http.HandlerFunc(provider.GetProjectHooksHandler))
	return r
}

func (s *suite) Test_GetProjectHooksHandler_happy_path(c *C) {
	var capturedProject string
	provider := &projectHandlerProvider{
		GetProjectHooks: func(project string) (types.Hooks, error) {
			capturedProject = project
			return types.Hooks{
				"test": map[string]string{},
			}, nil
		},
	}
	resp := s.testGET(c, s.getProjectHooksMuxWithProvider(provider), getProjectHooksTestURL)
	c.Assert(resp.StatusCode, Equals, 200)
	c.Assert(capturedProject, Equals, "project")
	c.Assert(resp.Header.Get("Content-Type"), Equals, "application/json")
	result := types.Hooks{}
	c.Assert(json.NewDecoder(resp.Body).Decode(&result), IsNil)
	c.Assert(result, HasLen, 1)
}

func (s *suite) Test_GetProjectHooksHandler_fails_if_GetProject_fails(c *C) {
	provider := &projectHandlerProvider{
		GetProjectHooks: func(project string) (types.Hooks, error) {
			return nil, types.NotFound
		},
	}
	resp := s.testGET(c, s.getProjectHooksMuxWithProvider(provider), getProjectHooksTestURL)
	c.Assert(resp.StatusCode, Equals, 404)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "")
}

/*
	AddProjectHandler

*/

func (s *suite) addProjectMuxWithProvider(provider *projectHandlerProvider) *mux.Router {
	r := mux.NewRouter()
	router := r.Methods("POST").Subrouter()
	router.Handle(AddProjectURL, http.HandlerFunc(provider.AddProjectHandler))
	return r
}

func (s *suite) Test_AddProjectHandler_happy_path(c *C) {
	provider := &projectHandlerProvider{
		AddProject: func(project *types.Project, username string) error {
			return nil
		},
	}
	data := types.NewProject("test")
	resp := s.testPOST(c, s.addProjectMuxWithProvider(provider), addProjectTestURL, data)
	c.Assert(resp.StatusCode, Equals, 200)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "")
}

func (s *suite) Test_AddProjectHandler_fails_if_invalid_json(c *C) {
	provider := &projectHandlerProvider{}
	resp := s.testPOST(c, s.addProjectMuxWithProvider(provider), addProjectTestURL, nil)
	c.Assert(resp.StatusCode, Equals, 400)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "Invalid JSON")
}

func (s *suite) Test_AddProjectHandler_fails_if_AddProject_fails(c *C) {
	provider := &projectHandlerProvider{
		AddProject: func(project *types.Project, username string) error {
			return types.AlreadyExists
		},
	}
	data := types.NewProject("test")
	resp := s.testPOST(c, s.addProjectMuxWithProvider(provider), addProjectTestURL, data)
	c.Assert(resp.StatusCode, Equals, 409)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "Resource already exists")
}

/*
	UpdateProjectHandler

*/

func (s *suite) updateProjectMuxWithProvider(provider *projectHandlerProvider) *mux.Router {
	r := mux.NewRouter()
	router := r.Methods("PUT").Subrouter()
	router.Handle(UpdateProjectURL, http.HandlerFunc(provider.UpdateProjectHandler))
	return r
}

func (s *suite) Test_UpdateProjectHandler_happy_path(c *C) {
	provider := &projectHandlerProvider{
		UpdateProject: func(project *types.Project) error {
			return nil
		},
	}
	data := types.NewProject("project")
	resp := s.testPUT(c, s.updateProjectMuxWithProvider(provider), updateProjectTestURL, data)
	c.Assert(resp.StatusCode, Equals, 201)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "")
}

func (s *suite) Test_UpdateProjectHandler_fails_if_invalid_json(c *C) {
	provider := &projectHandlerProvider{}
	resp := s.testPUT(c, s.updateProjectMuxWithProvider(provider), updateProjectTestURL, nil)
	c.Assert(resp.StatusCode, Equals, 400)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "Invalid JSON")
}

func (s *suite) Test_UpdateProjectHandler_fails_if_mux_name_doesnt_match_payload_name(c *C) {
	provider := &projectHandlerProvider{
		UpdateProject: func(project *types.Project) error {
			return types.AlreadyExists
		},
	}
	data := types.NewProject("wrong-name")
	resp := s.testPUT(c, s.updateProjectMuxWithProvider(provider), updateProjectTestURL, data)
	c.Assert(resp.StatusCode, Equals, 400)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "Project 'name' field doesn't correspond with URL")
}

func (s *suite) Test_UpdateProjectHandler_fails_if_UpdateProject_fails(c *C) {
	provider := &projectHandlerProvider{
		UpdateProject: func(project *types.Project) error {
			return types.AlreadyExists
		},
	}
	data := types.NewProject("project")
	resp := s.testPUT(c, s.updateProjectMuxWithProvider(provider), updateProjectTestURL, data)
	c.Assert(resp.StatusCode, Equals, 409)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "Resource already exists")
}

/*
	UpdateProjectHooksHandler

*/

func (s *suite) updateProjectHooksMuxWithProvider(provider *projectHandlerProvider) *mux.Router {
	r := mux.NewRouter()
	router := r.Methods("PUT").Subrouter()
	router.Handle(UpdateProjectHooksURL, http.HandlerFunc(provider.UpdateProjectHooksHandler))
	return r
}

func (s *suite) Test_UpdateProjectHooksHandler_happy_path(c *C) {
	provider := &projectHandlerProvider{
		UpdateProjectHooks: func(project string, hooks types.Hooks) error {
			return nil
		},
	}
	data := types.Hooks{}
	resp := s.testPUT(c, s.updateProjectHooksMuxWithProvider(provider), updateProjectHooksTestURL, data)
	c.Assert(resp.StatusCode, Equals, 201)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "")
}

func (s *suite) Test_UpdateProjectHooksHandler_fails_if_invalid_json(c *C) {
	provider := &projectHandlerProvider{}
	resp := s.testPUT(c, s.updateProjectHooksMuxWithProvider(provider), updateProjectHooksTestURL, nil)
	c.Assert(resp.StatusCode, Equals, 400)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "Invalid JSON")
}

func (s *suite) Test_UpdateProjectHooksHandler_fails_if_UpdateProjectHooks_fails(c *C) {
	provider := &projectHandlerProvider{
		UpdateProjectHooks: func(project string, hooks types.Hooks) error {
			return types.NotFound
		},
	}
	data := types.Hooks{}
	resp := s.testPUT(c, s.updateProjectHooksMuxWithProvider(provider), updateProjectHooksTestURL, data)
	c.Assert(resp.StatusCode, Equals, 404)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "")
}

/*
	HardDeleteProjectHandler
*/

func (s *suite) hardDeleteProjectMuxWithProvider(provider *projectHandlerProvider) *mux.Router {
	return s.GetMuxForHandler("DELETE", HardDeleteProjectHooksURL, provider.HardDeleteProjectHandler)
}

func (s *suite) Test_HardDeleteProjectHandler_happy_path(c *C) {
	var capturedProject string
	provider := &projectHandlerProvider{
		HardDeleteProject: func(project string) error {
			capturedProject = project
			return nil
		},
	}
	resp := s.testDELETE(c, s.hardDeleteProjectMuxWithProvider(provider), hardDeleteProjectHooksTestURL)
	s.ExpectSuccessResponse(c, resp, "")
	c.Assert(capturedProject, Equals, "project")
}

func (s *suite) Test_HardDeleteProjectHandler_fails_if_HardDeleteProject_fails(c *C) {
	provider := &projectHandlerProvider{
		HardDeleteProject: func(project string) error {
			return types.NotFound
		},
	}
	resp := s.testDELETE(c, s.hardDeleteProjectMuxWithProvider(provider), hardDeleteProjectHooksTestURL)
	s.ExpectErrorResponse(c, resp, 404, "")
}
