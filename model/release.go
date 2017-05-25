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
	"github.com/ankyra/escape-core"
	"github.com/ankyra/escape-core/parsers"
	"github.com/ankyra/escape-registry/dao"
	. "github.com/ankyra/escape-registry/dao/types"
	"strings"
)

func AddRelease(project, metadataJson string) error {
	metadata, err := core.NewReleaseMetadataFromJsonString(metadataJson)
	if err != nil {
		return NewUserError(err)
	}
	releaseId := metadata.GetReleaseId()
	parsed, err := parsers.ParseReleaseId(releaseId)
	if err != nil {
		return NewUserError(err)
	}
	if parsed.NeedsResolving() {
		return NewUserError(fmt.Errorf("Can't add release with unresolved version"))
	}
	release, err := dao.GetRelease(project, parsed.Name, releaseId)
	if err != nil && !dao.IsNotFound(err) {
		return err
	}
	if release != nil {
		return NewUserError(fmt.Errorf("Release %s already exists", releaseId))
	}
	_, err = dao.AddRelease(project, metadata)
	return err
}

func GetReleaseMetadata(project, releaseIdString string) (*core.ReleaseMetadata, error) {
	release, err := ResolveReleaseId(project, releaseIdString)
	if err != nil {
		return nil, err
	}
	return release.Metadata, nil
}

func ResolveReleaseId(project, releaseIdString string) (*Release, error) {
	releaseId, err := parsers.ParseReleaseId(releaseIdString)
	if err != nil {
		return nil, NewUserError(err)
	}
	if releaseId.Version == "latest" {
		version, err := getLastVersionForPrefix(project, releaseId.Name, "")
		if err != nil {
			return nil, NewUserError(err)
		}
		releaseId.Version = version.ToString()
	} else if strings.HasSuffix(releaseId.Version, ".@") {
		prefix := releaseId.Version[:len(releaseId.Version)-1]
		version, err := getLastVersionForPrefix(project, releaseId.Name, prefix)
		if err != nil {
			return nil, NewUserError(err)
		}
		releaseId.Version = prefix + version.ToString()
	}
	return dao.GetRelease(project, releaseId.Name, releaseId.ToString())
}
