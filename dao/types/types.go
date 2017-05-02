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
	"github.com/ankyra/escape-registry/shared"
)

type Metadata shared.ReleaseMetadata

type DAO interface {
	GetApplication(typ, name string) (ApplicationDAO, error)
	GetApplications() ([]ApplicationDAO, error)
	GetRelease(releaseId string) (ReleaseDAO, error)
	GetAllReleases() ([]ReleaseDAO, error)
	AddRelease(metadata Metadata) (ReleaseDAO, error)
	GetApplicationsByType(typ string) ([]string, error)
	GetReleaseTypes() ([]string, error)
}

type ApplicationDAO interface {
	GetType() string
	GetName() string
	FindAllVersions() ([]string, error)
}

type ReleaseDAO interface {
	GetApplication() ApplicationDAO
	GetVersion() string
	GetMetadata() Metadata

	GetPackageURIs() ([]string, error)
	AddPackageURI(uri string) error
}

var NotFound = fmt.Errorf("Not found")
var AlreadyExists = fmt.Errorf("Already exists")
