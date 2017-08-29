package model

import (
	"github.com/ankyra/escape-registry/dao"
	"github.com/ankyra/escape-registry/dao/types"
)

func GetDownstreamDependencies(project, name, version string) ([]*types.Dependency, error) {
	releaseId := name + "-" + version
	release, err := ResolveReleaseId(project, releaseId)
	if err != nil {
		return nil, err
	}
	return dao.GetDownstreamDependencies(release)
}
