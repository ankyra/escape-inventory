package sqlhelp

import (
	"database/sql"

	. "github.com/ankyra/escape-registry/dao/types"
)

func (s *SQLHelper) AddApplication(app *Application) error {
	return s.PrepareAndExecInsert(s.AddApplicationQuery,
		app.Name,
		app.Project,
		app.Description,
		app.LatestReleaseId,
		app.Logo)
}
func (s *SQLHelper) UpdateApplication(app *Application) error {
	return s.PrepareAndExecUpdate(s.UpdateApplicationQuery,
		app.Description,
		app.LatestReleaseId,
		app.Logo,
		app.Name,
		app.Project)
}

func (s *SQLHelper) GetApplications(project string) (map[string]*Application, error) {
	rows, err := s.PrepareAndQuery(s.GetApplicationsQuery, project)
	if err != nil {
		return nil, err
	}
	return s.scanApplications(rows)
}

func (s *SQLHelper) GetApplication(project, name string) (*Application, error) {
	rows, err := s.PrepareAndQuery(s.GetApplicationQuery, project, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		return s.scanApplication(rows)
	}
	return nil, NotFound
}

func (s *SQLHelper) scanApplication(rows *sql.Rows) (*Application, error) {
	var name, project, description, latestReleaseId, logo string
	if err := rows.Scan(&name, &project, &description, &latestReleaseId, &logo); err != nil {
		return nil, err
	}
	return &Application{
		Name:            name,
		Project:         project,
		Description:     description,
		LatestReleaseId: latestReleaseId,
		Logo:            logo,
	}, nil
}

func (s *SQLHelper) scanApplications(rows *sql.Rows) (map[string]*Application, error) {
	defer rows.Close()
	result := map[string]*Application{}
	for rows.Next() {
		app, err := s.scanApplication(rows)
		if err != nil {
			return nil, err
		}
		result[app.Name] = app
	}
	return result, nil
}
