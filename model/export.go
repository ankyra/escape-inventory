package model

import (
	"encoding/json"
	"github.com/ankyra/escape-registry/dao"
	"io"
)

func Export(w io.Writer) error {
	releases, err := dao.GetAllReleases()
	if err != nil {
		return err
	}
	result := []interface{}{}
	for _, r := range releases {
		metadata, err := r.GetMetadata().ToDict()
		if err != nil {
			return err
		}
		uris, err := r.GetPackageURIs()
		if err != nil {
			return err
		}
		metadata["URI"] = uris
		result = append(result, metadata)
	}
	enc := json.NewEncoder(w)
	if err := enc.Encode(result); err != nil {
		return err
	}
	return nil
}
