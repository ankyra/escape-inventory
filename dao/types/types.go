package types

import (
    "github.com/ankyra/escape-client/model/interfaces"
    "fmt"
)

type Metadata interfaces.ReleaseMetadata

type DAO interface {
    NewApplication(typ, name string) (ApplicationDAO, error)
    GetApplication(typ, name string) (ApplicationDAO, error)
    GetApplications() ([]ApplicationDAO, error)
    GetRelease(releaseId string) (ReleaseDAO, error)
}

type ApplicationDAO interface {
    GetType() string
    GetName() string

    FindAllVersions() ([]string, error)
    AddRelease(metadata Metadata) error
}

type ReleaseDAO interface {
    GetApplication() ApplicationDAO
    GetVersion() string
    GetMetadata() Metadata

    GetPackageURIs() ([]string, error)
    AddPackageURI(uri string) error
}

var NotFound = fmt.Errorf("Not found")
var AlreadyExists = fmt.Errorf("Already exists")
