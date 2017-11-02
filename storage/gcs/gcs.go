/*
Copyright 2017 Ankyra

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package gcs

import (
	"cloud.google.com/go/storage"
	"fmt"
	"github.com/ankyra/escape-core/parsers"
	"github.com/ankyra/escape-inventory/config"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"io"
	"strings"
)

type GoogleCloudStorageBackend struct {
	BucketString string
	Bucket       *storage.BucketHandle
	Client       *storage.Client
	Context      context.Context
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
	var client *storage.Client
	var err error
	if settings.Credentials == "" {
		client, err = storage.NewClient(ls.Context)
	} else {
		client, err = storage.NewClient(ls.Context, option.WithServiceAccountFile(settings.Credentials))
	}
	if err != nil {
		return err
	}
	ls.Client = client
	ls.Bucket = client.Bucket(ls.BucketString)
	return nil
}

func (ls *GoogleCloudStorageBackend) Upload(project string, releaseId *parsers.ReleaseId, pkg io.ReadSeeker) (string, error) {
	archive := strings.Join([]string{project, releaseId.Name, releaseId.ToString() + ".tgz"}, "/")
	writer := ls.Bucket.Object(archive).NewWriter(ls.Context)
	if _, err := io.Copy(writer, pkg); err != nil {
		return "", err
	}
	if err := writer.Close(); err != nil {
		return "", err
	}
	return "gcs://" + ls.BucketString + "/" + archive, nil
}

func (ls *GoogleCloudStorageBackend) Download(project, uri string) (io.Reader, error) {
	path := uri[len("gcs://"):]
	parts := strings.SplitN(path, "/", 2)
	bucket := ls.Client.Bucket(parts[0])
	return bucket.Object(parts[1]).NewReader(ls.Context)
}
