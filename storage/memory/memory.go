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

package memory

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/ankyra/escape-core/parsers"
	"github.com/ankyra/escape-inventory/config"
	"github.com/ankyra/escape-inventory/dao/types"
)

type InMemoryStorageBackend struct {
	URIs map[string]string
}

func NewInMemoryStorageBackend() *InMemoryStorageBackend {
	return &InMemoryStorageBackend{
		URIs: map[string]string{},
	}
}

func (m *InMemoryStorageBackend) Init(settings config.StorageSettings) error {
	return nil
}

func (m *InMemoryStorageBackend) Upload(project string, releaseId *parsers.ReleaseId, pkg io.ReadSeeker) (string, error) {
	uri := "mem://" + releaseId.Name + "-v" + releaseId.Version + ".tgz"
	data, err := ioutil.ReadAll(pkg)
	if err != nil {
		return "", err
	}
	m.URIs[uri] = string(data)
	return uri, nil
}

func (m *InMemoryStorageBackend) Download(project, uri string) (io.Reader, error) {
	data, exists := m.URIs[uri]
	if !exists {
		return nil, types.NotFound
	}
	return bytes.NewReader([]byte(data)), nil
}
