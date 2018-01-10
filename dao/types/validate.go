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

package types

import (
	"time"

	"github.com/ankyra/escape-core"
	"github.com/ankyra/escape-core/util"
	. "gopkg.in/check.v1"
)

func ValidateDAO(dao func() DAO, c *C) {
	Validate_AddRelease_Unique(dao(), c)
	Validate_AddRelease_Unique_per_project(dao(), c)
	Validate_AddRelease_Big_Metadata(dao(), c)
	Validate_GetRelease(dao(), c)
	Validate_GetRelease_NotFound(dao(), c)
	Validate_GetProjects(dao(), c)
	Validate_ProjectMetadata(dao(), c)
	Validate_HardDeleteProject(dao(), c)
	Validate_GetProjectsByGroups(dao(), c)
	Validate_ApplicationMetadata(dao(), c)
	Validate_GetDownstreamHooks(dao(), c)
	Validate_GetApplications(dao(), c)
	Validate_FindAllVersions(dao(), c)
	Validate_FindAllVersions_Empty(dao(), c)
	Validate_GetPackageURIs(dao(), c)
	Validate_AddPackageURI_Unique(dao(), c)
	Validate_GetAllReleases(dao(), c)
	Validate_ACL(dao(), c)
	Validate_GetReleasesWithoutProcessedDependencies(dao(), c)
	Validate_Dependencies(dao(), c)
	Validate_DependenciesByGroups(dao(), c)
	Validate_Metrics(dao(), c)
	Validate_Feed(dao(), c)
	Validate_WipeDatabase(dao(), c)
}

func addReleaseToProject(dao DAO, c *C, name, version, project string) *Release {
	app := NewApplication(project, name)
	dao.AddProject(NewProject(project))
	dao.AddApplication(app)
	metadataJson := `{"name": "` + name + `", "version": "` + version + `"}`
	metadata, err := core.NewReleaseMetadataFromJsonString(metadataJson)
	c.Assert(err, IsNil)
	release := NewRelease(app, metadata)
	release.UploadedBy = "123-123"
	release.UploadedAt = time.Unix(123, 0)
	c.Assert(dao.AddRelease(release), IsNil)
	return release
}

func addRelease(dao DAO, c *C, name, version string) *Release {
	return addReleaseToProject(dao, c, name, version, "_")
}

func Validate_AddRelease_Unique(dao DAO, c *C) {
	app := NewApplication("_", "dao-val")
	c.Assert(dao.AddProject(NewProject("_")), IsNil)
	c.Assert(dao.AddApplication(app), IsNil)
	metadataJson := `{"name": "dao-val", "version": "1"}`
	metadata, err := core.NewReleaseMetadataFromJsonString(metadataJson)
	c.Assert(err, IsNil)
	c.Assert(dao.AddRelease(NewRelease(app, metadata)), IsNil)
	c.Assert(dao.AddRelease(NewRelease(app, metadata)), Equals, AlreadyExists)
}

func Validate_AddRelease_Unique_per_project(dao DAO, c *C) {
	app1 := NewApplication("_", "dao-val")
	app2 := NewApplication("my-project", "dao-val")
	c.Assert(dao.AddProject(NewProject("_")), IsNil)
	c.Assert(dao.AddApplication(app1), IsNil)
	c.Assert(dao.AddProject(NewProject("my-project")), IsNil)
	c.Assert(dao.AddApplication(app2), IsNil)
	metadataJson := `{"name": "dao-val", "version": "1"}`
	metadata, err := core.NewReleaseMetadataFromJsonString(metadataJson)
	c.Assert(err, IsNil)
	c.Assert(dao.AddRelease(NewRelease(app1, metadata)), IsNil)
	c.Assert(dao.AddRelease(NewRelease(app2, metadata)), IsNil)
}

func Validate_AddRelease_Big_Metadata(dao DAO, c *C) {
	app1 := NewApplication("_", "dao-val")
	c.Assert(dao.AddProject(NewProject("_")), IsNil)
	c.Assert(dao.AddApplication(app1), IsNil)
	longValue := util.RandomString(70000)
	metadataJson := `{"name": "dao-val", "metadata": {"key": "` + longValue + `"}, "version": "1"}`
	metadata, err := core.NewReleaseMetadataFromJsonString(metadataJson)
	c.Assert(err, IsNil)
	c.Assert(dao.AddRelease(NewRelease(app1, metadata)), IsNil)
	release, err := dao.GetRelease("_", "dao-val", "dao-val-v1")
	c.Assert(err, IsNil)
	c.Assert(release.Metadata.Metadata["key"], Equals, longValue)
}

