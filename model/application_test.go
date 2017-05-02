package model

import (
	. "gopkg.in/check.v1"
)

type appSuite struct{}

var _ = Suite(&appSuite{})

func (s *appSuite) Test_GetMaxFromVersions_MorePreciseIsGreater(c *C) {
	versions := []string{"0", "0.0", "0.0.0"}
	maxVer := getMaxFromVersions(versions, "")
	c.Assert(maxVer.ToString(), Equals, "0.0.0")
}

func (s *appSuite) Test_GetMaxFromVersions_MorePreciseIsGreater2(c *C) {
	versions := []string{"0", "0.1"}
	maxVer := getMaxFromVersions(versions, "")
	c.Assert(maxVer.ToString(), Equals, "0.1")
}

func (s *appSuite) Test_GetMaxFromVersions_SmallerAndHigherBeatsLonger(c *C) {
	versions := []string{"0.0.1", "0.0.2", "0.0.3", "0.1"}
	maxVer := getMaxFromVersions(versions, "")
	c.Assert(maxVer.ToString(), Equals, "0.1")
}

func (s *appSuite) Test_GetMaxFromVersions_Prefix_Matching(c *C) {
	versions := []string{"0.0.1", "0.0.2", "0.0.3", "0.1"}
	maxVer := getMaxFromVersions(versions, "0.0.")
	c.Assert(maxVer.ToString(), Equals, "3")
}

func (s *appSuite) Test_GetNextVersion(c *C) {
	semver, err := GetNextVersion("archive-semver-test-latest", "")
	c.Assert(err, IsNil)
	c.Assert(semver, Equals, "0")

	err = AddRelease(`{"name": "semver-test", "type": "archive", "version": "0"}`)
	c.Assert(err, IsNil)
	err = AddRelease(`{"name": "semver-test", "type": "archive", "version": "0.1"}`)
	c.Assert(err, IsNil)
	semver, err = GetNextVersion("archive-semver-test-latest", "")
	c.Assert(err, IsNil)
	c.Assert(semver, Equals, "0.2")
}

func (s *appSuite) Test_GetNextVersion_With_Prefix(c *C) {
	semver, err := GetNextVersion("archive-semver2-test-latest", "")
	c.Assert(err, IsNil)
	c.Assert(semver, Equals, "0")

	err = AddRelease(`{"name": "semver2-test", "type": "archive", "version": "1"}`)
	c.Assert(err, IsNil)
	err = AddRelease(`{"name": "semver2-test", "type": "archive", "version": "0.1"}`)
	c.Assert(err, IsNil)
	semver, err = GetNextVersion("archive-semver2-test-latest", "0.")
	c.Assert(err, IsNil)
	c.Assert(semver, Equals, "0.2")
}
