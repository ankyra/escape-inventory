package model

import (
	"github.com/ankyra/escape-inventory/dao"
	"github.com/ankyra/escape-inventory/dao/types"
)

func AddCreateProjectFeedEvent(project, username string) error {
	event := types.NewCreateProjectEvent(project, username)
	return dao.AddFeedEvent(event)
}

func AddNewReleaseFeedEvent(project, name, version, username string) error {
	event := types.NewReleaseEvent(project, name, version, username)
	return dao.AddFeedEvent(event)
}
