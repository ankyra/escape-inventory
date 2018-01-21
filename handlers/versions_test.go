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

	core "github.com/ankyra/escape-core"
	"github.com/ankyra/escape-inventory/dao/types"
	"github.com/ankyra/escape-inventory/model"
	"github.com/gorilla/mux"
	. "gopkg.in/check.v1"
)

const (
	GetVersionURL          = "/api/v1/inventory/{project}/units/{name}/versions/{version}/"
	getVersionTestURL      = "/api/v1/inventory/project/units/name/versions/v1.0/"
	NextVersionURL         = "/api/v1/inventory/{project}/units/{name}/next-version"
	nextVersionTestURL     = "/api/v1/inventory/project/units/name/next-version?prefix=0.1"
	PreviousVersionURL     = "/api/v1/inventory/{project}/units/{name}/versions/{version}/previous/"
	previousVersionTestURL = "/api/v1/inventory/project/units/name/versions/v1.0/previous/"
	DiffURL                = "/api/v1/inventory/{project}/units/{name}/versions/{version}/diff/{diffWith}/"
	diffTestURL            = "/api/v1/inventory/project/units/name/versions/v1.0/diff/v1.1/"
)

/*
	GetVersionHandler

*/

func (s *suite) getVersionMuxWithProvider(provider *versionHandlerProvider) *mux.Router {
	r := mux.NewRouter()
	router := r.Methods("GET").Subrouter()
	router.Handle(GetVersionURL, http.HandlerFunc(provider.GetVersionHandler))
	return r
}

func (s *suite) Test_GetVersionHandler_happy_path(c *C) {
	var capturedProject, capturedName, capturedVersion string
	provider := &versionHandlerProvider{
		GetReleaseMetadata: func(project, name, version string) (*core.ReleaseMetadata, error) {
			capturedProject = project
			capturedName = name
			capturedVersion = version
			return core.NewReleaseMetadata("name", "v10.0"), nil
		},
	}
	resp := s.testGET(c, s.getVersionMuxWithProvider(provider), getVersionTestURL)
	c.Assert(resp.StatusCode, Equals, 200)
	c.Assert(resp.Header.Get("Content-Type"), Equals, "application/json")
	c.Assert(capturedProject, Equals, "project")
	c.Assert(capturedName, Equals, "name")
	c.Assert(capturedVersion, Equals, "v1.0")

	result := core.ReleaseMetadata{}
	c.Assert(json.NewDecoder(resp.Body).Decode(&result), IsNil)
	c.Assert(result.Name, Equals, "name")
	c.Assert(result.Version, Equals, "v10.0")
}

func (s *suite) Test_GetVersionHandler_happy_path_full(c *C) {
	var capturedProject, capturedName, capturedVersion string
	provider := &versionHandlerProvider{
		GetRelease: func(project, name, version string) (*model.ReleasePayload, error) {
			capturedProject = project
			capturedName = name
			capturedVersion = version
			return &model.ReleasePayload{Downloads: 2000}, nil
		},
	}
	resp := s.testGET(c, s.getVersionMuxWithProvider(provider), getVersionTestURL+"?full=true")
	c.Assert(resp.StatusCode, Equals, 200)
	c.Assert(resp.Header.Get("Content-Type"), Equals, "application/json")
	c.Assert(capturedProject, Equals, "project")
	c.Assert(capturedName, Equals, "name")
	c.Assert(capturedVersion, Equals, "v1.0")

	result := model.ReleasePayload{}
	c.Assert(json.NewDecoder(resp.Body).Decode(&result), IsNil)
	c.Assert(result.Downloads, Equals, 2000)
}

func (s *suite) Test_GetVersionHandler_fails_if_full_and_GetRelease_fails(c *C) {
	provider := &versionHandlerProvider{
		GetRelease: func(project, name, version string) (*model.ReleasePayload, error) {
			return nil, types.NotFound
		},
	}
	resp := s.testGET(c, s.getVersionMuxWithProvider(provider), getVersionTestURL+"?full=true")
	c.Assert(resp.StatusCode, Equals, 404)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "")
}

