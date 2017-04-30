package mem

import (
    . "github.com/ankyra/escape-registry/dao/types"
)


type mem_release struct {
    application *mem_application
    version string
    metadata Metadata
    packages []string
}

func newRelease(release Metadata, app *mem_application) *mem_release {
    return &mem_release{
        application: app,
        version: release.GetVersion(),
        metadata: release,
        packages: []string{},
    }
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

func (r *mem_release) GetPackageURIs() ([]string, error) {
    return r.packages, nil
}

func (r *mem_release) AddPackageURI(uri string) error {
    for _, u := range r.packages {
        if u == uri {
            return AlreadyExists
        }
    }
    r.packages = append(r.packages, uri)
    return nil
}
