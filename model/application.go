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
	"github.com/ankyra/escape-registry/dao"
	"github.com/ankyra/escape-registry/shared"
	"strings"
)

func GetNextVersion(releaseIdString, prefix string) (string, error) {
	latest, err := getLastVersionForPrefix(releaseIdString, prefix)
	if err != nil {
		if dao.IsNotFound(err) {
			return prefix + "0", nil
		}
		return "", NewUserError(err)
	}
	if err := latest.IncrementSmallest(); err != nil {
		return "", NewUserError(err)
	}
	return prefix + latest.ToString(), nil

}

func getLastVersionForPrefix(releaseIdString, prefix string) (*SemanticVersion, error) {
	releaseId, err := shared.ParseReleaseId(releaseIdString)
	if err != nil {
		return nil, NewUserError(err)
	}
	app, err := dao.GetApplication(releaseId.Type, releaseId.Name)
	if err != nil {
		return nil, NewUserError(err)
	}
	versions, err := app.FindAllVersions()
	if err != nil {
		return nil, err
	}
	return getMaxFromVersions(versions, prefix), nil
}

func getMaxFromVersions(versions []string, prefix string) *SemanticVersion {
	current := NewSemanticVersion("-1")
	for _, v := range versions {
		if strings.HasPrefix(v, prefix) {
			release_version := v[len(prefix):]
			newver := NewSemanticVersion(release_version)
			if current.LessOrEqual(newver) {
				current = newver
			}
		}
	}
	return current
}
