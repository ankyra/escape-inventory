package model

import (
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type semverSuite struct{}

var _ = Suite(&semverSuite{})

func (s *semverSuite) Test_LessOrEqual(c *C) {
	unit := NewSemanticVersion("0.0.3")
	c.Assert(unit.LessOrEqual(NewSemanticVersion("0.0.4.0")), Equals, true)
	c.Assert(unit.LessOrEqual(NewSemanticVersion("0.0.3.1")), Equals, true)
	c.Assert(unit.LessOrEqual(NewSemanticVersion("0.0.3.0")), Equals, true)
	c.Assert(unit.LessOrEqual(NewSemanticVersion("0.0.1")), Equals, false)
	c.Assert(unit.LessOrEqual(NewSemanticVersion("0.0.2")), Equals, false)
	c.Assert(unit.LessOrEqual(NewSemanticVersion("0.0.3")), Equals, true)
	c.Assert(unit.LessOrEqual(NewSemanticVersion("0.0.4")), Equals, true)
	c.Assert(unit.LessOrEqual(NewSemanticVersion("0.0.5")), Equals, true)
	c.Assert(unit.LessOrEqual(NewSemanticVersion("0.0")), Equals, false)
	c.Assert(unit.LessOrEqual(NewSemanticVersion("0.1")), Equals, true)
	c.Assert(unit.LessOrEqual(NewSemanticVersion("0")), Equals, false)
	c.Assert(unit.LessOrEqual(NewSemanticVersion("1")), Equals, true)
	c.Assert(unit.LessOrEqual(NewSemanticVersion("2")), Equals, true)
}

func (s *semverSuite) Test_IncrementSmallest(c *C) {
	unit := NewSemanticVersion("0.0.3")
	unit.IncrementSmallest()
	c.Assert(unit.ToString(), Equals, "0.0.4")
	unit = NewSemanticVersion("0.1")
	unit.IncrementSmallest()
	c.Assert(unit.ToString(), Equals, "0.2")
	unit = NewSemanticVersion("0")
	unit.IncrementSmallest()
	c.Assert(unit.ToString(), Equals, "1")
}
