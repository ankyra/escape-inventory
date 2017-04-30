package storage

import (
    "io"
    "fmt"
    "net/url"
    "github.com/ankyra/escape-registry/storage/local"
    "github.com/ankyra/escape-registry/storage/gcs"
    "github.com/ankyra/escape-registry/shared"
)

type StorageBackend interface {
    Upload(releaseId *shared.ReleaseId, pkg io.ReadSeeker) (string, error)
    Download(uri string) (io.ReadSeeker, error)
}

var storageBackends = map[string]StorageBackend{
    "file": local.NewLocalStorageBackend(),
    "gcs": gcs.NewGoogleCloudStorageBackend(),
}

var uploadBackend = "file"

func Upload(releaseId string, pkg io.ReadSeeker) (string, error) {
    backend, ok := storageBackends[uploadBackend]
    if !ok {
        return "", fmt.Errorf("Unknown scheme")
    }
    parsedReleaseId, err := shared.ParseReleaseId(releaseId)
    if err != nil {
        return "", err
    }
    return backend.Upload(parsedReleaseId, pkg)
}

func Download(uri string) (io.ReadSeeker, error) {
    u, err := url.Parse(uri)
    if err != nil {
        return nil, err
    }
    backend, ok := storageBackends[u.Scheme]
    if !ok {
        return nil, fmt.Errorf("Unknown scheme")
    }
    return backend.Download(uri)
}