func Validate_GetRelease(dao DAO, c *C) {
	addRelease(dao, c, "dao-val", "1")
	release, err := dao.GetRelease("_", "dao-val", "dao-val-v1")
	c.Assert(err, IsNil)
	c.Assert(release.Version, Equals, "1")
	c.Assert(release.Application.Name, Equals, "dao-val")
	c.Assert(release.Metadata.Version, Equals, "1")
	_, err = dao.GetRelease("other-project", "dao-val", "dao-val-v1")
	c.Assert(err, Equals, NotFound)
}

func Validate_GetRelease_NotFound(dao DAO, c *C) {
	_, err := dao.GetRelease("_", "archive-dao-val", "archive-dao-val-v1")
	c.Assert(err, Equals, NotFound)
}

func Validate_GetProjects(dao DAO, c *C) {
	empty, err := dao.GetProjects()
	c.Assert(err, IsNil)
	c.Assert(empty, HasLen, 0)

	c.Assert(dao.AddProject(NewProject("_")), IsNil)
	c.Assert(dao.AddProject(NewProject("project1")), Equals, nil)
	c.Assert(dao.AddProject(NewProject("project2")), Equals, nil)

	projects, err := dao.GetProjects()
	c.Assert(err, IsNil)
	c.Assert(projects, HasLen, 3)
	c.Assert(projects["_"].Name, Equals, "_")
	c.Assert(projects["project1"].Name, Equals, "project1")
	c.Assert(projects["project2"].Name, Equals, "project2")

	// default permission
	c.Assert(projects["_"].Permission, Equals, "admin")
	c.Assert(projects["project1"].Permission, Equals, "admin")
	c.Assert(projects["project2"].Permission, Equals, "admin")
}

func Validate_ProjectMetadata(dao DAO, c *C) {
	update := NewProject("test")
	update.Description = "yo"
	_, err := dao.GetProject("test")
	c.Assert(err, Equals, NotFound)
	hooks, err := dao.GetProjectHooks(update)
	c.Assert(err, Equals, NotFound)
	c.Assert(dao.UpdateProject(update), Equals, NotFound)
	c.Assert(dao.AddProject(update), IsNil)
	c.Assert(dao.AddProject(update), Equals, AlreadyExists)
	hooks, err = dao.GetProjectHooks(update)
	c.Assert(err, IsNil)

	project, err := dao.GetProject("test")
	c.Assert(err, IsNil)
	c.Assert(project.Name, Equals, "test")
	c.Assert(project.Description, Equals, "yo")

	update = NewProject("test")
	update.Description = "new description"
	c.Assert(dao.UpdateProject(update), IsNil)

	project, err = dao.GetProject("test")
	c.Assert(err, IsNil)
	c.Assert(project.Description, Equals, "new description")

	hooks, err = dao.GetProjectHooks(project)
	c.Assert(err, IsNil)
	c.Assert(hooks, HasLen, 0)

	newHooks := NewHooks()
	newHooks["slack"] = map[string]string{}
	newHooks["slack"]["url"] = "http://example.com"
	c.Assert(dao.SetProjectHooks(project, newHooks), IsNil)

	hooks, err = dao.GetProjectHooks(project)
	c.Assert(err, IsNil)
	c.Assert(hooks, HasLen, 1)
	c.Assert(hooks["slack"]["url"], Equals, "http://example.com")
}

