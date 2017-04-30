package mem

import (
	"testing"
	. "gopkg.in/check.v1"
    "github.com/ankyra/escape-registry/dao/types"
)

func Test(t *testing.T) { TestingT(t) }

type memSuite struct{}

var _ = Suite(&memSuite{})

func (s *memSuite) Test_DAO(c *C) {
    types.ValidateDAO(NewInMemoryDAO, c)
}
