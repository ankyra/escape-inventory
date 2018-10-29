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
	GetNamespacesURL = "/api/v1/inventory/"

	GetNamespaceURL     = "/api/v1/inventory/{namespace}/"
	getNamespaceTestURL = "/api/v1/inventory/namespace/"

	AddNamespaceURL     = "/api/v1/inventory/{namespace}/add-namespace"
	addNamespaceTestURL = "/api/v1/inventory/namespace/add-namespace"

	UpdateNamespaceURL     = "/api/v1/inventory/{namespace}/"
	updateNamespaceTestURL = "/api/v1/inventory/namespace/"

	GetNamespaceHooksURL     = "/api/v1/inventory/{namespace}/hooks/"
	getNamespaceHooksTestURL = "/api/v1/inventory/namespace/hooks/"

	UpdateNamespaceHooksURL     = "/api/v1/inventory/{namespace}/hooks/"
	updateNamespaceHooksTestURL = "/api/v1/inventory/namespace/hooks/"

	HardDeleteNamespaceHooksURL     = "/api/v1/inventory/{namespace}/hard-delete"
	hardDeleteNamespaceHooksTestURL = "/api/v1/inventory/namespace/hard-delete"
)

/*
	GetNamespacesHandler

*/

func (s *suite) getNamespacesMuxWithProvider(provider *namespaceHandlerProvider) *mux.Router {
	r := mux.NewRouter()
	router := r.Methods("GET").Subrouter()
	router.Handle(GetNamespacesURL, http.HandlerFunc(provider.GetNamespacesHandler))
	return r
}

func (s *suite) Test_GetNamespacesHandler_happy_path(c *C) {
	provider := &namespaceHandlerProvider{
		GetNamespaces: func() (map[string]*types.Project, error) {
			return map[string]*types.Project{
				"test": types.NewProject("test"),
			}, nil
		},
	}
	resp := s.testGET(c, s.getNamespacesMuxWithProvider(provider), GetNamespacesURL)
	c.Assert(resp.StatusCode, Equals, 200)
	c.Assert(resp.Header.Get("Content-Type"), Equals, "application/json")
	result := map[string]*types.Project{}
	c.Assert(json.NewDecoder(resp.Body).Decode(&result), IsNil)
	c.Assert(result, HasLen, 1)
	c.Assert(result["test"].Name, Equals, "test")
}

func (s *suite) Test_GetNamespacesHandler_fails_if_GetNamespaces_fails(c *C) {
	provider := &namespaceHandlerProvider{
		GetNamespaces: func() (map[string]*types.Project, error) {
			return nil, types.NotFound
		},
	}
	resp := s.testGET(c, s.getNamespacesMuxWithProvider(provider), GetNamespacesURL)
	c.Assert(resp.StatusCode, Equals, 404)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "")
}

/*
	GetNamespaceHandler

*/

func (s *suite) getNamespaceMuxWithProvider(provider *namespaceHandlerProvider) *mux.Router {
	r := mux.NewRouter()
	router := r.Methods("GET").Subrouter()
	router.Handle(GetNamespaceURL, http.HandlerFunc(provider.GetNamespaceHandler))
	return r
}

func (s *suite) Test_GetNamespaceHandler_happy_path(c *C) {
	var capturedNamespace string
	provider := &namespaceHandlerProvider{
		GetNamespace: func(namespace string) (*model.ProjectPayload, error) {
			capturedNamespace = namespace
			return &model.ProjectPayload{
				Project: types.NewProject("test"),
			}, nil
		},
	}
	resp := s.testGET(c, s.getNamespaceMuxWithProvider(provider), getNamespaceTestURL)
	c.Assert(resp.StatusCode, Equals, 200)
	c.Assert(capturedNamespace, Equals, "namespace")
	c.Assert(resp.Header.Get("Content-Type"), Equals, "application/json")
	result := model.ProjectPayload{}
	c.Assert(json.NewDecoder(resp.Body).Decode(&result), IsNil)
	c.Assert(result.Project.Name, Equals, "test")
}

func (s *suite) Test_GetNamespaceHandler_fails_if_GetNamespace_fails(c *C) {
	provider := &namespaceHandlerProvider{
		GetNamespace: func(namespace string) (*model.ProjectPayload, error) {
			return nil, types.NotFound
		},
	}
	resp := s.testGET(c, s.getNamespaceMuxWithProvider(provider), getNamespaceTestURL)
	c.Assert(resp.StatusCode, Equals, 404)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "")
}

