package mem

import (
	"strconv"

	. "github.com/ankyra/escape-inventory/dao/types"
)

func (a *dao) GetFeedPage(pageSize int) ([]*FeedEvent, error) {
	result := []*FeedEvent{}
	for i := len(a.events) - 1; i >= 0 && len(result) < pageSize; i-- {
		result = append(result, a.events[i])
	}
	return result, nil
}

func (a *dao) GetProjectFeedPage(project string, pageSize int) ([]*FeedEvent, error) {
	result := []*FeedEvent{}
	for i := len(a.events) - 1; i >= 0 && len(result) < pageSize; i-- {
		if a.events[i].Project == project {
			result = append(result, a.events[i])
		}
	}
	return result, nil
}

func (a *dao) GetApplicationFeedPage(project, application string, pageSize int) ([]*FeedEvent, error) {
	result := []*FeedEvent{}
	for i := len(a.events) - 1; i >= 0 && len(result) < pageSize; i-- {
		if a.events[i].Project == project && a.events[i].Application == application {
			result = append(result, a.events[i])
		}
	}
	return result, nil
}

func (a *dao) GetFeedPageByGroups(readGroups []string, pageSize int) ([]*FeedEvent, error) {
	added := map[string]bool{}
	result := []*FeedEvent{}
	for i := len(a.events) - 1; i >= 0 && len(result) < pageSize; i-- {
		allowedGroups, found := a.acls[a.events[i].Project]
		if !found {
			continue
		}
		for group, _ := range allowedGroups {
			found := group == "*"
			for _, g := range readGroups {
				if g == group {
					found = true
				}
			}
			if found {
				_, ok := added[a.events[i].ID]
				if !ok {
					added[a.events[i].ID] = true
					result = append(result, a.events[i])
				}
			}
		}
	}
	return result, nil
}

func (a *dao) AddFeedEvent(event *FeedEvent) error {
	if event.ID == "" {
		event.ID = strconv.Itoa(a.feedIDCounter)
		a.feedIDCounter++
	}
	a.events = append(a.events, event)
	return nil
}
