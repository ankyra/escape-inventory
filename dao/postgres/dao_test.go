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
	"io/ioutil"
	"os"
	"testing"

	"github.com/ankyra/escape-inventory/dao/types"
	_ "github.com/lib/pq"
	. "gopkg.in/check.v1"
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
			_, err = db.Exec(`TRUNCATE release CASCADE`)
			_, err = db.Exec(`TRUNCATE package CASCADE`)
			_, err = db.Exec(`TRUNCATE acl CASCADE`)
			_, err = db.Exec(`TRUNCATE application CASCADE`)
			_, err = db.Exec(`TRUNCATE project CASCADE`)
			_, err = db.Exec(`TRUNCATE release_dependency CASCADE`)
			_, err = db.Exec(`TRUNCATE subscriptions CASCADE`)

			// Create unit
			dao, err := NewPostgresDAO(url)
			c.Assert(err, IsNil)
			return dao
		}, c)
	} else {
		println("Postgres tests have not been enabled. Use ENABLE_POSTGRES_TESTS=1 to do so")
	}
}

func (s *memSuite) Test_DAO_migrate(c *C) {
	enableTest := false
	for _, e := range os.Environ() {
		if e == "ENABLE_POSTGRES_TESTS=1" {
			enableTest = true
		}
	}
	if !enableTest {
		return
	}
	testfiles := []string{"testdata/migration_test_dump.sql"}
	for _, file := range testfiles {
		// Empty database
		url := "postgres://postgres:@localhost/postgres?sslmode=disable"
		db, err := sql.Open("postgres", url)
		c.Assert(err, IsNil)
		_, err = db.Exec(`DROP TABLE release CASCADE`)
		_, err = db.Exec(`DROP TABLE package CASCADE`)
		_, err = db.Exec(`DROP TABLE acl CASCADE`)
		_, err = db.Exec(`DROP TABLE application CASCADE`)
		_, err = db.Exec(`DROP TABLE project CASCADE`)
		_, err = db.Exec(`DROP TABLE release_dependency CASCADE`)
		_, err = db.Exec(`DROP TABLE subscriptions CASCADE`)
		_, err = db.Exec(`DROP TABLE schema_migrations CASCADE`)

		bytes, err := ioutil.ReadFile(file)
		c.Assert(err, IsNil)
		_, err = db.Exec(string(bytes))
		c.Assert(err, IsNil)
		// Create unit
		dao, err := NewPostgresDAO(url)
		c.Assert(err, IsNil)

		release, err := dao.GetRelease("project", "test", "project/test-v1.0")
		c.Assert(err, IsNil)
		c.Assert(release.Metadata, Not(IsNil))
		c.Assert(release.Metadata.Description, Equals, "yo")
	}
}
