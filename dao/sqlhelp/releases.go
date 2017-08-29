package sqlhelp

import (
	"database/sql"

	core "github.com/ankyra/escape-core"
	. "github.com/ankyra/escape-registry/dao/types"
)

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
		return s.scanRelease(project, name, rows)
	}
	return nil, NotFound
}

func (s *SQLHelper) GetAllReleases() ([]*Release, error) {
	rows, err := s.PrepareAndQuery(s.GetAllReleasesQuery)
	if err != nil {
		return nil, err
	}
	return s.scanReleases(rows)
}

func (s *SQLHelper) GetAllReleasesWithoutProcessedDependencies() ([]*Release, error) {
	rows, err := s.PrepareAndQuery(s.GetAllReleasesWithoutProcessedDependenciesQuery)
	if err != nil {
		return nil, err
	}
	return s.scanReleases(rows)
}

func (s *SQLHelper) AddRelease(project string, release *core.ReleaseMetadata) (*Release, error) {
	err := s.PrepareAndExecInsert(s.AddReleaseQuery,
		project,
		release.Name,
		release.GetReleaseId(),
		release.Version,
		release.ToJson(),
	)
	if err != nil {
		return nil, err
	}
	return NewRelease(NewApplication(project, release.Name), release), nil
}

func (s *SQLHelper) UpdateRelease(release *Release) error {
	return s.PrepareAndExecUpdate(s.UpdateReleaseQuery,
		release.ProcessedDependencies,
		release.Application.Project,
		release.Application.Name,
		release.ReleaseId,
	)
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

func (s *SQLHelper) SetDependencies(release *Release, depends []*Dependency) error {
	for _, dep := range depends {
		err := s.PrepareAndExecInsert(s.InsertDependencyQuery,
			release.Application.Project,
			release.Application.Name,
			release.Version,
			dep.Project,
			dep.Application,
			dep.Version,
			dep.BuildScope,
			dep.DeployScope)
		if err != nil {
			return err
		}
	}
	return nil
}
func (s *SQLHelper) GetDependencies(release *Release) ([]*Dependency, error) {
	rows, err := s.PrepareAndQuery(s.GetDependenciesQuery, release.Application.Project, release.Application.Name, release.Version)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []*Dependency{}
	for rows.Next() {
		var depProject, depApplication, depVersion string
		var buildScope, deployScope bool
		if err := rows.Scan(&depProject, &depApplication, &depVersion, &buildScope, &deployScope); err != nil {
			return nil, err
		}
		result = append(result, &Dependency{
			Project:     depProject,
			Application: depApplication,
			Version:     depVersion,
			BuildScope:  buildScope,
			DeployScope: deployScope,
		})
	}
	return result, nil
}

func (s *SQLHelper) scanRelease(project, name string, rows *sql.Rows) (*Release, error) {
	var metadataJson string
	var processedDependencies bool
	if err := rows.Scan(&metadataJson, &processedDependencies); err != nil {
		return nil, err
	}
	metadata, err := core.NewReleaseMetadataFromJsonString(metadataJson)
	if err != nil {
		return nil, err
	}
	rel := NewRelease(NewApplication(project, name), metadata)
	rel.ProcessedDependencies = processedDependencies
	return rel, nil
}

func (s *SQLHelper) scanReleases(rows *sql.Rows) ([]*Release, error) {
	defer rows.Close()
	result := []*Release{}
	for rows.Next() {
		var project, metadataJson string
		var processedDependencies bool
		if err := rows.Scan(&project, &metadataJson, &processedDependencies); err != nil {
			return nil, err
		}
		metadata, err := core.NewReleaseMetadataFromJsonString(metadataJson)
		if err != nil {
			return nil, err
		}
		rel := NewRelease(NewApplication(project, metadata.Name), metadata)
		rel.ProcessedDependencies = processedDependencies
		result = append(result, rel)
	}
	return result, nil
}
