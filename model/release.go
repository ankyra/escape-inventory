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
	"github.com/ankyra/escape-registry/dao"
	. "github.com/ankyra/escape-registry/dao/types"
	"github.com/ankyra/escape-registry/shared"
	"strings"
)

func AddRelease(metadataJson string) error {
	metadata, err := shared.NewReleaseMetadataFromJsonString(metadataJson)
	if err != nil {
		return NewUserError(err)
	}
	releaseId := metadata.GetReleaseId()
	parsed, err := shared.ParseReleaseId(releaseId)
	if err != nil {
		return NewUserError(err)
	}
	if parsed.NeedsResolving() {
		return NewUserError(fmt.Errorf("Can't add release with unresolved version"))
	}
	release, err := dao.GetRelease(releaseId)
	if err != nil && !dao.IsNotFound(err) {
		return err
	}
	if release != nil {
		return NewUserError(fmt.Errorf("Release %s already exists", releaseId))
	}
	_, err = dao.AddRelease(metadata)
	return err
}

func GetReleaseMetadata(releaseIdString string) (Metadata, error) {
	release, err := ResolveReleaseId(releaseIdString)
	if err != nil {
		return nil, err
	}
	return release.GetMetadata(), nil
}

func ResolveReleaseId(releaseIdString string) (ReleaseDAO, error) {
	releaseId, err := shared.ParseReleaseId(releaseIdString)
	if err != nil {
		return nil, NewUserError(err)
	}
	if releaseId.Version == "latest" {
		version, err := getLastVersionForPrefix(releaseIdString, "")
		if err != nil {
			return nil, NewUserError(err)
		}
		releaseId.Version = version.ToString()
	} else if strings.HasSuffix(releaseId.Version, ".@") {
		prefix := releaseId.Version[:len(releaseId.Version)-1]
		version, err := getLastVersionForPrefix(releaseIdString, prefix)
		if err != nil {
			return nil, NewUserError(err)
		}
		releaseId.Version = prefix + version.ToString()
	}
	return dao.GetRelease(releaseId.ToString())
}
