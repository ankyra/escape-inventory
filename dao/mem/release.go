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

package mem

import (
	. "github.com/ankyra/escape-registry/dao/types"
)

type mem_release struct {
	application *mem_application
	version     string
	metadata    Metadata
	packages    []string
}

func newRelease(release Metadata, app *mem_application) *mem_release {
	return &mem_release{
		application: app,
		version:     release.GetVersion(),
		metadata:    release,
		packages:    []string{},
	}
}

func (r *mem_release) GetApplication() ApplicationDAO {
	return r.application
}

func (r *mem_release) GetVersion() string {
	return r.version
}

func (r *mem_release) GetMetadata() Metadata {
	return r.metadata
}

func (r *mem_release) GetPackageURIs() ([]string, error) {
	return r.packages, nil
}

func (r *mem_release) AddPackageURI(uri string) error {
	for _, u := range r.packages {
		if u == uri {
			return AlreadyExists
		}
	}
	r.packages = append(r.packages, uri)
	return nil
}
