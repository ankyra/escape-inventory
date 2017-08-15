package sqlhelp

import (
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
