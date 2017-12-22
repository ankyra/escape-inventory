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

	core "github.com/ankyra/escape-core"
	"github.com/ankyra/escape-inventory/dao/types"
	"github.com/gorilla/mux"
	. "gopkg.in/check.v1"
)

const (
	DiffURL     = "/api/v1/registry/{project}/units/{name}/versions/{version}/diff/{diffWith}/"
	diffTestURL = "/api/v1/registry/project/units/name/versions/v1.0/diff/v1.1/"
)

/*
	DiffHandler

*/

func (s *suite) diffMuxWithProvider(provider *diffHandlerProvider) *mux.Router {
	r := mux.NewRouter()
	router := r.Methods("GET").Subrouter()
	router.Handle(DiffURL, http.HandlerFunc(provider.DiffHandler))
	return r
}

func (s *suite) Test_DiffHandler_happy_path(c *C) {
	var capturedProject, capturedName, capturedVersion, capturedDiffWithVersion string
	provider := &diffHandlerProvider{
		Diff: func(project, name, version, diffWithVersion string) (map[string]map[string]core.Changes, error) {
			capturedProject = project
			capturedName = name
			capturedVersion = version
			capturedDiffWithVersion = diffWithVersion
			changes := map[string]map[string]core.Changes{
				"test": map[string]core.Changes{},
			}
			return changes, nil
		},
	}
	resp := s.testGET(c, s.diffMuxWithProvider(provider), diffTestURL)
	c.Assert(resp.StatusCode, Equals, 200)
	c.Assert(capturedProject, Equals, "project")
	c.Assert(capturedName, Equals, "name")
	c.Assert(capturedVersion, Equals, "v1.0")
	c.Assert(capturedDiffWithVersion, Equals, "v1.1")

	result := map[string]map[string]core.Changes{}
	c.Assert(json.NewDecoder(resp.Body).Decode(&result), IsNil)
	c.Assert(result, HasLen, 1)
}

func (s *suite) Test_DiffHandler_fails_if_Diff_fails(c *C) {
	provider := &diffHandlerProvider{
		Diff: func(project, name, version, diffWithVersion string) (map[string]map[string]core.Changes, error) {
			return nil, types.NotFound
		},
	}
	resp := s.testGET(c, s.diffMuxWithProvider(provider), diffTestURL)
	c.Assert(resp.StatusCode, Equals, 404)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "")
}
