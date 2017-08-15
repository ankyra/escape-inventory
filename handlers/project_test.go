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
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ankyra/escape-registry/dao"
	"github.com/ankyra/escape-registry/model"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type suite struct{}

var _ = Suite(&suite{})

func (s *suite) SetUpTest(c *C) {
	dao.TestSetup()
}

/*
	JSON TEST HELPERS
*/

func (s *suite) JsonNilBodyTest(c *C, handler http.HandlerFunc) {
	req := httptest.NewRequest("POST", "/url", nil)
	req.Body = nil
	w := httptest.NewRecorder()
	handler(w, req)
	resp := w.Result()
	c.Assert(resp.StatusCode, Equals, 400, Commentf("nil body should fail"))
}

func (s *suite) JsonEmptyBodyTest(c *C, handler http.HandlerFunc) {
	req := httptest.NewRequest("POST", "/api/v1/registry/add-project", nil)
	w := httptest.NewRecorder()
	handler(w, req)
	resp := w.Result()
	c.Assert(resp.StatusCode, Equals, 400, Commentf("empty body should fail"))
}

func (s *suite) JsonParseTesting(c *C, handler http.HandlerFunc) {
	s.JsonNilBodyTest(c, handler)
	s.JsonEmptyBodyTest(c, handler)
}

/*
	GET PROJECT
*/

func (s *suite) Test_GetProject_Not_Found(c *C) {
	req := httptest.NewRequest("GET", "/url", nil)
	w := httptest.NewRecorder()
	getProjectHandler(w, req, "test")

	resp := w.Result()
	c.Assert(resp.StatusCode, Equals, 404)
}

func (s *suite) Test_GetProject(c *C) {
	s.addProject(c, "test")
	req := httptest.NewRequest("GET", "/get", nil)
	w := httptest.NewRecorder()
	getProjectHandler(w, req, "test")

	resp := w.Result()
	c.Assert(resp.StatusCode, Equals, 200)

	result := model.ProjectPayload{}
	c.Assert(json.NewDecoder(resp.Body).Decode(&result), IsNil)

	c.Assert(result.Name, Equals, "test")
	c.Assert(result.Description, Equals, "")
	c.Assert(result.Units, HasLen, 0)
}

/*
	ADD PROJECT
*/

func (s *suite) Test_AddProject_Json_Tests(c *C) {
	s.JsonParseTesting(c, AddProjectHandler)
}

func (s *suite) Test_AddProject_Missing_Name(c *C) {
	body := `{}`
	req := httptest.NewRequest("POST", "/api/v1/registry/add-project", bytes.NewReader([]byte(body)))
	w := httptest.NewRecorder()
	AddProjectHandler(w, req)

	resp := w.Result()
	c.Assert(resp.StatusCode, Equals, 400)
}

func (s *suite) Test_AddProject_Already_Exists(c *C) {
	body := `{"name": "test"}`
	req := httptest.NewRequest("POST", "/api/v1/registry/add-project", bytes.NewReader([]byte(body)))
	w := httptest.NewRecorder()
	AddProjectHandler(w, req)

	resp := w.Result()
	c.Assert(resp.StatusCode, Equals, 200)

	req = httptest.NewRequest("POST", "/api/v1/registry/add-project", bytes.NewReader([]byte(body)))
	w = httptest.NewRecorder()
	AddProjectHandler(w, req)

	resp = w.Result()
	c.Assert(resp.StatusCode, Equals, 409)
}

/*
	ADD PROJECT HELPER
*/

func (s *suite) addProject(c *C, name string) {
	body := `{"name": "` + name + `"}`
	req := httptest.NewRequest("POST", "/api/v1/registry/add-project", bytes.NewReader([]byte(body)))
	w := httptest.NewRecorder()
	AddProjectHandler(w, req)
	resp := w.Result()
	c.Assert(resp.StatusCode, Equals, 200)
}

/*
	UPDATE PROJECT
*/

func (s *suite) Test_UpdateProject_Json_Test(c *C) {
	s.JsonParseTesting(c, UpdateProjectHandler)
}

func (s *suite) Test_UpdateProject_Missing_Name(c *C) {
	body := `{}`
	req := httptest.NewRequest("POST", "/api/v1/registry/update-project", bytes.NewReader([]byte(body)))
	w := httptest.NewRecorder()
	UpdateProjectHandler(w, req)

	resp := w.Result()
	c.Assert(resp.StatusCode, Equals, 400)
}

func (s *suite) Test_UpdateProject_Not_Found(c *C) {
	body := `{"name": "test"}`
	req := httptest.NewRequest("POST", "/api/v1/registry/update-project", bytes.NewReader([]byte(body)))
	w := httptest.NewRecorder()
	UpdateProjectHandler(w, req)

	resp := w.Result()
	c.Assert(resp.StatusCode, Equals, 404)
}

func (s *suite) Test_UpdateProject(c *C) {
	s.addProject(c, "test")
	body := `{"name": "test", "description": "test description"}`
	req := httptest.NewRequest("POST", "/api/v1/registry/update-project", bytes.NewReader([]byte(body)))
	w := httptest.NewRecorder()
	UpdateProjectHandler(w, req)

	resp := w.Result()
	c.Assert(resp.StatusCode, Equals, 201)
}
