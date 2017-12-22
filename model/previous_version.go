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
	core "github.com/ankyra/escape-core"
	"github.com/ankyra/escape-inventory/dao"
	"github.com/ankyra/escape-inventory/dao/types"
)

func GetPreviousReleaseMetadata(project, name, version string) (*core.ReleaseMetadata, error) {
	metadata, err := GetReleaseMetadata(project, name, version)
	if err != nil {
		return nil, err
	}
	prev, err := GetPreviousVersion(project, name, metadata.Version)
	if err != nil {
		return nil, err
	}
	return GetReleaseMetadata(project, name, prev)
}

func GetPreviousVersion(project, app, version string) (string, error) {
	prev, err := getPreviousVersion(project, app, version)
	if err != nil {
		return "", NewUserError(err)
	}
	if prev.Equals(NewSemanticVersion("-1")) {
		return "", types.NotFound
	}
	return prev.ToString(), nil

}

func getPreviousVersion(project, appName, version string) (*SemanticVersion, error) {
	app, err := dao.GetApplication(project, appName)
	if err != nil {
		return nil, NewUserError(err)
	}
	versions, err := dao.FindAllVersions(app)
	if err != nil {
		return nil, err
	}
	return getPrevVersion(versions, version), nil
}

func getPrevVersion(versions []string, version string) *SemanticVersion {
	prev := NewSemanticVersion(version)
	current := NewSemanticVersion("-1")
	for _, release_version := range versions {
		newver := NewSemanticVersion(release_version)
		if newver.LessOrEqual(prev) && !newver.Equals(prev) && current.LessOrEqual(newver) {
			current = newver
		}
	}
	return current
}
