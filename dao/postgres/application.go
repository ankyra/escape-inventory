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

import ()

type application_dao struct {
	project string
	name    string
	dao     *postgres_dao
}

func newApplicationDAO(project, name string, dao *postgres_dao) *application_dao {
	return &application_dao{
		project: project,
		name:    name,
		dao:     dao,
	}
}

func (a *application_dao) GetName() string {
	return a.name
}

func (a *application_dao) FindAllVersions() ([]string, error) {
	stmt, err := a.dao.db.Prepare("SELECT version FROM release WHERE project = $1 AND name = $2")
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(a.project, a.name)
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
