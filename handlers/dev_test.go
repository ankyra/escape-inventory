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
	"net/http/httptest"

	. "gopkg.in/check.v1"
)

func (s *suite) Test_DevHandler_WipeDatabase_CallsFunc(c *C) {
	req := httptest.NewRequest("GET", "/url", nil)
	w := httptest.NewRecorder()

	var called bool

	DevHandler{
		WipeDatabaseFunc: func() error {
			called = true
			return nil
		},
	}.WipeDatabase(w, req)

	c.Assert(called, Equals, true)

	c.Assert(w.Result().StatusCode, Equals, 200)
}
