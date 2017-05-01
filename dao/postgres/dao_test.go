package postgres

import (
    "os"
	"testing"
    "database/sql"
    _ "github.com/lib/pq"
	. "gopkg.in/check.v1"
    "github.com/ankyra/escape-registry/dao/types"
)

func Test(t *testing.T) { TestingT(t) }

type memSuite struct{}

var _ = Suite(&memSuite{})

func (s *memSuite) Test_DAO(c *C) {
    enableTest := false
    for _, e := range os.Environ() {
        if e == "ENABLE_POSTGRES_TESTS=1" {
            enableTest = true
        }
    }
    if enableTest {
        types.ValidateDAO(func() types.DAO {

            // Empty database
            url := "postgres://postgres:@localhost/postgres?sslmode=disable"
            db, err := sql.Open("postgres", url)
            c.Assert(err, IsNil)
            _, err = db.Exec(`DROP TABLE IF EXISTS release`)
            c.Assert(err, IsNil)
            _, err = db.Exec(`DROP TABLE IF EXISTS package`)
            c.Assert(err, IsNil)

            // Create unit
            dao, err := NewPostgresDAO(url)
            c.Assert(err, IsNil)
            return dao
        }, c)
    }
}
