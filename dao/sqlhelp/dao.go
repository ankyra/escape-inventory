package sqlhelp

import (
	"database/sql"
	"github.com/ankyra/escape-core"
	. "github.com/ankyra/escape-registry/dao/types"
)

type SQLHelper struct {
	DB                      *sql.DB
	GetProjectsQuery        string
	GetApplicationsQuery    string
	GetApplicationQuery     string
	FindAllVersionsQuery    string
	GetReleaseQuery         string
	GetAllReleasesQuery     string
	AddReleaseQuery         string
	GetPackageURIsQuery     string
	AddPackageURIQuery      string
	InsertACLQuery          string
	UpdateACLQuery          string
	DeleteACLQuery          string
	GetPermittedGroupsQuery string
	IsUniqueConstraintError func(error) bool
}

func (s *SQLHelper) GetProjects() ([]string, error) {
	rows, err := s.PrepareAndQuery(s.GetProjectsQuery)
	if err != nil {
		return nil, err
	}
	return s.ReadRowsIntoStringArray(rows)
}

func (s *SQLHelper) GetApplications(project string) ([]*Application, error) {
	rows, err := s.PrepareAndQuery(s.GetApplicationsQuery, project)
	if err != nil {
		return nil, err
	}
	apps, err := s.ReadRowsIntoStringArray(rows)
	if err != nil {
		return nil, err
	}
	result := []*Application{}
	for _, app := range apps {
		result = append(result, NewApplication(project, app))
	}
	return result, nil
}

func (s *SQLHelper) GetApplication(project, name string) (*Application, error) {
	rows, err := s.PrepareAndQuery(s.GetApplicationQuery, project, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		return NewApplication(project, name), nil
	}
	return nil, NotFound
}

func (s *SQLHelper) FindAllVersions(app *Application) ([]string, error) {
	rows, err := s.PrepareAndQuery(s.FindAllVersionsQuery, app.Project, app.Name)
	if err != nil {
		return nil, err
	}
	return s.ReadRowsIntoStringArray(rows)
}

func (s *SQLHelper) GetRelease(project, name, releaseId string) (*Release, error) {
	rows, err := s.PrepareAndQuery(s.GetReleaseQuery, project, name, releaseId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var metadataJson string
		if err := rows.Scan(&metadataJson); err != nil {
			return nil, err
		}
		metadata, err := core.NewReleaseMetadataFromJsonString(metadataJson)
		if err != nil {
			return nil, err
		}
		return NewRelease(NewApplication(project, name), metadata), nil
	}
	return nil, NotFound
}

func (s *SQLHelper) GetAllReleases() ([]*Release, error) {
	rows, err := s.PrepareAndQuery(s.GetAllReleasesQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []*Release{}
	for rows.Next() {
		var project, metadataJson string
		if err := rows.Scan(&project, &metadataJson); err != nil {
			return nil, err
		}
		metadata, err := core.NewReleaseMetadataFromJsonString(metadataJson)
		if err != nil {
			return nil, err
		}
		result = append(result, NewRelease(NewApplication(project, metadata.Name), metadata))
	}
	return result, nil
}

func (s *SQLHelper) AddRelease(project string, release *core.ReleaseMetadata) (*Release, error) {
	stmt, err := s.DB.Prepare(s.AddReleaseQuery)
	if err != nil {
		return nil, err
	}
	name := release.Name
	_, err = stmt.Exec(project, name, release.GetReleaseId(), release.Version, release.ToJson())
	if err != nil {
		if s.IsUniqueConstraintError(err) {
			return nil, AlreadyExists
		}
		return nil, err
	}
	return NewRelease(NewApplication(project, release.Name), release), nil
}

func (s *SQLHelper) GetPackageURIs(release *Release) ([]string, error) {
	rows, err := s.PrepareAndQuery(s.GetPackageURIsQuery, release.Application.Project, release.ReleaseId)
	if err != nil {
		return nil, err
	}
	return s.ReadRowsIntoStringArray(rows)
}

func (s *SQLHelper) AddPackageURI(release *Release, uri string) error {
	stmt, err := s.DB.Prepare(s.AddPackageURIQuery)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(release.Application.Project, release.ReleaseId, uri)
	if err != nil {
		if s.IsUniqueConstraintError(err) {
			return AlreadyExists
		}
		return err
	}
	return nil
}

func (s *SQLHelper) SetACL(project, group string, perm Permission) error {
	stmt, err := s.DB.Prepare(s.InsertACLQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(project, group, int(perm))
	if err != nil {
		if s.IsUniqueConstraintError(err) {
			stmt, err := s.DB.Prepare(s.UpdateACLQuery)
			if err != nil {
				return err
			}
			_, err = stmt.Exec(int(perm), project, group)
			return err
		}
		return err
	}
	return nil
}

func (s *SQLHelper) DeleteACL(project, group string) error {
	stmt, err := s.DB.Prepare(s.DeleteACLQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(project, group)
	return err
}

func (s *SQLHelper) GetPermittedGroups(project string, perm Permission) ([]string, error) {
	rows, err := s.PrepareAndQuery(s.GetPermittedGroupsQuery, project, int(perm))
	if err != nil {
		return nil, err
	}
	return s.ReadRowsIntoStringArray(rows)
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
