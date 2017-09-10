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
	"time"

	"github.com/ankyra/escape-core"
)

type Permission int

const ReadPermission = Permission(1)
const WritePermission = Permission(2)
const OwnerPermission = Permission(3)
const AdminPermission = Permission(4)

func (p Permission) String() string {
	if p == ReadPermission {
		return "read"
	} else if p == WritePermission {
		return "write"
	} else if p == OwnerPermission {
		return "owner"
	} else if p == AdminPermission {
		return "admin"
	}
	return "???"
}

type Hooks map[string]map[string]string

func NewHooks() Hooks {
	return map[string]map[string]string{}
}

type Project struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	OrgURL      string `json:"org_url"`
	Logo        string `json:"logo"`
	Hooks       Hooks  `json:"hooks"` // only used for view purposes
}

func NewProject(project string) *Project {
	return &Project{
		Name: project,
	}
}

type Application struct {
	Name          string    `json:"name"`
	Project       string    `json:"project"`
	Description   string    `json:"description"`
	LatestVersion string    `json:"latest_version"`
	Logo          string    `json:"logo"`
	UploadedBy    string    `json:"uploaded_by"`
	UploadedAt    time.Time `json:"uploaded_at"`
	Hooks         Hooks     `json:"hooks"` // only used for view purposes
}

func NewApplication(project, name string) *Application {
	return &Application{
		Name:    name,
		Project: project,
	}
}

type Release struct {
	Application           *Application
	ReleaseId             string
	Version               string
	Metadata              *core.ReleaseMetadata
	ProcessedDependencies bool
	Downloads             int
	UploadedBy            string
	UploadedAt            time.Time
}

type Dependency struct {
	Project     string `json:"project"`
	Application string `json:"name"`
	Version     string `json:"version"`
	BuildScope  bool   `json:"build"`
	DeployScope bool   `json:"deploy"`
	IsExtension bool   `json:"is_extension"`
}

func NewRelease(app *Application, metadata *core.ReleaseMetadata) *Release {
	return &Release{
		Application: app,
		ReleaseId:   metadata.GetReleaseId(),
		Version:     metadata.Version,
		Metadata:    metadata,
	}
}

type DAO interface {
	GetProject(project string) (*Project, error)
	AddProject(*Project) error
	UpdateProject(*Project) error
	GetProjects() (map[string]*Project, error)
	GetProjectsByGroups(readGroups []string) (map[string]*Project, error)
	GetProjectHooks(*Project) (Hooks, error)
	SetProjectHooks(*Project, Hooks) error

	GetApplication(project, name string) (*Application, error)
	AddApplication(app *Application) error
	UpdateApplication(app *Application) error
	GetApplications(project string) (map[string]*Application, error)
	FindAllVersions(application *Application) ([]string, error)
	GetApplicationHooks(*Application) (Hooks, error)
	SetApplicationHooks(*Application, Hooks) error

	GetRelease(project, name, releaseId string) (*Release, error)
	AddRelease(*Release) error
	UpdateRelease(*Release) error
	GetAllReleases() ([]*Release, error)
	GetAllReleasesWithoutProcessedDependencies() ([]*Release, error)

	SetDependencies(*Release, []*Dependency) error
	GetDependencies(*Release) ([]*Dependency, error)
	GetDownstreamDependencies(*Release) ([]*Dependency, error)
	GetDownstreamDependenciesByGroups(rel *Release, readGroups []string) ([]*Dependency, error)

	GetPackageURIs(release *Release) ([]string, error)
	AddPackageURI(release *Release, uri string) error

	SetACL(project, group string, perm Permission) error
	GetACL(project string) (map[string]Permission, error)
	DeleteACL(project, group string) error
	GetPermittedGroups(project string, perm Permission) ([]string, error)
}

var NotFound = fmt.Errorf("Not found")
var AlreadyExists = fmt.Errorf("Already exists")
