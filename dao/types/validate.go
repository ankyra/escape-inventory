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
	Validate_GetApplication(dao(), c)
	Validate_GetApplications(dao(), c)
	Validate_FindAllVersions(dao(), c)
	Validate_FindAllVersions_Empty(dao(), c)
	Validate_GetPackageURIs(dao(), c)
	Validate_AddPackageURI_Unique(dao(), c)
	Validate_GetAllReleases(dao(), c)
	Validate_ACL(dao(), c)
}

func addReleaseToProject(dao DAO, c *C, name, version, project string) *Release {
	metadataJson := `{"name": "` + name + `", "version": "` + version + `"}`
	metadata, err := core.NewReleaseMetadataFromJsonString(metadataJson)
	c.Assert(err, IsNil)
	result, err := dao.AddRelease(project, metadata)
	c.Assert(err, IsNil)
	return result
}

func addRelease(dao DAO, c *C, name, version string) *Release {
	return addReleaseToProject(dao, c, name, version, "_")
}

func Validate_AddRelease_Unique(dao DAO, c *C) {
	metadataJson := `{"name": "dao-val", "version": "1"}`
	metadata, err := core.NewReleaseMetadataFromJsonString(metadataJson)
	c.Assert(err, IsNil)
	_, err = dao.AddRelease("_", metadata)
	c.Assert(err, IsNil)
	_, err = dao.AddRelease("_", metadata)
	c.Assert(err, Equals, AlreadyExists)
}

func Validate_AddRelease_Unique_per_project(dao DAO, c *C) {
	metadataJson := `{"name": "dao-val", "version": "1"}`
	metadata, err := core.NewReleaseMetadataFromJsonString(metadataJson)
	c.Assert(err, IsNil)
	_, err = dao.AddRelease("_", metadata)
	c.Assert(err, IsNil)
	_, err = dao.AddRelease("my-project", metadata)
	c.Assert(err, IsNil)
}

func Validate_GetRelease(dao DAO, c *C) {
	metadataJson := `{"name": "dao-val", "version": "1"}`
	metadata, err := core.NewReleaseMetadataFromJsonString(metadataJson)
	c.Assert(err, IsNil)
	_, err = dao.AddRelease("_", metadata)
	c.Assert(err, IsNil)
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

	addReleaseToProject(dao, c, "test", "0.0.1", "_")
	addReleaseToProject(dao, c, "test", "0.0.1", "project1")
	addReleaseToProject(dao, c, "test", "0.0.1", "project2")
	addReleaseToProject(dao, c, "test", "0.0.2", "project2")

	projects, err := dao.GetProjects()
	c.Assert(err, IsNil)
	c.Assert(projects, HasLen, 3)
	c.Assert(projects, HasItem, "_")
	c.Assert(projects, HasItem, "project1")
	c.Assert(projects, HasItem, "project2")
}

