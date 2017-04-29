package model

import (
    "github.com/ankyra/escape-client/model/release"
    "github.com/ankyra/escape-registry/dao"
    . "github.com/ankyra/escape-registry/dao/types"
    "fmt"
)

func AddRelease(metadataJson string) error {
    metadata, err := release.NewReleaseMetadataFromJsonString(metadataJson)
    if err != nil {
        return err
    }
    releaseId := metadata.GetReleaseId()
    release, err := dao.GetRelease(releaseId)
    if err != nil && !dao.IsNotFound(err) {
        return err
    }
    if release != nil {
        return fmt.Errorf("Release %s already exists", releaseId)
    }
    return dao.AddRelease(metadata)
}

func GetReleaseMetadata(releaseId string) (Metadata, error) {
    release, err := dao.GetRelease(releaseId)
    if err != nil {
        return nil, err
    }
    return release.GetMetadata(), nil
}
