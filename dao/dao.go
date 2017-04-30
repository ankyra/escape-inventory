package dao

import (
    . "github.com/ankyra/escape-registry/dao/types"
    "github.com/ankyra/escape-registry/dao/mem"
)

var globalDAO = mem.NewInMemoryDAO()


func GetApplications() ([]ApplicationDAO, error) {
    return globalDAO.GetApplications()
}

func GetApplication(typ, name string) (ApplicationDAO, error) {
    return globalDAO.GetApplication(typ, name)
}

func GetRelease(releaseId string) (ReleaseDAO, error) {
    return globalDAO.GetRelease(releaseId)
}

func AddRelease(metadata Metadata) error {
    app, err := globalDAO.GetApplication(metadata.GetType(), metadata.GetName())
    if IsNotFound(err) {
        app, err = globalDAO.NewApplication(metadata.GetType(), metadata.GetName())
        if err != nil {
            return err
        }
    } else if err != nil {
        return err
    }
    return app.AddRelease(metadata)
}

func IsNotFound(err error) bool {
    return err == NotFound
}
