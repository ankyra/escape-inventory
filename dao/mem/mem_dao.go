package mem

import (
    . "github.com/ankyra/escape-registry/dao/types"
)

var applications = []ApplicationDAO{}
var releases = map[string]ReleaseDAO{}

type mem_dao struct { }
type mem_application struct {
    name string
    typ  string
}
type mem_release struct {
    application *mem_application
    version string
    metadata Metadata
}

func NewInMemoryDAO() DAO {
    return &mem_dao{}
}

func (a *mem_dao) GetApplications() ([]ApplicationDAO, error) {
    return applications, nil
}

func (a *mem_dao) GetApplication(typ, name string) (ApplicationDAO, error) {
    for _, app := range applications {
        if app.GetType() == typ && app.GetName() == name {
            return app, nil
        }
    }
    return nil, NotFound
}

func (a *mem_dao) GetRelease(releaseId string) (ReleaseDAO, error) {
    release, ok := releases[releaseId]
    if !ok {
        return nil, NotFound
    }
    return release, nil
}

func (a *mem_dao) AddRelease(release Metadata) error {
    app := &mem_application{
        name: release.GetName(),
        typ: release.GetType(),
    }
    applications = append(applications, app)
    result := &mem_release{
        application: app,
        version: release.GetVersion(),
        metadata: release,
    }
    releases[release.GetReleaseId()] = result
    return nil
}

func (a *mem_application) GetType() string {
    return a.typ
}

func (a *mem_application) GetName() string {
    return a.name
}

func (a *mem_application) FindAllVersions() ([]string, error) {
    versions := []string{}
    for _, r := range releases {
        if r.GetApplication() == a {
            versions = append(versions, r.GetVersion())
        }
    }
    return versions, nil
}

func (r *mem_release) GetApplication() ApplicationDAO {
    return r.application
}
func (r *mem_release) GetVersion() string {
    return r.version
}
func (r *mem_release) GetMetadata() Metadata {
    return r.metadata
}
