package sqlhelp

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	. "github.com/ankyra/escape-inventory/dao/types"
)

func (s *SQLHelper) GetFeedPage(pageSize int) ([]*FeedEvent, error) {
	fmt.Println(pageSize)
	rows, err := s.PrepareAndQuery(s.FeedEventPageQuery + strconv.Itoa(pageSize))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		return s.scanFeedEvents(rows)
	}
	return []*FeedEvent{}, nil
}

func (s *SQLHelper) GetProjectFeedPage(project string, pageSize int) ([]*FeedEvent, error) {
	return nil, nil
}

func (s *SQLHelper) GetFeedPageByGroups(readGroups []string, pageSize int) ([]*FeedEvent, error) {
	return nil, nil
}

func (s *SQLHelper) AddFeedEvent(event *FeedEvent) error {
	data, err := json.Marshal(event.Data)
	if err != nil {
		return err
	}
	return s.PrepareAndExecInsert(s.AddFeedEventQuery,
		event.Type,
		event.Username,
		event.Project,
		event.Timestamp.Unix(),
		string(data))
}

func (s *SQLHelper) scanFeedEvents(rows *sql.Rows) ([]*FeedEvent, error) {
	defer rows.Close()
	result := []*FeedEvent{}
	for rows.Next() {
		var eventType, username, project, data string
		var uploadedAt int64
		if err := rows.Scan(&eventType, &username, &project, &uploadedAt, &data); err != nil {
			return nil, err
		}
		ev := NewEvent(eventType, project)
		ev.Username = username
		ev.Timestamp = time.Unix(uploadedAt, 0)
		ev.Data = map[string]interface{}{}
		if err := json.Unmarshal([]byte(data), &ev.Data); err != nil {
			return nil, err
		}
		result = append(result, ev)
		fmt.Println(result)
	}
	return result, nil
}