/*
	GetNamespaceHooksHandler

*/

func (s *suite) getNamespaceHooksHandlerMuxWithProvider(provider *namespaceHandlerProvider) *mux.Router {
	r := mux.NewRouter()
	router := r.Methods("GET").Subrouter()
	router.Handle(GetNamespaceHooksURL, http.HandlerFunc(provider.GetNamespaceHooksHandler))
	return r
}

func (s *suite) Test_GetNamespaceHooksHandler_happy_path(c *C) {
	var capturedNamespace string
	provider := &namespaceHandlerProvider{
		GetNamespaceHooks: func(namespace string) (types.Hooks, error) {
			capturedNamespace = namespace
			return types.Hooks{
				"test": map[string]string{},
			}, nil
		},
	}
	resp := s.testGET(c, s.getNamespaceHooksHandlerMuxWithProvider(provider), getNamespaceHooksTestURL)
	c.Assert(resp.StatusCode, Equals, 200)
	c.Assert(capturedNamespace, Equals, "namespace")
	c.Assert(resp.Header.Get("Content-Type"), Equals, "application/json")
	result := types.Hooks{}
	c.Assert(json.NewDecoder(resp.Body).Decode(&result), IsNil)
	c.Assert(result, HasLen, 1)
}

func (s *suite) Test_GetNamespaceHooksHandler_fails_if_GetNamespaceHooks_fails(c *C) {
	provider := &namespaceHandlerProvider{
		GetNamespaceHooks: func(namespace string) (types.Hooks, error) {
			return nil, types.NotFound
		},
	}
	resp := s.testGET(c, s.getNamespaceHooksHandlerMuxWithProvider(provider), getNamespaceHooksTestURL)
	c.Assert(resp.StatusCode, Equals, 404)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "")
}

/*
	AddNamespaceHandler

*/

func (s *suite) addNamespaceMuxWithProvider(provider *namespaceHandlerProvider) *mux.Router {
	r := mux.NewRouter()
	router := r.Methods("POST").Subrouter()
	router.Handle(AddNamespaceURL, http.HandlerFunc(provider.AddNamespaceHandler))
	return r
}

func (s *suite) Test_AddNamespaceHandler_happy_path(c *C) {
	provider := &namespaceHandlerProvider{
		AddNamespace: func(namespace *types.Project, username string) error {
			return nil
		},
	}
	data := types.NewProject("test")
	resp := s.testPOST(c, s.addNamespaceMuxWithProvider(provider), addNamespaceTestURL, data)
	c.Assert(resp.StatusCode, Equals, 200)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "")
}

func (s *suite) Test_AddNamespaceHandler_fails_if_invalid_json(c *C) {
	provider := &namespaceHandlerProvider{}
	resp := s.testPOST(c, s.addNamespaceMuxWithProvider(provider), addNamespaceTestURL, nil)
	c.Assert(resp.StatusCode, Equals, 400)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "Invalid JSON")
}

func (s *suite) Test_AddNamespaceHandler_fails_if_AddNamespace_fails(c *C) {
	provider := &namespaceHandlerProvider{
		AddNamespace: func(namespace *types.Project, username string) error {
			return types.AlreadyExists
		},
	}
	data := types.NewProject("test")
	resp := s.testPOST(c, s.addNamespaceMuxWithProvider(provider), addNamespaceTestURL, data)
	c.Assert(resp.StatusCode, Equals, 409)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "Resource already exists")
}

/*
	UpdateNamespaceHandler

*/

func (s *suite) updateNamespaceMuxWithProvider(provider *namespaceHandlerProvider) *mux.Router {
	r := mux.NewRouter()
	router := r.Methods("PUT").Subrouter()
	router.Handle(UpdateNamespaceURL, http.HandlerFunc(provider.UpdateNamespaceHandler))
	return r
}

func (s *suite) Test_UpdateNamespaceHandler_happy_path(c *C) {
	provider := &namespaceHandlerProvider{
		UpdateNamespace: func(namespace *types.Project) error {
			return nil
		},
	}
	data := types.NewProject("namespace")
	resp := s.testPUT(c, s.updateNamespaceMuxWithProvider(provider), updateNamespaceTestURL, data)
	c.Assert(resp.StatusCode, Equals, 201)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "")
}

func (s *suite) Test_UpdateNamespaceHandler_fails_if_invalid_json(c *C) {
	provider := &namespaceHandlerProvider{}
	resp := s.testPUT(c, s.updateNamespaceMuxWithProvider(provider), updateNamespaceTestURL, nil)
	c.Assert(resp.StatusCode, Equals, 400)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "Invalid JSON")
}

