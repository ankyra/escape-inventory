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

package sqlite

import (
	"os"
	"testing"

	"github.com/ankyra/escape-inventory/dao/types"
	. "gopkg.in/check.v1"
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
