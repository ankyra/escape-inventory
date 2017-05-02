package mem

import (
	"github.com/ankyra/escape-registry/dao/types"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type memSuite struct{}

var _ = Suite(&memSuite{})

func (s *memSuite) Test_DAO(c *C) {
	types.ValidateDAO(NewInMemoryDAO, c)
}
