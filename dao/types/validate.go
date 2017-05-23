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
	Validate_GetApplication(dao(), c)
	Validate_GetApplications(dao(), c)
	Validate_FindAllVersions(dao(), c)
	Validate_FindAllVersions_Empty(dao(), c)
	Validate_GetPackageURIs(dao(), c)
	Validate_AddPackageURI_Unique(dao(), c)
	Validate_GetAllReleases(dao(), c)
	Validate_ACL(dao(), c)
}

func addReleaseToProject(dao DAO, c *C, name, version, project string) ReleaseDAO {
	metadataJson := `{"name": "` + name + `", "version": "` + version + `"}`
	metadata, err := core.NewReleaseMetadataFromJsonString(metadataJson)
	c.Assert(err, IsNil)
	result, err := dao.AddRelease(project, metadata)
	c.Assert(err, IsNil)
	return result
}

func addRelease(dao DAO, c *C, name, version string) ReleaseDAO {
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
	c.Assert(release.GetVersion(), Equals, "1")
	c.Assert(release.GetApplication().GetName(), Equals, "dao-val")
	c.Assert(release.GetMetadata().GetVersion(), Equals, "1")
	_, err = dao.GetRelease("other-project", "dao-val", "dao-val-v1")
	c.Assert(err, Equals, NotFound)
}

func Validate_GetRelease_NotFound(dao DAO, c *C) {
	_, err := dao.GetRelease("_", "archive-dao-val", "archive-dao-val-v1")
	c.Assert(err, Equals, NotFound)
}

func Validate_GetApplication(dao DAO, c *C) {
	_, err := dao.GetApplication("_", "dao-val")
	c.Assert(err, Equals, NotFound)
	addRelease(dao, c, "dao-val", "0.0.1")
	app, err := dao.GetApplication("_", "dao-val")
	c.Assert(err, IsNil)
	c.Assert(app.GetName(), Equals, "dao-val")
}

func Validate_GetApplications(dao DAO, c *C) {
	addRelease(dao, c, "archive-dao-archive", "0.1")
	addRelease(dao, c, "archive-dao-archive", "0.2")
	addRelease(dao, c, "ansible-dao-ansible", "0.1")
	addRelease(dao, c, "ansible-dao-ansible", "0.2")
	addReleaseToProject(dao, c, "other-project", "0.3", "other-rpoject")
	applications, err := dao.GetApplications("_")
	c.Assert(err, IsNil)
	var archive, ansible ApplicationDAO
	for _, app := range applications {
		if app.GetName() == "ansible-dao-ansible" {
			ansible = app
		} else if app.GetName() == "archive-dao-archive" {
			archive = app
		} else {
			c.Fail()
		}
	}
	c.Assert(applications, HasLen, 2)
	c.Assert(archive, Not(IsNil))
	c.Assert(ansible, Not(IsNil))
	c.Assert(archive.GetName(), Equals, "archive-dao-archive")
	c.Assert(ansible.GetName(), Equals, "ansible-dao-ansible")
}

func Validate_FindAllVersions(dao DAO, c *C) {
	addRelease(dao, c, "dao-val", "0.0.1")
	addRelease(dao, c, "dao-val", "0.0.2")
	addReleaseToProject(dao, c, "dao-val", "0.0.3", "other-project")
	app, err := dao.GetApplication("_", "dao-val")
	c.Assert(err, IsNil)
	versions, err := app.FindAllVersions()
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
	versions, err := app.FindAllVersions()
	c.Assert(err, IsNil)
	c.Assert(len(versions), Equals, 1)
}

func Validate_GetPackageURIs(dao DAO, c *C) {
	release := addRelease(dao, c, "dao-val", "1")
	_ = addReleaseToProject(dao, c, "dao-val", "1", "other-project")
	err := release.AddPackageURI("file:///test.txt")
	c.Assert(err, IsNil)
	err = release.AddPackageURI("gcs:///test.txt")
	c.Assert(err, IsNil)

	release, err = dao.GetRelease("_", "dao-val", "dao-val-v1")
	c.Assert(err, IsNil)

	uris, err := release.GetPackageURIs()
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
	uris, err = release.GetPackageURIs()
	c.Assert(err, IsNil)
	c.Assert(uris, HasLen, 0)
}

func Validate_AddPackageURI_Unique(dao DAO, c *C) {
	metadataJson := `{"name": "dao-val", "version": "1"}`
	metadata, err := core.NewReleaseMetadataFromJsonString(metadataJson)
	c.Assert(err, IsNil)
	release, err := dao.AddRelease("_", metadata)
	c.Assert(err, IsNil)
	err = release.AddPackageURI("file:///test.txt")
	c.Assert(err, IsNil)
	err = release.AddPackageURI("file:///test.txt")
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
	err = dao.SetACL("_", "admin", ReadAndWritePermission)
	c.Assert(err, IsNil)

	groups, err := dao.GetPermittedGroups("_", ReadPermission)
	c.Assert(err, IsNil)
	c.Assert(groups, HasLen, 2)
	c.Assert(groups, HasItem, "*")
	c.Assert(groups, HasItem, "admin")

	groups, err = dao.GetPermittedGroups("_", WritePermission)
	c.Assert(err, IsNil)
	c.Assert(groups, HasLen, 1)
	c.Assert(groups, HasItem, "admin")

	err = dao.DeleteACL("_", "*")
	c.Assert(err, IsNil)

	err = dao.DeleteACL("doesnt-exist", "*")
	c.Assert(err, IsNil)

	groups, err = dao.GetPermittedGroups("_", ReadPermission)
	c.Assert(err, IsNil)
	c.Assert(groups, HasLen, 1)
	c.Assert(groups, DeepEquals, []string{"admin"})

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
