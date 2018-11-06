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

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ankyra/escape-inventory/cmd"
	"github.com/ankyra/escape-inventory/config"
	"github.com/ankyra/escape-inventory/dao"
	"github.com/ankyra/escape-inventory/dao/types"
	"github.com/ankyra/escape-inventory/model"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type suite struct{}

var _ = Suite(&suite{})

var handler http.Handler
var rr *httptest.ResponseRecorder

func (s *suite) SetUpTest(c *C) {
	dao.TestSetup()
	config, _ := config.NewConfig([]string{})
	handler = cmd.GetHandler(getMux(config))
	rr = httptest.NewRecorder()
}

const (
	registerEndpoint = "/api/v1/inventory/my-project/register"

	getProjectEndpoint    = "/api/v1/inventory/" + applicationsTestProject + "/"
	addProjectEndpoint    = "/api/v1/inventory/test/add-project"
	updateProjectEndpoint = "/api/v1/inventory/test/"
	getProjectsEndpoints  = "/api/v1/inventory/"

	getProjectUnitsEndpoint       = "/api/v1/inventory/" + applicationsTestProject + "/units/"
	getProjectUnitEndpoint        = "/api/v1/inventory/" + applicationsTestProject + "/units/my-app/"
	applicationsTestProject       = "applications-test-prj"
	applicationsEndpointNoProject = "/api/v1/inventory/doesnt-exist/"

	applicationVersionsTestProject = "versions-test-prj"
	applicationVersionsEndpoint    = "/api/v1/inventory/" + applicationVersionsTestProject + "/units/my-app/versions/"
	applicationVersionsNoProject   = "/api/v1/inventory/doesnt-exist/units/my-app/"
	applicationVersionsNoApp       = "/api/v1/inventory/versions-test/units/doesnt-exist/"

	nextVersionProject  = "next-version-prj"
	nextVersionEndpoint = "/api/v1/inventory/" + nextVersionProject + "/units/my-app/next-version"

	getVersionProject        = "get-version-prj"
	getVersionEndpoint       = "/api/v1/inventory/" + getVersionProject + "/units/my-app/versions/v1/"
	getLatestVersionEndpoint = "/api/v1/inventory/" + getVersionProject + "/units/my-app/versions/latest/"
	getAutoVersionEndpoint   = "/api/v1/inventory/" + getVersionProject + "/units/my-app/versions/v0.0.@/"
	getPreviousEndpoint      = "/api/v1/inventory/" + getVersionProject + "/units/my-app/versions/v0.0.2/previous/"
	getPreviousEndpoint2     = "/api/v1/inventory/" + getVersionProject + "/units/my-app/versions/v0.0.1/previous/"
	getPreviousEndpoint3     = "/api/v1/inventory/" + getVersionProject + "/units/my-app/versions/v0.0.@/previous/"
	getDiffEndpoint          = "/api/v1/inventory/" + getVersionProject + "/units/my-app/versions/v0.0.2/diff/"
	getDiffWithEndpoint      = "/api/v1/inventory/" + getVersionProject + "/units/my-app/versions/v0.0.3/diff/v0.0.1/"

	importEndpoint           = "/api/v1/internal/import"
	importGetVersionEndpoint = "/api/v1/inventory/_/units/my-app/versions/v1/"
	exportProject            = "export-prj"
	exportEndpoint           = "/api/v1/internal/export"

	metricsEndpoint = "/metrics"
	healthEndpoint  = "/health"
)

func testRequest(c *C, req *http.Request, expectedStatus int) {
	handler.ServeHTTP(rr, req)
	c.Assert(rr.Code, DeepEquals, expectedStatus, Commentf("%s Responded with body '%s'", req.URL, rr.Body.String()))
}

func (s *suite) addRelease(c *C, project, version string) {
	rr = httptest.NewRecorder()
	body := bytes.NewReader([]byte(`{"name": "my-app", "version": "` + version + `", "project": "` + project + `"}`))
	req, _ := http.NewRequest("POST", "/api/v1/inventory/"+project+"/register", body)
	testRequest(c, req, 200)
	rr = httptest.NewRecorder()
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

	result := map[string]map[string]string{}
	err := json.Unmarshal([]byte(rr.Body.String()), &result)
	c.Assert(err, IsNil)
	c.Assert(result["project1"]["name"], Equals, "project1")
	c.Assert(result["project2"]["name"], Equals, "project2")
}

func (s *suite) Test_GetProject(c *C) {
	s.addRelease(c, applicationsTestProject, "1")
	s.addRelease(c, applicationsTestProject, "2")
	req, _ := http.NewRequest("GET", getProjectEndpoint, nil)
	testRequest(c, req, http.StatusOK)

	result := model.NamespacePayload{}
	err := json.Unmarshal([]byte(rr.Body.String()), &result)
	c.Assert(err, IsNil)
	c.Assert(result.Name, Equals, applicationsTestProject)
	c.Assert(result.Units["my-app"].Name, Equals, "my-app")
}

