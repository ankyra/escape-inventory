package sqlhelp

import (
	. "github.com/ankyra/escape-inventory/dao/types"
)

func (s *SQLHelper) GetUserMetrics(userID string) (*Metrics, error) {
	err := s.PrepareAndExecInsertIgnoreDups(s.CreateUserIDMetricsQuery, userID)
	if err != nil {
		return nil, err
	}
	rows, err := s.PrepareAndQuery(s.GetMetricsByUserIDQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var projectCount int
		if err := rows.Scan(&projectCount); err != nil {
			return nil, err
		}
		return NewMetrics(projectCount), nil
	}
	return nil, NotFound
}

func (s *SQLHelper) SetUserMetrics(userID string, previous, new *Metrics) error {
	if previous.ProjectCount != new.ProjectCount {
		return s.PrepareAndExecUpdate(s.SetProjectCountMetricForUser, userID, previous.ProjectCount, new.ProjectCount)
	}
	return nil
}
