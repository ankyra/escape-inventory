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
	"reflect"
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
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	OrgURL         string   `json:"org_url"`
	Logo           string   `json:"logo"`
	Hooks          Hooks    `json:"hooks,omitempty"` // only used for view purposes
	Permission     string   `json:"permission"`      // only used for view purposes
	MatchingGroups []string `json:"-"`               // used to work out highest permission in model
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

func NewDependency(project, name, version string) *Dependency {
	return &Dependency{
		Project:     project,
		Application: name,
		Version:     version,
	}
}

func NewRelease(app *Application, metadata *core.ReleaseMetadata) *Release {
	return &Release{
		Application: app,
		ReleaseId:   metadata.GetReleaseId(),
		Version:     metadata.Version,
		Metadata:    metadata,
	}
}

type Metrics struct {
	ProjectCount int
}

func NewMetrics(projectCount int) *Metrics {
	return &Metrics{
		ProjectCount: projectCount,
	}
}

const FeedPageSize = 7

type FeedEvent struct {
	ID        string                 `json:"id"` // set by DAO
	Type      string                 `json:"type"`
	Username  string                 `json:"user"`
	Project   string                 `json:"project"`
	Timestamp time.Time              `json:"time"`
	Data      map[string]interface{} `json:"data"`
}

func NewEvent(typ, project, username string) *FeedEvent {
	return &FeedEvent{
		Type:      typ,
		Project:   project,
		Username:  username,
		Timestamp: time.Now(),
		Data:      map[string]interface{}{},
	}
}
func NewEventWithData(typ, project, username string, data map[string]interface{}) *FeedEvent {
	ev := NewEvent(typ, project, username)
	ev.Data = data
	return ev
}

func NewCreateProjectEvent(project, username string) *FeedEvent {
	return NewEvent("CREATE_PROJECT", project, username)
}

func NewReleaseEvent(project, name, version, uploadedBy string) *FeedEvent {
	data := map[string]interface{}{
		"name":        name,
		"version":     version,
		"uploaded_by": uploadedBy,
	}
	return NewEventWithData("NEW_RELEASE", project, uploadedBy, data)
}

func (f *FeedEvent) Equals(other *FeedEvent) bool {
	fTimestamp := f.Timestamp.Truncate(time.Second)
	otherTimestamp := other.Timestamp.Truncate(time.Second)
	return f.Type == other.Type &&
		f.Username == other.Username &&
		f.Project == other.Project &&
		fTimestamp == otherTimestamp &&
		reflect.DeepEqual(f.Data, other.Data)
}

type DAO interface {
	GetProject(project string) (*Project, error)
	AddProject(*Project) error
	HardDeleteProject(project string) error
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
	GetDownstreamHooks(*Application) ([]*Hooks, error)
	SetApplicationSubscribesToUpdatesFrom(*Application, []*Application) error

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

	GetUserMetrics(username string) (*Metrics, error)
	SetUserMetrics(username string, previous, new *Metrics) error

	GetFeedPage(pageSize int) ([]*FeedEvent, error)
	GetProjectFeedPage(project string, pageSize int) ([]*FeedEvent, error)
	GetFeedPageByGroups(readGroups []string, pageSize int) ([]*FeedEvent, error)
	AddFeedEvent(event *FeedEvent) error

	WipeDatabase() error
}

var NotFound = fmt.Errorf("Not found")
var AlreadyExists = fmt.Errorf("Already exists")
var LimitError = fmt.Errorf("Plan limit exceeded")
var Unauthorized = fmt.Errorf("Unauthorized")
