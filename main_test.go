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

package main

import (
	"bytes"
	"encoding/json"
	"github.com/ankyra/escape-registry/cmd"
	"github.com/ankyra/escape-registry/dao"
	. "gopkg.in/check.v1"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type suite struct{}

var _ = Suite(&suite{})

var handler http.Handler
var rr *httptest.ResponseRecorder

func (s *suite) SetUpTest(c *C) {
	dao.TestSetup()
	handler = cmd.GetHandler(getMux())
	rr = httptest.NewRecorder()
}

const (
	registerEndpoint = "/a/my-project/register"

	getProjectsEndpoints = "/a/"

	applicationsTestProject       = "applications-test-prj"
	applicationsEndpoint          = "/a/" + applicationsTestProject + "/"
	applicationsEndpointNoProject = "/a/doesnt-exist/"

	applicationVersionsTestProject = "versions-test-prj"
	applicationVersionsEndpoint    = "/a/" + applicationVersionsTestProject + "/my-app/"
	applicationVersionsNoProject   = "/a/doesnt-exist/my-app/"
	applicationVersionsNoApp       = "/a/versions-test/doesnt-exist/"

	nextVersionProject  = "next-version-prj"
	nextVersionEndpoint = "/a/" + nextVersionProject + "/my-app/next-version"

	getVersionProject        = "get-version-prj"
	getVersionEndpoint       = "/a/" + getVersionProject + "/my-app/v1/"
	getLatestVersionEndpoint = "/a/" + getVersionProject + "/my-app/latest/"
	getAutoVersionEndpoint   = "/a/" + getVersionProject + "/my-app/v0.0.@/"
	getPreviousEndpoint      = "/a/" + getVersionProject + "/my-app/v0.0.2/previous/"
	getPreviousEndpoint2     = "/a/" + getVersionProject + "/my-app/v0.0.1/previous/"
	getPreviousEndpoint3     = "/a/" + getVersionProject + "/my-app/v0.0.@/previous/"
	getDiffEndpoint          = "/a/" + getVersionProject + "/my-app/v0.0.2/diff/"
	getDiffWithEndpoint      = "/a/" + getVersionProject + "/my-app/v0.0.3/diff/v0.0.1/"

	importEndpoint           = "/import"
	importGetVersionEndpoint = "/a/_/my-app/v1/"
	exportProject            = "export-prj"
	exportEndpoint           = "/export"

	metricsEndpoint = "/metrics"
	healthEndpoint  = "/health"
)

func testRequest(c *C, req *http.Request, expectedStatus int) {
	handler.ServeHTTP(rr, req)
	c.Assert(rr.Code, DeepEquals, expectedStatus, Commentf("%s Responded with body '%s'", req.URL, rr.Body.String()))
}

func (s *suite) addRelease(c *C, project, version string) {
	rr = httptest.NewRecorder()
	body := bytes.NewReader([]byte(`{"name": "my-app", "version": "` + version + `"}`))
	req, _ := http.NewRequest("POST", "/a/"+project+"/register", body)
	testRequest(c, req, 200)
	rr = httptest.NewRecorder()
}

func (s *suite) Test_Register_fails_with_empty_body(c *C) {
	req, _ := http.NewRequest("POST", registerEndpoint, nil)
	testRequest(c, req, 400)
}

func (s *suite) Test_Register_fails_with_invalid_json(c *C) {
	body := bytes.NewReader([]byte("hello"))
	req, _ := http.NewRequest("POST", registerEndpoint, body)
	testRequest(c, req, 400)
}

func (s *suite) Test_Register_fails_if_required_fields_are_missing(c *C) {
	cases := []string{
		`{"name": "missing-version"}`,
		`{"version": "1"}`,
		`{}`,
		`[]`,
		`null`,
		`12`,
	}
	for _, testCase := range cases {
		body := bytes.NewReader([]byte(testCase))
		req, _ := http.NewRequest("POST", registerEndpoint, body)
		testRequest(c, req, 400)
		rr = httptest.NewRecorder()
	}
}

func (s *suite) Test_Register_success_with_minimal_metadata(c *C) {
	body := bytes.NewReader([]byte(`{"name": "my-app", "version": "1"}`))
	req, _ := http.NewRequest("POST", registerEndpoint, body)
	testRequest(c, req, 200)
}

func (s *suite) Test_GetApplications_empty_list(c *C) {
	req, _ := http.NewRequest("GET", applicationsEndpointNoProject, nil)
	testRequest(c, req, 404)
}

func (s *suite) Test_GetProjects(c *C) {
	s.addRelease(c, "project1", "1")
	s.addRelease(c, "project2", "2")
	req, _ := http.NewRequest("GET", getProjectsEndpoints, nil)
	testRequest(c, req, http.StatusOK)
	c.Assert(rr.Body.String(), Equals, `["project1","project2"]`)
}

func (s *suite) Test_GetApplications(c *C) {
	s.addRelease(c, applicationsTestProject, "1")
	s.addRelease(c, applicationsTestProject, "2")
	req, _ := http.NewRequest("GET", applicationsEndpoint, nil)
	testRequest(c, req, http.StatusOK)
	c.Assert(rr.Body.String(), Equals, `["my-app"]`)
}

type hasItemChecker struct{}

var HasItem = &hasItemChecker{}

func (*hasItemChecker) Info() *CheckerInfo {
	return &CheckerInfo{Name: "HasItem", Params: []string{"obtained", "expected to have item"}}
}
func (*hasItemChecker) Check(params []interface{}, names []string) (bool, string) {
	obtained := params[0]
	expectedItem := params[1]
	switch obtained.(type) {
	case []interface{}:
		for _, v := range obtained.([]interface{}) {
			if v == expectedItem {
				return true, ""
			}
		}
	case []string:
		for _, v := range obtained.([]string) {
			if v == expectedItem {
				return true, ""
			}
		}
	default:
		return false, "Unexpected type."
	}
	return false, "Item not found"
}

func (s *suite) Test_GetVersions(c *C) {
	s.addRelease(c, applicationVersionsTestProject, "1")
	s.addRelease(c, applicationVersionsTestProject, "2")
	req, _ := http.NewRequest("GET", applicationVersionsEndpoint, nil)
	testRequest(c, req, http.StatusOK)
	result := []string{}
	err := json.Unmarshal([]byte(rr.Body.String()), &result)
	c.Assert(err, IsNil)
	c.Assert(result, HasLen, 2)
	c.Assert(result, HasItem, "1")
	c.Assert(result, HasItem, "2")
}

func (s *suite) Test_GetVersions_fails_if_app_not_found(c *C) {
	s.addRelease(c, applicationVersionsTestProject, "1")
	s.addRelease(c, applicationVersionsTestProject, "2")
	req, _ := http.NewRequest("GET", applicationVersionsNoApp, nil)
	testRequest(c, req, 404)
	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", applicationVersionsNoProject, nil)
	testRequest(c, req, 404)
}

func (s *suite) Test_NextVersion(c *C) {
	req, _ := http.NewRequest("GET", nextVersionEndpoint, nil)
	testRequest(c, req, http.StatusOK)
	c.Assert(rr.Body.String(), Equals, `0`)

	s.addRelease(c, nextVersionProject, "0")
	req, _ = http.NewRequest("GET", nextVersionEndpoint, nil)
	testRequest(c, req, http.StatusOK)
	c.Assert(rr.Body.String(), Equals, `1`)

	s.addRelease(c, nextVersionProject, "10")
	req, _ = http.NewRequest("GET", nextVersionEndpoint, nil)
	testRequest(c, req, http.StatusOK)
	c.Assert(rr.Body.String(), Equals, `11`)
}

func (s *suite) Test_NextVersion_with_prefix(c *C) {
	req, _ := http.NewRequest("GET", nextVersionEndpoint+"?prefix=0.0.", nil)
	testRequest(c, req, http.StatusOK)
	c.Assert(rr.Body.String(), Equals, `0.0.0`)

	s.addRelease(c, nextVersionProject, "0.1.0")
	s.addRelease(c, nextVersionProject, "0.0.0")
	req, _ = http.NewRequest("GET", nextVersionEndpoint+"?prefix=0.0.", nil)
	testRequest(c, req, http.StatusOK)
	c.Assert(rr.Body.String(), Equals, `0.0.1`)
}

func (s *suite) Test_GetVersion(c *C) {
	s.addRelease(c, getVersionProject, "1")
	req, _ := http.NewRequest("GET", getVersionEndpoint, nil)
	testRequest(c, req, http.StatusOK)
	result := map[string]interface{}{}
	err := json.Unmarshal([]byte(rr.Body.String()), &result)
	c.Assert(err, IsNil)
	c.Assert(result["name"], Equals, "my-app")
	c.Assert(result["version"], Equals, "1")
}

func (s *suite) Test_GetVersion_Resolves_latest(c *C) {
	s.addRelease(c, getVersionProject, "0.0.1")
	s.addRelease(c, getVersionProject, "0.0.2")
	req, _ := http.NewRequest("GET", getLatestVersionEndpoint, nil)
	testRequest(c, req, http.StatusOK)
	result := map[string]interface{}{}
	err := json.Unmarshal([]byte(rr.Body.String()), &result)
	c.Assert(err, IsNil)
	c.Assert(result["name"], Equals, "my-app")
	c.Assert(result["version"], Equals, "0.0.2")
}

func (s *suite) Test_GetVersion_Resolves_auto_version(c *C) {
	s.addRelease(c, getVersionProject, "0.0.1")
	s.addRelease(c, getVersionProject, "0.0.2")
	req, _ := http.NewRequest("GET", getAutoVersionEndpoint, nil)
	testRequest(c, req, http.StatusOK)
	result := map[string]interface{}{}
	err := json.Unmarshal([]byte(rr.Body.String()), &result)
	c.Assert(err, IsNil)
	c.Assert(result["name"], Equals, "my-app")
	c.Assert(result["version"], Equals, "0.0.2")
}

func (s *suite) Test_GetPreviousVersion(c *C) {
	s.addRelease(c, getVersionProject, "0.0.1")
	s.addRelease(c, getVersionProject, "0.0.2")
	req, _ := http.NewRequest("GET", getPreviousEndpoint, nil)
	testRequest(c, req, http.StatusOK)
	result := map[string]interface{}{}
	err := json.Unmarshal([]byte(rr.Body.String()), &result)
	c.Assert(err, IsNil)
	c.Assert(result["name"], Equals, "my-app")
	c.Assert(result["version"], Equals, "0.0.1")

	req, _ = http.NewRequest("GET", getPreviousEndpoint2, nil)
	rr = httptest.NewRecorder()
	testRequest(c, req, http.StatusNotFound)

	req, _ = http.NewRequest("GET", getPreviousEndpoint3, nil)
	rr = httptest.NewRecorder()
	testRequest(c, req, http.StatusOK)
	err = json.Unmarshal([]byte(rr.Body.String()), &result)
	c.Assert(err, IsNil)
	c.Assert(result["name"], Equals, "my-app")
	c.Assert(result["version"], Equals, "0.0.1")
}

func (s *suite) Test_Diff(c *C) {
	s.addRelease(c, getVersionProject, "0.0.1")
	s.addRelease(c, getVersionProject, "0.0.2")
	s.addRelease(c, getVersionProject, "0.0.3")
	req, _ := http.NewRequest("GET", getDiffEndpoint, nil)
	testRequest(c, req, http.StatusOK)
	result := []string{}
	err := json.Unmarshal([]byte(rr.Body.String()), &result)
	c.Assert(err, IsNil)
	c.Assert(result, HasLen, 1)
	c.Assert(result[0], Equals, "Change Version from '0.0.1' to '0.0.2'")

	req, _ = http.NewRequest("GET", getDiffWithEndpoint, nil)
	rr = httptest.NewRecorder()
	testRequest(c, req, http.StatusOK)
	err = json.Unmarshal([]byte(rr.Body.String()), &result)
	c.Assert(err, IsNil)
	c.Assert(result, HasLen, 1)
	c.Assert(result[0], Equals, "Change Version from '0.0.1' to '0.0.3'")
}

func (s *suite) Test_GetVersion_fails_if_app_doesnt_exist(c *C) {
	req, _ := http.NewRequest("GET", getVersionEndpoint, nil)
	testRequest(c, req, 404)
}

func (s *suite) Test_GetVersion_fails_if_version_doesnt_exist(c *C) {
	versions := []string{
		"v1",
		"v1.0",
		"v100.100.100",
		"latest",
		"@",
		"v@",
		"v0.0.@",
	}
	for _, v := range versions {
		req, _ := http.NewRequest("GET", "/a/"+getVersionProject+"/my-app/"+v+"/", nil)
		testRequest(c, req, 404)
		rr = httptest.NewRecorder()
	}
}

func (s *suite) Test_GetVersion_fails_if_version_format_invalid(c *C) {
	versions := []string{
		"1",
		"v12asdpokasdk",
		"null",
		"1.0",
	}
	for _, v := range versions {
		req, _ := http.NewRequest("GET", "/a/"+getVersionProject+"/my-app/"+v+"/", nil)
		testRequest(c, req, 400)
		rr = httptest.NewRecorder()
	}
}

func (s *suite) Test_Import(c *C) {
	body := bytes.NewReader([]byte(`[
        {"name": "my-app", "version": "1", "project": "_"},
        {"name": "my-app", "version": "1", "project": "ankyra"},
        {"name": "my-app", "version": "2", "project": "ankyra"}
    ]
    `))
	req, _ := http.NewRequest("POST", importEndpoint, body)
	testRequest(c, req, 200)
	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", importGetVersionEndpoint, nil)
	testRequest(c, req, http.StatusOK)
}

func (s *suite) Test_Import_fails_with_malformed_json(c *C) {
	body := []string{
		"{}",
		"[{}]",
		"12",
		`"string"`,
	}
	for _, b := range body {
		req, _ := http.NewRequest("POST", importEndpoint, bytes.NewReader([]byte(b)))
		testRequest(c, req, 400)
		rr = httptest.NewRecorder()
	}
}
func (s *suite) Test_Import_fails_with_empty_body(c *C) {
	req, _ := http.NewRequest("POST", importEndpoint, nil)
	testRequest(c, req, 400)
}

func (s *suite) Test_Export(c *C) {
	s.addRelease(c, exportProject, "1")
	s.addRelease(c, exportProject, "2")
	req, _ := http.NewRequest("GET", exportEndpoint, nil)
	testRequest(c, req, http.StatusOK)
	result := []map[string]interface{}{}
	err := json.Unmarshal([]byte(rr.Body.String()), &result)
	c.Assert(err, IsNil)
	c.Assert(result, HasLen, 2)
	var first, second map[string]interface{}
	for _, release := range result {
		if release["version"].(string) == "1" {
			first = release
		}
		if release["version"].(string) == "2" {
			second = release
		}
	}
	c.Assert(first, Not(IsNil))
	c.Assert(second, Not(IsNil))
	c.Assert(first["name"], DeepEquals, "my-app")
	c.Assert(first["project"], DeepEquals, exportProject)
	c.Assert(first["URI"], DeepEquals, []interface{}{})
}

func (s *suite) Test_Metrics(c *C) {
	req, _ := http.NewRequest("GET", metricsEndpoint, nil)
	testRequest(c, req, 200)
}

func (s *suite) Test_HealthCheck(c *C) {
	req, _ := http.NewRequest("GET", healthEndpoint, nil)
	testRequest(c, req, 200)
}
