package sqlhelp

import (
	. "github.com/ankyra/escape-inventory/dao/types"
)

func (s *SQLHelper) GetFeedPage(pageSize int) ([]*FeedEvent, error) {
	return nil, nil
}

func (s *SQLHelper) GetProjectFeedPage(project string, pageSize int) ([]*FeedEvent, error) {
	return nil, nil
}

func (s *SQLHelper) GetFeedPageByGroups(readGroups []string, pageSize int) ([]*FeedEvent, error) {
	return nil, nil
}

func (s *SQLHelper) AddFeedEvent(event *FeedEvent) error {
	return nil
}
