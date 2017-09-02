package sqlhelp

import (
	"database/sql"
	"strconv"
	"strings"

	. "github.com/ankyra/escape-registry/dao/types"
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
func (s *SQLHelper) GetDownstreamDependenciesByGroups(release *Release, readGroups []string) ([]*Dependency, error) {
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
			insertMarks = append(insertMarks, "$"+strconv.Itoa(i+4))
		} else {
			insertMarks = append(insertMarks, "?")
		}
	}
	query := s.GetDownstreamDependenciesByGroupsQuery
	if len(readGroups) == 1 {
		query += " = " + insertMarks[0]
	} else {
		query += "IN (" + strings.Join(insertMarks, ", ") + ")"
	}
	interfaceGroups := []interface{}{
		release.Application.Project,
		release.Application.Name,
		release.Version,
	}
	for _, g := range readGroups {
		interfaceGroups = append(interfaceGroups, g)
	}
	rows, err := s.PrepareAndQuery(query, interfaceGroups...)
	if err != nil {
		return nil, err
	}
	return s.scanDependencies(rows)
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
