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
	"github.com/ankyra/escape-core"
	"github.com/ankyra/escape-registry/config"
	"github.com/ankyra/escape-registry/dao/mem"
	"github.com/ankyra/escape-registry/dao/postgres"
	"github.com/ankyra/escape-registry/dao/sqlite"
	. "github.com/ankyra/escape-registry/dao/types"
)

var globalDAO = mem.NewInMemoryDAO()

func LoadFromConfig(conf *config.Config) error {
	if conf.Database == "" {
		return fmt.Errorf("Missing database configuration variable")
	} else if conf.Database == "memory" {
		globalDAO = mem.NewInMemoryDAO()
		return nil
	} else if conf.Database == "sqlite" {
		dao, err := sqlite.NewSQLiteDAO(conf.DatabaseSettings.Path)
		if err != nil {
			return err
		}
		globalDAO = dao
		return nil
	} else if conf.Database == "postgres" {
		dao, err := postgres.NewPostgresDAO(conf.DatabaseSettings.PostgresUrl)
		if err != nil {
			return err
		}
		globalDAO = dao
		return nil
	}
	return fmt.Errorf("Unknown database backend: %s", conf.Database)
}
func TestSetup() {
	globalDAO = mem.NewInMemoryDAO()
}

func GetApplications(project string) ([]*Application, error) {
	return globalDAO.GetApplications(project)
}

func GetApplication(project, name string) (*Application, error) {
	return globalDAO.GetApplication(project, name)
}

func FindAllVersions(app *Application) ([]string, error) {
	return globalDAO.FindAllVersions(app)
}

func GetRelease(project, name, releaseId string) (*Release, error) {
	return globalDAO.GetRelease(project, name, releaseId)
}

func AddRelease(project string, metadata *core.ReleaseMetadata) (*Release, error) {
	return globalDAO.AddRelease(project, metadata)
}

func GetPermittedGroups(project string, perm Permission) ([]string, error) {
	return globalDAO.GetPermittedGroups(project, perm)
}

func SetACL(project, group string, perm Permission) error {
	return globalDAO.SetACL(project, group, perm)
}

func DeleteACL(project, group string) error {
	return globalDAO.DeleteACL(project, group)
}

func GetPackageURIs(r *Release) ([]string, error) {
	return globalDAO.GetPackageURIs(r)
}

func AddPackageURI(r *Release, uri string) error {
	return globalDAO.AddPackageURI(r, uri)
}

// TODO: Rename to export releases
func GetAllReleases() ([]*Release, error) {
	return globalDAO.GetAllReleases()
}

func IsNotFound(err error) bool {
	return err == NotFound
}
func IsAlreadyExists(err error) bool {
	return err == AlreadyExists
}
