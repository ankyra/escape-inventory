package model

import (
    "strings"
    "fmt"
    "github.com/ankyra/escape-registry/shared"
    "github.com/ankyra/escape-registry/dao"
    . "github.com/ankyra/escape-registry/dao/types"
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

func GetReleaseMetadata(releaseIdString string) (Metadata, error) {
    release, err := ResolveReleaseId(releaseIdString)
    if err != nil {
        return nil, err
    }
    return release.GetMetadata(), nil
}

func ResolveReleaseId(releaseIdString string) (ReleaseDAO, error) {
    releaseId, err := shared.ParseReleaseId(releaseIdString)
    if err != nil {
        return nil, NewUserError(err)
    }
    if releaseId.Version == "latest" {
        version, err := getLastVersionForPrefix(releaseIdString, "")
        if err != nil {
            return nil, NewUserError(err)
        }
        releaseId.Version = version.ToString()
    } else if strings.HasSuffix(releaseId.Version, ".@") {
        prefix := releaseId.Version[:len(releaseId.Version) - 1]
        version, err := getLastVersionForPrefix(releaseIdString, prefix)
        if err != nil {
            return nil, NewUserError(err)
        }
        releaseId.Version = prefix + version.ToString()
    }
    return dao.GetRelease(releaseId.ToString())
}
