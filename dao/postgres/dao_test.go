/*
Copyright 2017 Ankyra

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

package postgres

import (
	"database/sql"
	"github.com/ankyra/escape-registry/dao/types"
	_ "github.com/lib/pq"
	. "gopkg.in/check.v1"
	"os"
	"testing"
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
			_, err = db.Exec(`DROP TABLE IF EXISTS acl`)
			c.Assert(err, IsNil)

			// Create unit
			dao, err := NewPostgresDAO(url)
			c.Assert(err, IsNil)
			return dao
		}, c)
	} else {
		println("Postgres tests have not been enabled. Use ENABLE_POSTGRES_TESTS=1 to do so")
	}
}
