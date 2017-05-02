package shared

import (
	"testing"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type releaseIdSuite struct{}

var _ = Suite(&releaseIdSuite{})

func (s *releaseIdSuite) Test_ReleaseId_Happy_Path(c *C) {
    id, err := ParseReleaseId("type-name-v1.0")
    c.Assert(err, IsNil)
	c.Assert(id.Type, Equals, "type")
	c.Assert(id.Name, Equals, "name")
	c.Assert(id.Version, Equals, "1.0")
}

func (s *releaseIdSuite) Test_ReleaseId_Can_Have_Dashes(c *C) {
    id, err := ParseReleaseId("type-name-with-dashes-v1.0")
    c.Assert(err, IsNil)
	c.Assert(id.Name, Equals, "name-with-dashes")
}

func (s *releaseIdSuite) Test_ReleaseId_Parse_Latest1(c *C) {
    id, err := ParseReleaseId("type-name-latest")
    c.Assert(err, IsNil)
	c.Assert(id.Version, Equals, "latest")
}

func (s *releaseIdSuite) Test_ReleaseId_Parse_Latest2(c *C) {
    id, err := ParseReleaseId("type-name-@")
    c.Assert(err, IsNil)
	c.Assert(id.Version, Equals, "latest")
}

func (s *releaseIdSuite) Test_ReleaseId_Parse_Latest3(c *C) {
    id, err := ParseReleaseId("type-name-v@")
    c.Assert(err, IsNil)
	c.Assert(id.Version, Equals, "latest")
}


func (s *releaseIdSuite) Test_ReleaseId_Parse_Version(c *C) {
    id, err := ParseReleaseId("type-name-v1.0")
    c.Assert(err, IsNil)
	c.Assert(id.Version, Equals, "1.0")
}

func (s *releaseIdSuite) Test_ReleaseId_Invalid_Format1(c *C) {
    _, err := ParseReleaseId("type")
    c.Assert(err.Error(), Equals, "Invalid release format: type")
}
func (s *releaseIdSuite) Test_ReleaseId_Invalid_Format2(c *C) {
    _, err := ParseReleaseId("type-name")
    c.Assert(err.Error(), Equals, "Invalid release format: type-name")
}
func (s *releaseIdSuite) Test_ReleaseId_Missing_Version(c *C) {
    _, err := ParseReleaseId("type-name-nope")
    c.Assert(err.Error(), Equals, "Invalid version string in release ID 'type-name-nope': nope")
}
func (s *releaseIdSuite) Test_ReleaseId_Invalid_Version(c *C) {
    _, err := ParseReleaseId("type-name-vnope")
    c.Assert(err.Error(), Equals, "Invalid release ID 'type-name-vnope': Invalid version format: nope")
}

func (s *releaseIdSuite) Test_ValidateVersion(c *C) {
    c.Assert(validateVersion("latest"), IsNil)
    c.Assert(validateVersion("0"), IsNil)
    c.Assert(validateVersion("10"), IsNil)
    c.Assert(validateVersion("0.0"), IsNil)
    c.Assert(validateVersion("0.10"), IsNil)
    c.Assert(validateVersion("0.0.0"), IsNil)
    c.Assert(validateVersion("0.0.10"), IsNil)
    c.Assert(validateVersion("0.@"), IsNil)
    c.Assert(validateVersion("0.0.@"), IsNil)
}

func (s *releaseIdSuite) Test_ValidateVersion_Error(c *C) {
    c.Assert(validateVersion("whatsthisnow"), Not(IsNil))
    c.Assert(validateVersion("nope"), Not(IsNil))
    c.Assert(validateVersion("0.test"), Not(IsNil))
    c.Assert(validateVersion("0.0.test"), Not(IsNil))
    c.Assert(validateVersion("0.0.latest"), Not(IsNil))
    c.Assert(validateVersion("0-0"), Not(IsNil))
    c.Assert(validateVersion("0_0"), Not(IsNil))
    c.Assert(validateVersion("0@"), Not(IsNil))
    c.Assert(validateVersion("0.0@"), Not(IsNil))
}
