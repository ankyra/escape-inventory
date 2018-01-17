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

package ql

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ankyra/escape-core/util"
	"github.com/ankyra/escape-inventory/dao/sqlhelp"
	. "github.com/ankyra/escape-inventory/dao/types"
	_ "github.com/cznic/ql/driver"
	"github.com/mattes/migrate"
	ql_migrate "github.com/mattes/migrate/database/ql"
	"github.com/mattes/migrate/source/go-bindata"
)

func NewQLDAO(path string) (DAO, error) {

	err := startupCheckDir(path)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("ql", path)
	if err != nil {
		return nil, fmt.Errorf("Couldn't open ql storage backend '%s': %s", path, err.Error())
	}

	driver, err := ql_migrate.WithInstance(db, &ql_migrate.Config{})
	s, err := bindata.WithInstance(bindata.Resource(AssetNames(),
		func(name string) ([]byte, error) {
			return Asset(name)
		}))
	if err != nil {
		return nil, fmt.Errorf("Couldn't initialise ql storage backend '%s' [1]: %s", path, err.Error())
	}
	m, err := migrate.NewWithInstance("go-bindata", s, "ql", driver)
	if err != nil {
		return nil, fmt.Errorf("Couldn't initialise migrations for ql storage backend '%s' [1]: %s", path, err.Error())
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("Couldn't apply migrations to ql storage backend '%s' [1]: %s", path, err.Error())
	}

	return &sqlhelp.SQLHelper{
		DB: db,
		UseNumericInsertMarks: true,
		GetProjectQuery:       `SELECT name, description, orgURL, logo FROM project WHERE name = $1`,
		AddProjectQuery:       `INSERT INTO project(name, description, orgURL, logo) VALUES ($1, $2, $3, $4)`,
		UpdateProjectQuery:    `UPDATE project SET name = $1, description = $2, orgURL = $3, logo = $4 WHERE name = $5`,
		GetProjectsQuery:      `SELECT name, description, orgURL, logo FROM project`,
		GetProjectsByGroupsQuery: `SELECT project.name, project.description, project.orgURL, project.logo, acl.group_name 
								     FROM project, acl
									 WHERE acl.project = project.name
									 AND acl.group_name `,
		GetProjectHooksQuery: `SELECT hooks FROM project WHERE name = $1`,
		SetProjectHooksQuery: `UPDATE project SET hooks = $1 WHERE name = $2`,

		GetApplicationQuery: `SELECT name, project, description, latest_version, logo, uploaded_by, uploaded_at 
								  FROM application WHERE project = $1 AND name = $2`,
		AddApplicationQuery: `INSERT INTO application(name, project, description, latest_version, logo, uploaded_by, uploaded_at)
						      VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		UpdateApplicationQuery: `UPDATE application SET description = $1, latest_version = $2, logo = $3, uploaded_by = $4, uploaded_at = $5 
								 WHERE name = $6 AND project = $7`,
		GetApplicationsQuery: `SELECT name, project, description, latest_version, logo, uploaded_by, uploaded_at
								  FROM application WHERE project = $1`,
		GetApplicationHooksQuery: `SELECT hooks FROM application WHERE project = $1 AND name = $2`,
		SetApplicationHooksQuery: `UPDATE application SET hooks = $1 WHERE project = $2 AND name = $3`,
		DeleteSubscriptionsQuery: `DELETE FROM subscriptions WHERE project = $1 AND name = $2`,
		AddSubscriptionQuery:     `INSERT INTO subscriptions (project, name, subscription_project, subscription_name) VALUES ($1, $2, $3, $4);`,
		GetDownstreamSubscriptionsQuery: `SELECT application.hooks 
											FROM application, subscriptions
											WHERE subscriptions.project = application.project 
											AND subscriptions.name = application.name
											AND subscriptions.subscription_project = $1 
								  			AND subscriptions.subscription_name = $2`,

		AddReleaseQuery: "INSERT INTO release(project, name, release_id, version, metadata, uploaded_by, uploaded_at) VALUES($1, $2, $3, $4, $5, $6, $7)",
		GetReleaseQuery: `SELECT metadata, processed_dependencies, downloads, uploaded_by, uploaded_at
						  FROM release 
						  WHERE project = $1 AND name = $2 AND release_id = $3`,
		UpdateReleaseQuery:                              `UPDATE release SET processed_dependencies = $1, downloads = $2 WHERE project = $3 AND name = $4 AND release_id = $5`,
		GetAllReleasesQuery:                             "SELECT project, metadata, processed_dependencies, downloads, uploaded_by, uploaded_at FROM release",
		GetAllReleasesWithoutProcessedDependenciesQuery: `SELECT project, metadata, processed_dependencies, downloads, uploaded_by, uploaded_at FROM release WHERE processed_dependencies = false`,
		FindAllVersionsQuery:                            "SELECT version FROM release WHERE project = $1 AND name = $2",

		GetPackageURIsQuery: "SELECT uri FROM package WHERE project = $1 AND release_id = $2",
		AddPackageURIQuery:  "INSERT INTO package (project, release_id, uri) VALUES ($1, $2, $3)",

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
		GetDownstreamDependenciesByGroupsQuery: `SELECT release_dependency.project, release_dependency.name, release_dependency.version, 
													release_dependency.build_scope, release_dependency.deploy_scope, release_dependency.is_extension
							   						FROM release_dependency, acl
													WHERE release_dependency.project = acl.project
													AND release_dependency.dep_project = $1 
													AND release_dependency.dep_name = $2 AND release_dependency.dep_version = $3
													AND acl.group_name `,
		GetACLQuery:                  "SELECT group_name, permission FROM acl WHERE project = $1",
		InsertACLQuery:               "INSERT INTO acl(project, group_name, permission) VALUES($1, $2, $3)",
		UpdateACLQuery:               "UPDATE acl SET permission = $1 WHERE project = $2 AND group_name = $3",
		DeleteACLQuery:               "DELETE FROM acl WHERE project = $1 AND group_name = $2",
		GetPermittedGroupsQuery:      "SELECT group_name FROM acl WHERE project = $1 AND (permission >= $2)",
		CreateUserIDMetricsQuery:     `INSERT INTO metrics(user_id) VALUES($1)`,
		GetMetricsByUserIDQuery:      `SELECT project_count FROM metrics WHERE user_id = $1`,
		SetProjectCountMetricForUser: `UPDATE metrics SET project_count = $3 WHERE user_id = $1 AND project_count = $2`,
		AddFeedEventQuery:            `INSERT INTO feed_events(event_type, username, project, timestamp, data) VALUES ($1, $2, $3, $4, $5)`,
		FeedEventPageQuery: `SELECT id() as id, event_type, username, project, timestamp, data 
							 FROM feed_events ORDER BY id DESC LIMIT $1`,
		ProjectFeedEventPageQuery: `SELECT id() as id, event_type, username, project, timestamp, data 
									 FROM feed_events WHERE project = $1 ORDER BY id DESC LIMIT $2`,
		FeedEventsByGroupsPageQuery: `SELECT id(f) as id, f.event_type, f.username, f.project, f.timestamp, f.data 
									  FROM feed_events AS f, acl
									  WHERE f.project = acl.project
									  AND acl.group_name `,
		HardDeleteProjectFeedEventsQuery:          `DELETE FROM feed_events WHERE project = $1`,
		HardDeleteProjectACLQuery:                 `DELETE FROM acl WHERE project = $1`,
		HardDeleteProjectPackageURIsQuery:         `DELETE FROM package WHERE project = $1`,
		HardDeleteProjectUnitSubscriptions:        `DELETE FROM subscriptions WHERE project = $1`,
		HardDeleteProjectReleaseDependenciesQuery: `DELETE FROM release_dependency WHERE project = $1`,
		HardDeleteProjectReleasesQuery:            `DELETE FROM release WHERE project = $1`,
		HardDeleteProjectApplicationsQuery:        `DELETE FROM application WHERE project = $1`,
		HardDeleteProjectQuery:                    `DELETE FROM project WHERE name = $1 `,
		WipeDatabaseFunc: func(s *sqlhelp.SQLHelper) error {
			queries := []string{
				`TRUNCATE TABLE release`,
				`TRUNCATE TABLE package`,
				`TRUNCATE TABLE acl`,
				`TRUNCATE TABLE application`,
				`TRUNCATE TABLE project`,
				`TRUNCATE TABLE release_dependency`,
				`TRUNCATE TABLE subscriptions`,
				`TRUNCATE TABLE metrics`,
				`TRUNCATE TABLE feed_events`,
			}

			for _, query := range queries {
				if err := s.PrepareAndExec(query); err != nil {
					fmt.Println(err)
					return err
				}
			}

			return nil
		},
		IsUniqueConstraintError: func(err error) bool {
			return strings.Contains(err.Error(), "duplicate value")
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

	permissionTestFileName := ""

	for len(permissionTestFileName) == 0 {
		fileName := escapeDir + "." + util.RandomString(6)
		_, err = os.Stat(fileName)
		if os.IsNotExist(err) {
			permissionTestFileName = fileName
		}
	}

	err = ioutil.WriteFile(permissionTestFileName, []byte(""), 0644)
	if err != nil {
		return fmt.Errorf("Couldn't write to %s : %s ", escapeDir, err.Error())
	}

	err = os.Remove(permissionTestFileName)
	if err != nil {
		return fmt.Errorf("Couldn't remove file from %s : %s ", escapeDir, err.Error())
	}

	return nil
}
