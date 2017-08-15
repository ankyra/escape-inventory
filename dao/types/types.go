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

type Project struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	OrgURL      string `json:"org_url"`
	Logo        string `json:"logo"`
}

func NewProject(project string) *Project {
	return &Project{
		Name: project,
	}
}

type Application struct {
	Name            string
	Project         string
	Description     string
	LatestReleaseId string
	Logo            string
}

func NewApplication(project, name string) *Application {
	return &Application{
		Name:    name,
		Project: project,
	}
}

type Release struct {
	Application *Application
	ReleaseId   string
	Version     string
	Metadata    *core.ReleaseMetadata
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

	AddApplication(app *Application) error
	UpdateApplication(app *Application) error
	GetApplications(project string) ([]*Application, error)
	GetApplication(project, name string) (*Application, error)
	FindAllVersions(application *Application) ([]string, error)

	GetRelease(project, name, releaseId string) (*Release, error)
	AddRelease(project string, metadata *core.ReleaseMetadata) (*Release, error)
	GetAllReleases() ([]*Release, error)
	GetPackageURIs(release *Release) ([]string, error)
	AddPackageURI(release *Release, uri string) error

	SetACL(project, group string, perm Permission) error
	GetACL(project string) (map[string]Permission, error)
	DeleteACL(project, group string) error
	GetPermittedGroups(project string, perm Permission) ([]string, error)
}

var NotFound = fmt.Errorf("Not found")
var AlreadyExists = fmt.Errorf("Already exists")
