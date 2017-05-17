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
	"github.com/ankyra/escape-core/parsers"
	"github.com/ankyra/escape-registry/dao"
	"github.com/ankyra/escape-registry/dao/types"
	"github.com/ankyra/escape-registry/storage"
	"io"
	"log"
)

func UploadPackage(releaseId string, pkg io.ReadSeeker) error {
	parsed, err := parsers.ParseReleaseId(releaseId)
	if err != nil {
		return NewUserError(err)
	}
	if parsed.NeedsResolving() {
		return NewUserError(fmt.Errorf("Can't upload package against unresolved version"))
	}
	release, err := dao.GetRelease(releaseId)
	if err != nil {
		return NewUserError(err)
	}
	uri, err := storage.Upload(releaseId, pkg)
	if err != nil {
		return err
	}
	return release.AddPackageURI(uri)
}

func GetDownloadReadSeeker(releaseId string) (io.Reader, error) {
	release, err := ResolveReleaseId(releaseId)
	if err != nil {
		return nil, err
	}
	uris, err := release.GetPackageURIs()
	if err != nil {
		return nil, err
	}
	for _, uri := range uris {
		reader, err := storage.Download(uri)
		if err == nil {
			return reader, nil
		}
		log.Printf("Warn: %s\n", err.Error())
	}
	return nil, types.NotFound
}
