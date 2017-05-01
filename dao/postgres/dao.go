package postgres

import (
    "fmt"
    "database/sql"
    _ "github.com/lib/pq"
    . "github.com/ankyra/escape-registry/dao/types"
    "github.com/ankyra/escape-registry/shared"
)

type postgres_dao struct {
    db *sql.DB
}

func NewPostgresDAO(url string) (DAO, error) {
    db, err := sql.Open("postgres", url)
    if err != nil {
        return nil, fmt.Errorf("Couldn't open Postgres storage backend '%s': %s", url, err.Error())
    }
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS release (
            typ varchar(32), 
            name varchar(128), 
            release_id varchar(256),
            version varchar(32),
            metadata text,
            project varchar(32),
            PRIMARY KEY(typ, name, version, project)
        )`)
    if err != nil {
        return nil, fmt.Errorf("Couldn't initialise Postgres storage backend '%s' [1]: %s", url, err.Error())
    }
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS package (
            release_id varchar(256), 
            uri varchar(256), 
            PRIMARY KEY(release_id, uri)
        )`)
    if err != nil {
        return nil, fmt.Errorf("Couldn't initialise Postgres storage backend '%s' [2]: %s", url, err.Error())
    }
    return &postgres_dao{
        db: db,
    }, nil
}

func (a *postgres_dao) GetApplications() ([]ApplicationDAO, error) {
    stmt, err := a.db.Prepare("SELECT DISTINCT(typ), name FROM release")
    if err != nil {
        return nil, err
    }
    rows, err := stmt.Query()
    if err != nil {
        return nil, err
    }
    defer rows.Close()
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

func (a *postgres_dao) GetReleaseTypes() ([]string, error) {
    stmt, err := a.db.Prepare("SELECT DISTINCT(typ) FROM release")
    if err != nil {
        return nil, err
    }
    rows, err := stmt.Query()
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    result := []string{}
    for rows.Next() {
        var typ string
        if err := rows.Scan(&typ); err != nil {
            return nil, err
        }
        result = append(result, typ)
    }
    return result, nil
}

func (a *postgres_dao) GetApplicationsByType(typ string) ([]string, error) {
    stmt, err := a.db.Prepare("SELECT DISTINCT(name) FROM release WHERE typ = $1")
    if err != nil {
        return nil, err
    }
    rows, err := stmt.Query(typ)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    result := []string{}
    for rows.Next() {
        var name string
        if err := rows.Scan(&name); err != nil {
            return nil, err
        }
        result = append(result, name)
    }
    return result, nil
}

func (a *postgres_dao) GetApplication(typ, name string) (ApplicationDAO, error) {
    stmt, err := a.db.Prepare("SELECT name FROM release WHERE typ = $1 AND name = $2")
    if err != nil {
        return nil, err
    }
    rows, err := stmt.Query(typ, name)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    for rows.Next() {
        return newApplicationDAO(typ, name, a), nil
    }
    return nil, NotFound
}

func (a *postgres_dao) GetRelease(releaseId string) (ReleaseDAO, error) {
    stmt, err := a.db.Prepare("SELECT metadata FROM release WHERE release_id = $1")
    if err != nil {
        return nil, err
    }
    rows, err := stmt.Query(releaseId)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    for rows.Next() {
        var metadataJson string
        if err := rows.Scan(&metadataJson); err != nil {
            return nil, err
        }
        metadata, err := shared.NewReleaseMetadataFromJsonString(metadataJson)
        if err != nil {
            return nil, err
        }
        return newRelease(metadata, a), nil
    }
    return nil, NotFound
}

func (a *postgres_dao) GetAllReleases() ([]ReleaseDAO, error) {
    result := []ReleaseDAO{}
    stmt, err := a.db.Prepare("SELECT metadata FROM release")
    if err != nil {
        return nil, err
    }
    rows, err := stmt.Query()
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    for rows.Next() {
        var metadataJson string
        if err := rows.Scan(&metadataJson); err != nil {
            return nil, err
        }
        metadata, err := shared.NewReleaseMetadataFromJsonString(metadataJson)
        if err != nil {
            return nil, err
        }
        result = append(result, newRelease(metadata, a))
    }
    return result, nil
}

func (a *postgres_dao) AddRelease(release Metadata) (ReleaseDAO, error) {
    releaseDao := newRelease(release, a)
    return releaseDao.Save()
}
