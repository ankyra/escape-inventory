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
	"net/http"

	"github.com/gorilla/mux"
	. "gopkg.in/check.v1"
)

func (s *suite) Test_HealthCheckHandler(c *C) {
	r := mux.NewRouter()
	router := r.Methods("GET").Subrouter()
	router.Handle("/", http.HandlerFunc(HealthCheckHandler))
	resp := s.testGET(c, r, "/")
	c.Assert(resp.StatusCode, Equals, 200)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "")
}
