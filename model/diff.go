package model

import (
	core "github.com/ankyra/escape-core"
)

func Diff(project, name, version, diffWith string) (map[string]map[string]core.Changes, error) {
	metadata, err := GetReleaseMetadata(project, name, version)
	if err != nil {
		return nil, err
	}
	if diffWith == "" {
		prev, err := GetPreviousVersion(project, name, metadata.Version)
		if err != nil {
			return nil, err
		}
		diffWith = prev
	}
	previousMetadata, err := GetReleaseMetadata(project, name, diffWith)
	if err != nil {
		return nil, err
	}
	return core.Diff(previousMetadata, metadata).Collapse(), nil
}
