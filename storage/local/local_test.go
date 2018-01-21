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

package local

import (
	"github.com/ankyra/escape-core/parsers"
	. "gopkg.in/check.v1"
	"io/ioutil"
	"os"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type localSuite struct{}

var _ = Suite(&localSuite{})

const test_local_storage_path = "/tmp/escape-test/"
const test_data_path = "/tmp/escape_test_data.txt"

var test_data = []byte("Hello!")

func (s *localSuite) Test_Local_Storage_Backend_Upload(c *C) {
	backend := NewLocalStorageBackendWithStoragePath(test_local_storage_path)
	err := ioutil.WriteFile(test_data_path, test_data, 0644)
	c.Assert(err, IsNil)
	pkg, err := os.Open(test_data_path)
	c.Assert(err, IsNil)
	releaseId, err := parsers.ParseReleaseId("archive-upload-test-v1")
	c.Assert(err, IsNil)
	uri, err := backend.Upload("project", releaseId, pkg)
	c.Assert(err, IsNil)
	c.Assert(uri, Equals, "file:///tmp/escape-test/project/archive-upload-test/archive-upload-test-v1.tgz")
	os.RemoveAll(test_data_path)
	os.RemoveAll(test_local_storage_path)
}

func (s *localSuite) Test_Local_Storage_Backend_Download(c *C) {
	backend := NewLocalStorageBackendWithStoragePath(test_local_storage_path)
	err := ioutil.WriteFile(test_data_path, test_data, 0644)
	c.Assert(err, IsNil)
	pkg, err := os.Open(test_data_path)
	c.Assert(err, IsNil)
	releaseId, err := parsers.ParseReleaseId("archive-upload-test-v1")
	c.Assert(err, IsNil)
	uri, err := backend.Upload("project", releaseId, pkg)
	c.Assert(err, IsNil)
	seeker, err := backend.Download("project", uri)

	payload, err := ioutil.ReadAll(seeker)
	c.Assert(err, IsNil)
	c.Assert(payload, DeepEquals, test_data)

	os.RemoveAll(test_data_path)
	os.RemoveAll(test_local_storage_path)
}
