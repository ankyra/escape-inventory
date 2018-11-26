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

package dao

import (
	"fmt"

	core "github.com/ankyra/escape-core"
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

func GetNamespace(namespace string) (*Project, error) {
	return GlobalDAO.GetNamespace(namespace)
}

func AddNamespace(namespace *Project) error {
	return GlobalDAO.AddNamespace(namespace)
}

func UpdateNamespace(namespace *Project) error {
	return GlobalDAO.UpdateNamespace(namespace)
}

func GetNamespaces() (map[string]*Project, error) {
	return GlobalDAO.GetNamespaces()
}

func GetNamespacesByNames(namespaces []string) (map[string]*Project, error) {
	return GlobalDAO.GetNamespacesByNames(namespaces)
}

func GetNamespacesForUser(namespaces []string) (map[string]*Project, error) {
	return GlobalDAO.GetNamespacesForUser(namespaces)
}

func GetNamespacesFilteredBy(f *NamespacesFilter) (map[string]*Project, error) {
	return GlobalDAO.GetNamespacesFilteredBy(f)
}

func GetNamespaceHooks(namespace *Project) (Hooks, error) {
	return GlobalDAO.GetNamespaceHooks(namespace)
}

func SetNamespaceHooks(namespace *Project, hooks Hooks) error {
	return GlobalDAO.SetNamespaceHooks(namespace, hooks)
}

func HardDeleteNamespace(namespace string) error {
	return GlobalDAO.HardDeleteNamespace(namespace)
}

func GetApplications(namespace string) (map[string]*Application, error) {
	return GlobalDAO.GetApplications(namespace)
}

func AddApplication(app *Application) error {
	return GlobalDAO.AddApplication(app)
}

func UpdateApplication(app *Application) error {
	return GlobalDAO.UpdateApplication(app)
}

func GetApplication(namespace, name string) (*Application, error) {
	return GlobalDAO.GetApplication(namespace, name)
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

func GetRelease(namespace, name, releaseId string) (*Release, error) {
	return GlobalDAO.GetRelease(namespace, name, releaseId)
}
func GetReleaseByTag(namespace, name, tag string) (*Release, error) {
	return GlobalDAO.GetReleaseByTag(namespace, name, tag)
}
func TagRelease(release *Release, tag string) error {
	return GlobalDAO.TagRelease(release, tag)
}

func GetProviders(providerName string) (map[string]*MinimalReleaseMetadata, error) {
	return GlobalDAO.GetProviders(providerName)
}

func GetProvidersFilteredBy(providerName string, f *ProvidersFilter) (map[string]*MinimalReleaseMetadata, error) {
	return GlobalDAO.GetProvidersFilteredBy(providerName, f)
}

func RegisterProviders(release *core.ReleaseMetadata) error {
	return GlobalDAO.RegisterProviders(release)
}

func FindAllVersions(app *Application) ([]string, error) {
	return GlobalDAO.FindAllVersions(app)
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

func GetDownstreamDependenciesFilteredBy(r *Release, f *DownstreamDependenciesFilter) ([]*Dependency, error) {
	return GlobalDAO.GetDownstreamDependenciesFilteredBy(r, f)
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
