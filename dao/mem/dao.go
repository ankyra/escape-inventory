package mem

import (
    . "github.com/ankyra/escape-registry/dao/types"
)

type mem_dao struct {
    applications []ApplicationDAO
    releases map[string]ReleaseDAO
}

func NewInMemoryDAO() DAO {
    return &mem_dao{
        applications: []ApplicationDAO{},
        releases: map[string]ReleaseDAO{},
    }
}

func (a *mem_dao) GetApplications() ([]ApplicationDAO, error) {
    return a.applications, nil
}

func (a *mem_dao) GetApplication(typ, name string) (ApplicationDAO, error) {
    for _, app := range a.applications {
        if app.GetType() == typ && app.GetName() == name {
            return app, nil
        }
    }
    return nil, NotFound
}

func (a *mem_dao) NewApplication(typ, name string) (ApplicationDAO, error) {
    for _, app := range a.applications {
        if app.GetType() == typ && app.GetName() == name {
            return nil, AlreadyExists
        }
    }
    result := newApplication(typ, name, a)
    a.applications = append(a.applications, result)
    return result, nil
}

func (a *mem_dao) GetRelease(releaseId string) (ReleaseDAO, error) {
    release, ok := a.releases[releaseId]
    if !ok {
        return nil, NotFound
    }
    return release, nil
}
