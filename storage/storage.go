package storage

import (
    "io"
    "fmt"
    "net/url"
    "github.com/ankyra/escape-registry/storage/local"
    "github.com/ankyra/escape-registry/storage/gcs"
    "github.com/ankyra/escape-registry/shared"
    "github.com/ankyra/escape-registry/config"
)

type StorageBackend interface {
    Init(settings config.StorageSettings) error
    Upload(releaseId *shared.ReleaseId, pkg io.ReadSeeker) (string, error)
    Download(uri string) (io.ReadSeeker, error)
}

var storageBackends = map[string]StorageBackend{
    "local": local.NewLocalStorageBackend(),
    "gcs": gcs.NewGoogleCloudStorageBackend(),
}

var uploadBackend = "local"

func LoadFromConfig(conf *config.Config) error {
    switch conf.StorageBackend {
    case "":
        return fmt.Errorf("Missing storage backend configuration variable")
    case "local", "gcs":
        backend, _ := storageBackends[conf.StorageBackend]
        err := backend.Init(conf.StorageSettings)
        if err != nil {
            return fmt.Errorf("Could not initialize '%s' storage backend: %s", conf.StorageBackend, err.Error())
        }
        uploadBackend = conf.StorageBackend
        return nil
    }
    return fmt.Errorf("Unknown storage backend: %s", conf.StorageBackend)
}

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
