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

package dao

import (
	"fmt"

	"github.com/ankyra/escape-inventory/config"
	"github.com/ankyra/escape-inventory/dao/mem"
	"github.com/ankyra/escape-inventory/dao/postgres"
	"github.com/ankyra/escape-inventory/dao/ql"
	. "github.com/ankyra/escape-inventory/dao/types"
)

var GlobalDAO = mem.NewInMemoryDAO()

func LoadFromConfig(conf *config.Config) error {
	if conf.Database == "" {
		return fmt.Errorf("Missing database configuration variable")
	} else if conf.Database == "memory" {
		GlobalDAO = mem.NewInMemoryDAO()
		return nil
	} else if conf.Database == "ql" {
		dao, err := ql.NewQLDAO(conf.DatabaseSettings.Path)
		if err != nil {
			return err
		}
		GlobalDAO = dao
		return nil
	} else if conf.Database == "postgres" {
		dao, err := postgres.NewPostgresDAO(conf.DatabaseSettings.PostgresUrl)
		if err != nil {
			return err
		}
		GlobalDAO = dao
		return nil
	}

	return fmt.Errorf("Unknown database backend: %s", conf.Database)
}
func TestSetup() {
	GlobalDAO = mem.NewInMemoryDAO()
}

func GetProject(project string) (*Project, error) {
	return GlobalDAO.GetProject(project)
}

func AddProject(project *Project) error {
	return GlobalDAO.AddProject(project)
}

func UpdateProject(project *Project) error {
	return GlobalDAO.UpdateProject(project)
}

func GetProjects() (map[string]*Project, error) {
	return GlobalDAO.GetProjects()
}

func GetProjectsByGroups(readGroups []string) (map[string]*Project, error) {
	return GlobalDAO.GetProjectsByGroups(readGroups)
}

func GetProjectHooks(project *Project) (Hooks, error) {
	return GlobalDAO.GetProjectHooks(project)
}

func SetProjectHooks(project *Project, hooks Hooks) error {
	return GlobalDAO.SetProjectHooks(project, hooks)
}

func GetApplications(project string) (map[string]*Application, error) {
	return GlobalDAO.GetApplications(project)
}

func AddApplication(app *Application) error {
	return GlobalDAO.AddApplication(app)
}

func UpdateApplication(app *Application) error {
	return GlobalDAO.UpdateApplication(app)
}

func GetApplication(project, name string) (*Application, error) {
	return GlobalDAO.GetApplication(project, name)
}

func GetApplicationHooks(app *Application) (Hooks, error) {
	return GlobalDAO.GetApplicationHooks(app)
}

func SetApplicationHooks(app *Application, hooks Hooks) error {
	return GlobalDAO.SetApplicationHooks(app, hooks)
}

func GetDownstreamHooks(app *Application) ([]*Hooks, error) {
	return GlobalDAO.GetDownstreamHooks(app)
}
func SetApplicationSubscribesToUpdatesFrom(app *Application, upstream []*Application) error {
	return GlobalDAO.SetApplicationSubscribesToUpdatesFrom(app, upstream)
}

func AddRelease(release *Release) error {
	return GlobalDAO.AddRelease(release)
}
func UpdateRelease(release *Release) error {
	return GlobalDAO.UpdateRelease(release)
}

func GetRelease(project, name, releaseId string) (*Release, error) {
	return GlobalDAO.GetRelease(project, name, releaseId)
}
func FindAllVersions(app *Application) ([]string, error) {
	return GlobalDAO.FindAllVersions(app)
}

func GetPermittedGroups(project string, perm Permission) ([]string, error) {
	return GlobalDAO.GetPermittedGroups(project, perm)
}

func SetACL(project, group string, perm Permission) error {
	return GlobalDAO.SetACL(project, group, perm)
}

func GetACL(project string) (map[string]Permission, error) {
	return GlobalDAO.GetACL(project)
}

func DeleteACL(project, group string) error {
	return GlobalDAO.DeleteACL(project, group)
}

func GetPackageURIs(r *Release) ([]string, error) {
	return GlobalDAO.GetPackageURIs(r)
}

func AddPackageURI(r *Release, uri string) error {
	return GlobalDAO.AddPackageURI(r, uri)
}

func SetDependencies(r *Release, deps []*Dependency) error {
	return GlobalDAO.SetDependencies(r, deps)
}

func GetDependencies(r *Release) ([]*Dependency, error) {
	return GlobalDAO.GetDependencies(r)
}

func GetDownstreamDependencies(r *Release) ([]*Dependency, error) {
	return GlobalDAO.GetDownstreamDependencies(r)
}

func GetDownstreamDependenciesByGroups(r *Release, readGroups []string) ([]*Dependency, error) {
	return GlobalDAO.GetDownstreamDependenciesByGroups(r, readGroups)
}

func GetAllReleases() ([]*Release, error) {
	return GlobalDAO.GetAllReleases()
}
func GetAllReleasesWithoutProcessedDependencies() ([]*Release, error) {
	return GlobalDAO.GetAllReleasesWithoutProcessedDependencies()
}

func GetUserMetrics(username string) (*Metrics, error) {
	return GlobalDAO.GetUserMetrics(username)
}

func SetUserMetrics(username string, previous, new *Metrics) error {
	return GlobalDAO.SetUserMetrics(username, previous, new)
}

func GetFeedPage(pageSize int) ([]*FeedEvent, error) {
	return GlobalDAO.GetFeedPage(pageSize)
}

func GetProjectFeedPage(project string, pageSize int) ([]*FeedEvent, error) {
	return GlobalDAO.GetProjectFeedPage(project, pageSize)
}

func GetFeedPageByGroups(readGroups []string, pageSize int) ([]*FeedEvent, error) {
	return GlobalDAO.GetFeedPageByGroups(readGroups, pageSize)
}

func AddFeedEvent(event *FeedEvent) error {
	return GlobalDAO.AddFeedEvent(event)
}

func IsNotFound(err error) bool {
	return err == NotFound
}
func IsAlreadyExists(err error) bool {
	return err == AlreadyExists
}
func IsLimitError(err error) bool {
	return err == LimitError
}
func IsUnauthorized(err error) bool {
	return err == Unauthorized
}