func Validate_HardDeleteProject(dao DAO, c *C) {
	prj := NewProject("_")
	otherPrj := NewProject("other-prj")
	c.Assert(dao.AddProject(prj), IsNil)
	c.Assert(dao.AddProject(otherPrj), IsNil)
	c.Assert(dao.AddApplication(NewApplication("_", "yoooo")), IsNil)
	release := addRelease(dao, c, "dao-val", "1")
	c.Assert(dao.AddPackageURI(release, "http://example.com"), IsNil)
	deps := []*Dependency{
		NewDependency("asdoias", "eroijaeirjo", "1.0"),
		NewDependency("_", "yo", "1.0"),
	}
	c.Assert(dao.SetDependencies(release, deps), IsNil)
	addRelease(dao, c, "dao-val", "2")
	dao.SetACL("_", "group1", ReadPermission)
	dao.AddFeedEvent(NewCreateProjectEvent("_", "admin"))
	hooks := NewHooks()
	hooks["test"] = map[string]string{
		"wut": "wat",
	}
	c.Assert(dao.SetProjectHooks(prj, hooks), IsNil)
	app := NewApplication("_", "dao-val")
	c.Assert(dao.SetApplicationHooks(app, hooks), IsNil)
	otherApp := NewApplication("other-prj", "yo-yo")
	c.Assert(dao.AddApplication(otherApp), IsNil)
	c.Assert(dao.SetApplicationHooks(otherApp, hooks), IsNil)
	c.Assert(dao.SetApplicationSubscribesToUpdatesFrom(otherApp, []*Application{app}), IsNil)
	c.Assert(dao.SetApplicationSubscribesToUpdatesFrom(app, []*Application{otherApp}), IsNil)

	// Before Delete
	_, err := dao.GetProject("_")
	c.Assert(err, IsNil)
	c.Assert(dao.UpdateProject(prj), IsNil)
	prjs, err := dao.GetProjects()
	c.Assert(err, IsNil)
	c.Assert(prjs, HasLen, 2)
	hooks, err = dao.GetProjectHooks(prj)
	c.Assert(err, IsNil)
	c.Assert(hooks, HasLen, 1)
	_, err = dao.GetApplication("_", "dao-val")
	c.Assert(err, IsNil)
	apps, err := dao.GetApplications("_")
	c.Assert(err, IsNil)
	c.Assert(apps, HasLen, 2)
	appHooks, err := dao.GetApplicationHooks(app)
	c.Assert(err, IsNil)
	c.Assert(appHooks, HasLen, 1)
	c.Assert(dao.AddApplication(NewApplication("_", "dao-val")), Equals, AlreadyExists)
	r, err := dao.GetRelease("_", "dao-val", "dao-val-v1")
	c.Assert(err, IsNil)
	c.Assert(r, DeepEquals, release)
	rels, err := dao.GetAllReleases()
	c.Assert(err, IsNil)
	c.Assert(rels, HasLen, 2)
	pkgs, err := dao.GetPackageURIs(release)
	c.Assert(err, IsNil)
	c.Assert(pkgs, HasLen, 1)
	down, err := dao.GetDownstreamHooks(app)
	c.Assert(err, IsNil)
	c.Assert(down, HasLen, 1)
	down, err = dao.GetDownstreamHooks(otherApp)
	c.Assert(err, IsNil)
	c.Assert(down, HasLen, 1)
	deps, err = dao.GetDependencies(release)
	c.Assert(err, IsNil)
	c.Assert(deps, HasLen, 2)
	perm, err := dao.GetACL("_")
	c.Assert(err, IsNil)
	c.Assert(perm["group1"], Equals, ReadPermission)
	feed, err := dao.GetProjectFeedPage("_", 1)
	c.Assert(err, IsNil)
	c.Assert(feed, HasLen, 1)
	c.Assert(dao.AddProject(prj), Equals, AlreadyExists)

	// Delete
	c.Assert(dao.HardDeleteProject("_"), IsNil)

	// After Delete
	_, err = dao.GetProject("_")
	c.Assert(err, Equals, NotFound)
	c.Assert(dao.UpdateProject(prj), Equals, NotFound)
	prjs, err = dao.GetProjects()
	c.Assert(err, IsNil)
	c.Assert(prjs, HasLen, 1)
	_, err = dao.GetProjectHooks(prj)
	c.Assert(err, Equals, NotFound)
	_, err = dao.GetApplicationHooks(app)
	c.Assert(err, Equals, NotFound)
	_, err = dao.GetApplication("_", "dao-val")
	c.Assert(err, Equals, NotFound)
	apps, err = dao.GetApplications("_")
	c.Assert(err, IsNil)
	c.Assert(apps, HasLen, 0)
	_, err = dao.GetRelease("_", "dao-val", "dao-val-v1")
	c.Assert(err, Equals, NotFound)
	pkgs, err = dao.GetPackageURIs(release)
	c.Assert(err, IsNil)
	c.Assert(pkgs, HasLen, 0)
	deps, err = dao.GetDependencies(release)
	c.Assert(err, IsNil)
	c.Assert(deps, HasLen, 0)
	rels, err = dao.GetAllReleases()
	c.Assert(err, IsNil)
	c.Assert(rels, HasLen, 0)
	down, err = dao.GetDownstreamHooks(app)
	c.Assert(err, IsNil)
	c.Assert(down, HasLen, 1)
	down, err = dao.GetDownstreamHooks(otherApp)
	c.Assert(err, IsNil)
	c.Assert(down, HasLen, 0)
	perm, err = dao.GetACL("_")
	c.Assert(err, IsNil)
	c.Assert(perm, HasLen, 0)
	feed, err = dao.GetProjectFeedPage("_", 1)
	c.Assert(err, IsNil)
	c.Assert(feed, HasLen, 0)

	// Re-adding
	c.Assert(dao.AddProject(prj), IsNil)
	c.Assert(dao.AddApplication(NewApplication("_", "dao-val")), IsNil)
	addRelease(dao, c, "dao-val", "1")

	// After re-adding
	rels, err = dao.GetAllReleases()
	c.Assert(err, IsNil)
	c.Assert(rels, HasLen, 1)
	apps, err = dao.GetApplications("_")
	c.Assert(err, IsNil)
	c.Assert(apps, HasLen, 1)
	hooks, err = dao.GetProjectHooks(prj)
	c.Assert(err, IsNil)
	c.Assert(hooks, HasLen, 0)
	appHooks, err = dao.GetApplicationHooks(app)
	c.Assert(err, IsNil)
	c.Assert(appHooks, HasLen, 0)
	pkgs, err = dao.GetPackageURIs(release)
	c.Assert(err, IsNil)
	c.Assert(pkgs, HasLen, 0)
	deps, err = dao.GetDependencies(release)
	c.Assert(err, IsNil)
	c.Assert(deps, HasLen, 0)
	down, err = dao.GetDownstreamHooks(app)
	c.Assert(err, IsNil)
	c.Assert(down, HasLen, 1)
	down, err = dao.GetDownstreamHooks(otherApp)
	c.Assert(err, IsNil)
	c.Assert(down, HasLen, 0)
}