func (s *suite) Test_GetProjectUnits(c *C) {
	s.addRelease(c, applicationsTestProject, "1")
	s.addRelease(c, applicationsTestProject, "2")
	req, _ := http.NewRequest("GET", getProjectUnitsEndpoint, nil)
	testRequest(c, req, http.StatusOK)

	result := map[string]*types.Application{}
	err := json.Unmarshal([]byte(rr.Body.String()), &result)
	c.Assert(err, IsNil)
	c.Assert(result, HasLen, 1)
	c.Assert(result["my-app"], Not(IsNil))
	c.Assert(result["my-app"].Name, Equals, "my-app")
	c.Assert(result["my-app"].Project, Equals, applicationsTestProject)
	c.Assert(result["my-app"].Description, Equals, "")
	c.Assert(result["my-app"].Logo, Equals, "")
	c.Assert(result["my-app"].LatestVersion, Equals, "2")
}

func (s *suite) Test_GetProjectUnit(c *C) {
	s.addRelease(c, applicationsTestProject, "1")
	s.addRelease(c, applicationsTestProject, "2")
	req, _ := http.NewRequest("GET", getProjectUnitEndpoint, nil)
	testRequest(c, req, http.StatusOK)

	result := model.ApplicationPayload{}
	err := json.Unmarshal([]byte(rr.Body.String()), &result)
	fmt.Println(rr.Body.String())
	c.Assert(err, IsNil)
	c.Assert(result, Not(IsNil))
	c.Assert(result.Name, Equals, "my-app")
	c.Assert(result.Project, Equals, applicationsTestProject)
	c.Assert(result.Description, Equals, "")
	c.Assert(result.Logo, Equals, "")
	c.Assert(result.LatestVersion, Equals, "2")
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
	result := map[string]map[string][]map[string]interface{}{}
	err := json.Unmarshal([]byte(rr.Body.String()), &result)
	c.Assert(err, IsNil)
	c.Assert(result, HasLen, 1)
	c.Assert(result["Version"]["change"][0]["PreviousValue"], Equals, "0.0.1")
	c.Assert(result["Version"]["change"][0]["NewValue"], Equals, "0.0.2")

	req, _ = http.NewRequest("GET", getDiffWithEndpoint, nil)
	rr = httptest.NewRecorder()
	testRequest(c, req, http.StatusOK)
	err = json.Unmarshal([]byte(rr.Body.String()), &result)
	c.Assert(err, IsNil)
	c.Assert(result, HasLen, 1)
	c.Assert(result["Version"]["change"][0]["PreviousValue"], Equals, "0.0.1")
	c.Assert(result["Version"]["change"][0]["NewValue"], Equals, "0.0.3")
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
		req, _ := http.NewRequest("GET", "/api/v1/inventory/"+getVersionProject+"/units/my-app/versions/"+v+"/", nil)
		testRequest(c, req, 404)
		rr = httptest.NewRecorder()
	}
}

func (s *suite) Test_GetVersion_fails_if_version_format_invalid(c *C) {
	versions := []string{
		"asd1@",
		"v12asdpokasdk",
		"null",
		"1.-",
	}
	for _, v := range versions {
		req, _ := http.NewRequest("GET", "/api/v1/inventory/"+getVersionProject+"/units/my-app/versions/"+v+"/", nil)
		testRequest(c, req, 400)
		rr = httptest.NewRecorder()
	}
}

func (s *suite) Test_Metrics(c *C) {
	req, _ := http.NewRequest("GET", metricsEndpoint, nil)
	testRequest(c, req, 200)
}

func (s *suite) Test_HealthCheck(c *C) {
	req, _ := http.NewRequest("GET", healthEndpoint, nil)
	testRequest(c, req, 200)
}

func (s *suite) Test_HardDeleteProject(c *C) {
	s.addRelease(c, "project1", "1")
	s.addRelease(c, "project2", "2")
	req, _ := http.NewRequest("GET", getProjectsEndpoints, nil)
	rr = httptest.NewRecorder()
	testRequest(c, req, http.StatusOK)
	result := map[string]map[string]string{}
	err := json.Unmarshal([]byte(rr.Body.String()), &result)
	c.Assert(err, IsNil)
	c.Assert(result, HasLen, 2)

	req, _ = http.NewRequest("DELETE", "/api/v1/inventory/project1/hard-delete", nil)
	rr = httptest.NewRecorder()
	testRequest(c, req, 200)

	req, _ = http.NewRequest("GET", getProjectsEndpoints, nil)
	rr = httptest.NewRecorder()
	testRequest(c, req, http.StatusOK)
	result = map[string]map[string]string{}
	err = json.Unmarshal([]byte(rr.Body.String()), &result)
	c.Assert(err, IsNil)
	c.Assert(result, HasLen, 1)
}
