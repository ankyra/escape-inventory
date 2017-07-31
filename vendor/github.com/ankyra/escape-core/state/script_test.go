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

package state

import (
	"github.com/ankyra/escape-core"
	"github.com/ankyra/escape-core/script"
	"github.com/ankyra/escape-core/variables"
	. "gopkg.in/check.v1"
)

type scriptSuite struct{}

var _ = Suite(&scriptSuite{})

func (s *scriptSuite) SetUpTest(c *C) {
	var err error
	prj, err := NewProjectStateFromFile("prj", "testdata/project_script.json", nil)
	c.Assert(err, IsNil)
	env := prj.GetEnvironmentStateOrMakeNew("dev")
	depl = env.GetOrCreateDeploymentState("archive-release")
	fullDepl = env.GetOrCreateDeploymentState("archive-full")
	dep := env.GetOrCreateDeploymentState("archive-release-with-deps")
	deplWithDeps = dep.GetDeploymentOrMakeNew("deploy", "archive-release")
}

func (s *scriptSuite) Test_ToScript(c *C) {
	metadata := core.NewReleaseMetadata("test", "1.0")
	metadata.Metadata["value"] = "yo"
	input, err := variables.NewVariableFromString("user_level", "string")
	c.Assert(err, IsNil)
	metadata.AddInputVariable(input)
	metadata.AddOutputVariable(input)
	unit := NewStateCompiler(nil).CompileState(depl, metadata, "deploy", true)
	dicts := map[string][]string{
		"inputs":   []string{"user_level"},
		"outputs":  []string{"user_level"},
		"metadata": []string{"value"},
	}
	test_helper_check_script_environment(c, unit, dicts, "archive-release")
}

func (s *scriptSuite) Test_ToScript_doesnt_include_variable_that_are_not_defined_in_release_metadata(c *C) {
	metadata := core.NewReleaseMetadata("test", "1.0")
	unit := NewStateCompiler(nil).CompileState(depl, metadata, "deploy", true)
	dicts := map[string][]string{
		"inputs":   []string{},
		"outputs":  []string{},
		"metadata": []string{},
	}
	test_helper_check_script_environment(c, unit, dicts, "archive-release")
}

func (s *scriptSuite) Test_ToScriptEnvironment_adds_dependencies(c *C) {
	resolver := newResolverFromMap(map[string]*core.ReleaseMetadata{
		"test-v1.0": core.NewReleaseMetadata("test", "1.0"),
	})
	metadata := core.NewReleaseMetadata("test", "1.0")
	metadata.SetDependencies([]string{"test-v1.0"})
	env, err := ToScriptEnvironment(fullDepl, metadata, "build", resolver)
	c.Assert(err, IsNil)
	c.Assert(script.IsDictAtom((*env)["$"]), Equals, true)
	dict := script.ExpectDictAtom((*env)["$"])
	dicts := map[string][]string{
		"inputs":   []string{},
		"outputs":  []string{},
		"metadata": []string{},
	}
	test_helper_check_script_environment(c, dict["this"], dicts, "archive-full")
	test_helper_check_script_environment(c, dict["test-v1.0"], dicts, "_/test")
}

func (s *scriptSuite) Test_ToScriptEnvironment_honours_variable_context(c *C) {
	resolver := newResolverFromMap(map[string]*core.ReleaseMetadata{
		"test-v1.0": core.NewReleaseMetadata("test", "1.0"),
	})
	metadata := core.NewReleaseMetadata("test", "1.0")
	metadata.SetDependencies([]string{"test-v1.0"})
	metadata.SetVariableInContext("test", "test-v1.0")
	env, err := ToScriptEnvironment(fullDepl, metadata, "build", resolver)
	c.Assert(err, IsNil)
	c.Assert(script.IsDictAtom((*env)["$"]), Equals, true)
	dict := script.ExpectDictAtom((*env)["$"])
	dicts := map[string][]string{
		"inputs":   []string{},
		"outputs":  []string{},
		"metadata": []string{},
	}
	test_helper_check_script_environment(c, dict["this"], dicts, "archive-full")
	test_helper_check_script_environment(c, dict["test-v1.0"], dicts, "_/test")
	test_helper_check_script_environment(c, dict["test"], dicts, "_/test")
}

func (s *scriptSuite) Test_ToScriptEnvironment_ignores_missing_variables_in_variable_context(c *C) {
	resolver := newResolverFromMap(map[string]*core.ReleaseMetadata{
		"test-v1.0": core.NewReleaseMetadata("test", "1.0"),
	})
	metadata := core.NewReleaseMetadata("test", "1.0")
	metadata.SetDependencies([]string{"test-v1.0"})
	metadata.SetVariableInContext("test", "doesnt-exist-1.0")
	env, err := ToScriptEnvironment(fullDepl, metadata, "build", resolver)
	c.Assert(err, IsNil)
	c.Assert(script.IsDictAtom((*env)["$"]), Equals, true)
	dict := script.ExpectDictAtom((*env)["$"])
	c.Assert(dict["test"], IsNil)
}

