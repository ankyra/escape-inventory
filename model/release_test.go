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

package model

import (
	"fmt"

	"github.com/ankyra/escape-inventory/dao"
	"github.com/ankyra/escape-inventory/dao/types"
	. "gopkg.in/check.v1"
)

type releaseSuite struct{}

var _ = Suite(&releaseSuite{})

func (s *releaseSuite) SetUpTest(c *C) {
	dao.TestSetup()
}

func (s *releaseSuite) Test_AddRelease_Missing_Field_Name(c *C) {
	_, err := AddRelease("_", `{"version": "0"}`)
	c.Assert(err, Not(IsNil))
}

func (s *releaseSuite) Test_AddRelease_Missing_Field_Version(c *C) {
	_, err := AddRelease("_", `{"name": "asdaiasd"}`)
	c.Assert(err, Not(IsNil))
}

func (s *releaseSuite) Test_AddRelease_GetMetadata(c *C) {
	_, err := AddRelease("_", `{"name": "rel-test", "version": "0"}`)
	c.Assert(err, IsNil)
	metadata, err := GetReleaseMetadata("_", "rel-test", "v0")
	c.Assert(err, IsNil)
	c.Assert(metadata.GetReleaseId(), Equals, "rel-test-v0")
}

func (s *releaseSuite) Test_GetMetadataNotFound(c *C) {
	_, err := GetReleaseMetadata("_", "whatiojasdiofjasd-test", "v0")
	c.Assert(dao.IsNotFound(err), Equals, true)
}

func (s *releaseSuite) Test_AddRelease_Creates_Namespace_Metadata(c *C) {
	_, err := dao.GetNamespace("test")
	c.Assert(err, Equals, types.NotFound)

	_, err = AddRelease("test", `{"name": "rel-test", "version": "0"}`)
	c.Assert(err, IsNil)

	prj, err := dao.GetNamespace("test")
	c.Assert(err, IsNil)
	c.Assert(prj.Name, Equals, "test")
}

func (s *releaseSuite) Test_AddRelease_Creates_Project_Metadata_fails_if_invalid_namespace_name(c *C) {
	_, err := dao.GetNamespace("invalid$")
	c.Assert(err, Equals, types.NotFound)

	_, err = AddRelease("invalid$", `{"name": "rel-test", "version": "0"}`)
	c.Assert(err, DeepEquals, NewUserError(fmt.Errorf("Invalid name 'invalid$'")))
}

func (s *releaseSuite) Test_AddRelease_Creates_Application_Metadata(c *C) {
	_, err := dao.GetApplication("test", "up-test")
	c.Assert(err, Equals, types.NotFound)

	_, err = AddRelease("test", `{"name": "up-test", "version": "0", "description": "testing", "project": "test"}`)
	c.Assert(err, IsNil)

	app, err := dao.GetApplication("test", "up-test")
	c.Assert(err, IsNil)
	c.Assert(app.Name, Equals, "up-test")
	c.Assert(app.Project, Equals, "test")
	c.Assert(app.Description, Equals, "testing")
	c.Assert(app.LatestVersion, Equals, "0")
}

func (s *releaseSuite) Test_AddRelease_Creates_Application_Metadata_fails_if_invalid_app_name(c *C) {
	_, err := dao.GetApplication("test", "up/test")
	c.Assert(err, Equals, types.NotFound)

	_, err = AddRelease("test", `{"name": "up/test", "version": "0", "description": "testing", "project": "test"}`)
	c.Assert(err, DeepEquals, NewUserError(fmt.Errorf("Invalid name 'up/test'")))
}

func (s *releaseSuite) Test_AddRelease_Updates_Application_Metadata(c *C) {
	_, err := dao.GetApplication("test", "up-test")
	c.Assert(err, Equals, types.NotFound)

	_, err = AddRelease("test", `{"name": "up-test", "version": "0", "description": "testing", "project": "test"}`)
	c.Assert(err, IsNil)

	app, err := dao.GetApplication("test", "up-test")
	c.Assert(err, IsNil)
	c.Assert(app.Name, Equals, "up-test")
	c.Assert(app.Project, Equals, "test")
	c.Assert(app.Description, Equals, "testing")
	c.Assert(app.LatestVersion, Equals, "0")

	_, err = AddRelease("test", `{"name": "up-test", "version": "1", "description": "updated", "project": "test"}`)
	c.Assert(err, IsNil)

	app, err = dao.GetApplication("test", "up-test")
	c.Assert(err, IsNil)
	c.Assert(app.Name, Equals, "up-test")
	c.Assert(app.Project, Equals, "test")
	c.Assert(app.Description, Equals, "updated")
	c.Assert(app.LatestVersion, Equals, "1")
}

func (s *releaseSuite) Test_AddRelease_Processes_Dependencies(c *C) {
	_, err := AddRelease("test", `{"name": "up-test", "version": "1", "project": "test"}`)
	c.Assert(err, IsNil)
	release, err := ResolveReleaseId("test", "up-test", "v1")
	c.Assert(err, IsNil)
	c.Assert(release.ProcessedDependencies, Equals, true)
}

func (s *releaseSuite) Test_AddRelease_Processes_Dependencies_2(c *C) {
	_, err := AddRelease("test", `{"name": "up-test", "version": "1", "project": "test", 
								   "depends": [{
									   "release_id": "yo/test-v1"
								   }, {
									   "release_id": "yoyo/tester-v2",
									   "scopes": ["build"]
								   }],
									 "extends": [{
										 "release_id": "aight/ai-v3"
									 }]
								 }`)
	c.Assert(err, IsNil)
	release, err := ResolveReleaseId("test", "up-test", "v1")
	c.Assert(err, IsNil)
	c.Assert(release.ProcessedDependencies, Equals, true)
	deps, err := dao.GetDependencies(release)
	c.Assert(err, IsNil)
	c.Assert(deps, HasLen, 3)
	c.Assert(deps[0].Project, Equals, "yo")
	c.Assert(deps[0].Application, Equals, "test")
	c.Assert(deps[0].Version, Equals, "1")
	c.Assert(deps[0].BuildScope, Equals, true)
	c.Assert(deps[0].DeployScope, Equals, true)
	c.Assert(deps[1].Project, Equals, "yoyo")
	c.Assert(deps[1].Application, Equals, "tester")
	c.Assert(deps[1].Version, Equals, "2")
	c.Assert(deps[1].BuildScope, Equals, true)
	c.Assert(deps[1].DeployScope, Equals, false)
	c.Assert(deps[2].Project, Equals, "aight")
	c.Assert(deps[2].Application, Equals, "ai")
	c.Assert(deps[2].Version, Equals, "3")
	c.Assert(deps[2].BuildScope, Equals, true)
	c.Assert(deps[2].DeployScope, Equals, true)
}

func (s *releaseSuite) Test_AddRelease_Registers_Providers(c *C) {
	_, err := AddRelease("test", `{"name": "up-test", "version": "1", "project": "test", 
								   "provides": [{
										 "name": "provider-test"
									 }]
								 }`)
	c.Assert(err, IsNil)
	providers, err := dao.GetProviders("provider-test")
	c.Assert(err, IsNil)
	c.Assert(providers, HasLen, 1)
	c.Assert(providers["test/up-test-v1"], Not(IsNil))
	c.Assert(providers["test/up-test-v1"].Version, Equals, "1")
}
