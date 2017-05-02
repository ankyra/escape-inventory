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

func GetApplications() ([]ApplicationDAO, error) {
	return globalDAO.GetApplications()
}

func GetApplication(typ, name string) (ApplicationDAO, error) {
	return globalDAO.GetApplication(typ, name)
}

func GetRelease(releaseId string) (ReleaseDAO, error) {
	return globalDAO.GetRelease(releaseId)
}
func GetAllReleases() ([]ReleaseDAO, error) {
	return globalDAO.GetAllReleases()
}

func AddRelease(metadata Metadata) (ReleaseDAO, error) {
	return globalDAO.AddRelease(metadata)
}
func GetReleaseTypes() ([]string, error) {
	return globalDAO.GetReleaseTypes()
}
func GetApplicationsByType(typ string) ([]string, error) {
	return globalDAO.GetApplicationsByType(typ)
}

func IsNotFound(err error) bool {
	return err == NotFound
}
func IsAlreadyExists(err error) bool {
	return err == AlreadyExists
}
