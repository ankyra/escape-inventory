package dao

import (
    . "github.com/ankyra/escape-registry/dao/types"
    "github.com/ankyra/escape-registry/dao/mem"
    "github.com/ankyra/escape-registry/dao/sqlite"
    "github.com/ankyra/escape-registry/config"
    "fmt"
)

var globalDAO = mem.NewInMemoryDAO()

func LoadFromConfig(conf *config.Config) (error) {
    if conf.Database == "" {
        return fmt.Errorf("Missing database configuration variable")
    } else if conf.Database == "memory" {
        globalDAO = mem.NewInMemoryDAO()
        return nil
    } else if conf.Database == "sqlite" {
        dao, err := sqlite.NewSQLiteDAO(conf.DatabaseSettings.Path)
        if err != nil {
            return err
        }
        globalDAO = dao
        return nil
    }
    return fmt.Errorf("Unknown database backend: %s", conf.Database)
}

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
func IsAlreadyExists(err error) bool {
    return err == AlreadyExists
}