func (s *suite) Test_GetVersionHandler_fails_if_not_full_and_GetReleaseMetadata_fails(c *C) {
	provider := &versionHandlerProvider{
		GetReleaseMetadata: func(project, name, version string) (*core.ReleaseMetadata, error) {
			return nil, types.NotFound
		},
	}
	resp := s.testGET(c, s.getVersionMuxWithProvider(provider), getVersionTestURL)
	c.Assert(resp.StatusCode, Equals, 404)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "")
}

/*
	NextVersionHandler

*/

func (s *suite) nextVersionMuxWithProvider(provider *versionHandlerProvider) *mux.Router {
	r := mux.NewRouter()
	router := r.Methods("GET").Subrouter()
	router.Handle(NextVersionURL, http.HandlerFunc(provider.NextVersionHandler))
	return r
}

func (s *suite) Test_NextVersionHandler_happy_path(c *C) {
	var capturedProject, capturedName, capturedPrefix string
	provider := &versionHandlerProvider{
		GetNextVersion: func(project, name, prefix string) (string, error) {
			capturedProject = project
			capturedName = name
			capturedPrefix = prefix
			return "0.9", nil
		},
	}
	resp := s.testGET(c, s.nextVersionMuxWithProvider(provider), nextVersionTestURL)
	c.Assert(resp.StatusCode, Equals, 200)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(capturedProject, Equals, "project")
	c.Assert(capturedName, Equals, "name")
	c.Assert(capturedPrefix, Equals, "0.1")
	c.Assert(string(body), Equals, "0.9")
}

func (s *suite) Test_NextVersionHandler_fails_if_GetNextVersion_fails(c *C) {
	provider := &versionHandlerProvider{
		GetNextVersion: func(project, name, prefix string) (string, error) {
			return "", types.NotFound
		},
	}
	resp := s.testGET(c, s.nextVersionMuxWithProvider(provider), nextVersionTestURL)
	c.Assert(resp.StatusCode, Equals, 404)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "")
}

/*
	PreviousVersionHandler
*/

func (s *suite) previousVersionMuxWithProvider(provider *versionHandlerProvider) *mux.Router {
	r := mux.NewRouter()
	router := r.Methods("GET").Subrouter()
	router.Handle(PreviousVersionURL, http.HandlerFunc(provider.PreviousVersionHandler))
	return r
}

func (s *suite) Test_PreviousVersionHandler_happy_path(c *C) {
	var capturedProject, capturedName, capturedVersion string
	provider := &versionHandlerProvider{
		GetPreviousVersion: func(project, name, version string) (*core.ReleaseMetadata, error) {
			capturedProject = project
			capturedName = name
			capturedVersion = version
			return core.NewReleaseMetadata("name", "v0.9"), nil
		},
	}
	resp := s.testGET(c, s.previousVersionMuxWithProvider(provider), previousVersionTestURL)
	c.Assert(resp.StatusCode, Equals, 200)
	c.Assert(resp.Header.Get("Content-Type"), Equals, "application/json")
	c.Assert(capturedProject, Equals, "project")
	c.Assert(capturedName, Equals, "name")
	c.Assert(capturedVersion, Equals, "v1.0")

	result := core.ReleaseMetadata{}
	c.Assert(json.NewDecoder(resp.Body).Decode(&result), IsNil)
	c.Assert(result.Name, Equals, "name")
	c.Assert(result.Version, Equals, "v0.9")
}

func (s *suite) Test_PreviousVersionHandler_fails_if_GetNextVersion_fails(c *C) {
	provider := &versionHandlerProvider{
		GetPreviousVersion: func(project, name, prefix string) (*core.ReleaseMetadata, error) {
			return nil, types.NotFound
		},
	}
	resp := s.testGET(c, s.previousVersionMuxWithProvider(provider), previousVersionTestURL)
	c.Assert(resp.StatusCode, Equals, 404)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "")
}

/*
	DiffHandler
*/

func (s *suite) diffMuxWithProvider(provider *versionHandlerProvider) *mux.Router {
	r := mux.NewRouter()
	router := r.Methods("GET").Subrouter()
	router.Handle(DiffURL, http.HandlerFunc(provider.DiffHandler))
	return r
}

func (s *suite) Test_DiffHandler_happy_path(c *C) {
	var capturedProject, capturedName, capturedVersion, capturedDiffWithVersion string
	provider := &versionHandlerProvider{
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
	provider := &versionHandlerProvider{
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
