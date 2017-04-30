package local

import (
    "os"
    "io/ioutil"
	"testing"
	. "gopkg.in/check.v1"
    "github.com/ankyra/escape-registry/shared"
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
    releaseId, err := shared.ParseReleaseId("archive-upload-test-v1")
    c.Assert(err, IsNil)
    uri, err := backend.Upload(releaseId, pkg)
    c.Assert(err, IsNil)
    c.Assert(uri, Equals, "file:///tmp/escape-test/archive/upload-test/archive-upload-test-v1.tgz")
    os.RemoveAll(test_data_path)
    os.RemoveAll(test_local_storage_path)
}

func (s *localSuite) Test_Local_Storage_Backend_Download(c *C) {
    backend := NewLocalStorageBackendWithStoragePath(test_local_storage_path)
    err := ioutil.WriteFile(test_data_path, test_data, 0644)
    c.Assert(err, IsNil)
    pkg, err := os.Open(test_data_path)
    c.Assert(err, IsNil)
    releaseId, err := shared.ParseReleaseId("archive-upload-test-v1")
    c.Assert(err, IsNil)
    uri, err := backend.Upload(releaseId, pkg)
    c.Assert(err, IsNil)
    seeker, err := backend.Download(uri)

    payload, err := ioutil.ReadAll(seeker)
    c.Assert(err, IsNil)
    c.Assert(payload, DeepEquals, test_data)

    os.RemoveAll(test_data_path)
    os.RemoveAll(test_local_storage_path)
}
