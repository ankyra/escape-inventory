package sqlite

import (
	"github.com/ankyra/escape-registry/dao/types"
	. "gopkg.in/check.v1"
	"os"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type sqliteSuite struct{}

var _ = Suite(&sqliteSuite{})

func (s *sqliteSuite) Test_DAO(c *C) {
	types.ValidateDAO(func() types.DAO {
		os.RemoveAll("./test.db")
		dao, err := NewSQLiteDAO("./test.db")
		c.Assert(err, IsNil)
		return dao
	}, c)
	os.RemoveAll("./test.db")
}
