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
	"strings"
	"time"

	"github.com/ankyra/escape-core"
	"github.com/ankyra/escape-core/parsers"
	"github.com/ankyra/escape-inventory/dao"
	. "github.com/ankyra/escape-inventory/dao/types"
)

func ensureProjectExists(project, username string) error {
	prj, err := dao.GetProject(project)
	if err == nil {
		return nil
	}
	if err != NotFound {
		return err
	}
	if err := core.ValidateProjectName(project); err != nil {
		return NewUserError(err)
	}
	prj = NewProject(project)
	if err := dao.AddProject(prj); err != nil {
		return err
	}
	return AddCreateProjectFeedEvent(project, username)
}

func updateApp(app *Application, metadata *core.ReleaseMetadata, byUser string, uploadedAt time.Time) {
	app.Description = metadata.Description
	app.Logo = metadata.Logo
	app.LatestVersion = metadata.Version
	if byUser != "" {
		app.UploadedBy = byUser
	}
	if uploadedAt != time.Unix(0, 0) {
		app.UploadedAt = uploadedAt
	}
}

func ensureApplicationExists(project, byUser string, metadata *core.ReleaseMetadata, uploadAt time.Time) error {
	name := metadata.Name
	app, err := dao.GetApplication(project, name)
	if err != nil && err != NotFound {
		return err
	} else if err == nil {
		updateApp(app, metadata, byUser, uploadAt)
		return dao.UpdateApplication(app)
	}
	app = NewApplication(project, name)
	updateApp(app, metadata, byUser, uploadAt)
	if err := dao.AddApplication(app); err != nil {
		return err
	}
	return AddNewApplicationFeedEvent(project, name, byUser)
}

func AddRelease(project, metadataJson string) (*core.ReleaseMetadata, error) {
	return AddReleaseByUser(project, metadataJson, "")
}

func AddReleaseByUser(project, metadataJson, uploadUser string) (*core.ReleaseMetadata, error) {
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
		return nil, NewUserError(fmt.Errorf("Release format version v%s is not supported (this Inventory supports up to v%s)", metadata.ApiVersion, core.CurrentApiVersion))
	}
	release, err := dao.GetRelease(project, parsed.Name, releaseId)
	if err != nil && !dao.IsNotFound(err) {
		return nil, err
	}
	if release != nil {
		return nil, NewUserError(fmt.Errorf("Release %s already exists", releaseId))
	}
	if err := ensureProjectExists(project, uploadUser); err != nil {
		return nil, err
	}
	result := NewRelease(NewApplication(project, metadata.Name), metadata)
	result.UploadedBy = uploadUser
	result.UploadedAt = time.Now()
	if err := ensureApplicationExists(project, result.UploadedBy, metadata, result.UploadedAt); err != nil {
		return nil, err
	}
	if err := dao.AddRelease(result); err != nil {
		return nil, err
	}
	if err := dao.RegisterProviders(metadata); err != nil {
		return nil, err
	}
	err = AddNewReleaseFeedEvent(project, metadata.Name, metadata.Version, uploadUser)
	if err != nil {
		return nil, err
	}
	return result.Metadata, ProcessDependencies(result)
}

func ProcessDependencies(release *Release) error {
	deps := []*Dependency{}
	apps := []*Application{}
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
		apps = append(apps, NewApplication(parsed.Project, parsed.Name))
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
		apps = append(apps, NewApplication(parsed.Project, parsed.Name))
	}
	if err := dao.SetDependencies(release, deps); err != nil {
		return err
	}
	if err := dao.SetApplicationSubscribesToUpdatesFrom(release.Application, apps); err != nil {
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
	Release    *core.ReleaseMetadata `json:"release"`
	Versions   []string              `json:"versions"`
	IsLatest   bool                  `json:"is_latest"`
	Downloads  int                   `json:"downloads"`
	UploadedBy string                `json:"uploaded_by"`
	UploadedAt time.Time             `json:"uploaded_at"`
}

func GetRelease(project, name, version string) (*ReleasePayload, error) {
	releaseIdString := name + "-" + version
	release, err := ResolveReleaseId(project, releaseIdString)
	if err != nil {
		return nil, err
	}
	versions, err := GetApplicationVersions(project, release.Application.Name)
	if err != nil {
		return nil, err
	}
	maxVersion := getMaxFromVersions(versions, "")
	isLatest := maxVersion.ToString() == release.Version
	return &ReleasePayload{
		Release:    release.Metadata,
		Versions:   versions,
		IsLatest:   isLatest,
		Downloads:  release.Downloads,
		UploadedBy: release.UploadedBy,
		UploadedAt: release.UploadedAt,
	}, nil
}

func GetReleaseMetadata(project, name, version string) (*core.ReleaseMetadata, error) {
	releaseIdString := name + "-v" + version
	if strings.HasPrefix(version, "v") || version == "latest" {
		releaseIdString = name + "-" + version
	}
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
