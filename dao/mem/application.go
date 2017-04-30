package mem

import (
    . "github.com/ankyra/escape-registry/dao/types"
)

type mem_application struct {
    dao *mem_dao
    typ  string
    name string
    releases map[string]ReleaseDAO
}

func newApplication(typ, name string, dao *mem_dao) *mem_application {
    return &mem_application{
        dao: dao,
        typ: typ,
        name: name,
        releases: map[string]ReleaseDAO{},
    }
}

func (a *mem_application) GetType() string {
    return a.typ
}

func (a *mem_application) GetName() string {
    return a.name
}

func (a *mem_application) FindAllVersions() ([]string, error) {
    versions := []string{}
    for _, r := range a.releases {
        if r.GetApplication() == a {
            versions = append(versions, r.GetVersion())
        }
    }
    return versions, nil
}

