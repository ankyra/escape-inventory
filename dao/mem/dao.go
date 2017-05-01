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

func (a *mem_dao) GetReleaseTypes() ([]string, error) {
    types := map[string]bool{}
    for _, app := range a.applications {
        types[app.GetType()] = true
    }
    result := []string{}
    for typ, _ := range types {
        result = append(result, typ)
    }
    return result, nil
}

func (a *mem_dao) GetApplication(typ, name string) (ApplicationDAO, error) {
    for _, app := range a.applications {
        if app.GetType() == typ && app.GetName() == name {
            return app, nil
        }
    }
    return nil, NotFound
}

func (a *mem_dao) GetRelease(releaseId string) (ReleaseDAO, error) {
    release, ok := a.releases[releaseId]
    if !ok {
        return nil, NotFound
    }
    return release, nil
}

func (a *mem_dao) GetAllReleases() ([]ReleaseDAO, error) {
    result := []ReleaseDAO{}
    for _, rel := range a.releases {
        result = append(result, rel)
    }
    return result, nil
}

func (a *mem_dao) AddRelease(release Metadata) (ReleaseDAO, error) {
    key := release.GetReleaseId()
    _, alreadyExists := a.releases[key]
    if alreadyExists {
        return nil, AlreadyExists
    }
    var application *mem_application
    typ := release.GetType()
    name := release.GetName()
    for _, app := range a.applications {
        if app.GetName() == name && app.GetType() == typ {
            application = app.(*mem_application)
        }
    }
    if application == nil {
        application = newApplication(typ, name, a)
        a.applications = append(a.applications, application)
    }
    a.releases[key] = newRelease(release, application)
    application.releases[key] = a.releases[key]
    return a.releases[key], nil
}