func Validate_GetProjectsByGroups(dao DAO, c *C) {
	anon := []string{}
	oneGroup := []string{"project1"}
	allGroups := []string{"project1", "project2"}
	cases := [][]string{anon, oneGroup, allGroups}

	for _, testCase := range cases {
		empty, err := dao.GetProjectsByGroups(testCase)
		c.Assert(err, IsNil)
		c.Assert(empty, HasLen, 0)
	}

	c.Assert(dao.AddProject(NewProject("_")), Equals, nil)
	c.Assert(dao.AddProject(NewProject("project1")), Equals, nil)
	c.Assert(dao.AddProject(NewProject("project2")), Equals, nil)

	c.Assert(dao.SetACL("_", "*", ReadPermission), IsNil)

	for _, testCase := range cases {
		projects, err := dao.GetProjectsByGroups(testCase)
		c.Assert(err, IsNil)
		c.Assert(projects, HasLen, 1, Commentf("%s should have one group, got %v", testCase, projects))
		c.Assert(projects["_"].Name, Equals, "_")
		c.Assert(projects["_"].Permission, Equals, "admin")
		c.Assert(projects["_"].MatchingGroups, DeepEquals, []string{"*"})
	}

	c.Assert(dao.SetACL("project1", "project1", ReadPermission), IsNil)
	c.Assert(dao.SetACL("_", "project1", WritePermission), IsNil)

	projects, err := dao.GetProjectsByGroups(anon)
	c.Assert(err, IsNil)
	c.Assert(projects, HasLen, 1)

	projects, err = dao.GetProjectsByGroups(oneGroup)
	c.Assert(err, IsNil)
	c.Assert(projects, HasLen, 2)
	c.Assert(projects["_"].Name, Equals, "_")
	c.Assert(projects["_"].MatchingGroups, HasLen, 2)
	c.Assert(projects["_"].MatchingGroups, HasItem, "project1")
	c.Assert(projects["_"].MatchingGroups, HasItem, "*")
	c.Assert(projects["project1"].Name, Equals, "project1")
	c.Assert(projects["project1"].MatchingGroups, DeepEquals, []string{"project1"})
	c.Assert(projects["_"].Permission, Equals, "admin")
	c.Assert(projects["project1"].Permission, Equals, "admin")

	projects, err = dao.GetProjectsByGroups(allGroups)
	c.Assert(err, IsNil)
	c.Assert(projects, HasLen, 2)
	c.Assert(projects["_"].Name, Equals, "_")
	c.Assert(projects["_"].MatchingGroups, HasLen, 2)
	c.Assert(projects["_"].MatchingGroups, HasItem, "project1")
	c.Assert(projects["_"].MatchingGroups, HasItem, "*")
	c.Assert(projects["project1"].Name, Equals, "project1")
	c.Assert(projects["project1"].MatchingGroups, DeepEquals, []string{"project1"})
	c.Assert(projects["_"].Permission, Equals, "admin")
	c.Assert(projects["project1"].Permission, Equals, "admin")

	c.Assert(dao.SetACL("project1", "project2", WritePermission), IsNil)
	c.Assert(dao.SetACL("project1", "project2", ReadPermission), IsNil)
	c.Assert(dao.SetACL("project2", "project2", ReadPermission), IsNil)

	projects, err = dao.GetProjectsByGroups(anon)
	c.Assert(err, IsNil)
	c.Assert(projects, HasLen, 1)

	projects, err = dao.GetProjectsByGroups(oneGroup)
	c.Assert(err, IsNil)
	c.Assert(projects, HasLen, 2)

	projects, err = dao.GetProjectsByGroups(allGroups)
	c.Assert(err, IsNil)
	c.Assert(projects, HasLen, 3)
	c.Assert(projects["_"].Name, Equals, "_")
	c.Assert(projects["project1"].Name, Equals, "project1")
	c.Assert(projects["project2"].Name, Equals, "project2")
	c.Assert(projects["_"].MatchingGroups, HasLen, 2)
	c.Assert(projects["_"].MatchingGroups, HasItem, "project1")
	c.Assert(projects["_"].MatchingGroups, HasItem, "*")
	c.Assert(projects["project1"].MatchingGroups, HasLen, 2)
	c.Assert(projects["project1"].MatchingGroups, HasItem, "project1")
	c.Assert(projects["project1"].MatchingGroups, HasItem, "project2")
	c.Assert(projects["project2"].MatchingGroups, DeepEquals, []string{"project2"})
	c.Assert(projects["_"].Permission, Equals, "admin")
	c.Assert(projects["project1"].Permission, Equals, "admin")
	c.Assert(projects["project2"].Permission, Equals, "admin")
}

