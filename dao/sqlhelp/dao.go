package sqlhelp

import (
	"database/sql"

	. "github.com/ankyra/escape-registry/dao/types"
)

type SQLHelper struct {
	DB                      *sql.DB
	UseNumericInsertMarks   bool
	IsUniqueConstraintError func(error) bool

	GetProjectQuery          string
	AddProjectQuery          string
	UpdateProjectQuery       string
	GetProjectsQuery         string
	GetProjectsByGroupsQuery string

	AddApplicationQuery    string
	UpdateApplicationQuery string
	GetApplicationsQuery   string
	GetApplicationQuery    string

	FindAllVersionsQuery  string
	GetReleaseQuery       string
	GetAllReleasesQuery   string
	InsertDependencyQuery string
	GetDependenciesQuery  string
	AddReleaseQuery       string
	GetPackageURIsQuery   string
	AddPackageURIQuery    string

	GetACLQuery             string
	InsertACLQuery          string
	UpdateACLQuery          string
	DeleteACLQuery          string
	GetPermittedGroupsQuery string
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

func (s *SQLHelper) PrepareAndExecInsert(query string, arg ...interface{}) error {
	stmt, err := s.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(arg...)
	if err != nil {
		if s.IsUniqueConstraintError(err) {
			return AlreadyExists
		}
	}
	return err
}

func (s *SQLHelper) PrepareAndExecUpdate(query string, arg ...interface{}) error {
	stmt, err := s.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	result, err := stmt.Exec(arg...)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return NotFound
	}
	return err
}
