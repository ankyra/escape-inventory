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
	"github.com/ankyra/escape-core"
	. "github.com/ankyra/escape-registry/dao/types"
	pq "github.com/lib/pq"
)

type dao struct {
	db *sql.DB
}

func NewPostgresDAO(url string) (DAO, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("Couldn't open Postgres storage backend '%s': %s", url, err.Error())
	}
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS release (
            name varchar(128), 
            release_id varchar(256),
            version varchar(32),
            metadata text,
            project varchar(32),
            PRIMARY KEY(name, version, project)
        )`)
	if err != nil {
		return nil, fmt.Errorf("Couldn't initialise Postgres storage backend '%s' [1]: %s", url, err.Error())
	}
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS package (
            project varchar(32),
            release_id varchar(256), 
            uri varchar(256), 
            PRIMARY KEY(release_id, uri, project)
        )`)
	if err != nil {
		return nil, fmt.Errorf("Couldn't initialise Postgres storage backend '%s' [2]: %s", url, err.Error())
	}
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS acl (
            project varchar(32),
            group_name varchar(256),
            permission varchar(1), 
            PRIMARY KEY(project, group_name)
        )`)
	if err != nil {
		return nil, fmt.Errorf("Couldn't initialise Postgres storage backend '%s' [2]: %s", url, err.Error())
	}
	return &dao{
		db: db,
	}, nil
}

func (a *dao) GetApplications(project string) ([]*Application, error) {
	stmt, err := a.db.Prepare("SELECT DISTINCT(name) FROM release WHERE project = $1")
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(project)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []*Application{}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		result = append(result, NewApplication(project, name))
	}
	return result, nil
}

func (a *dao) GetApplication(project, name string) (*Application, error) {
	stmt, err := a.db.Prepare("SELECT name FROM release WHERE project = $1 AND name = $2")
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(project, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		return NewApplication(project, name), nil
	}
	return nil, NotFound
}

func (a *dao) FindAllVersions(app *Application) ([]string, error) {
	stmt, err := a.db.Prepare("SELECT version FROM release WHERE project = $1 AND name = $2")
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(app.Project, app.Name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []string{}
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		result = append(result, version)
	}
	return result, nil
}

func (a *dao) GetRelease(project, name, releaseId string) (*Release, error) {
	stmt, err := a.db.Prepare("SELECT metadata FROM release WHERE project = $1 AND release_id = $2")
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(project, releaseId)
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
		return NewRelease(NewApplication(project, name), metadata), nil
	}
	return nil, NotFound
}

func (a *dao) GetAllReleases() ([]*Release, error) {
	result := []*Release{}
	stmt, err := a.db.Prepare("SELECT project, metadata FROM release")
	if err != nil {
		return nil, err
	}
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
		result = append(result, NewRelease(NewApplication(project, metadata.GetName()), metadata))
	}
	return result, nil
}

func (a *dao) AddRelease(project string, release *core.ReleaseMetadata) (*Release, error) {
	stmt, err := a.db.Prepare(`
        INSERT INTO release(project, name, release_id, version, metadata) VALUES($1, $2, $3, $4, $5)`)
	if err != nil {
		return nil, err
	}
	name := release.GetName()
	_, err = stmt.Exec(project, name, release.GetReleaseId(), release.GetVersion(), release.ToJson())
	if err != nil {
		if err.(*pq.Error).Code.Name() == "unique_violation" {
			return nil, AlreadyExists
		}
		return nil, err
	}
	return NewRelease(NewApplication(project, release.GetName()), release), nil
}

func (a *dao) GetPackageURIs(release *Release) ([]string, error) {
	stmt, err := a.db.Prepare("SELECT uri FROM package WHERE project = $1 AND release_id = $2")
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(release.Application.Project, release.ReleaseId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []string{}
	for rows.Next() {
		var uri string
		if err := rows.Scan(&uri); err != nil {
			return nil, err
		}
		result = append(result, uri)
	}
	return result, nil
}

func (a *dao) AddPackageURI(release *Release, uri string) error {
	stmt, err := a.db.Prepare("INSERT INTO package (project, release_id, uri) VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(release.Application.Project, release.ReleaseId, uri)
	if err != nil {
		if err.(*pq.Error).Code.Name() == "unique_violation" {
			return AlreadyExists
		}
		return err
	}
	return nil
}

func (a *dao) SetACL(project, group string, perm Permission) error {
	stmt, err := a.db.Prepare(`INSERT INTO acl(project, group_name, permission) VALUES($1, $2, $3)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(project, group, string(perm))
	if err != nil {
		if err.(*pq.Error).Code.Name() == "unique_violation" {
			stmt, err := a.db.Prepare(`UPDATE acl SET permission = $1 WHERE project = $2 AND group_name = $3`)
			if err != nil {
				return err
			}
			_, err = stmt.Exec(project, group, string(perm))
			return err
		}
		return err
	}
	return nil
}
func (a *dao) DeleteACL(project, group string) error {
	stmt, err := a.db.Prepare("DELETE FROM acl WHERE project = $1 AND group_name = $2")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(project, group)
	return err
}
func (a *dao) GetPermittedGroups(project string, perm Permission) ([]string, error) {
	result := []string{}
	stmt, err := a.db.Prepare("SELECT group_name FROM acl WHERE project = $1 AND (permission = $2 OR permission = $3)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(project, string(perm), string(ReadAndWritePermission))
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
