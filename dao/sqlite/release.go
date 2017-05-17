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
	"github.com/ankyra/escape-core"
	. "github.com/ankyra/escape-registry/dao/types"
	sqlite3 "github.com/mattn/go-sqlite3"
)

type release_dao struct {
	dao       *sql_dao
	releaseId string
	version   string
	metadata  *core.ReleaseMetadata
}

func newRelease(release *core.ReleaseMetadata, dao *sql_dao) *release_dao {
	return &release_dao{
		dao:       dao,
		releaseId: release.GetReleaseId(),
		version:   release.GetVersion(),
		metadata:  release,
	}
}

func (r *release_dao) GetApplication() ApplicationDAO {
	return newApplicationDAO(
		r.metadata.GetName(),
		r.dao,
	)
}

func (r *release_dao) GetVersion() string {
	return r.version
}

func (r *release_dao) GetMetadata() *core.ReleaseMetadata {
	return r.metadata
}

func (r *release_dao) GetPackageURIs() ([]string, error) {
	stmt, err := r.dao.db.Prepare("SELECT uri FROM package WHERE release_id = ?")
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(r.releaseId)
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

func (r *release_dao) AddPackageURI(uri string) error {
	stmt, err := r.dao.db.Prepare("INSERT INTO package (release_id, uri) VALUES (?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(r.releaseId, uri)
	if err != nil {
		if err.(sqlite3.Error).Code == sqlite3.ErrConstraint {
			return AlreadyExists
		}
		return err
	}
	return nil
}

func (r *release_dao) Save() (ReleaseDAO, error) {
	stmt, err := r.dao.db.Prepare(`
        INSERT INTO release(project, name, release_id, version, metadata) VALUES(?, ?, ?, ?, ?)`)
	if err != nil {
		return nil, err
	}
	project := ""
	name := r.metadata.GetName()
	_, err = stmt.Exec(project, name, r.releaseId, r.version, r.metadata.ToJson())
	if err != nil {
		if err.(sqlite3.Error).Code == sqlite3.ErrConstraint {
			return nil, AlreadyExists
		}
		return nil, err
	}
	return r, nil
}
