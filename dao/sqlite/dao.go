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
	"github.com/ankyra/escape-core"
	. "github.com/ankyra/escape-registry/dao/types"
	sqlite3 "github.com/mattn/go-sqlite3"
)

type sql_dao struct {
	db *sql.DB
}

func NewSQLiteDAO(path string) (DAO, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("Couldn't open SQLite storage backend '%s': %s", path, err.Error())
	}
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS release (
            name string, 
            release_id string,
            version string,
            metadata string,
            project string,
            PRIMARY KEY(name, version, project)
        )`)
	if err != nil {
		return nil, fmt.Errorf("Couldn't initialise SQLite storage backend '%s' [1]: %s", path, err.Error())
	}
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS package (
            project string,
            release_id string, 
            uri string, 
            PRIMARY KEY(project, release_id, uri)
        )`)
	if err != nil {
		return nil, fmt.Errorf("Couldn't initialise SQLite storage backend '%s' [2]: %s", path, err.Error())
	}
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS acl (
            project string,
            group_name string, 
            permission varchar(1),
            PRIMARY KEY(project, group_name)
        )`)
	if err != nil {
		return nil, fmt.Errorf("Couldn't initialise SQLite storage backend '%s' [3]: %s", path, err.Error())
	}
	return &sql_dao{
		db: db,
	}, nil
}

func (a *sql_dao) GetApplications(project string) ([]ApplicationDAO, error) {
	stmt, err := a.db.Prepare("SELECT DISTINCT(name) FROM release WHERE project = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(project)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []ApplicationDAO{}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		result = append(result, newApplicationDAO(project, name, a))
	}
	return result, nil
}

func (a *sql_dao) GetApplication(project, name string) (ApplicationDAO, error) {
	stmt, err := a.db.Prepare("SELECT name FROM release WHERE project = ? AND name = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(project, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		return newApplicationDAO(project, name, a), nil
	}
	return nil, NotFound
}

func (a *sql_dao) GetRelease(project, name, releaseId string) (ReleaseDAO, error) {
	stmt, err := a.db.Prepare("SELECT metadata FROM release WHERE project = ? AND name = ? AND release_id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(project, name, releaseId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var metadataJson string
		if err := rows.Scan(&metadataJson); err != nil {
			return nil, err
		}
		metadata, err := core.NewReleaseMetadataFromJsonString(metadataJson)
		if err != nil {
			return nil, err
		}
		return newRelease(project, metadata, a), nil
	}
	return nil, NotFound
}

func (a *sql_dao) GetAllReleases() ([]ReleaseDAO, error) {
	result := []ReleaseDAO{}
	stmt, err := a.db.Prepare("SELECT project, metadata FROM release")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var project, metadataJson string
		if err := rows.Scan(&project, &metadataJson); err != nil {
			return nil, err
		}
		metadata, err := core.NewReleaseMetadataFromJsonString(metadataJson)
		if err != nil {
			return nil, err
		}
		result = append(result, newRelease(project, metadata, a))
	}
	return result, nil
}

func (a *sql_dao) AddRelease(project string, release *core.ReleaseMetadata) (ReleaseDAO, error) {
	releaseDao := newRelease(project, release, a)
	return releaseDao.Save()
}

func (a *sql_dao) SetACL(project, group string, perm Permission) error {
	stmt, err := a.db.Prepare(`INSERT INTO acl(project, group_name, permission) VALUES(?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(project, group, perm)
	if err != nil {
		if err.(sqlite3.Error).Code == sqlite3.ErrConstraint {
			stmt, err := a.db.Prepare(`UPDATE acl SET permission = ? WHERE project = ? AND group_name = ?`)
			if err != nil {
				return err
			}
			_, err = stmt.Exec(project, group, perm)
			return err
		}
		return err
	}
	return nil
}

func (a *sql_dao) DeleteACL(project, group string) error {
	stmt, err := a.db.Prepare("DELETE FROM acl WHERE project = ? AND group_name = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(project, group)
	return err
}

func (a *sql_dao) GetPermittedGroups(project string, perm Permission) ([]string, error) {
	result := []string{}
	stmt, err := a.db.Prepare("SELECT group_name FROM acl WHERE project = ? AND (permission = ? OR permission = ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(project, perm, ReadAndWritePermission)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var group string
		if err := rows.Scan(&group); err != nil {
			return nil, err
		}
		result = append(result, group)
	}
	return result, nil
}