func Validate_ApplicationMetadata(dao DAO, c *C) {
	app := NewApplication("project", "name")
	update := NewApplication("project", "name")
	update.Description = "Test"
	update.UploadedBy = "123-123"
	update.UploadedAt = time.Unix(123, 0)

	_, err := dao.GetApplication("project", "name")
	c.Assert(err, Equals, NotFound)
	hooks, err := dao.GetApplicationHooks(app)
	c.Assert(err, Equals, NotFound)

	c.Assert(dao.AddProject(NewProject("project")), Equals, nil)
	hooks, err = dao.GetApplicationHooks(app)
	c.Assert(err, Equals, NotFound)
	c.Assert(hooks, HasLen, 0)
	c.Assert(dao.UpdateApplication(update), Equals, NotFound)
	c.Assert(dao.AddApplication(app), IsNil)
	hooks, err = dao.GetApplicationHooks(app)
	c.Assert(err, IsNil)
	c.Assert(hooks, HasLen, 0)
	c.Assert(dao.AddApplication(app), Equals, AlreadyExists)

	app, err = dao.GetApplication("project", "name")
	c.Assert(err, IsNil)
	c.Assert(app.Name, Equals, "name")
	c.Assert(app.Project, Equals, "project")
	c.Assert(app.Description, Equals, "")

	c.Assert(dao.UpdateApplication(update), IsNil)

	app, err = dao.GetApplication("project", "name")
	c.Assert(err, IsNil)
	c.Assert(app.Name, Equals, "name")
	c.Assert(app.Project, Equals, "project")
	c.Assert(app.Description, Equals, "Test")
	c.Assert(app.UploadedBy, Equals, "123-123")
	c.Assert(app.UploadedAt, Equals, time.Unix(123, 0))

	newHooks := NewHooks()
	newHooks["slack"] = map[string]string{}
	newHooks["slack"]["url"] = "http://example.com"
	c.Assert(dao.SetApplicationHooks(app, newHooks), IsNil)

	hooks, err = dao.GetApplicationHooks(app)
	c.Assert(err, IsNil)
	c.Assert(hooks, HasLen, 1)
	c.Assert(hooks["slack"]["url"], Equals, "http://example.com")
}

func Validate_GetDownstreamHooks(dao DAO, c *C) {
	app1 := NewApplication("project", "app1")
	app2 := NewApplication("project", "app2")
	app3 := NewApplication("project", "app3")
	c.Assert(dao.AddProject(NewProject("project")), Equals, nil)
	c.Assert(dao.AddApplication(app1), IsNil)
	c.Assert(dao.AddApplication(app2), IsNil)
	c.Assert(dao.AddApplication(app3), IsNil)
	c.Assert(dao.SetApplicationSubscribesToUpdatesFrom(app1, []*Application{app2, app3}), IsNil)
	c.Assert(dao.SetApplicationSubscribesToUpdatesFrom(app2, []*Application{app3}), IsNil)
	hooks1, err := dao.GetDownstreamHooks(app1)
	c.Assert(err, IsNil)
	hooks2, err := dao.GetDownstreamHooks(app2)
	c.Assert(err, IsNil)
	hooks3, err := dao.GetDownstreamHooks(app3)
	c.Assert(err, IsNil)
	c.Assert(hooks1, HasLen, 0)
	c.Assert(hooks2, HasLen, 1)
	c.Assert(hooks3, HasLen, 2)
}

