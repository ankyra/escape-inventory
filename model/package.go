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

package model

import (
	"fmt"
	"io"
	"log"

	"github.com/ankyra/escape-core/parsers"
	"github.com/ankyra/escape-inventory/dao"
	"github.com/ankyra/escape-inventory/dao/types"
	"github.com/ankyra/escape-inventory/storage"
)

type storageProvider struct {
	Upload   func(project, releaseId string, pkg io.ReadSeeker) (string, error)
	Download func(project, uri string) (io.Reader, error)
}

func newStorageProvider() *storageProvider {
	return &storageProvider{
		Upload:   storage.Upload,
		Download: storage.Download,
	}
}

func UploadPackage(project, releaseId string, pkg io.ReadSeeker) error {
	return newStorageProvider().UploadPackage(project, releaseId, pkg)
}

func GetDownloadReadSeeker(project, releaseId string) (io.Reader, error) {
	return newStorageProvider().GetDownloadReadSeeker(project, releaseId)
}

func (s *storageProvider) UploadPackage(project, releaseId string, pkg io.ReadSeeker) error {
	parsed, err := parsers.ParseReleaseId(releaseId)
	if err != nil {
		return NewUserError(err)
	}
	if parsed.NeedsResolving() {
		return NewUserError(fmt.Errorf("Can't upload package against unresolved version '%s/%s'", project, releaseId))
	}
	release, err := dao.GetRelease(project, parsed.Name, releaseId)
	if err != nil {
		return NewUserError(err)
	}
	uri, err := s.Upload(project, releaseId, pkg)
	if err != nil {
		return err
	}
	return dao.AddPackageURI(release, uri)
}

func (s *storageProvider) GetDownloadReadSeeker(project, releaseId string) (io.Reader, error) {
	release, err := ResolveReleaseId(project, releaseId)
	if err != nil {
		return nil, err
	}
	uris, err := dao.GetPackageURIs(release)
	if err != nil {
		return nil, err
	}
	lastError := types.NotFound
	for _, uri := range uris {
		reader, err := s.Download(project, uri)
		if err == nil {
			release.Downloads += 1
			return reader, dao.UpdateRelease(release)
		}
		lastError = err
		log.Printf("Warn: %s\n", err.Error())
	}
	return nil, lastError
}
