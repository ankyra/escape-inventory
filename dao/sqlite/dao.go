package sqlite

import (
    "fmt"
    "database/sql"
    "github.com/ankyra/escape-client/model/release"
    _ "github.com/mattn/go-sqlite3"
    . "github.com/ankyra/escape-registry/dao/types"
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
            typ string, 
            name string, 
            release_id string,
            version string,
            metadata string,
            project string,
            PRIMARY KEY(typ, name, version, project)
        )`)
    if err != nil {
        return nil, fmt.Errorf("Couldn't initialise SQLite storage backend '%s' [1]: %s", path, err.Error())
    }
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS package (
            release_id string, 
            uri string, 
            PRIMARY KEY(release_id, uri)
        )`)
    if err != nil {
        return nil, fmt.Errorf("Couldn't initialise SQLite storage backend '%s' [2]: %s", path, err.Error())
    }
    return &sql_dao{
        db: db,
    }, nil
}

func (a *sql_dao) GetApplications() ([]ApplicationDAO, error) {
    stmt, err := a.db.Prepare("SELECT DISTINCT(typ), name FROM release")
    if err != nil {
        return nil, err
    }
    rows, err := stmt.Query()
    if err != nil {
        return nil, err
    }
    result := []ApplicationDAO{}
    for rows.Next() {
        var typ, name string
        if err := rows.Scan(&typ, &name); err != nil {
            return nil, err
        }
        result = append(result, newApplicationDAO(typ, name, a))
    }
    return result, nil
}

func (a *sql_dao) GetApplication(typ, name string) (ApplicationDAO, error) {
    stmt, err := a.db.Prepare("SELECT name FROM release WHERE typ = ? AND name = ?")
    if err != nil {
        return nil, err
    }
    rows, err := stmt.Query(typ, name)
    if err != nil {
        return nil, err
    }
    for rows.Next() {
        return newApplicationDAO(typ, name, a), nil
    }
    return nil, NotFound
}

func (a *sql_dao) GetRelease(releaseId string) (ReleaseDAO, error) {
    stmt, err := a.db.Prepare("SELECT metadata FROM release WHERE release_id = ?")
    if err != nil {
        return nil, err
    }
    rows, err := stmt.Query(releaseId)
    if err != nil {
        return nil, err
    }
    for rows.Next() {
        var metadataJson string
        if err := rows.Scan(&metadataJson); err != nil {
            return nil, err
        }
        metadata, err := release.NewReleaseMetadataFromJsonString(metadataJson)
        if err != nil {
            return nil, err
        }
        return newRelease(metadata, a), nil
    }
    return nil, NotFound
}

func (a *sql_dao) GetAllReleases() ([]ReleaseDAO, error) {
    result := []ReleaseDAO{}
    stmt, err := a.db.Prepare("SELECT metadata FROM release")
    if err != nil {
        return nil, err
    }
    rows, err := stmt.Query()
    if err != nil {
        return nil, err
    }
    for rows.Next() {
        var metadataJson string
        if err := rows.Scan(&metadataJson); err != nil {
            return nil, err
        }
        metadata, err := release.NewReleaseMetadataFromJsonString(metadataJson)
        if err != nil {
            return nil, err
        }
        result = append(result, newRelease(metadata, a))
    }
    return result, nil
}

func (a *sql_dao) AddRelease(release Metadata) (ReleaseDAO, error) {
    releaseDao := newRelease(release, a)
    return releaseDao.Save()
}
