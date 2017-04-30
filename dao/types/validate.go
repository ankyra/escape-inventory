package types

import (
	. "gopkg.in/check.v1"
    "github.com/ankyra/escape-client/model/release"
)


func ValidateDAO(dao func() DAO, c *C) {
    Validate_NewApplication(dao(), c)
    Validate_NewApplication_Unique(dao(), c)
    Validate_GetApplication(dao(), c)
    Validate_GetApplications(dao(), c)
    Validate_FindAllVersions(dao(), c)
    Validate_FindAllVersions_Empty(dao(), c)
    Validate_AddRelease_Unique(dao(), c)
    Validate_GetRelease(dao(), c)
    Validate_GetRelease_NotFound(dao(), c)
    Validate_GetPackageURIs(dao(), c)
    Validate_GetAllReleases(dao(), c)
    Validate_AddPackageURI_Unique(dao(), c)
}

func Validate_NewApplication(dao DAO, c *C) {
    app, err := dao.NewApplication("archive", "dao-val")
    c.Assert(err, IsNil)
    c.Assert(app.GetType(), Equals, "archive")
    c.Assert(app.GetName(), Equals, "dao-val")
}

func Validate_NewApplication_Unique(dao DAO, c *C) {
    _, err := dao.NewApplication("archive", "dao-val")
    c.Assert(err, IsNil)
    _, err = dao.NewApplication("archive", "dao-val")
    c.Assert(err, Equals, AlreadyExists)
}

func Validate_GetApplication(dao DAO, c *C) {
    _, err := dao.GetApplication("archive", "dao-val")
    c.Assert(err, Equals, NotFound)
    _, err = dao.NewApplication("archive", "dao-val")
    c.Assert(err, IsNil)
    app, err := dao.GetApplication("archive", "dao-val")
    c.Assert(err, IsNil)
    c.Assert(app.GetType(), Equals, "archive")
    c.Assert(app.GetName(), Equals, "dao-val")
}

func Validate_GetApplications(dao DAO, c *C) {
    _, err := dao.NewApplication("archive", "dao-archive")
    c.Assert(err, IsNil)
    _, err = dao.NewApplication("ansible", "dao-ansible")
    c.Assert(err, IsNil)
    applications, err := dao.GetApplications()
    c.Assert(err, IsNil)
    var archive, ansible ApplicationDAO
    for _, app := range applications {
        if app.GetType() == "archive" {
            archive = app
        } else if app.GetType() == "ansible" {
            ansible = app
        } else {
            c.Fail()
        }
    }
    c.Assert(archive, Not(IsNil))
    c.Assert(ansible, Not(IsNil))
    c.Assert(archive.GetName(), Equals, "dao-archive")
    c.Assert(ansible.GetName(), Equals, "dao-ansible")
}

func Validate_FindAllVersions(dao DAO, c *C) {
    app, err := dao.NewApplication("archive", "dao-val")
    c.Assert(err, IsNil)

    metadataJson := `{"name": "dao-val", "type": "archive", "version": "0.0.1"}`
    metadata, err := release.NewReleaseMetadataFromJsonString(metadataJson)
    c.Assert(err, IsNil)
    err = app.AddRelease(metadata)
    c.Assert(err, IsNil)

    metadataJson = `{"name": "dao-val", "type": "archive", "version": "0.0.2"}`
    metadata, err = release.NewReleaseMetadataFromJsonString(metadataJson)
    c.Assert(err, IsNil)
    err = app.AddRelease(metadata)
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
    app, err := dao.NewApplication("archive", "dao-val")
    c.Assert(err, IsNil)
    versions, err := app.FindAllVersions()
    c.Assert(err, IsNil)
    c.Assert(len(versions), Equals, 0)
}

func Validate_AddRelease_Unique(dao DAO, c *C) {
    app, err := dao.NewApplication("archive", "dao-val")
    c.Assert(err, IsNil)
    metadataJson := `{"name": "dao-val", "type": "archive", "version": "1"}`
    metadata, err := release.NewReleaseMetadataFromJsonString(metadataJson)
    c.Assert(err, IsNil)
    err = app.AddRelease(metadata)
    c.Assert(err, IsNil)
    err = app.AddRelease(metadata)
    c.Assert(err, Equals, AlreadyExists)
}

func Validate_GetRelease(dao DAO, c *C) {
    app, err := dao.NewApplication("archive", "dao-val")
    c.Assert(err, IsNil)
    metadataJson := `{"name": "dao-val", "type": "archive", "version": "1"}`
    metadata, err := release.NewReleaseMetadataFromJsonString(metadataJson)
    c.Assert(err, IsNil)
    err = app.AddRelease(metadata)
    c.Assert(err, IsNil)
    release, err := dao.GetRelease("archive-dao-val-v1")
    c.Assert(err, IsNil)
    c.Assert(release.GetVersion(), Equals, "1")
    c.Assert(release.GetApplication().GetName(), Equals, "dao-val")
    c.Assert(release.GetApplication().GetType(), Equals, "archive")
    c.Assert(release.GetMetadata().GetType(), Equals, "archive")
    c.Assert(release.GetMetadata().GetVersion(), Equals, "1")
}

func Validate_GetRelease_NotFound(dao DAO, c *C) {
    _, err := dao.GetRelease("archive-dao-val-v1")
    c.Assert(err, Equals, NotFound)
}

func Validate_GetPackageURIs(dao DAO, c *C) {
    app, err := dao.NewApplication("archive", "dao-val")
    c.Assert(err, IsNil)
    metadataJson := `{"name": "dao-val", "type": "archive", "version": "1"}`
    metadata, err := release.NewReleaseMetadataFromJsonString(metadataJson)
    c.Assert(err, IsNil)
    err = app.AddRelease(metadata)
    c.Assert(err, IsNil)
    release, err := dao.GetRelease("archive-dao-val-v1")
    c.Assert(err, IsNil)
    err = release.AddPackageURI("file:///test.txt")
    c.Assert(err, IsNil)
    err = release.AddPackageURI("gcs:///test.txt")
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
}

func Validate_AddPackageURI_Unique(dao DAO, c *C) {
    app, err := dao.NewApplication("archive", "dao-val")
    c.Assert(err, IsNil)
    metadataJson := `{"name": "dao-val", "type": "archive", "version": "1"}`
    metadata, err := release.NewReleaseMetadataFromJsonString(metadataJson)
    c.Assert(err, IsNil)
    err = app.AddRelease(metadata)
    c.Assert(err, IsNil)
    release, err := dao.GetRelease("archive-dao-val-v1")
    c.Assert(err, IsNil)
    err = release.AddPackageURI("file:///test.txt")
    c.Assert(err, IsNil)
    err = release.AddPackageURI("file:///test.txt")
    c.Assert(err, Equals, AlreadyExists)
}

func Validate_GetAllReleases(dao DAO, c *C) {
    app, err := dao.NewApplication("archive", "dao-val")
    c.Assert(err, IsNil)
    metadataJson := `{"name": "dao-val", "type": "archive", "version": "1"}`
    metadata, err := release.NewReleaseMetadataFromJsonString(metadataJson)
    c.Assert(err, IsNil)
    err = app.AddRelease(metadata)
    c.Assert(err, IsNil)
    releases, err := dao.GetAllReleases()
    c.Assert(err, IsNil)
    c.Assert(releases, HasLen, 1)
}
