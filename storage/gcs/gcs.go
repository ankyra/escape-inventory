package gcs

import (
    "io"
    "fmt"
    "strings"
    "github.com/ankyra/escape-registry/shared"
    "github.com/ankyra/escape-registry/config"
    "cloud.google.com/go/storage"
    "golang.org/x/net/context"
)

type GoogleCloudStorageBackend struct {
    BucketString string
    Bucket *storage.BucketHandle
    Context context.Context
}

func NewGoogleCloudStorageBackend() *GoogleCloudStorageBackend {
    return &GoogleCloudStorageBackend{}
}

func (ls *GoogleCloudStorageBackend) Init(settings config.StorageSettings) error {
    if settings.Bucket == "" {
        return fmt.Errorf("Missing storage_settings.bucket configuration variable")
    }
    ls.BucketString = settings.Bucket
    ls.Context = context.Background()
    client, err := storage.NewClient(ls.Context)
    if err != nil {
        return err
    }
    ls.Bucket = client.Bucket(ls.BucketString)
    return nil
}

func (ls *GoogleCloudStorageBackend) Upload(releaseId *shared.ReleaseId, pkg io.ReadSeeker) (string, error) {
    archive := strings.Join([]string{releaseId.Type, releaseId.Name, releaseId.ToString() + ".tgz"}, "/")
    writer := ls.Bucket.Object(archive).NewWriter(ls.Context)
    if _, err := io.Copy(writer, pkg); err != nil {
        return "", err
    }
    if err := writer.Close(); err != nil {
        return "", err
    }
    return "gs://" + ls.BucketString + "/" + archive, nil
}

func (ls *GoogleCloudStorageBackend) Download(uri string) (io.ReadSeeker, error) {
    return nil, nil
}

