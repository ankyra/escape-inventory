package model

import (
	. "gopkg.in/check.v1"
    "github.com/ankyra/escape-registry/dao"
)

type releaseSuite struct{}

var _ = Suite(&releaseSuite{})

func (s *releaseSuite) Test_AddRelease_Missing_Field_Name(c *C) {
    err := AddRelease(`{"type": "archive", "version": "0"}`)
    c.Assert(err, Not(IsNil))
}
func (s *releaseSuite) Test_AddRelease_Missing_Field_Type(c *C) {
    err := AddRelease(`{"name": "aiosdjioasjdoas", "version": "0"}`)
    c.Assert(err, Not(IsNil))
}
func (s *releaseSuite) Test_AddRelease_Missing_Field_Version(c *C) {
    err := AddRelease(`{"type": "archive", "name": "asdaiasd"}`)
    c.Assert(err, Not(IsNil))
}

func (s *releaseSuite) Test_AddRelease_GetMetadata(c *C) {
    err := AddRelease(`{"name": "rel-test", "type": "archive", "version": "0"}`)
    c.Assert(err, IsNil)
    metadata, err := GetReleaseMetadata("archive-rel-test-v0")
    c.Assert(err, IsNil)
    c.Assert(metadata.GetReleaseId(), Equals, "archive-rel-test-v0")
}

func (s *releaseSuite) Test_GetMetadataNotFound(c *C) {
    _, err := GetReleaseMetadata("archive-whatiojasdiofjasd-test-v0")
    c.Assert(dao.IsNotFound(err), Equals, true)
}
