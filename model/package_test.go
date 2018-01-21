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
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/ankyra/escape-inventory/dao"

	. "gopkg.in/check.v1"
)

func (s *appSuite) Test_UploadPackage_happy_path(c *C) {
	dao.TestSetup()
	storage := &storageProvider{
		Upload: func(project, releaseId string, pkg io.ReadSeeker) (string, error) {
			return "mem://" + project + "/" + releaseId + ".tgz", nil
		},
	}

	_, err := AddRelease("project", `{"name": "name", "version": "1.0.0"}`)
	c.Assert(err, IsNil)
	pkg := bytes.NewReader([]byte("package data"))
	err = storage.UploadPackage("project", "name-v1.0.0", pkg)
	c.Assert(err, IsNil)

	release, err := dao.GetRelease("project", "name", "name-v1.0.0")
	c.Assert(err, IsNil)
	uris, err := dao.GetPackageURIs(release)
	c.Assert(err, IsNil)
	c.Assert(uris, HasLen, 1)
	c.Assert(uris[0], Equals, "mem://project/name-v1.0.0.tgz")
}

func (s *appSuite) Test_UploadPackage_fails_on_invalid_release_id(c *C) {
	err := UploadPackage("project", "asdoijasdoijasd", nil)
	c.Assert(err, Not(IsNil))
	c.Assert(err.Error(), Equals, "Invalid release format: asdoijasdoijasd")
}

func (s *appSuite) Test_UploadPackage_fails_on_release_id_that_needs_resolving(c *C) {
	err := UploadPackage("project", "name-latest", nil)
	c.Assert(err, Not(IsNil))
	c.Assert(err.Error(), Equals, "Can't upload package against unresolved version 'project/name-latest'")
}

func (s *appSuite) Test_UploadPackage_fails_if_release_not_found(c *C) {
	dao.TestSetup()
	err := UploadPackage("project", "name-v1.0.0", nil)
	c.Assert(err, Not(IsNil))
	c.Assert(err.Error(), Equals, "Not found")
}

func (s *appSuite) Test_UploadPackage_fails_if_upload_fails(c *C) {
	dao.TestSetup()
	storage := &storageProvider{
		Upload: func(project, releaseId string, pkg io.ReadSeeker) (string, error) {
			return "", errors.New("error uploading")
		},
	}

	_, err := AddRelease("project", `{"name": "name", "version": "1.0.0"}`)
	c.Assert(err, IsNil)
	pkg := bytes.NewReader([]byte("package data"))
	err = storage.UploadPackage("project", "name-v1.0.0", pkg)
	c.Assert(err, Not(IsNil))
	c.Assert(err.Error(), Equals, "error uploading")
}

/*
	DOWNLOAD
*/

func (s *appSuite) Test_DownloadPackage_happy_path(c *C) {
	dao.TestSetup()
	storage := &storageProvider{
		Download: func(project, uri string) (io.Reader, error) {
			return bytes.NewReader([]byte("package data")), nil
		},
	}
	_, err := AddRelease("project", `{"name": "name", "version": "1.0.0"}`)
	c.Assert(err, IsNil)
	release, err := dao.GetRelease("project", "name", "name-v1.0.0")
	c.Assert(err, IsNil)
	err = dao.AddPackageURI(release, "mem://project/name-v1.0.0.tar.gz")
	c.Assert(err, IsNil)
	reader, err := storage.GetDownloadReadSeeker("project", "name-v1.0.0")
	c.Assert(err, IsNil)
	data, err := ioutil.ReadAll(reader)
	c.Assert(err, IsNil)
	c.Assert(string(data), Equals, "package data")
}

func (s *appSuite) Test_DownloadPackage_fails_if_release_not_found(c *C) {
	dao.TestSetup()
	_, err := GetDownloadReadSeeker("project", "name-v1.0.0")
	c.Assert(err, Not(IsNil))
	c.Assert(err.Error(), Equals, "Not found")
}

func (s *appSuite) Test_DownloadPackage_fails_if_uri_not_found(c *C) {
	dao.TestSetup()
	_, err := AddRelease("project", `{"name": "name", "version": "1.0.0"}`)
	c.Assert(err, IsNil)
	_, err = GetDownloadReadSeeker("project", "name-v1.0.0")
	c.Assert(err, Not(IsNil))
	c.Assert(err.Error(), Equals, "Not found")
}

func (s *appSuite) Test_DownloadPackage_fails_if_download_fails(c *C) {
	dao.TestSetup()
	storage := &storageProvider{
		Download: func(project, uri string) (io.Reader, error) {
			return nil, fmt.Errorf("Download error")
		},
	}
	_, err := AddRelease("project", `{"name": "name", "version": "1.0.0"}`)
	c.Assert(err, IsNil)
	release, err := dao.GetRelease("project", "name", "name-v1.0.0")
	c.Assert(err, IsNil)
	err = dao.AddPackageURI(release, "mem://project/name-v1.0.0.tar.gz")
	c.Assert(err, IsNil)
	_, err = storage.GetDownloadReadSeeker("project", "name-v1.0.0")
	c.Assert(err, Not(IsNil))
	c.Assert(err.Error(), Equals, "Download error")
}
