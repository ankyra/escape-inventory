package model

import (
    "encoding/json"
    "github.com/ankyra/escape-registry/dao"
    "github.com/ankyra/escape-client/model/release"
)


func Import(releases []map[string]interface{}) error {
    for _, rel := range releases {
        metadataJson, err := json.Marshal(rel)
        if err != nil {
            return err
        }
        metadata, err := release.NewReleaseMetadataFromJsonString(string(metadataJson))
        if err != nil {
            return err
        }
        releaseDAO, err := dao.AddRelease(metadata)
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
