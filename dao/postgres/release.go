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
	"github.com/ankyra/escape-core"
	. "github.com/ankyra/escape-registry/dao/types"
	"github.com/lib/pq"
)

type release_dao struct {
	dao       *postgres_dao
	project   string
	releaseId string
	version   string
	metadata  *core.ReleaseMetadata
}

func newRelease(project string, release *core.ReleaseMetadata, dao *postgres_dao) *release_dao {
	return &release_dao{
		dao:       dao,
		project:   project,
		releaseId: release.GetReleaseId(),
		version:   release.GetVersion(),
		metadata:  release,
	}
}

func (r *release_dao) GetApplication() ApplicationDAO {
	return newApplicationDAO(
		r.project,
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
	stmt, err := r.dao.db.Prepare("SELECT uri FROM package WHERE project = $1 AND release_id = $2")
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(r.project, r.releaseId)
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
	stmt, err := r.dao.db.Prepare("INSERT INTO package (project, release_id, uri) VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(r.project, r.releaseId, uri)
	if err != nil {
		if err.(*pq.Error).Code.Name() == "unique_violation" {
			return AlreadyExists
		}
		return err
	}
	return nil
}

func (r *release_dao) Save() (ReleaseDAO, error) {
	stmt, err := r.dao.db.Prepare(`
        INSERT INTO release(project, name, release_id, version, metadata) VALUES($1, $2, $3, $4, $5)`)
	if err != nil {
		return nil, err
	}
	name := r.metadata.GetName()
	_, err = stmt.Exec(r.project, name, r.releaseId, r.version, r.metadata.ToJson())
	if err != nil {
		if err.(*pq.Error).Code.Name() == "unique_violation" {
			return nil, AlreadyExists
		}
		return nil, err
	}
	return r, nil
}