func Validate_ProjectMetadata(dao DAO, c *C) {
	_, err := dao.GetProject("test")
	c.Assert(err, Equals, NotFound)

	update := NewProject("test")
	update.Description = "yo"
	c.Assert(dao.UpdateProject(update), Equals, NotFound)
	c.Assert(dao.AddProject(update), IsNil)

	project, err := dao.GetProject("test")
	c.Assert(err, IsNil)
	c.Assert(project.Name, Equals, "test")
	c.Assert(project.Description, Equals, "yo")

	c.Assert(dao.AddProject(update), Equals, AlreadyExists)

	update = NewProject("test")
	update.Description = "new description"
	c.Assert(dao.UpdateProject(update), IsNil)

	project, err = dao.GetProject("test")
	c.Assert(err, IsNil)
	c.Assert(project.Description, Equals, "new description")
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

	addReleaseToProject(dao, c, "test", "0.0.1", "_")
	addReleaseToProject(dao, c, "test", "0.0.1", "project1")
	addReleaseToProject(dao, c, "test", "0.0.1", "project2")
	addReleaseToProject(dao, c, "test", "0.0.2", "project2")

	c.Assert(dao.SetACL("_", "*", ReadPermission), IsNil)

	for _, testCase := range cases {
		projects, err := dao.GetProjectsByGroups(testCase)
		c.Assert(err, IsNil)
		c.Assert(projects, HasLen, 1, Commentf("%s should have one group, got %v", testCase, projects))
		c.Assert(projects, HasItem, "_")
	}

	c.Assert(dao.SetACL("project1", "project1", ReadPermission), IsNil)

	projects, err := dao.GetProjectsByGroups(anon)
	c.Assert(err, IsNil)
	c.Assert(projects, HasLen, 1)

	projects, err = dao.GetProjectsByGroups(oneGroup)
	c.Assert(err, IsNil)
	c.Assert(projects, HasLen, 2)
	c.Assert(projects, HasItem, "_")
	c.Assert(projects, HasItem, "project1")

	projects, err = dao.GetProjectsByGroups(allGroups)
	c.Assert(err, IsNil)
	c.Assert(projects, HasLen, 2)
	c.Assert(projects, HasItem, "_")
	c.Assert(projects, HasItem, "project1")

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
	c.Assert(projects, HasItem, "_")
	c.Assert(projects, HasItem, "project1")
	c.Assert(projects, HasItem, "project2")
}

func Validate_GetApplication(dao DAO, c *C) {
	_, err := dao.GetApplication("_", "dao-val")
	c.Assert(err, Equals, NotFound)
	addRelease(dao, c, "dao-val", "0.0.1")
	app, err := dao.GetApplication("_", "dao-val")
	c.Assert(err, IsNil)
	c.Assert(app.Name, Equals, "dao-val")
}

func Validate_GetApplications(dao DAO, c *C) {
	addRelease(dao, c, "archive-dao-archive", "0.1")
	addRelease(dao, c, "archive-dao-archive", "0.2")
	addRelease(dao, c, "ansible-dao-ansible", "0.1")
	addRelease(dao, c, "ansible-dao-ansible", "0.2")
	addReleaseToProject(dao, c, "other-project", "0.3", "other-rpoject")
	applications, err := dao.GetApplications("_")
	c.Assert(err, IsNil)
	var archive, ansible *Application
	for _, app := range applications {
		if app.Name == "ansible-dao-ansible" {
			ansible = app
		} else if app.Name == "archive-dao-archive" {
			archive = app
		} else {
			c.Fail()
		}
	}
	c.Assert(applications, HasLen, 2)
	c.Assert(archive, Not(IsNil))
	c.Assert(ansible, Not(IsNil))
	c.Assert(archive.Name, Equals, "archive-dao-archive")
	c.Assert(ansible.Name, Equals, "ansible-dao-ansible")
}

func Validate_FindAllVersions(dao DAO, c *C) {
	addRelease(dao, c, "dao-val", "0.0.1")
	addRelease(dao, c, "dao-val", "0.0.2")
	addReleaseToProject(dao, c, "dao-val", "0.0.3", "other-project")
	app, err := dao.GetApplication("_", "dao-val")
	c.Assert(err, IsNil)
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
	addRelease(dao, c, "dao-val", "0.1")
	app, err := dao.GetApplication("_", "dao-val")
	c.Assert(err, IsNil)
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
	metadataJson := `{"name": "dao-val", "version": "1"}`
	metadata, err := core.NewReleaseMetadataFromJsonString(metadataJson)
	c.Assert(err, IsNil)
	release, err := dao.AddRelease("_", metadata)
	c.Assert(err, IsNil)
	err = dao.AddPackageURI(release, "file:///test.txt")
	c.Assert(err, IsNil)
	err = dao.AddPackageURI(release, "file:///test.txt")
	c.Assert(err, Equals, AlreadyExists)
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
