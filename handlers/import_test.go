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
	"io/ioutil"
	"net/http/httptest"
	"strings"

	. "gopkg.in/check.v1"
)

func (s *suite) Test_ImportHandler_ImportReleases_CallsFunc(c *C) {
	req := httptest.NewRequest("GET", "/url", strings.NewReader(`[{"test": 1}]`))
	w := httptest.NewRecorder()

	var called bool
	var input []map[string]interface{}

	ImportHandler{
		ImportReleasesFunc: func(releases []map[string]interface{}) error {
			called = true
			input = releases
			return nil
		},
	}.ImportReleases(w, req)

	c.Assert(called, Equals, true)
	c.Assert(input, NotNil)
	c.Assert(input, HasLen, 1)
	c.Assert(input[0]["test"], Equals, float64((1)))

	c.Assert(w.Result().StatusCode, Equals, 200)

	body, _ := ioutil.ReadAll(w.Result().Body)
	c.Assert(string(body), Equals, "")

}

func (s *suite) Test_ImportHandler_ImportReleases_Errors_NoBody(c *C) {
	req := httptest.NewRequest("GET", "/url", nil)
	w := httptest.NewRecorder()

	ImportHandler{
		ImportReleasesFunc: func(releases []map[string]interface{}) error {
			return nil
		},
	}.ImportReleases(w, req)

	c.Assert(w.Result().StatusCode, Equals, 400)

	body, _ := ioutil.ReadAll(w.Result().Body)
	c.Assert(string(body), Equals, "unexpected end of JSON input")
}

func (s *suite) Test_ImportHandler_ImportReleases_Errors_IncorrectBodySchema(c *C) {
	req := httptest.NewRequest("GET", "/url", strings.NewReader(`{}`))
	w := httptest.NewRecorder()

	ImportHandler{
		ImportReleasesFunc: func(releases []map[string]interface{}) error {
			return nil
		},
	}.ImportReleases(w, req)

	c.Assert(w.Result().StatusCode, Equals, 400)

	body, _ := ioutil.ReadAll(w.Result().Body)
	c.Assert(string(body), Equals, "json: cannot unmarshal object into Go value of type []map[string]interface {}")
}
