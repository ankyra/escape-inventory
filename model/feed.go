package model

import (
	"github.com/ankyra/escape-inventory/dao/types"
)

func GetFeedPage(page int) ([]*types.FeedEvent, error) {
	return []*types.FeedEvent{
		types.NewReleaseEvent("_", "escape", "0.1.1", "admin@service-account"),
		types.NewReleaseEvent("_", "escape", "0.1.0", "admin@service-account"),
		types.NewCreateProjectEvent("ankyra"),
		types.NewCreateProjectEvent("_"),
	}, nil
}

func GetProjectFeed(project string, page int) ([]*types.FeedEvent, error) {
	return []*types.FeedEvent{
		types.NewReleaseEvent("_", "escape", "0.1.1", "admin@service-account"),
		types.NewReleaseEvent("_", "escape", "0.1.0", "admin@service-account"),
		types.NewCreateProjectEvent("ankyra"),
		types.NewCreateProjectEvent("_"),
	}, nil
}
