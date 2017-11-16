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

package ql

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/ankyra/escape-inventory/dao/types"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type qlSuite struct{}

var _ = Suite(&qlSuite{})

func (s *qlSuite) Test_DAO(c *C) {
	os.Mkdir("testdata", os.ModePerm)
	var dbName string
	types.ValidateDAO(func() types.DAO {
		dbName := fmt.Sprintf("./testdata/%s.db", randomString())
		os.RemoveAll(dbName)
		dao, err := NewQLDAO(dbName)
		c.Assert(err, IsNil)
		return dao
	}, c)
	os.RemoveAll(dbName)
	os.RemoveAll("testdata")
}

func randomString() string {
	rand.Seed(time.Now().UnixNano())

	var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, 6)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}
