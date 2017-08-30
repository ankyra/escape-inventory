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
	"strings"

	"github.com/ankyra/escape-core"
	"github.com/ankyra/escape-core/parsers"
	"github.com/ankyra/escape-registry/dao"
	. "github.com/ankyra/escape-registry/dao/types"
)

func ensureProjectExists(project string) error {
	prj, err := dao.GetProject(project)
	if err == nil {
		return nil
	}
	if err != NotFound {
		return err
	}
	prj = NewProject(project)
	return dao.AddProject(prj)
}

func ensureApplicationExists(project string, metadata *core.ReleaseMetadata) error {
	name := metadata.Name
	app, err := dao.GetApplication(project, name)
	if err == nil {
		app.Description = metadata.Description
		app.Logo = metadata.Logo
		app.LatestVersion = metadata.Version
		return dao.UpdateApplication(app)
	}
	if err != NotFound {
		return err
	}
	app = NewApplication(project, name)
	app.Description = metadata.Description
	app.Logo = metadata.Logo
	app.LatestVersion = metadata.Version
	return dao.AddApplication(app)
}

func AddRelease(project, metadataJson string) (*core.ReleaseMetadata, error) {
	metadata, err := core.NewReleaseMetadataFromJsonString(metadataJson)
	if err != nil {
		return nil, NewUserError(err)
	}
	releaseId := metadata.GetReleaseId()
	parsed, err := parsers.ParseReleaseId(releaseId)
	if err != nil {
		return nil, NewUserError(err)
	}
	if parsed.NeedsResolving() {
		return nil, NewUserError(fmt.Errorf("Can't add release with unresolved version"))
	}
	if metadata.ApiVersion > core.CurrentApiVersion {
		return nil, NewUserError(fmt.Errorf("Release format version v%s is not supported (this registry supports up to v%s)", metadata.ApiVersion, core.CurrentApiVersion))
	}
	release, err := dao.GetRelease(project, parsed.Name, releaseId)
	if err != nil && !dao.IsNotFound(err) {
		return nil, err
	}
	if release != nil {
		return nil, NewUserError(fmt.Errorf("Release %s already exists", releaseId))
	}
	if err := ensureProjectExists(project); err != nil {
		return nil, err
	}
	if err := ensureApplicationExists(project, metadata); err != nil {
		return nil, err
	}
	result, err := dao.AddRelease(project, metadata)
	if err != nil {
		return nil, err
	}
	return result.Metadata, ProcessDependencies(result)
}

func ProcessDependencies(release *Release) error {
	deps := []*Dependency{}
	for _, dep := range release.Metadata.Depends {
		parsed, err := parsers.ParseQualifiedReleaseId(dep.ReleaseId)
		if err != nil {
			return fmt.Errorf("Couldn't parse dependency: %s", err.Error())
		}
		d := Dependency{
			Project:     parsed.Project,
			Application: parsed.Name,
			Version:     parsed.Version,
			BuildScope:  dep.InScope("build"),
			DeployScope: dep.InScope("deploy"),
			IsExtension: false,
		}
		deps = append(deps, &d)
	}
	for _, ext := range release.Metadata.Extends {
		parsed, err := parsers.ParseQualifiedReleaseId(ext.ReleaseId)
		if err != nil {
			return fmt.Errorf("Couldn't parse dependency: %s", err.Error())
		}
		d := Dependency{
			Project:     parsed.Project,
			Application: parsed.Name,
			Version:     parsed.Version,
			BuildScope:  true,
			DeployScope: true,
			IsExtension: true,
		}
		deps = append(deps, &d)
	}
	if err := dao.SetDependencies(release, deps); err != nil {
		return err
	}
	release.ProcessedDependencies = true
	return dao.UpdateRelease(release)
}

func ProcessUnprocessedReleases() error {
	releases, err := dao.GetAllReleasesWithoutProcessedDependencies()
	if err != nil {
		return err
	}
	for _, release := range releases {
		if err := ProcessDependencies(release); err != nil {
			return err
		}
	}
	return nil
}

type ReleasePayload struct {
	Release   *core.ReleaseMetadata `json:"release"`
	Downloads int                   `json:"downloads"`
}

func GetRelease(project, releaseIdString string) (*ReleasePayload, error) {
	release, err := ResolveReleaseId(project, releaseIdString)
	if err != nil {
		return nil, err
	}
	return &ReleasePayload{
		Release:   release.Metadata,
		Downloads: release.Downloads,
	}, nil
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
