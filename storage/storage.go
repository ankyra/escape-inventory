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

package storage

import (
	"fmt"
	"github.com/ankyra/escape-core/parsers"
	"github.com/ankyra/escape-registry/config"
	"github.com/ankyra/escape-registry/storage/gcs"
	"github.com/ankyra/escape-registry/storage/local"
	"io"
	"net/url"
)

type StorageBackend interface {
	Init(settings config.StorageSettings) error
	Upload(releaseId *parsers.ReleaseId, pkg io.ReadSeeker) (string, error)
	Download(uri string) (io.Reader, error)
}

var storageBackends = map[string]StorageBackend{
	"local": local.NewLocalStorageBackend(),
	"file":  local.NewLocalStorageBackend(),
	"gcs":   gcs.NewGoogleCloudStorageBackend(),
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
	parsedReleaseId, err := parsers.ParseReleaseId(releaseId)
	if err != nil {
		return "", err
	}
	return backend.Upload(parsedReleaseId, pkg)
}

func Download(uri string) (io.Reader, error) {
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
