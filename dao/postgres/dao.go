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
	pq "github.com/lib/pq"
)

var schema = `
CREATE TABLE IF NOT EXISTS release (
    name varchar(128), 
    release_id varchar(256),
    version varchar(32),
    metadata text,
    project varchar(32),
    PRIMARY KEY(name, version, project)
);

CREATE TABLE IF NOT EXISTS package (
    project varchar(32),
    release_id varchar(256), 
    uri varchar(256), 
    PRIMARY KEY(release_id, uri, project)
);

CREATE TABLE IF NOT EXISTS acl (
    project varchar(32),
    group_name varchar(256),
    permission int, 
    PRIMARY KEY(project, group_name)
);
`

func NewPostgresDAO(url string) (DAO, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("Couldn't open Postgres storage backend '%s': %s", url, err.Error())
	}
	_, err = db.Exec(schema)
	if err != nil {
		return nil, fmt.Errorf("Couldn't initialise Postgres storage backend '%s' [1]: %s", url, err.Error())
	}
	return &sqlhelp.SQLHelper{
		DB:                      db,
		GetApplicationsQuery:    "SELECT DISTINCT(name) FROM release WHERE project = $1",
		GetApplicationQuery:     "SELECT name FROM release WHERE project = $1 AND name = $2",
		FindAllVersionsQuery:    "SELECT version FROM release WHERE project = $1 AND name = $2",
		GetReleaseQuery:         "SELECT metadata FROM release WHERE project = $1 AND name = $2 AND release_id = $3",
		GetAllReleasesQuery:     "SELECT project, metadata FROM release",
		AddReleaseQuery:         "INSERT INTO release(project, name, release_id, version, metadata) VALUES($1, $2, $3, $4, $5)",
		GetPackageURIsQuery:     "SELECT uri FROM package WHERE project = $1 AND release_id = $2",
		AddPackageURIQuery:      "INSERT INTO package (project, release_id, uri) VALUES ($1, $2, $3)",
		InsertACLQuery:          "INSERT INTO acl(project, group_name, permission) VALUES ($1, $2, $3)",
		UpdateACLQuery:          "UPDATE acl SET permission = $1 WHERE project = $2 AND group_name = $3",
		DeleteACLQuery:          "DELETE FROM acl WHERE project = $1 AND group_name = $2",
		GetPermittedGroupsQuery: "SELECT group_name FROM acl WHERE project = $1 AND (permission = $2 OR permission = $3)",
		IsUniqueConstraintError: func(err error) bool {
			_, typeOk := err.(*pq.Error)
			return typeOk && err.(*pq.Error).Code.Name() == "unique_violation"
		},
	}, nil
}
