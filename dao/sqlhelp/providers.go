package sqlhelp

import (
	"database/sql"
	"strconv"
	"strings"

	core "github.com/ankyra/escape-core"
	. "github.com/ankyra/escape-inventory/dao/types"
)

func (s *SQLHelper) GetProviders(providerName string) (map[string]*MinimalReleaseMetadata, error) {
	rows, err := s.PrepareAndQuery(s.GetProviderReleasesQuery, providerName)
	if err != nil {
		return nil, err
	}
	return s.scanMinimalReleaseMetadata(rows)
}

func (s *SQLHelper) scanMinimalReleaseMetadata(rows *sql.Rows) (map[string]*MinimalReleaseMetadata, error) {
	result := map[string]*MinimalReleaseMetadata{}
	defer rows.Close()
	for rows.Next() {
		var project, application, version, description string
		if err := rows.Scan(&project, &application, &version, &description); err != nil {
			return nil, err
		}
		mini := &MinimalReleaseMetadata{
			Project:     project,
			Application: application,
			Version:     version,
			Description: description,
		}
		result[mini.GetReleaseId()] = mini
	}
	return result, nil
}

func (s *SQLHelper) GetProvidersByGroups(providerName string, groups []string) (map[string]*MinimalReleaseMetadata, error) {
	starFound := false
	for _, g := range groups {
		if g == "*" {
			starFound = true
			break
		}
	}
	if !starFound {
		groups = append(groups, "*")
	}
	insertMarks := []string{}
	for i, _ := range groups {
		if s.UseNumericInsertMarks {
			insertMarks = append(insertMarks, "$"+strconv.Itoa(i+2))
		} else {
			insertMarks = append(insertMarks, "?")
		}
	}
	query := s.GetProviderReleasesByGroupsQuery
	if len(groups) == 1 {
		query += " = " + insertMarks[0]
	} else {
		query += "IN (" + strings.Join(insertMarks, ", ") + ")"
	}
	interfaceGroups := []interface{}{providerName}
	for _, g := range groups {
		interfaceGroups = append(interfaceGroups, g)
	}
	rows, err := s.PrepareAndQuery(query, interfaceGroups...)
	if err != nil {
		return nil, err
	}
	return s.scanMinimalReleaseMetadata(rows)
}

func (s *SQLHelper) RegisterProviders(release *core.ReleaseMetadata) error {
	rows, err := s.PrepareAndQuery(s.GetProvidersForReleaseQuery, release.Project, release.Name)
	if err != nil {
		return err
	}
	defer rows.Close()
	current := map[string]string{}
	for rows.Next() {
		var provider, version string
		if err := rows.Scan(&provider, &version); err != nil {
			return err
		}
		current[provider] = version
	}

	for _, provider := range release.Provides {
		currentVersion, ok := current[provider.Name]
		if !ok {

			err := s.PrepareAndExecInsert(s.SetProviderQuery,
				release.Project,
				release.Name,
				release.Version,
				release.Description,
				provider.Name)
			if err != nil {
				return err
			}
		}

		currentV := core.NewSemanticVersion(currentVersion)
		newV := core.NewSemanticVersion(release.Version)
		if !newV.LessOrEqual(currentV) {
			err := s.PrepareAndExecUpdate(s.UpdateProviderQuery,
				release.Project,
				release.Name,
				release.Version,
				release.Description,
				provider.Name,
			)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
