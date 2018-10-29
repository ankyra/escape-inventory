package model

import (
	core "github.com/ankyra/escape-core"
)

func Diff(namespace, name, version, diffWith string) (map[string]map[string]core.Changes, error) {
	metadata, err := GetReleaseMetadata(namespace, name, version)
	if err != nil {
		return nil, err
	}
	if diffWith == "" {
		prev, err := GetPreviousVersion(namespace, name, metadata.Version)
		if err != nil {
			return nil, err
		}
		diffWith = prev
	}
	previousMetadata, err := GetReleaseMetadata(namespace, name, diffWith)
	if err != nil {
		return nil, err
	}
	return core.Diff(previousMetadata, metadata).Collapse(), nil
}
