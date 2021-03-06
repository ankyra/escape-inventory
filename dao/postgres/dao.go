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

package postgres

import (
	"database/sql"
	"fmt"

	"github.com/ankyra/escape-inventory/dao/sqlhelp"
	. "github.com/ankyra/escape-inventory/dao/types"
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
		UseNumericInsertMarks:     true,
		GetProjectQuery:           `SELECT name, description, orgURL, logo, is_public FROM project WHERE name = $1`,
		AddProjectQuery:           `INSERT INTO project(name, description, orgURL, logo, is_public) VALUES ($1, $2, $3, $4, $5)`,
		UpdateProjectQuery:        `UPDATE project SET name = $1, description = $2, orgURL = $3, logo = $4, is_public = $6 WHERE name = $5`,
		GetProjectsQuery:          `SELECT name, description, orgURL, logo, is_public FROM project`,
		GetNamespacesByNamesQuery: `SELECT name, description, orgURL, logo, is_public FROM project WHERE name`,
		GetNamespacesForUserQuery: `SELECT name, description, orgURL, logo, is_public FROM project WHERE is_public = true`,
		GetProjectHooksQuery:      `SELECT hooks FROM project WHERE name = $1`,
		SetProjectHooksQuery:      `UPDATE project SET hooks = $1 WHERE name = $2`,

		GetApplicationQuery: `SELECT name, project, description, latest_version, logo, uploaded_by, uploaded_at 
							     FROM application WHERE project = $1 AND name = $2`,
		AddApplicationQuery: `INSERT INTO application(name, project, description, latest_version, logo, uploaded_by, uploaded_at)
						      VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		UpdateApplicationQuery: `UPDATE application 
								 SET description = $1, latest_version = $2, logo = $3,
                                     uploaded_by = $4, uploaded_at = $5
								 WHERE name = $6 AND project = $7`,
		GetApplicationsQuery: `SELECT name, project, description, latest_version, logo, uploaded_by, uploaded_at
								  FROM application WHERE project = $1`,
		GetApplicationHooksQuery: `SELECT hooks FROM application WHERE project = $1 AND name = $2`,
		SetApplicationHooksQuery: `UPDATE application SET hooks = $1 WHERE project = $2 AND name = $3`,
		DeleteSubscriptionsQuery: `DELETE FROM subscriptions WHERE project = $1 AND name = $2`,
		AddSubscriptionQuery:     `INSERT INTO subscriptions (project, name, subscription_project, subscription_name) VALUES ($1, $2, $3, $4);`,
		GetDownstreamSubscriptionsQuery: `SELECT app.hooks FROM 
								subscriptions AS sub
								JOIN application AS app 
								ON sub.project = app.project AND sub.name = app.name
								WHERE sub.subscription_project = $1 
								  AND sub.subscription_name = $2`,

		AddReleaseQuery: `INSERT INTO 
                          release(project, name, release_id, version, metadata, uploaded_by, uploaded_at) 
                          VALUES($1, $2, $3, $4, $5, $6, $7)`,
		UpdateReleaseQuery:                              `UPDATE release SET processed_dependencies = $1, downloads = $2 WHERE project = $3 AND name = $4 AND release_id = $5`,
		GetReleaseQuery:                                 `SELECT metadata, processed_dependencies, downloads, uploaded_by, uploaded_at FROM release WHERE project = $1 AND name = $2 AND release_id = $3`,
		GetAllReleasesQuery:                             "SELECT project, metadata, processed_dependencies, downloads, uploaded_by, uploaded_at FROM release",
		GetAllReleasesWithoutProcessedDependenciesQuery: `SELECT project, metadata, processed_dependencies, downloads, uploaded_by, uploaded_at FROM release WHERE processed_dependencies = 'false'`,
		FindAllVersionsQuery:                            "SELECT version FROM release WHERE project = $1 AND name = $2",

		GetReleaseByTagQuery: `SELECT metadata, processed_dependencies, downloads, uploaded_by, uploaded_at 
							    FROM release, release_tags AS rt 
								WHERE rt.project = $1 AND rt.application = $2 AND rt.tag = $3 
								  AND rt.version = release.version 
								  AND rt.project = release.project 
								  AND rt.application = release.name`,
		AddReleaseTagQuery:    `INSERT INTO release_tags(project, application, tag, version) VALUES ($1, $2, $3, $4)`,
		UpdateReleaseTagQuery: `UPDATE release_tags SET version = $4 WHERE project = $1 AND application = $2 AND tag = $3`,

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

		GetPackageURIsQuery: "SELECT uri FROM package WHERE project = $1 AND release_id = $2",
		AddPackageURIQuery:  "INSERT INTO package (project, release_id, uri) VALUES ($1, $2, $3)",

		CreateUserIDMetricsQuery:                  `INSERT INTO metrics(user_id) VALUES($1)`,
		GetMetricsByUserIDQuery:                   `SELECT project_count FROM metrics WHERE user_id = $1`,
		SetProjectCountMetricForUser:              `UPDATE metrics SET project_count = $3 WHERE user_id = $1 AND project_count = $2`,
		GetProviderReleasesQuery:                  `SELECT project, application, version, description FROM providers WHERE provider = $1`,
		GetProvidersForReleaseQuery:               `SELECT provider, version FROM providers WHERE project = $1 AND application = $2`,
		SetProviderQuery:                          `INSERT INTO providers(project, application, version, description, provider) VALUES ($1, $2, $3, $4, $5)`,
		UpdateProviderQuery:                       `UPDATE providers SET version = $3, description = $4 WHERE project = $1 AND application = $2 AND provider = $5`,
		HardDeleteProjectPackageURIsQuery:         `DELETE FROM package WHERE project = $1`,
		HardDeleteProjectUnitSubscriptions:        `DELETE FROM subscriptions WHERE project = $1`,
		HardDeleteProjectReleaseDependenciesQuery: `DELETE FROM release_dependency WHERE project = $1`,
		HardDeleteProjectReleasesQuery:            `DELETE FROM release WHERE project = $1`,
		HardDeleteProjectApplicationsQuery:        `DELETE FROM application WHERE project = $1`,
		HardDeleteProjectQuery:                    `DELETE FROM project WHERE name = $1 `,
		WipeDatabaseFunc: func(s *sqlhelp.SQLHelper) error {
			queries := []string{
				`TRUNCATE release CASCADE`,
				`TRUNCATE package CASCADE`,
				`TRUNCATE acl CASCADE`,
				`TRUNCATE application CASCADE`,
				`TRUNCATE project CASCADE`,
				`TRUNCATE release_dependency CASCADE`,
				`TRUNCATE subscriptions CASCADE`,
				`TRUNCATE metrics CASCADE`,
				`TRUNCATE providers CASCADE`,
			}

			for _, query := range queries {
				if err := s.PrepareAndExec(query); err != nil {
					return err
				}
			}

			return nil
		},
		IsUniqueConstraintError: func(err error) bool {
			_, typeOk := err.(*pq.Error)
			return typeOk && err.(*pq.Error).Code.Name() == "unique_violation"
		},
	}, nil
}
