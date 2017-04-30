package sqlite

import (
)

type application_dao struct {
    typ string
    name string
    dao *sql_dao
}

func newApplicationDAO(typ, name string, dao *sql_dao) *application_dao {
    return &application_dao{
        typ: typ,
        name: name,
        dao: dao,
    }
}

func (a *application_dao) GetType() string {
    return a.typ
}

func (a *application_dao) GetName() string {
    return a.name
}

func (a *application_dao) FindAllVersions() ([]string, error) {
    stmt, err := a.dao.db.Prepare("SELECT version FROM release WHERE typ = ? AND name = ?")
    if err != nil {
        return nil, err
    }
    rows, err := stmt.Query(a.typ, a.name)
    if err != nil {
        return nil, err
    }
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
