package sqlhelp

import (
	"database/sql"

	. "github.com/ankyra/escape-inventory/dao/types"
)

type SQLHelper struct {
	DB                      *sql.DB
	UseNumericInsertMarks   bool
	WipeDatabaseFunc        func(*SQLHelper) error
	IsUniqueConstraintError func(error) bool

	GetProjectQuery          string
	AddProjectQuery          string
	UpdateProjectQuery       string
	GetProjectsQuery         string
	GetProjectsByGroupsQuery string
	GetProjectHooksQuery     string
	SetProjectHooksQuery     string

	AddApplicationQuery      string
	UpdateApplicationQuery   string
	GetApplicationsQuery     string
	GetApplicationQuery      string
	GetApplicationHooksQuery string
	SetApplicationHooksQuery string

	DeleteSubscriptionsQuery        string
	AddSubscriptionQuery            string
	GetDownstreamSubscriptionsQuery string

	AddReleaseQuery                                 string
	UpdateReleaseQuery                              string
	GetReleaseQuery                                 string
	GetAllReleasesQuery                             string
	GetAllReleasesWithoutProcessedDependenciesQuery string
	FindAllVersionsQuery                            string

	InsertDependencyQuery                  string
	GetDependenciesQuery                   string
	GetDownstreamDependenciesQuery         string
	GetDownstreamDependenciesByGroupsQuery string

	GetPackageURIsQuery string
	AddPackageURIQuery  string

	GetACLQuery             string
	InsertACLQuery          string
	UpdateACLQuery          string
	DeleteACLQuery          string
	GetPermittedGroupsQuery string

	CreateUserIDMetricsQuery     string
	GetMetricsByUserIDQuery      string
	SetProjectCountMetricForUser string

	FeedEventPageQuery string
	AddFeedEventQuery  string
}

func (s *SQLHelper) ReadRowsIntoStringArray(rows *sql.Rows) ([]string, error) {
	defer rows.Close()
	result := []string{}
	for rows.Next() {
		var arg string
		if err := rows.Scan(&arg); err != nil {
			return nil, err
		}
		result = append(result, arg)
	}
	return result, nil
}

func (s *SQLHelper) PrepareAndQuery(query string, arg ...interface{}) (*sql.Rows, error) {
	stmt, err := s.DB.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Query(arg...)
}

func (s *SQLHelper) PrepareAndExec(query string, arg ...interface{}) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(query, arg...)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *SQLHelper) PrepareAndExecInsert(query string, arg ...interface{}) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(query, arg...)

	if err != nil {
		tx.Rollback()
		if s.IsUniqueConstraintError(err) {
			return AlreadyExists
		}
		return err
	}
	return tx.Commit()
}

func (s *SQLHelper) PrepareAndExecInsertIgnoreDups(query string, arg ...interface{}) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(query, arg...)
	if err != nil {
		tx.Rollback()
		if s.IsUniqueConstraintError(err) {
			return nil
		}
		return err
	}
	return tx.Commit()
}

func (s *SQLHelper) PrepareAndExecUpdate(query string, arg ...interface{}) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	result, err := tx.Exec(query, arg...)
	if err != nil {
		tx.Rollback()
		if s.IsUniqueConstraintError(err) {
			return AlreadyExists
		}
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}
	if rowsAffected == 0 {
		tx.Rollback()
		return NotFound
	}
	return tx.Commit()
}

func (s *SQLHelper) WipeDatabase() error {
	return s.WipeDatabaseFunc(s)
}
