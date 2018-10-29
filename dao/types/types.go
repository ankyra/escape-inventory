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

// Mainly used for view purposes
type MinimalReleaseMetadata struct {
	Application string `json:"application"`
	Project     string `json:"project"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

func (m *MinimalReleaseMetadata) GetReleaseId() string {
	return m.Project + "/" + m.Application + "-v" + m.Version
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
	ID          string                 `json:"id"` // set by DAO
	Type        string                 `json:"type"`
	Username    string                 `json:"user"`
	Project     string                 `json:"project"`
	Application string                 `json:"application"`
	Timestamp   time.Time              `json:"time"`
	Data        map[string]interface{} `json:"data"`
}

func NewEvent(typ, namespace, username string) *FeedEvent {
	return &FeedEvent{
		Type:        typ,
		Project:     namespace,
		Application: "",
		Username:    username,
		Timestamp:   time.Now(),
		Data:        map[string]interface{}{},
	}
}
func NewEventWithData(typ, namespace, username string, data map[string]interface{}) *FeedEvent {
	ev := NewEvent(typ, namespace, username)
	ev.Data = data
	return ev
}

func NewCreateProjectEvent(namespace, username string) *FeedEvent {
	return NewEvent("CREATE_PROJECT", namespace, username)
}

func NewUserAddedToProjectEvent(namespace, username, addedByUser string) *FeedEvent {
	data := map[string]interface{}{
		"added_username": username,
	}
	return NewEventWithData("USER_ADDED_TO_PROJECT", namespace, addedByUser, data)
}

func NewUserRemovedFromProjectEvent(namespace, username, removedByUser string) *FeedEvent {
	data := map[string]interface{}{
		"removed_username": username,
	}
	return NewEventWithData("USER_REMOVED_FROM_PROJECT", namespace, removedByUser, data)
}

func NewCreateApplicationEvent(namespace, application, username string) *FeedEvent {
	ev := NewEvent("CREATE_APPLICATION", namespace, username)
	ev.Application = application
	return ev
}

func NewReleaseEvent(namespace, name, version, uploadedBy string) *FeedEvent {
	data := map[string]interface{}{
		"name":        name,
		"version":     version,
		"uploaded_by": uploadedBy,
	}
	ev := NewEventWithData("NEW_RELEASE", namespace, uploadedBy, data)
	ev.Application = name
	return ev
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
	GetNamespace(namespace string) (*Project, error)
	AddNamespace(*Project) error
	HardDeleteNamespace(namespace string) error
	UpdateNamespace(*Project) error
	GetNamespaces() (map[string]*Project, error)
	GetNamespacesByGroups(readGroups []string) (map[string]*Project, error)
	GetNamespaceHooks(*Project) (Hooks, error)
	SetNamespaceHooks(*Project, Hooks) error

	GetApplication(namespace, name string) (*Application, error)
	AddApplication(app *Application) error
	UpdateApplication(app *Application) error
	GetApplications(namespace string) (map[string]*Application, error)
	FindAllVersions(application *Application) ([]string, error)
	GetApplicationHooks(*Application) (Hooks, error)
	SetApplicationHooks(*Application, Hooks) error
	GetDownstreamHooks(*Application) ([]*Hooks, error)
	SetApplicationSubscribesToUpdatesFrom(*Application, []*Application) error

	GetRelease(namespace, name, releaseId string) (*Release, error)
	AddRelease(*Release) error
	UpdateRelease(*Release) error
	GetAllReleases() ([]*Release, error)
	GetAllReleasesWithoutProcessedDependencies() ([]*Release, error)

	SetDependencies(*Release, []*Dependency) error
	GetDependencies(*Release) ([]*Dependency, error)
	GetDownstreamDependencies(*Release) ([]*Dependency, error)
	GetDownstreamDependenciesByGroups(rel *Release, readGroups []string) ([]*Dependency, error)

	GetProviders(providerName string) (map[string]*MinimalReleaseMetadata, error)
	GetProvidersByGroups(providerName string, groups []string) (map[string]*MinimalReleaseMetadata, error)
	RegisterProviders(release *core.ReleaseMetadata) error

	GetPackageURIs(release *Release) ([]string, error)
	AddPackageURI(release *Release, uri string) error

	SetACL(namespace, group string, perm Permission) error
	GetACL(namespace string) (map[string]Permission, error)
	DeleteACL(namespace, group string) error
	GetPermittedGroups(namespace string, perm Permission) ([]string, error)

	GetUserMetrics(username string) (*Metrics, error)
	SetUserMetrics(username string, previous, new *Metrics) error

	WipeDatabase() error
}

var NotFound = fmt.Errorf("Not found")
var AlreadyExists = fmt.Errorf("Already exists")
var LimitError = fmt.Errorf("Plan limit exceeded")
var Unauthorized = fmt.Errorf("Unauthorized")
