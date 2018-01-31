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

func AddNewApplicationFeedEvent(project, name, username string) error {
	event := types.NewCreateApplicationEvent(project, name, username)
	return dao.AddFeedEvent(event)
}

func AddNewUserAddedToProjectFeedEvent(project, username, addedByUser string) error {
	event := types.NewUserAddedToProjectEvent(project, username, addedByUser)
	return dao.AddFeedEvent(event)
}

func AddNewUserRemovedFromProjectFeedEvent(project, username, removedByUser string) error {
	event := types.NewUserRemovedFromProjectEvent(project, username, removedByUser)
	return dao.AddFeedEvent(event)
}
