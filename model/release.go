package model

import (
    "github.com/ankyra/escape-registry/shared"
    "github.com/ankyra/escape-registry/dao"
    . "github.com/ankyra/escape-registry/dao/types"
    "fmt"
)

func AddRelease(metadataJson string) error {
    metadata, err := shared.NewReleaseMetadataFromJsonString(metadataJson)
    if err != nil {
        return NewUserError(err)
    }
    releaseId := metadata.GetReleaseId()
    release, err := dao.GetRelease(releaseId)
    if err != nil && !dao.IsNotFound(err) {
        return err
    }
    if release != nil {
        return NewUserError(fmt.Errorf("Release %s already exists", releaseId))
    }
    _, err = dao.AddRelease(metadata)
    return err
}

func GetReleaseMetadata(releaseId string) (Metadata, error) {
    release, err := dao.GetRelease(releaseId)
    if err != nil {
        return nil, NewUserError(err)
    }
    return release.GetMetadata(), nil
}
