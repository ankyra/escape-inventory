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
	DownstreamURL          = "/api/v1/inventory/{project}/units/{name}/versions/{version}/downstream"
	downstreamTestURL      = "/api/v1/inventory/project/units/name/versions/v1.0.0/downstream"
	DependencyGraphURL     = "/api/v1/inventory/{project}/units/{name}/versions/{version}/dependency-graph"
	dependencyGraphTestURL = "/api/v1/inventory/project/units/name/versions/v1.0.0/dependency-graph"
)

/*
	DownstreamHandler
*/

func (s *suite) downstreamMuxWithProvider(provider *dependencyHandlerProvider) *mux.Router {
	r := mux.NewRouter()
	router := r.Methods("GET").Subrouter()
	router.Handle(DownstreamURL, http.HandlerFunc(provider.DownstreamHandler))
	return r
}

func (s *suite) Test_DownstreamHandler_happy_path(c *C) {
	provider := &dependencyHandlerProvider{
		GetDownstreamDependencies: func(project, name, version string) ([]*types.Dependency, error) {
			deps := []*types.Dependency{types.NewDependency("prj", "dep", "1.0")}
			return deps, nil
		},
	}
	resp := s.testGET(c, s.downstreamMuxWithProvider(provider), downstreamTestURL)
	c.Assert(resp.StatusCode, Equals, 200)
	result := []*types.Dependency{}
	c.Assert(json.NewDecoder(resp.Body).Decode(&result), IsNil)
	c.Assert(result, HasLen, 1)
	c.Assert(result[0], DeepEquals, types.NewDependency("prj", "dep", "1.0"))
	c.Assert(resp.Header.Get("Content-Type"), Equals, "application/json")
}

func (s *suite) Test_DownstreamHandler_fails_if_GetDownStreamDependencies_fails(c *C) {
	provider := &dependencyHandlerProvider{
		GetDownstreamDependencies: func(project, name, version string) ([]*types.Dependency, error) {
			return nil, types.NotFound
		},
	}
	resp := s.testGET(c, s.downstreamMuxWithProvider(provider), downstreamTestURL)
	c.Assert(resp.StatusCode, Equals, 404)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "")
}

/*
	DependencyGraphHandler
*/

func (s *suite) dependencyGraphMuxWithProvider(provider *dependencyHandlerProvider) *mux.Router {
	r := mux.NewRouter()
	router := r.Methods("GET").Subrouter()
	router.Handle(DependencyGraphURL, http.HandlerFunc(provider.DependencyGraphHandler))
	return r
}

func (s *suite) Test_DependencyGraphHandler_happy_path(c *C) {
	provider := &dependencyHandlerProvider{
		GetDependencyGraph: func(project, name, version string, downstreamFunc model.DownstreamDependenciesResolver) (*model.DependencyGraph, error) {
			return nil, nil
		},
	}
	resp := s.testGET(c, s.dependencyGraphMuxWithProvider(provider), dependencyGraphTestURL)
	c.Assert(resp.StatusCode, Equals, 200)
	result := model.DependencyGraph{}
	c.Assert(json.NewDecoder(resp.Body).Decode(&result), IsNil)
	c.Assert(result.Nodes, HasLen, 0)
	c.Assert(result.Edges, HasLen, 0)
	c.Assert(resp.Header.Get("Content-Type"), Equals, "application/json")
}

func (s *suite) Test_DependencyGraphHandler_fails_if_GetDependencyGraph_fails(c *C) {
	provider := &dependencyHandlerProvider{
		GetDependencyGraph: func(project, name, version string, downstreamFunc model.DownstreamDependenciesResolver) (*model.DependencyGraph, error) {
			return nil, types.NotFound
		},
	}
	resp := s.testGET(c, s.dependencyGraphMuxWithProvider(provider), dependencyGraphTestURL)
	c.Assert(resp.StatusCode, Equals, 404)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "")
}
