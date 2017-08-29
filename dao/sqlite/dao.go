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

package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/ankyra/escape-registry/dao/sqlhelp"
	. "github.com/ankyra/escape-registry/dao/types"
	"github.com/mattes/migrate"
	sqlite_migrate "github.com/mattes/migrate/database/sqlite3"
	"github.com/mattes/migrate/source/go-bindata"
	sqlite3 "github.com/mattn/go-sqlite3"
)

func NewSQLiteDAO(path string) (DAO, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("Couldn't open SQLite storage backend '%s': %s", path, err.Error())
	}
	driver, err := sqlite_migrate.WithInstance(db, &sqlite_migrate.Config{})
	s, err := bindata.WithInstance(bindata.Resource(AssetNames(),
		func(name string) ([]byte, error) {
			return Asset(name)
		}))
	if err != nil {
		return nil, fmt.Errorf("Couldn't initialise SQLite storage backend '%s' [1]: %s", path, err.Error())
	}
	m, err := migrate.NewWithInstance("go-bindata", s, "sqlite3", driver)
	if err != nil {
		return nil, fmt.Errorf("Couldn't initialise migrations for SQLite storage backend '%s' [1]: %s", path, err.Error())
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("Couldn't apply migrations to SQLite storage backend '%s' [1]: %s", path, err.Error())
	}
	return &sqlhelp.SQLHelper{
		DB:                 db,
		GetProjectQuery:    `SELECT name, description, orgURL, logo FROM project WHERE name = ?`,
		AddProjectQuery:    `INSERT INTO project(name, description, orgURL, logo) VALUES (?, ?, ?, ?)`,
		UpdateProjectQuery: `UPDATE project SET name = ?, description = ?, orgURL = ?, logo = ? WHERE name = ?`,
		GetProjectsQuery:   `SELECT name, description, orgURL, logo FROM project`,
		GetProjectsByGroupsQuery: `SELECT p.name, p.description, p.orgURL, p.logo 
								     FROM project AS p
									 JOIN acl ON p.name = acl.project
									 WHERE group_name `,

		GetApplicationQuery: `SELECT name, project, description, latest_version, logo 
								  FROM application WHERE project = ? AND name = ?`,
		AddApplicationQuery: `INSERT INTO application(name, project, description, latest_version, logo)
						      VALUES (?, ?, ?, ?, ?)`,
		UpdateApplicationQuery: `UPDATE application SET description = ?, latest_version = ?, logo = ? 
								 WHERE name = ? AND project = ?`,
		GetApplicationsQuery: `SELECT name, project, description, latest_version, logo 
								  FROM application WHERE project = ?`,

		AddReleaseQuery: "INSERT INTO release(project, name, release_id, version, metadata) VALUES(?, ?, ?, ?, ?)",
		GetReleaseQuery: `SELECT metadata, processed_dependencies 
						  FROM release 
						  WHERE project = ? AND name = ? AND release_id = ?`,
		UpdateReleaseQuery:                              `UPDATE release SET processed_dependencies = ? WHERE project = ? AND name = ? AND release_id = ?`,
		GetAllReleasesQuery:                             "SELECT project, metadata, processed_dependencies FROM release",
		GetAllReleasesWithoutProcessedDependenciesQuery: `SELECT project, metadata, processed_dependencies FROM release WHERE processed_dependencies = 'false'`,
		FindAllVersionsQuery:                            "SELECT version FROM release WHERE project = ? AND name = ?",

		GetPackageURIsQuery: "SELECT uri FROM package WHERE project = ? AND release_id = ?",
		AddPackageURIQuery:  "INSERT INTO package (project, release_id, uri) VALUES (?, ?, ?)",

		InsertDependencyQuery: `INSERT INTO release_dependency(project, name, version,
										dep_project, dep_name, dep_version,
										build_scope, deploy_scope)
								VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		GetDependenciesQuery: `SELECT dep_project, dep_name, dep_version, 
									  build_scope, deploy_scope 
							   FROM release_dependency 
							   WHERE project = ? AND name = ? AND version = ?`,
		GetDownstreamDependenciesQuery: `SELECT project, name, version, 
									  build_scope, deploy_scope 
							   FROM release_dependency 
							   WHERE dep_project = ? AND dep_name = ? AND dep_version = ?`,
		GetACLQuery:             "SELECT group_name, permission FROM acl WHERE project = ?",
		InsertACLQuery:          "INSERT INTO acl(project, group_name, permission) VALUES(?, ?, ?)",
		UpdateACLQuery:          "UPDATE acl SET permission = ? WHERE project = ? AND group_name = ?",
		DeleteACLQuery:          "DELETE FROM acl WHERE project = ? AND group_name = ?",
		GetPermittedGroupsQuery: "SELECT group_name FROM acl WHERE project = ? AND (permission >= ?)",
		IsUniqueConstraintError: func(err error) bool {
			return err.(sqlite3.Error).Code == sqlite3.ErrConstraint
		},
	}, nil
}
