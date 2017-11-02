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

package model

import (
	"encoding/json"
	"io"

	"github.com/ankyra/escape-inventory/dao"
)

func Export(w io.Writer) error {
	releases, err := dao.GetAllReleases()
	if err != nil {
		return err
	}
	result := []interface{}{}
	for _, r := range releases {
		metadata, err := r.Metadata.ToDict()
		if err != nil {
			return err
		}
		uris, err := dao.GetPackageURIs(r)
		if err != nil {
			return err
		}
		metadata["project"] = r.Application.Project
		metadata["URI"] = uris
		result = append(result, metadata)
	}
	enc := json.NewEncoder(w)
	if err := enc.Encode(result); err != nil {
		return err
	}
	return nil
}
