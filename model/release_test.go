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

package model

import (
	"github.com/ankyra/escape-registry/dao"
	. "gopkg.in/check.v1"
)

type releaseSuite struct{}

var _ = Suite(&releaseSuite{})

func (s *releaseSuite) Test_AddRelease_Missing_Field_Name(c *C) {
	err := AddRelease(`{"version": "0"}`)
	c.Assert(err, Not(IsNil))
}
func (s *releaseSuite) Test_AddRelease_Missing_Field_Version(c *C) {
	err := AddRelease(`{"name": "asdaiasd"}`)
	c.Assert(err, Not(IsNil))
}

func (s *releaseSuite) Test_AddRelease_GetMetadata(c *C) {
	err := AddRelease(`{"name": "rel-test", "version": "0"}`)
	c.Assert(err, IsNil)
	metadata, err := GetReleaseMetadata("rel-test-v0")
	c.Assert(err, IsNil)
	c.Assert(metadata.GetReleaseId(), Equals, "rel-test-v0")
}

func (s *releaseSuite) Test_GetMetadataNotFound(c *C) {
	_, err := GetReleaseMetadata("whatiojasdiofjasd-test-v0")
	c.Assert(dao.IsNotFound(err), Equals, true)
}