func Validate_GetApplications(dao DAO, c *C) {
	app1 := NewApplication("_", "archive")
	app1.Description = "archive stuff"
	app2 := NewApplication("_", "ansible")
	app2.Description = "ansible stuff"
	app3 := NewApplication("other-project", "whatever")
	c.Assert(dao.AddProject(NewProject("_")), IsNil)
	c.Assert(dao.AddProject(NewProject("other-project")), IsNil)
	c.Assert(dao.AddApplication(app1), IsNil)
	c.Assert(dao.AddApplication(app2), IsNil)
	c.Assert(dao.AddApplication(app3), IsNil)
	applications, err := dao.GetApplications("_")
	c.Assert(err, IsNil)
	archive := applications["archive"]
	ansible := applications["ansible"]
	c.Assert(applications, HasLen, 2)
	c.Assert(archive, Not(IsNil))
	c.Assert(ansible, Not(IsNil))
	c.Assert(archive.Name, Equals, "archive")
	c.Assert(archive.Description, Equals, "archive stuff")
	c.Assert(ansible.Name, Equals, "ansible")
	c.Assert(ansible.Description, Equals, "ansible stuff")
}

func Validate_FindAllVersions(dao DAO, c *C) {
	c.Assert(dao.AddProject(NewProject("_")), IsNil)
	app := NewApplication("_", "dao-val")
	c.Assert(dao.AddApplication(app), IsNil)
	addRelease(dao, c, "dao-val", "0.0.1")
	addRelease(dao, c, "dao-val", "0.0.2")
	addReleaseToProject(dao, c, "dao-val", "0.0.3", "other-project")
	versions, err := dao.FindAllVersions(app)
	c.Assert(err, IsNil)
	c.Assert(len(versions), Equals, 2)
	var firstFound, secondFound bool
	for _, v := range versions {
		if v == "0.0.1" {
			firstFound = true
		} else if v == "0.0.2" {
			secondFound = true
		} else {
			c.Fail()
		}
	}
	c.Assert(firstFound, Equals, true)
	c.Assert(secondFound, Equals, true)
}

func Validate_FindAllVersions_Empty(dao DAO, c *C) {
	c.Assert(dao.AddProject(NewProject("_")), IsNil)
	app := NewApplication("_", "dao-val")
	c.Assert(dao.AddApplication(app), IsNil)
	addRelease(dao, c, "dao-val", "0.1")
	versions, err := dao.FindAllVersions(app)
	c.Assert(err, IsNil)
	c.Assert(len(versions), Equals, 1)
}

func Validate_GetPackageURIs(dao DAO, c *C) {
	release := addRelease(dao, c, "dao-val", "1")
	_ = addReleaseToProject(dao, c, "dao-val", "1", "other-project")
	err := dao.AddPackageURI(release, "file:///test.txt")
	c.Assert(err, IsNil)
	err = dao.AddPackageURI(release, "gcs:///test.txt")
	c.Assert(err, IsNil)

	release, err = dao.GetRelease("_", "dao-val", "dao-val-v1")
	c.Assert(err, IsNil)

	uris, err := dao.GetPackageURIs(release)
	c.Assert(err, IsNil)
	var fileFound, gcsFound bool
	for _, uri := range uris {
		if uri == "file:///test.txt" {
			fileFound = true
		}
		if uri == "gcs:///test.txt" {
			gcsFound = true
		}
	}
	c.Assert(fileFound, Equals, true)
	c.Assert(gcsFound, Equals, true)

	release, err = dao.GetRelease("other-project", "dao-val", "dao-val-v1")
	c.Assert(err, IsNil)
	uris, err = dao.GetPackageURIs(release)
	c.Assert(err, IsNil)
	c.Assert(uris, HasLen, 0)
}

func Validate_AddPackageURI_Unique(dao DAO, c *C) {
	release := addRelease(dao, c, "dao-val", "1")
	c.Assert(dao.AddPackageURI(release, "file:///test.txt"), IsNil)
	c.Assert(dao.AddPackageURI(release, "file:///test.txt"), Equals, AlreadyExists)
}

func Validate_GetAllReleases(dao DAO, c *C) {
	addRelease(dao, c, "dao-val", "0.1")
	addRelease(dao, c, "dao-val", "0.2")
	releases, err := dao.GetAllReleases()
	c.Assert(err, IsNil)
	c.Assert(releases, HasLen, 2)
}

