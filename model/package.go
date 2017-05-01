package model

import (
    "io"
    "fmt"
    "log"
    "github.com/ankyra/escape-registry/dao"
    "github.com/ankyra/escape-registry/storage"
)

func UploadPackage(releaseId string, pkg io.ReadSeeker) error {
    release, err := dao.GetRelease(releaseId)
    if err != nil {
        return NewUserError(err)
    }
    uri, err := storage.Upload(releaseId, pkg)
    if err != nil {
        return err
    }
    return release.AddPackageURI(uri)
}


func GetDownloadReadSeeker(releaseId string) (io.Reader, error) {
    release, err := dao.GetRelease(releaseId)
    if err != nil {
        return nil, NewUserError(err)
    }
    uris, err := release.GetPackageURIs()
    if err != nil {
        return nil, err
    }
    for _, uri := range uris {
        log.Println(uri)
        reader, err := storage.Download(uri)
        if err == nil {
            return reader, nil
        }
        log.Printf("Warn: %s\n", err.Error())
    }
    return nil, fmt.Errorf("No valid URI found")
}