func (s *scriptSuite) Test_ToScriptEnvironment_doesnt_add_dependencies_that_are_not_in_metadata(c *C) {
	resolver := newResolverFromMap(map[string]*core.ReleaseMetadata{})
	metadata := core.NewReleaseMetadata("test", "1.0")
	env, err := ToScriptEnvironment(fullDepl, metadata, "build", resolver)
	c.Assert(err, IsNil)
	c.Assert(script.IsDictAtom((*env)["$"]), Equals, true)
	dict := script.ExpectDictAtom((*env)["$"])
	dicts := map[string][]string{
		"inputs":   []string{},
		"outputs":  []string{},
		"metadata": []string{},
	}
	test_helper_check_script_environment(c, dict["this"], dicts, "archive-full")
	c.Assert(dict["test-v1.0"], IsNil)
}

func (s *scriptSuite) Test_ToScriptEnvironment_fails_if_deployment_state_is_missing(c *C) {
	metadata := core.NewReleaseMetadata("test", "1.0")
	_, err := ToScriptEnvironment(nil, metadata, "build", nil)
	c.Assert(err, Not(IsNil))
}

func (s *scriptSuite) Test_ToScriptEnvironment_fails_if_dependency_metadata_is_missing(c *C) {
	resolver := newResolverFromMap(map[string]*core.ReleaseMetadata{})
	metadata := core.NewReleaseMetadata("test", "1.0")
	metadata.SetDependencies([]string{"archive-dep-v1.0"})
	_, err := ToScriptEnvironment(fullDepl, metadata, "build", resolver)
	c.Assert(err, Not(IsNil))
}

func (s *scriptSuite) Test_ToScriptEnvironment_adds_consumers(c *C) {
	resolver := newResolverFromMap(map[string]*core.ReleaseMetadata{
		"archive-full-v1.0": core.NewReleaseMetadata("test", "1.0"),
	})
	metadata := core.NewReleaseMetadata("test", "1.0")
	metadata.SetConsumes([]string{"test"})
	depl.SetProvider("build", "test", "archive-full")
	env, err := ToScriptEnvironment(depl, metadata, "build", resolver)
	c.Assert(err, IsNil)
	c.Assert(script.IsDictAtom((*env)["$"]), Equals, true)
	dict := script.ExpectDictAtom((*env)["$"])
	dicts := map[string][]string{
		"inputs":   []string{},
		"outputs":  []string{},
		"metadata": []string{},
	}
	test_helper_check_script_environment(c, dict["this"], dicts, "archive-release")
	test_helper_check_script_environment(c, dict["test"], dicts, "archive-full")
}

func (s *scriptSuite) Test_ToScriptEnvironment_fails_if_missing_provider_state(c *C) {
	resolver := newResolverFromMap(map[string]*core.ReleaseMetadata{})
	metadata := core.NewReleaseMetadata("test", "1.0")
	metadata.SetConsumes([]string{"test"})
	_, err := ToScriptEnvironment(depl, metadata, "build", resolver)
	c.Assert(err, Not(IsNil))
}

func (s *scriptSuite) Test_ToScriptEnvironment_fails_if_missing_provider_metadata(c *C) {
	resolver := newResolverFromMap(map[string]*core.ReleaseMetadata{})
	metadata := core.NewReleaseMetadata("test", "1.0")
	metadata.SetConsumes([]string{"test"})
	depl.SetProvider("build", "test", "archive-full")
	_, err := ToScriptEnvironment(depl, metadata, "build", resolver)
	c.Assert(err, Not(IsNil))
}

func (s *scriptSuite) Test_ToScriptEnvironment_fails_if_missing_provider_state_in_environment(c *C) {
	resolver := newResolverFromMap(map[string]*core.ReleaseMetadata{
		"archive-full-v1.0": core.NewReleaseMetadata("test", "1.0"),
	})
	metadata := core.NewReleaseMetadata("test", "1.0")
	metadata.SetConsumes([]string{"test"})
	depl.SetProvider("build", "test", "this-doesnt-exist")
	_, err := ToScriptEnvironment(depl, metadata, "build", resolver)
	c.Assert(err, Not(IsNil))
}

func test_helper_check_script_environment(c *C, unit script.Script, dicts map[string][]string, name string) {
	c.Assert(script.IsDictAtom(unit), Equals, true)
	dict := script.ExpectDictAtom(unit)
	strings := map[string]string{
		"version":     "1.0",
		"description": "",
		"logo":        "",
		"release":     "test-v1.0",
		"id":          "_/test-v1.0",
		"name":        "test",
		"branch":      "",
		"revision":    "",
		"project":     "project_name",
		"environment": "dev",
		"deployment":  name,
	}
	for key, val := range strings {
		c.Assert(script.IsStringAtom(dict[key]), Equals, true, Commentf("Expecting %s to be of type string, but was %T", key, dict[key]))
		c.Assert(script.ExpectStringAtom(dict[key]), Equals, val, Commentf("Expecting '%s' to be '%s'. Got '%s'", key, val, script.ExpectStringAtom(dict[key])))
	}
	for key, keys := range dicts {
		c.Assert(script.IsDictAtom(dict[key]), Equals, true, Commentf("Expecting %s to be of type dict, but was %T", key, dict[key]))
		d := script.ExpectDictAtom(dict[key])
		c.Assert(d, HasLen, len(keys), Commentf("Expecting %d values in %s dict.", len(keys), key))
		for _, k := range keys {
			c.Assert(script.IsStringAtom(d[k]), Equals, true, Commentf("Expecting %s to be of type string, but was %T", k, d[k]))
		}
	}
}
