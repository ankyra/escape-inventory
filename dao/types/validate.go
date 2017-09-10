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
	. "gopkg.in/check.v1"
)

func ValidateDAO(dao func() DAO, c *C) {
	Validate_AddRelease_Unique(dao(), c)
	Validate_AddRelease_Unique_per_project(dao(), c)
	Validate_GetRelease(dao(), c)
	Validate_GetRelease_NotFound(dao(), c)
	Validate_GetProjects(dao(), c)
	Validate_ProjectMetadata(dao(), c)
	Validate_GetProjectsByGroups(dao(), c)
	Validate_ApplicationMetadata(dao(), c)
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
	}

	c.Assert(dao.SetACL("project1", "project1", ReadPermission), IsNil)

	projects, err := dao.GetProjectsByGroups(anon)
	c.Assert(err, IsNil)
	c.Assert(projects, HasLen, 1)

	projects, err = dao.GetProjectsByGroups(oneGroup)
	c.Assert(err, IsNil)
	c.Assert(projects, HasLen, 2)
	c.Assert(projects["_"].Name, Equals, "_")
	c.Assert(projects["project1"].Name, Equals, "project1")

	projects, err = dao.GetProjectsByGroups(allGroups)
	c.Assert(err, IsNil)
	c.Assert(projects, HasLen, 2)
	c.Assert(projects["_"].Name, Equals, "_")
	c.Assert(projects["project1"].Name, Equals, "project1")

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
}

func Validate_ApplicationMetadata(dao DAO, c *C) {
	app := NewApplication("project", "name")
	update := NewApplication("project", "name")
	update.Description = "Test"
	update.UploadedBy = "123-123"
	update.UploadedAt = time.Unix(123, 0)

	_, err := dao.GetApplication("project", "name")
	c.Assert(err, Equals, NotFound)

	c.Assert(dao.AddProject(NewProject("project")), Equals, nil)
	c.Assert(dao.UpdateApplication(update), Equals, NotFound)
	c.Assert(dao.AddApplication(app), IsNil)
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