func Validate_ACL(dao DAO, c *C) {
	err := dao.SetACL("_", "*", ReadPermission)
	c.Assert(err, IsNil)
	err = dao.SetACL("_", "writer", WritePermission)
	c.Assert(err, IsNil)
	err = dao.SetACL("_", "admin", AdminPermission)
	c.Assert(err, IsNil)

	groups, err := dao.GetPermittedGroups("_", ReadPermission)
	c.Assert(err, IsNil)
	c.Assert(groups, HasLen, 3)
	c.Assert(groups, HasItem, "*")
	c.Assert(groups, HasItem, "writer")
	c.Assert(groups, HasItem, "admin")

	groups, err = dao.GetPermittedGroups("_", WritePermission)
	c.Assert(err, IsNil)
	c.Assert(groups, HasLen, 2)
	c.Assert(groups, HasItem, "writer")
	c.Assert(groups, HasItem, "admin")

	members, err := dao.GetACL("_")
	c.Assert(err, IsNil)
	c.Assert(members["*"], Equals, ReadPermission)
	c.Assert(members["writer"], Equals, WritePermission)
	c.Assert(members["admin"], Equals, AdminPermission)

	err = dao.DeleteACL("_", "*")
	c.Assert(err, IsNil)

	err = dao.DeleteACL("doesnt-exist", "*")
	c.Assert(err, IsNil)

	groups, err = dao.GetPermittedGroups("_", ReadPermission)
	c.Assert(err, IsNil)
	c.Assert(groups, HasLen, 2)
	c.Assert(groups, HasItem, "writer")
	c.Assert(groups, HasItem, "admin")

	groups, err = dao.GetPermittedGroups("doesnt-exist", ReadPermission)
	c.Assert(err, IsNil)
	c.Assert(groups, DeepEquals, []string{})
}

func Validate_GetReleasesWithoutProcessedDependencies(dao DAO, c *C) {
	releases, err := dao.GetAllReleasesWithoutProcessedDependencies()
	c.Assert(err, IsNil)
	c.Assert(releases, HasLen, 0)
	release := addRelease(dao, c, "dao-val", "1")
	c.Assert(release.ProcessedDependencies, Equals, false)
	releases, err = dao.GetAllReleasesWithoutProcessedDependencies()
	c.Assert(err, IsNil)
	c.Assert(releases, HasLen, 1)
	release.ProcessedDependencies = true
	release.Downloads = 14
	c.Assert(dao.UpdateRelease(release), IsNil)
	releases, err = dao.GetAllReleasesWithoutProcessedDependencies()
	c.Assert(err, IsNil)
	c.Assert(releases, HasLen, 0)
	release, err = dao.GetRelease("_", "dao-val", "dao-val-v1")
	c.Assert(err, IsNil)
	c.Assert(release.ProcessedDependencies, Equals, true)
	c.Assert(release.Downloads, Equals, 14)
	c.Assert(release.UploadedBy, Equals, "123-123")
	c.Assert(release.UploadedAt, Equals, time.Unix(123, 0))
}

func Validate_Dependencies(dao DAO, c *C) {
	release := addRelease(dao, c, "dao-val", "1")
	deps, err := dao.GetDependencies(release)
	c.Assert(err, IsNil)
	c.Assert(deps, HasLen, 0)

	app := NewApplication("_", "dao-parent")
	dao.AddApplication(app)
	metadataJson := `{"name": "dao-parent", "version": "1", "depends": [{"id": "_/dao-val-v1", "scopes": ["build"]}]}`
	metadata, err := core.NewReleaseMetadataFromJsonString(metadataJson)
	c.Assert(err, IsNil)
	releaseParent := NewRelease(app, metadata)
	c.Assert(dao.AddRelease(releaseParent), IsNil)
	dependencies := []*Dependency{
		&Dependency{
			Project:     "_",
			Application: "dao-val",
			Version:     "1",
			BuildScope:  true,
			DeployScope: false,
		},
	}
	c.Assert(dao.SetDependencies(releaseParent, dependencies), IsNil)
	deps, err = dao.GetDependencies(releaseParent)
	c.Assert(err, IsNil)
	c.Assert(deps, DeepEquals, dependencies)

	downstream := []*Dependency{
		&Dependency{
			Project:     "_",
			Application: "dao-parent",
			Version:     "1",
			BuildScope:  true,
			DeployScope: false,
		},
	}
	ds, err := dao.GetDownstreamDependencies(release)
	c.Assert(err, IsNil)
	c.Assert(ds, HasLen, 1)
	c.Assert(ds[0], DeepEquals, downstream[0])
}

