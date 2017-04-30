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
func GetAllReleases() ([]ReleaseDAO, error) {
    return globalDAO.GetAllReleases()
}

func AddRelease(metadata Metadata) (ReleaseDAO, error) {
    return globalDAO.AddRelease(metadata)
}

func IsNotFound(err error) bool {
    return err == NotFound
}
