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

import "time"

type ApplicationsDAO interface {
	GetApplication(namespace, name string) (*Application, error)
	AddApplication(app *Application) error
	UpdateApplication(app *Application) error
	GetApplications(namespace string) (map[string]*Application, error)
	FindAllVersions(application *Application) ([]string, error)
	GetApplicationHooks(*Application) (Hooks, error)
	SetApplicationHooks(*Application, Hooks) error
	GetDownstreamHooks(*Application) ([]*Hooks, error)
	SetApplicationSubscribesToUpdatesFrom(*Application, []*Application) error
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
