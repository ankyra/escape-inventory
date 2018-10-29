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

package model

import (
	"fmt"

	"github.com/ankyra/escape-inventory/dao"
	"github.com/ankyra/escape-inventory/dao/types"
	. "gopkg.in/check.v1"
)

type suite struct{}

var _ = Suite(&suite{})

func (s *suite) SetUpTest(c *C) {
	dao.TestSetup()
}

func (s *suite) Test_AddNamespace_fails_if_no_namespacet_name(c *C) {
	p := types.NewProject("")
	c.Assert(AddNamespace(p, "username"), DeepEquals, NewUserError(fmt.Errorf("Missing name")))
}

func (s *suite) Test_AddNamespace_fails_with_invalid_namespace_name(c *C) {
	cases := []string{
		"a$",
		"a'",
		"a*",
		"a^",
		"a@",
		"123",
		"a\"",
		"test/",
	}
	for _, name := range cases {
		p := types.NewProject(name)
		c.Assert(AddNamespace(p, "username"), DeepEquals, NewUserError(fmt.Errorf("Invalid name '%s'", name)))
	}
}
