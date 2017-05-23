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

package types

import (
	"fmt"
	"github.com/ankyra/escape-core"
)

type DAO interface {
	GetApplications(project string) ([]ApplicationDAO, error)
	GetApplication(project, name string) (ApplicationDAO, error)
	GetRelease(project, releaseId string) (ReleaseDAO, error)
	AddRelease(project string, metadata *core.ReleaseMetadata) (ReleaseDAO, error)
	GetAllReleases() ([]ReleaseDAO, error)
}

type ApplicationDAO interface {
	GetName() string
	FindAllVersions() ([]string, error)
}

type ReleaseDAO interface {
	GetApplication() ApplicationDAO
	GetVersion() string
	GetMetadata() *core.ReleaseMetadata

	GetPackageURIs() ([]string, error)
	AddPackageURI(uri string) error
}

var NotFound = fmt.Errorf("Not found")
var AlreadyExists = fmt.Errorf("Already exists")