func (s *suite) Test_UpdateNamespaceHandler_fails_if_mux_name_doesnt_match_payload_name(c *C) {
	provider := &namespaceHandlerProvider{
		UpdateNamespace: func(namespace *types.Project) error {
			return types.AlreadyExists
		},
	}
	data := types.NewProject("wrong-name")
	resp := s.testPUT(c, s.updateNamespaceMuxWithProvider(provider), updateNamespaceTestURL, data)
	c.Assert(resp.StatusCode, Equals, 400)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "Namespace 'name' field doesn't correspond with URL")
}

func (s *suite) Test_UpdateNamespaceHandler_fails_if_UpdateNamespace_fails(c *C) {
	provider := &namespaceHandlerProvider{
		UpdateNamespace: func(namespace *types.Project) error {
			return types.AlreadyExists
		},
	}
	data := types.NewProject("namespace")
	resp := s.testPUT(c, s.updateNamespaceMuxWithProvider(provider), updateNamespaceTestURL, data)
	c.Assert(resp.StatusCode, Equals, 409)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "Resource already exists")
}

/*
	UpdateNamespaceHooksHandler

*/

func (s *suite) updateNamespaceHooksHandlerMuxWithProvider(provider *namespaceHandlerProvider) *mux.Router {
	r := mux.NewRouter()
	router := r.Methods("PUT").Subrouter()
	router.Handle(UpdateNamespaceHooksURL, http.HandlerFunc(provider.UpdateNamespaceHooksHandler))
	return r
}

func (s *suite) Test_UpdateNamespaceHooksHandler_happy_path(c *C) {
	provider := &namespaceHandlerProvider{
		UpdateNamespaceHooks: func(namespace string, hooks types.Hooks) error {
			return nil
		},
	}
	data := types.Hooks{}
	resp := s.testPUT(c, s.updateNamespaceHooksHandlerMuxWithProvider(provider), updateNamespaceHooksTestURL, data)
	c.Assert(resp.StatusCode, Equals, 201)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "")
}

func (s *suite) Test_UpdateNamespaceHooksHandler_fails_if_invalid_json(c *C) {
	provider := &namespaceHandlerProvider{}
	resp := s.testPUT(c, s.updateNamespaceHooksHandlerMuxWithProvider(provider), updateNamespaceHooksTestURL, nil)
	c.Assert(resp.StatusCode, Equals, 400)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "Invalid JSON")
}

func (s *suite) Test_UpdateNamespaceHooksHandler_fails_if_UpdateNamespaceHooks_fails(c *C) {
	provider := &namespaceHandlerProvider{
		UpdateNamespaceHooks: func(namespace string, hooks types.Hooks) error {
			return types.NotFound
		},
	}
	data := types.Hooks{}
	resp := s.testPUT(c, s.updateNamespaceHooksHandlerMuxWithProvider(provider), updateNamespaceHooksTestURL, data)
	c.Assert(resp.StatusCode, Equals, 404)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "")
}

/*
	HardDeleteNamespaceHandler
*/

func (s *suite) hardDeleteNamespaceHandlerMuxWithProvider(provider *namespaceHandlerProvider) *mux.Router {
	return s.GetMuxForHandler("DELETE", HardDeleteNamespaceHooksURL, provider.HardDeleteNamespaceHandler)
}

func (s *suite) Test_HardDeleteNamespaceHandler_happy_path(c *C) {
	var capturedNamespace string
	provider := &namespaceHandlerProvider{
		HardDeleteNamespace: func(namespace string) error {
			capturedNamespace = namespace
			return nil
		},
	}
	resp := s.testDELETE(c, s.hardDeleteNamespaceHandlerMuxWithProvider(provider), hardDeleteNamespaceHooksTestURL)
	s.ExpectSuccessResponse(c, resp, "")
	c.Assert(capturedNamespace, Equals, "namespace")
}

func (s *suite) Test_HardDeleteNamespaceHandler_fails_if_HardDeleteNamespace_fails(c *C) {
	provider := &namespaceHandlerProvider{
		HardDeleteNamespace: func(namespace string) error {
			return types.NotFound
		},
	}
	resp := s.testDELETE(c, s.hardDeleteNamespaceHandlerMuxWithProvider(provider), hardDeleteNamespaceHooksTestURL)
	s.ExpectErrorResponse(c, resp, 404, "")
}
