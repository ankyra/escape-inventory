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
	sqlite3 "github.com/mattn/go-sqlite3"
)

var schema = `
CREATE TABLE IF NOT EXISTS release (
    name string, 
    release_id string,
    version string,
    metadata string,
    project string,
    PRIMARY KEY(name, version, project)
);

CREATE TABLE IF NOT EXISTS package (
    project string,
    release_id string, 
    uri string, 
    PRIMARY KEY(project, release_id, uri)
);

CREATE TABLE IF NOT EXISTS acl (
    project string,
    group_name string, 
    permission varchar(1),
    PRIMARY KEY(project, group_name)
);
`

func NewSQLiteDAO(path string) (DAO, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("Couldn't open SQLite storage backend '%s': %s", path, err.Error())
	}
	_, err = db.Exec(schema)
	if err != nil {
		return nil, fmt.Errorf("Couldn't initialise SQLite storage backend '%s' [1]: %s", path, err.Error())
	}
	return &sqlhelp.SQLHelper{
		DB:                      db,
		GetApplicationsQuery:    "SELECT DISTINCT(name) FROM release WHERE project = ?",
		GetApplicationQuery:     "SELECT name FROM release WHERE project = ? AND name = ?",
		FindAllVersionsQuery:    "SELECT version FROM release WHERE project = ? AND name = ?",
		GetReleaseQuery:         "SELECT metadata FROM release WHERE project = ? AND name = ? AND release_id = ?",
		GetAllReleasesQuery:     "SELECT project, metadata FROM release",
		AddReleaseQuery:         "INSERT INTO release(project, name, release_id, version, metadata) VALUES(?, ?, ?, ?, ?)",
		GetPackageURIsQuery:     "SELECT uri FROM package WHERE project = ? AND release_id = ?",
		AddPackageURIQuery:      "INSERT INTO package (project, release_id, uri) VALUES (?, ?, ?)",
		InsertACLQuery:          "INSERT INTO acl(project, group_name, permission) VALUES(?, ?, ?)",
		UpdateACLQuery:          "UPDATE acl SET permission = ? WHERE project = ? AND group_name = ?",
		DeleteACLQuery:          "DELETE FROM acl WHERE project = ? AND group_name = ?",
		GetPermittedGroupsQuery: "SELECT group_name FROM acl WHERE project = ? AND (permission = ? OR permission = ?)",
		IsUniqueConstraintError: func(err error) bool {
			return err.(sqlite3.Error).Code == sqlite3.ErrConstraint
		},
	}, nil
}
