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
	. "gopkg.in/check.v1"
)

type importSuite struct{}

var _ = Suite(&importSuite{})

func (s *importSuite) Test_Import_Empty(c *C) {
	releases := []map[string]interface{}{}
	err := Import(releases)
	c.Assert(err, IsNil)
}

func (s *importSuite) Test_Import(c *C) {
	releases := []map[string]interface{}{
		map[string]interface{}{
			"type":    "archive",
			"name":    "import-test",
			"version": "1",
		},
	}
	err := Import(releases)
	c.Assert(err, IsNil)
	metadata, err := GetReleaseMetadata("archive-import-test-v1")
	c.Assert(err, IsNil)
	c.Assert(metadata.GetName(), Equals, "import-test")
}

func (s *importSuite) Test_Import_Ignore_Existing(c *C) {
	releases := []map[string]interface{}{
		map[string]interface{}{
			"type":    "archive",
			"name":    "import-exists-test",
			"version": "1",
		},
		map[string]interface{}{
			"type":    "archive",
			"name":    "import-exists-test",
			"version": "1",
		},
	}
	err := Import(releases)
	c.Assert(err, IsNil)
	metadata, err := GetReleaseMetadata("archive-import-exists-test-v1")
	c.Assert(err, IsNil)
	c.Assert(metadata.GetName(), Equals, "import-exists-test")
}
