package sqlhelp

import (
	"database/sql"

	. "github.com/ankyra/escape-inventory/dao/types"
)

func (s *SQLHelper) GetAllReleasesWithoutProcessedDependencies() ([]*Release, error) {
	rows, err := s.PrepareAndQuery(s.GetAllReleasesWithoutProcessedDependenciesQuery)
	if err != nil {
		return nil, err
	}
	return s.scanReleases(rows)
}

func (s *SQLHelper) SetDependencies(release *Release, depends []*Dependency) error {
	for _, dep := range depends {
		err := s.PrepareAndExecInsertIgnoreDups(s.InsertDependencyQuery,
			release.Application.Project,
			release.Application.Name,
			release.Version,
			dep.Project,
			dep.Application,
			dep.Version,
			dep.BuildScope,
			dep.DeployScope,
			dep.IsExtension)
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
	return s.scanDependencies(rows)
}

func (s *SQLHelper) GetDownstreamDependencies(release *Release) ([]*Dependency, error) {
	rows, err := s.PrepareAndQuery(s.GetDownstreamDependenciesQuery, release.Application.Project, release.Application.Name, release.Version)
	if err != nil {
		return nil, err
	}
	return s.scanDependencies(rows)
}

func (s *SQLHelper) GetDownstreamDependenciesFilteredBy(release *Release, query *DownstreamDependenciesFilter) ([]*Dependency, error) {
	return nil, nil
}

func (s *SQLHelper) scanDependencies(rows *sql.Rows) ([]*Dependency, error) {
	defer rows.Close()
	result := []*Dependency{}
	for rows.Next() {
		var depProject, depApplication, depVersion string
		var buildScope, deployScope, isExtension bool
		if err := rows.Scan(&depProject, &depApplication, &depVersion, &buildScope, &deployScope, &isExtension); err != nil {
			return nil, err
		}
		result = append(result, &Dependency{
			Project:     depProject,
			Application: depApplication,
			Version:     depVersion,
			BuildScope:  buildScope,
			DeployScope: deployScope,
			IsExtension: isExtension,
		})
	}
	return result, nil
}