func Validate_DependenciesByGroups(dao DAO, c *C) {
	c.Assert(dao.SetACL("_", "cheeky-group", ReadPermission), IsNil)

	release := addRelease(dao, c, "dao-val", "1")
	deps, err := dao.GetDependencies(release)
	c.Assert(err, IsNil)
	c.Assert(deps, HasLen, 0)

	app := NewApplication("_", "dao-parent")
	dao.AddApplication(app)
	metadataJson := `{"name": "dao-parent", "version": "1", "depends": [{"id": "_/dao-val-v1", "scopes": ["build"]}]}`
	metadata, err := core.NewReleaseMetadataFromJsonString(metadataJson)
	c.Assert(err, IsNil)
	releaseParent := NewRelease(app, metadata)
	c.Assert(dao.AddRelease(releaseParent), IsNil)
	dependencies := []*Dependency{
		&Dependency{
			Project:     "_",
			Application: "dao-val",
			Version:     "1",
			BuildScope:  true,
			DeployScope: false,
		},
	}
	c.Assert(dao.SetDependencies(releaseParent, dependencies), IsNil)
	ds, err := dao.GetDownstreamDependenciesByGroups(release, []string{})
	c.Assert(err, IsNil)
	c.Assert(ds, HasLen, 0)

	ds, err = dao.GetDownstreamDependenciesByGroups(release, []string{"cheeky-group"})
	c.Assert(err, IsNil)
	c.Assert(ds, HasLen, 1)
}

func Validate_Metrics(dao DAO, c *C) {
	metrics, err := dao.GetUserMetrics("test-user")
	c.Assert(err, IsNil)
	c.Assert(metrics, Not(IsNil))
	c.Assert(metrics.ProjectCount, Equals, 0)

	newMetrics := Metrics{
		ProjectCount: 3,
	}
	err = dao.SetUserMetrics("test-user", metrics, &newMetrics)
	c.Assert(err, IsNil)
	obtained, err := dao.GetUserMetrics("test-user")
	c.Assert(err, IsNil)
	c.Assert(obtained, Not(IsNil))
	c.Assert(obtained.ProjectCount, Equals, 3)

	err = dao.SetUserMetrics("yo-i-dont-exist", metrics, &newMetrics)
	c.Assert(err, Equals, NotFound)
}

func Validate_Feed(dao DAO, c *C) {
	events := []*FeedEvent{
		NewCreateProjectEvent("test1", "user"),
		NewReleaseEvent("test2", "app", "1.0", "user"),
		NewCreateProjectEvent("test2", "user"),
		NewCreateProjectEvent("test3", "user"),
		NewCreateProjectEvent("test4", "user"),
		NewCreateProjectEvent("test5", "user"),
	}
	for _, ev := range events {
		c.Assert(dao.AddFeedEvent(ev), IsNil)
	}
	received, err := dao.GetFeedPage(4)
	c.Assert(err, IsNil)
	c.Assert(received, HasLen, 4)
	expected := [][]*FeedEvent{
		[]*FeedEvent{received[0], events[5]},
		[]*FeedEvent{received[1], events[4]},
		[]*FeedEvent{received[2], events[3]},
		[]*FeedEvent{received[3], events[2]},
	}
	for i, exp := range expected {
		c.Assert(exp[0].Equals(exp[1]), Equals, true, Commentf("%d item '%s' expected; was '%s'", i, exp[1], exp[0]))
	}

	received, err = dao.GetProjectFeedPage("test2", 4)
	c.Assert(err, IsNil)
	c.Assert(received, HasLen, 2)
	c.Assert(received[0].Equals(events[2]), Equals, true, Commentf("expected '%s'; was '%s'", events[2], received[0]))
	c.Assert(received[1].Equals(events[1]), Equals, true, Commentf("expected '%s'; was '%s'", events[1], received[1]))

	c.Assert(dao.SetACL("test1", "test1", ReadPermission), IsNil)
	received, err = dao.GetFeedPageByGroups([]string{"test1"}, 4)
	c.Assert(err, IsNil)
	c.Assert(received, HasLen, 1)
	c.Assert(received[0].Equals(events[0]), Equals, true, Commentf("expected '%s'; was '%s'", events[0], received[0]))
}

func Validate_WipeDatabase(dao DAO, c *C) {
	addReleaseToProject(dao, c, "test", "1.0.0", "test-project")
	dao.WipeDatabase()

	projects, _ := dao.GetProjects()
	c.Assert(projects, HasLen, 0)
}

type hasItemChecker struct{}

var HasItem = &hasItemChecker{}

func (*hasItemChecker) Info() *CheckerInfo {
	return &CheckerInfo{Name: "HasItem", Params: []string{"obtained", "expected to have item"}}
}
func (*hasItemChecker) Check(params []interface{}, names []string) (bool, string) {
	obtained := params[0]
	expectedItem := params[1]
	switch obtained.(type) {
	case []interface{}:
		for _, v := range obtained.([]interface{}) {
			if v == expectedItem {
				return true, ""
			}
		}
	case []string:
		for _, v := range obtained.([]string) {
			if v == expectedItem {
				return true, ""
			}
		}
	default:
		return false, "Unexpected type."
	}
	return false, "Item not found"
}
