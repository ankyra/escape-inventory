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
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/ankyra/escape-inventory/dao/sqlhelp"
	. "github.com/ankyra/escape-inventory/dao/types"
	"github.com/mattes/migrate"
	sqlite_migrate "github.com/mattes/migrate/database/sqlite3"
	"github.com/mattes/migrate/source/go-bindata"
	sqlite3 "github.com/mattn/go-sqlite3"
)

func NewSQLiteDAO(path string) (DAO, error) {

	err := startupCheckDir(path)
	if err != nil {
		return nil, err
	}

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
		GetProjectHooksQuery: `SELECT hooks FROM project WHERE name = ?`,
		SetProjectHooksQuery: `UPDATE project SET hooks = ? WHERE name = ?`,

		GetApplicationQuery: `SELECT name, project, description, latest_version, logo, uploaded_by, uploaded_at 
								  FROM application WHERE project = ? AND name = ?`,
		AddApplicationQuery: `INSERT INTO application(name, project, description, latest_version, logo)
						      VALUES (?, ?, ?, ?, ?)`,
		UpdateApplicationQuery: `UPDATE application SET description = ?, latest_version = ?, logo = ?, uploaded_by = ?, uploaded_at = ? 
								 WHERE name = ? AND project = ?`,
		GetApplicationsQuery: `SELECT name, project, description, latest_version, logo, uploaded_by, uploaded_at
								  FROM application WHERE project = ?`,
		GetApplicationHooksQuery: `SELECT hooks FROM application WHERE project = ? AND name = ?`,
		SetApplicationHooksQuery: `UPDATE application SET hooks = ? WHERE project = ? AND name = ?`,
		DeleteSubscriptionsQuery: `DELETE FROM subscriptions WHERE project = ? AND name = ?`,
		AddSubscriptionQuery:     `INSERT INTO subscriptions (project, name, subscription_project, subscription_name) VALUES (?, ?, ?, ?);`,
		GetDownstreamSubscriptionsQuery: `SELECT app.hooks FROM 
								subscriptions AS sub
								JOIN application AS app 
								ON sub.project = app.project AND sub.name = app.name
								WHERE sub.subscription_project = ? 
								  AND sub.subscription_name = ?`,

		AddReleaseQuery: "INSERT INTO release(project, name, release_id, version, metadata, uploaded_by, uploaded_at) VALUES(?, ?, ?, ?, ?, ?, ?)",
		GetReleaseQuery: `SELECT metadata, processed_dependencies, downloads, uploaded_by, uploaded_at
						  FROM release 
						  WHERE project = ? AND name = ? AND release_id = ?`,
		UpdateReleaseQuery:                              `UPDATE release SET processed_dependencies = ?, downloads = ? WHERE project = ? AND name = ? AND release_id = ?`,
		GetAllReleasesQuery:                             "SELECT project, metadata, processed_dependencies, downloads, uploaded_by, uploaded_at FROM release",
		GetAllReleasesWithoutProcessedDependenciesQuery: `SELECT project, metadata, processed_dependencies, downloads, uploaded_by, uploaded_at FROM release WHERE processed_dependencies = 'false'`,
		FindAllVersionsQuery:                            "SELECT version FROM release WHERE project = ? AND name = ?",

		GetPackageURIsQuery: "SELECT uri FROM package WHERE project = ? AND release_id = ?",
		AddPackageURIQuery:  "INSERT INTO package (project, release_id, uri) VALUES (?, ?, ?)",

		InsertDependencyQuery: `INSERT INTO release_dependency(project, name, version,
										dep_project, dep_name, dep_version,
										build_scope, deploy_scope, is_extension)
								VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		GetDependenciesQuery: `SELECT dep_project, dep_name, dep_version, 
									  build_scope, deploy_scope, is_extension
							   FROM release_dependency 
							   WHERE project = ? AND name = ? AND version = ?`,
		GetDownstreamDependenciesQuery: `SELECT project, name, version, 
									  build_scope, deploy_scope, is_extension
							   FROM release_dependency 
							   WHERE dep_project = ? AND dep_name = ? AND dep_version = ?`,
		GetDownstreamDependenciesByGroupsQuery: `SELECT r.project, r.name, r.version, 
									  r.build_scope, r.deploy_scope, r.is_extension
							   FROM release_dependency AS r
							   JOIN acl ON r.project = acl.project
							   WHERE r.dep_project = ? AND r.dep_name = ? AND r.dep_version = ?
							   AND group_name `,
		GetACLQuery:             "SELECT group_name, permission FROM acl WHERE project = ?",
		InsertACLQuery:          "INSERT INTO acl(project, group_name, permission) VALUES(?, ?, ?)",
		UpdateACLQuery:          "UPDATE acl SET permission = ? WHERE project = ? AND group_name = ?",
		DeleteACLQuery:          "DELETE FROM acl WHERE project = ? AND group_name = ?",
		GetPermittedGroupsQuery: "SELECT group_name FROM acl WHERE project = ? AND (permission >= ?)",
		WipeDatabaseFunc: func(s *sqlhelp.SQLHelper) error {
			queries := []string{
				`DELETE FROM release`,
				`DELETE FROM package`,
				`DELETE FROM acl`,
				`DELETE FROM application`,
				`DELETE FROM project`,
				`DELETE FROM release_dependency`,
				`DELETE FROM subscriptions`,
			}

			for _, query := range queries {
				if err := s.PrepareAndExec(query); err != nil {
					return err
				}
			}

			return nil
		},
		IsUniqueConstraintError: func(err error) bool {
			return err.(sqlite3.Error).Code == sqlite3.ErrConstraint
		},
	}, nil
}

func startupCheckDir(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("Could not build absolute path from %s :%s", path, err.Error())
	}

	escapeDir, _ := filepath.Split(absPath)
	_, err = os.Stat(escapeDir)
	if os.IsNotExist(err) {
		return fmt.Errorf("Directory does not exist %s :%s", escapeDir, err.Error())
	}

	rand.Seed(time.Now().UnixNano())

	var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, 6)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	permissionTestFileName := string(b)

	err = ioutil.WriteFile(escapeDir+permissionTestFileName, []byte(""), 0644)
	if err != nil {
		return fmt.Errorf("Couldn't write to %s : %s ", escapeDir, err.Error())
	}

	err = os.Remove(escapeDir + permissionTestFileName)
	if err != nil {
		return fmt.Errorf("Couldn't remove file from %s : %s ", escapeDir, err.Error())
	}

	return nil
}
