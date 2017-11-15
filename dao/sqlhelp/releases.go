package sqlhelp

import (
	"database/sql"
	"time"

	core "github.com/ankyra/escape-core"
	. "github.com/ankyra/escape-inventory/dao/types"
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

func (s *SQLHelper) AddRelease(release *Release) error {
	return s.PrepareAndExecInsert(s.AddReleaseQuery,
		release.Application.Project,
		release.Application.Name,
		release.Metadata.GetReleaseId(),
		release.Version,
		release.Metadata.ToJson(),
		release.UploadedBy,
		release.UploadedAt.Unix(),
	)
}

func (s *SQLHelper) UpdateRelease(release *Release) error {
	return s.PrepareAndExecUpdate(s.UpdateReleaseQuery,
		release.ProcessedDependencies,
		release.Downloads,
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
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(s.AddPackageURIQuery, release.Application.Project, release.ReleaseId, uri)
	if err != nil {
		if s.IsUniqueConstraintError(err) {
			return AlreadyExists
		}
		return err
	}
	return tx.Commit()
}

func (s *SQLHelper) scanRelease(project, name string, rows *sql.Rows) (*Release, error) {
	var metadataJson, uploadedBy string
	var processedDependencies bool
	var downloads int
	var uploadedAt int64
	if err := rows.Scan(&metadataJson, &processedDependencies, &downloads, &uploadedBy, &uploadedAt); err != nil {
		return nil, err
	}
	metadata, err := core.NewReleaseMetadataFromJsonString(metadataJson)
	if err != nil {
		return nil, err
	}
	rel := NewRelease(NewApplication(project, name), metadata)
	rel.ProcessedDependencies = processedDependencies
	rel.Downloads = downloads
	rel.UploadedBy = uploadedBy
	rel.UploadedAt = time.Unix(uploadedAt, 0)
	return rel, nil
}

func (s *SQLHelper) scanReleases(rows *sql.Rows) ([]*Release, error) {
	defer rows.Close()
	result := []*Release{}
	for rows.Next() {
		var project, metadataJson, uploadedBy string
		var processedDependencies bool
		var downloads int
		var uploadedAt int64
		if err := rows.Scan(&project, &metadataJson, &processedDependencies, &downloads, &uploadedBy, &uploadedAt); err != nil {
			return nil, err
		}
		metadata, err := core.NewReleaseMetadataFromJsonString(metadataJson)
		if err != nil {
			return nil, err
		}
		rel := NewRelease(NewApplication(project, metadata.Name), metadata)
		rel.ProcessedDependencies = processedDependencies
		rel.Downloads = downloads
		rel.UploadedBy = uploadedBy
		rel.UploadedAt = time.Unix(uploadedAt, 0)
		result = append(result, rel)
	}
	return result, nil
}
