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
	"time"

	"github.com/ankyra/escape-core"
	"github.com/ankyra/escape-core/parsers"
	"github.com/ankyra/escape-inventory/dao"
	. "github.com/ankyra/escape-inventory/dao/types"
	"github.com/ankyra/escape-middleware/errors"
)

func ensureNamespaceExists(namespace, username string) error {
	prj, err := dao.GetNamespace(namespace)
	if err == nil {
		return nil
	}
	if err != NotFound {
		return err
	}
	if err := core.ValidateProjectName(namespace); err != nil {
		return NewUserError(err)
	}
	prj = NewProject(namespace)
	return dao.AddNamespace(prj)
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

func ensureApplicationExists(namespace, byUser string, metadata *core.ReleaseMetadata, uploadAt time.Time) error {
	name := metadata.Name
	app, err := dao.GetApplication(namespace, name)
	if err != nil && err != NotFound {
		return err
	} else if err == nil {
		updateApp(app, metadata, byUser, uploadAt)
		return dao.UpdateApplication(app)
	}
	app = NewApplication(namespace, name)
	updateApp(app, metadata, byUser, uploadAt)
	return dao.AddApplication(app)
}

func AddRelease(namespace, metadataJson string) (*core.ReleaseMetadata, error) {
	return AddReleaseByUser(namespace, metadataJson, "")
}

func AddReleaseByUser(namespace, metadataJson, uploadUser string) (*core.ReleaseMetadata, error) {
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
		return nil, NewUserError(fmt.Errorf("Release format version v%d is not supported (this Inventory supports up to v%d)", metadata.ApiVersion, core.CurrentApiVersion))
	}
	release, err := dao.GetRelease(namespace, parsed.Name, releaseId)
	if err != nil && !dao.IsNotFound(err) {
		return nil, err
	}
	if release != nil {
		return nil, NewUserError(fmt.Errorf("Release %s already exists", releaseId))
	}
	if err := ensureNamespaceExists(namespace, uploadUser); err != nil {
		return nil, err
	}
	result := NewRelease(NewApplication(namespace, metadata.Name), metadata)
	result.UploadedBy = uploadUser
	result.UploadedAt = time.Now()
	if err := ensureApplicationExists(namespace, result.UploadedBy, metadata, result.UploadedAt); err != nil {
		return nil, err
	}
	if err := dao.AddRelease(result); err != nil {
		return nil, err
	}
	if err := dao.RegisterProviders(metadata); err != nil {
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

func GetRelease(namespace, name, version string) (*ReleasePayload, error) {
	release, err := ResolveReleaseId(namespace, name, version)
	if err != nil {
		return nil, err
	}
	versions, err := GetApplicationVersions(namespace, release.Application.Name)
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

func GetReleaseMetadata(namespace, name, version string) (*core.ReleaseMetadata, error) {
	release, err := ResolveReleaseId(namespace, name, version)
	if err != nil {
		return nil, err
	}
	return release.Metadata, nil
}

func ResolveReleaseId(namespace, application, versionQuery string) (*Release, error) {
	vq, err := parsers.ParseVersionQuery(versionQuery)
	if err != nil {
		return nil, NewUserError(err)
	}
	if vq.LatestVersion {
		version, err := getLastVersionForPrefix(namespace, application, "")
		if err != nil {
			return nil, NewUserError(err)
		}
		versionQuery = version.ToString()
	} else if vq.VersionPrefix != "" {
		version, err := getLastVersionForPrefix(namespace, application, vq.VersionPrefix)
		if err != nil {
			return nil, NewUserError(err)
		}
		versionQuery = vq.VersionPrefix + version.ToString()
	} else if vq.SpecificVersion != "" {
		versionQuery = vq.SpecificVersion
	} else if vq.SpecificTag != "" {
		return dao.GetReleaseByTag(namespace, application, vq.SpecificTag)
	} else {
		return nil, errors.NewUserError("Unsupported version query")
	}
	return dao.GetRelease(namespace, application, application+"-v"+versionQuery)
}

func TagRelease(namespace, application, releaseId, tag string) error {
	parsed, err := parsers.ParseQualifiedReleaseId(releaseId)
	if err != nil {
		return NewUserError(err)
	}
	if !parsers.IsValidTag(tag) {
		return NewUserError(fmt.Errorf("The tag '%s' is invalid/not supported.", tag))
	}
	if parsed.Project != namespace {
		return NewUserError(fmt.Errorf("Mismatch between namespace from URL and namespace from release_id."))
	}
	if parsed.Name != application {
		return NewUserError(fmt.Errorf("Mismatch between application from URL and application from release_id."))
	}
	release, err := ResolveReleaseId(namespace, application, parsed.Version)
	if err != nil {
		return err
	}
	return dao.TagRelease(release, tag)
}
