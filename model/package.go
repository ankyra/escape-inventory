package model

import (
    "io"
    "fmt"
    "github.com/ankyra/escape-registry/dao"
    "github.com/ankyra/escape-registry/storage"
)

func UploadPackage(releaseId string, pkg io.ReadSeeker) error {
    release, err := dao.GetRelease(releaseId)
    if err != nil {
        return err
    }
    uri, err := storage.Upload(releaseId, pkg)
    if err != nil {
        return err
    }
    return release.AddPackageURI(uri)
}


func GetDownloadReadSeeker(releaseId string) (io.ReadSeeker, error) {
    release, err := dao.GetRelease(releaseId)
    if err != nil {
        return nil, err
    }
    uris, err := release.GetPackageURIs()
    if err != nil {
        return nil, err
    }
    for _, uri := range uris {
        seeker, err := storage.Download(uri)
        if err != nil {
            return seeker, err
        }
    }
    return nil, fmt.Errorf("No valid URI found")
}
