package sqlhelp

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	. "github.com/ankyra/escape-inventory/dao/types"
)

func (s *SQLHelper) GetFeedPage(pageSize int) ([]*FeedEvent, error) {
	rows, err := s.PrepareAndQuery(s.FeedEventPageQuery, pageSize)
	if err != nil {
		return nil, err
	}
	return s.scanFeedEvents(rows)
}

func (s *SQLHelper) GetProjectFeedPage(project string, pageSize int) ([]*FeedEvent, error) {
	rows, err := s.PrepareAndQuery(s.ProjectFeedEventPageQuery, project, pageSize)
	if err != nil {
		return nil, err
	}
	return s.scanFeedEvents(rows)
}

func (s *SQLHelper) GetFeedPageByGroups(readGroups []string, pageSize int) ([]*FeedEvent, error) {
	starFound := false
	for _, g := range readGroups {
		if g == "*" {
			starFound = true
			break
		}
	}
	if !starFound {
		readGroups = append(readGroups, "*")
	}
	insertMarks := []string{}
	for i, _ := range readGroups {
		if s.UseNumericInsertMarks {
			insertMarks = append(insertMarks, "$"+strconv.Itoa(i+1))
		} else {
			insertMarks = append(insertMarks, "?")
		}
	}
	query := s.FeedEventsByGroupsPageQuery
	if len(readGroups) == 1 {
		query += " = " + insertMarks[0]
	} else {
		query += "IN (" + strings.Join(insertMarks, ", ") + ")"
	}
	interfaceGroups := []interface{}{}
	for _, g := range readGroups {
		interfaceGroups = append(interfaceGroups, g)
	}
	query += " ORDER BY id DESC LIMIT " + strconv.Itoa(pageSize)
	rows, err := s.PrepareAndQuery(query, interfaceGroups...)
	if err != nil {
		return nil, err
	}
	return s.scanFeedEvents(rows)
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
		var id, eventType, username, project, data string
		var uploadedAt int64
		if err := rows.Scan(&id, &eventType, &username, &project, &uploadedAt, &data); err != nil {
			return nil, err
		}
		ev := &FeedEvent{}
		ev.ID = id
		ev.Type = eventType
		ev.Project = project
		ev.Username = username
		ev.Timestamp = time.Unix(uploadedAt, 0)
		ev.Data = map[string]interface{}{}
		if err := json.Unmarshal([]byte(data), &ev.Data); err != nil {
			return nil, err
		}
		result = append(result, ev)
	}
	return result, nil
}
