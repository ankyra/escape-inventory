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
	"io"
	"log"

	"github.com/ankyra/escape-core/parsers"
	"github.com/ankyra/escape-inventory/dao"
	"github.com/ankyra/escape-inventory/dao/types"
	"github.com/ankyra/escape-inventory/storage"
)

type storageProvider struct {
	Upload   func(namespace, releaseId string, pkg io.ReadSeeker) (string, error)
	Download func(namespace, uri string) (io.Reader, error)
}

func newStorageProvider() *storageProvider {
	return &storageProvider{
		Upload:   storage.Upload,
		Download: storage.Download,
	}
}

func UploadPackage(namespace, releaseId string, pkg io.ReadSeeker) error {
	return newStorageProvider().UploadPackage(namespace, releaseId, pkg)
}

func GetDownloadReadSeeker(namespace, releaseId string) (io.Reader, error) {
	return newStorageProvider().GetDownloadReadSeeker(namespace, releaseId)
}

func (s *storageProvider) UploadPackage(namespace, releaseId string, pkg io.ReadSeeker) error {
	parsed, err := parsers.ParseReleaseId(releaseId)
	if err != nil {
		return NewUserError(err)
	}
	if parsed.NeedsResolving() {
		return NewUserError(fmt.Errorf("Can't upload package against unresolved version '%s/%s'", namespace, releaseId))
	}
	release, err := dao.GetRelease(namespace, parsed.Name, releaseId)
	if err != nil {
		return NewUserError(err)
	}
	uri, err := s.Upload(namespace, releaseId, pkg)
	if err != nil {
		return err
	}
	return dao.AddPackageURI(release, uri)
}

func (s *storageProvider) GetDownloadReadSeeker(namespace, releaseId string) (io.Reader, error) {
	release, err := ResolveReleaseId(namespace, releaseId)
	if err != nil {
		return nil, err
	}
	uris, err := dao.GetPackageURIs(release)
	if err != nil {
		return nil, err
	}
	lastError := types.NotFound
	for _, uri := range uris {
		reader, err := s.Download(namespace, uri)
		if err == nil {
			release.Downloads += 1
			return reader, dao.UpdateRelease(release)
		}
		lastError = err
		log.Printf("Warn: %s\n", err.Error())
	}
	return nil, lastError
}
