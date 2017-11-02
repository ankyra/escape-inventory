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
	"net/http/httptest"

	"github.com/ankyra/escape-inventory/dao/types"

	. "gopkg.in/check.v1"
)

func (s *suite) Test_GetApplication_NotFound(c *C) {
	req := httptest.NewRequest("GET", "/url", nil)
	w := httptest.NewRecorder()
	getApplicationHandler(w, req, "prj", "name")

	resp := w.Result()
	c.Assert(resp.StatusCode, Equals, 404)
}

func (s *suite) Test_GetApplications(c *C) {
	s.addProject(c, "test")
	req := httptest.NewRequest("GET", "/get", nil)
	w := httptest.NewRecorder()
	getApplicationsHandler(w, req, "test")

	resp := w.Result()
	c.Assert(resp.StatusCode, Equals, 200)

	result := map[string]*types.Application{}
	c.Assert(json.NewDecoder(resp.Body).Decode(&result), IsNil)

	c.Assert(result, HasLen, 0)
}
