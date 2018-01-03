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
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ankyra/escape-inventory/dao"
	"github.com/ankyra/escape-inventory/dao/types"
	"github.com/ankyra/escape-inventory/model"
	"github.com/gorilla/mux"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type suite struct{}

var _ = Suite(&suite{})

func (s *suite) SetUpTest(c *C) {
	dao.TestSetup()
}

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

func (s *suite) Test_HandleError_nil_error(c *C) {
	rr := httptest.NewRecorder()
	HandleError(rr, nil, nil)
	c.Assert(rr.Code, Equals, 500)
	c.Assert(rr.Body.String(), Equals, "")
}

func (s *suite) Test_HandleError_not_found(c *C) {
	rr := httptest.NewRecorder()
	HandleError(rr, nil, types.NotFound)
	c.Assert(rr.Code, Equals, 404)
	c.Assert(rr.Body.String(), Equals, "")
}

func (s *suite) Test_HandleError_already_exists(c *C) {
	rr := httptest.NewRecorder()
	HandleError(rr, nil, types.AlreadyExists)
	c.Assert(rr.Code, Equals, 409)
	c.Assert(rr.Body.String(), Equals, "Resource already exists")
}

func (s *suite) Test_HandleError_limit_reached(c *C) {
	rr := httptest.NewRecorder()
	HandleError(rr, nil, types.LimitError)
	c.Assert(rr.Code, Equals, 402)
	c.Assert(rr.Body.String(), Equals, "Plan limit exceeded")
}

func (s *suite) Test_HandleError_user_error(c *C) {
	rr := httptest.NewRecorder()
	HandleError(rr, nil, model.NewUserError(errors.New("ouch")))
	c.Assert(rr.Code, Equals, 400)
	c.Assert(rr.Body.String(), Equals, "ouch")
}

func (s *suite) Test_HandleError_server_error(c *C) {
	rr := httptest.NewRecorder()
	HandleError(rr, nil, errors.New("server error"))
	c.Assert(rr.Code, Equals, 500)
	c.Assert(rr.Body.String(), Equals, "")
}

func (s *suite) Test_ReadJsonBodyOrFail(c *C) {
	rr := httptest.NewRecorder()
	body := []byte(`{"test": "uh"}`)
	req, err := http.NewRequest("POST", "something", bytes.NewReader(body))
	c.Assert(err, IsNil)
	var result struct{ Test string }
	ReadJsonBodyOrFail(rr, req, &result)
	c.Assert(rr.Code, Equals, 200)
	c.Assert(result.Test, Equals, "uh")
}

func (s *suite) Test_ReadJsonBodyOrFail_empty_body(c *C) {
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "something", nil)
	c.Assert(err, IsNil)
	ReadJsonBodyOrFail(rr, req, nil)
	c.Assert(rr.Code, Equals, 400)
	c.Assert(rr.Body.String(), Equals, "Empty body")
}

func (s *suite) Test_ReadJsonBodyOrFail_decode_failure(c *C) {
	rr := httptest.NewRecorder()
	body := []byte(`{"test": "uh"}`)
	req, err := http.NewRequest("POST", "something", bytes.NewReader(body))
	c.Assert(err, IsNil)
	var result string
	ReadJsonBodyOrFail(rr, req, &result)
	c.Assert(rr.Code, Equals, 400)
	c.Assert(rr.Body.String(), Equals, "Invalid JSON")
}
