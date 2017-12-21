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
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/ankyra/escape-inventory/dao/types"
	"github.com/gorilla/mux"
	. "gopkg.in/check.v1"
)

func (s *suite) testGET(c *C, r *mux.Router, url string) *http.Response {
	var reader io.Reader = nil
	req := httptest.NewRequest("GET", url, reader)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w.Result()
}

func (s *suite) testPUT(c *C, r *mux.Router, url string, data interface{}) *http.Response {
	var reader io.Reader = nil
	if data != nil {
		payload, err := json.Marshal(data)
		c.Assert(err, IsNil)
		reader = bytes.NewReader(payload)
	}
	req := httptest.NewRequest("PUT", url, reader)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w.Result()
}

func (s *suite) testPOST(c *C, r *mux.Router, url string, data interface{}) *http.Response {
	var reader io.Reader = nil
	if data != nil {
		payload, err := json.Marshal(data)
		c.Assert(err, IsNil)
		reader = bytes.NewReader(payload)
	}
	req := httptest.NewRequest("POST", url, reader)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w.Result()
}

func (s *suite) testPOST_file(c *C, r *mux.Router, url, path string) *http.Response {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	if path != "" {
		fileWriter, err := bodyWriter.CreateFormFile("file", path)
		c.Assert(err, IsNil)
		fh, err := os.Open(path)
		c.Assert(err, IsNil)
		_, err = io.Copy(fileWriter, fh)
		c.Assert(err, IsNil)
	}
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	req, err := http.NewRequest("POST", url, bodyBuf)
	c.Assert(err, IsNil)
	req.Header.Add("Content-Type", contentType)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w.Result()
}

func (s *suite) Test_JsonSuccess(c *C) {
	rr := httptest.NewRecorder()
	result := map[string]string{"test": "yo"}
	JsonSuccess(rr, result)
	c.Assert(rr.Code, Equals, 200)
	c.Assert(rr.Body.String(), Equals, "{\"test\":\"yo\"}\n")
}

func (s *suite) Test_ErrorOrJsonSuccess_success(c *C) {
	rr := httptest.NewRecorder()
	result := map[string]string{"test": "yo"}
	ErrorOrJsonSuccess(rr, nil, result, nil)
	c.Assert(rr.Code, Equals, 200)
	c.Assert(rr.Body.String(), Equals, "{\"test\":\"yo\"}\n")
}

func (s *suite) Test_ErrorOrJsonSuccess_error(c *C) {
	rr := httptest.NewRecorder()
	result := map[string]string{"test": "yo"}
	ErrorOrJsonSuccess(rr, nil, result, types.AlreadyExists)
	c.Assert(rr.Code, Equals, 409)
	c.Assert(rr.Body.String(), Equals, "Resource already exists")
}
