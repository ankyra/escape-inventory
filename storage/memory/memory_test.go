/*
Copyright 2017, 2018 Ankyra

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

package memory

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/ankyra/escape-core/parsers"
	"github.com/ankyra/escape-inventory/config"
	"github.com/ankyra/escape-inventory/dao/types"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type suite struct{}

var _ = Suite(&suite{})

func (s *suite) Test_InMemoryStorageBackend(c *C) {
	uri := "mem://namespace/name-v1.0.tgz"
	unit := NewInMemoryStorageBackend()
	unit.Init(config.StorageSettings{})

	_, err := unit.Download("namespace", uri)
	c.Assert(err, Equals, types.NotFound)

	pkg := bytes.NewReader([]byte("package data"))
	releaseId, err := parsers.ParseReleaseId("namespace/name-v1.0")
	c.Assert(err, IsNil)
	str, err := unit.Upload("namespace", releaseId, pkg)
	c.Assert(err, IsNil)
	c.Assert(str, Equals, uri)

	data, err := unit.Download("namespace", uri)
	c.Assert(err, IsNil)
	bytes, err := ioutil.ReadAll(data)
	c.Assert(err, IsNil)
	c.Assert(string(bytes), Equals, "package data")
}
