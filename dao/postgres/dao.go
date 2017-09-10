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

package postgres

import (
	"database/sql"
	"fmt"

	"github.com/ankyra/escape-registry/dao/sqlhelp"
	. "github.com/ankyra/escape-registry/dao/types"
	"github.com/lib/pq"
	"github.com/mattes/migrate"
	_ "github.com/mattes/migrate/database/postgres"
	"github.com/mattes/migrate/source/go-bindata"
)

func NewPostgresDAO(url string) (DAO, error) {
	s, err := bindata.WithInstance(bindata.Resource(AssetNames(),
		func(name string) ([]byte, error) {
			return Asset(name)
		}))
	if err != nil {
		return nil, fmt.Errorf("Couldn't open Postgres storage backend '%s' [3]: %s", url, err.Error())
	}
	m, err := migrate.NewWithSourceInstance("go-bindata", s, url)
	if err != nil {
		return nil, fmt.Errorf("Couldn't open migrations for Postgres storage backend '%s' [3]: %s", url, err.Error())
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("Couldn't apply Postgres migrations to '%s' [4]: %s", url, err.Error())
	}
	sourceError, dbError := m.Close()
	if sourceError != nil {
		return nil, fmt.Errorf("Couldn't close Postgres migrations to '%s' [4]: %s", url, sourceError)
	}
	if dbError != nil {
		return nil, fmt.Errorf("Couldn't close Postgres migrations to '%s' [4]: %s", url, dbError)
	}
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("Couldn't open Postgres storage backend '%s': %s", url, err.Error())
	}
	return &sqlhelp.SQLHelper{
		DB: db,
		UseNumericInsertMarks: true,
		GetProjectQuery:       `SELECT name, description, orgURL, logo FROM project WHERE name = $1`,
		AddProjectQuery:       `INSERT INTO project(name, description, orgURL, logo) VALUES ($1, $2, $3, $4)`,
		UpdateProjectQuery:    `UPDATE project SET name = $1, description = $2, orgURL = $3, logo = $4 WHERE name = $5`,
		GetProjectsQuery:      `SELECT name, description, orgURL, logo FROM project`,
		GetProjectsByGroupsQuery: `SELECT p.name, p.description, p.orgURL, p.logo 
								     FROM project AS p
									 JOIN acl ON p.name = acl.project
									 WHERE group_name `,
		GetProjectHooksQuery: `SELECT hooks FROM project WHERE name = $1`,
		SetProjectHooksQuery: `UPDATE project SET hooks = $1 WHERE name = $2`,

		GetApplicationQuery: `SELECT name, project, description, latest_version, logo, uploaded_by, uploaded_at 
							     FROM application WHERE project = $1 AND name = $2`,
		AddApplicationQuery: `INSERT INTO application(name, project, description, latest_version, logo)
						      VALUES ($1, $2, $3, $4, $5)`,
		UpdateApplicationQuery: `UPDATE application 
								 SET description = $1, latest_version = $2, logo = $3,
                                     uploaded_by = $4, uploaded_at = $5
								 WHERE name = $6 AND project = $7`,
		GetApplicationsQuery: `SELECT name, project, description, latest_version, logo, uploaded_by, uploaded_at
								  FROM application WHERE project = $1`,
		GetApplicationHooksQuery: `SELECT hooks FROM application WHERE project = $1 AND name = $2`,
		SetApplicationHooksQuery: `UPDATE application SET hooks = $1 WHERE project = $2 AND name = $3`,

		AddReleaseQuery: `INSERT INTO 
                          release(project, name, release_id, version, metadata, uploaded_by, uploaded_at) 
                          VALUES($1, $2, $3, $4, $5, $6, $7)`,
		UpdateReleaseQuery:                              `UPDATE release SET processed_dependencies = $1, downloads = $2 WHERE project = $3 AND name = $4 AND release_id = $5`,
		GetReleaseQuery:                                 `SELECT metadata, processed_dependencies, downloads, uploaded_by, uploaded_at FROM release WHERE project = $1 AND name = $2 AND release_id = $3`,
		GetAllReleasesQuery:                             "SELECT project, metadata, processed_dependencies, downloads, uploaded_by, uploaded_at FROM release",
		GetAllReleasesWithoutProcessedDependenciesQuery: `SELECT project, metadata, processed_dependencies, downloads, uploaded_by, uploaded_at FROM release WHERE processed_dependencies = 'false'`,
		FindAllVersionsQuery:                            "SELECT version FROM release WHERE project = $1 AND name = $2",

		InsertDependencyQuery: `INSERT INTO release_dependency(project, name, version,
										dep_project, dep_name, dep_version,
										build_scope, deploy_scope, is_extension)
								VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		GetDependenciesQuery: `SELECT dep_project, dep_name, dep_version, 
									  build_scope, deploy_scope, is_extension
							   FROM release_dependency 
							   WHERE project = $1 AND name = $2 AND version = $3`,
		GetDownstreamDependenciesQuery: `SELECT project, name, version, 
							   build_scope, deploy_scope, is_extension
							   FROM release_dependency 
							   WHERE dep_project = $1 AND dep_name = $2 AND dep_version = $3`,
		GetDownstreamDependenciesByGroupsQuery: `SELECT r.project, r.name, r.version, 
							   r.build_scope, r.deploy_scope, r.is_extension
							   FROM release_dependency AS r
							   JOIN acl ON r.project = acl.project
							   WHERE dep_project = $1 AND dep_name = $2 AND dep_version = $3
							   AND group_name `,

		GetPackageURIsQuery: "SELECT uri FROM package WHERE project = $1 AND release_id = $2",
		AddPackageURIQuery:  "INSERT INTO package (project, release_id, uri) VALUES ($1, $2, $3)",

		GetACLQuery:             "SELECT group_name, permission FROM acl WHERE project = $1",
		InsertACLQuery:          "INSERT INTO acl(project, group_name, permission) VALUES ($1, $2, $3)",
		UpdateACLQuery:          "UPDATE acl SET permission = $1 WHERE project = $2 AND group_name = $3",
		DeleteACLQuery:          "DELETE FROM acl WHERE project = $1 AND group_name = $2",
		GetPermittedGroupsQuery: "SELECT group_name FROM acl WHERE project = $1 AND (permission >= $2)",
		IsUniqueConstraintError: func(err error) bool {
			_, typeOk := err.(*pq.Error)
			return typeOk && err.(*pq.Error).Code.Name() == "unique_violation"
		},
	}, nil
}
