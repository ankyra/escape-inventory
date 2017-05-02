package model

import (
	"encoding/json"
	"fmt"
	"github.com/ankyra/escape-registry/dao"
	"github.com/ankyra/escape-registry/shared"
)

func Import(releases []map[string]interface{}) error {
	for _, rel := range releases {
		metadataJson, err := json.Marshal(rel)
		if err != nil {
			return NewUserError(fmt.Errorf("Could not parse JSON: %s", err.Error()))
		}
		metadata, err := shared.NewReleaseMetadataFromJsonString(string(metadataJson))
		if err != nil {
			return NewUserError(fmt.Errorf("Could not get metadata from JSON: %s", err.Error()))
		}
		releaseDAO, err := dao.AddRelease(metadata)
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
