package types

import (
    "github.com/ankyra/escape-client/model/interfaces"
    "fmt"
)

type Metadata interfaces.ReleaseMetadata

type DAO interface {
    GetApplications() ([]ApplicationDAO, error)
    GetApplication(typ, name string) (ApplicationDAO, error)
    GetRelease(releaseId string) (ReleaseDAO, error)
    AddRelease(metadata Metadata) error
}

type ApplicationDAO interface {
    GetType() string
    GetName() string
    FindAllVersions() ([]string, error)
}

type ReleaseDAO interface {
    GetApplication() ApplicationDAO
    GetVersion() string
    GetMetadata() Metadata
}

type PackageDAO interface {
    GetRelease() ReleaseDAO
    GetURI() string
}

var NotFound = fmt.Errorf("Not found")

