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
	"fmt"
	"github.com/ankyra/escape-core"
	"github.com/ankyra/escape-registry/dao"
)

func Import(releases []map[string]interface{}) error {
	for _, rel := range releases {
		metadataJson, err := json.Marshal(rel)
		if err != nil {
			return NewUserError(fmt.Errorf("Could not parse JSON: %s", err.Error()))
		}
		metadata, err := core.NewReleaseMetadataFromJsonString(string(metadataJson))
		if err != nil {
			return NewUserError(fmt.Errorf("Could not get metadata from JSON: %s", err.Error()))
		}
		releaseDAO, err := dao.AddRelease("_", metadata)
		if dao.IsAlreadyExists(err) {
			continue
		}
		if err != nil {
			return err
		}
		uris, ok := rel["URI"]
		if ok {
			uriList := uris.([]interface{})
			for _, uri := range uriList {
				if err := releaseDAO.AddPackageURI(uri.(string)); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
