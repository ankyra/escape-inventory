package shared

import (
	"fmt"
	"regexp"
	"strings"
)

type ReleaseId struct {
	Type    string
	Name    string
	Version string
}

func ParseReleaseId(releaseId string) (*ReleaseId, error) {
	split := strings.Split(releaseId, "-")
	if len(split) < 3 { // type-build-version
		return nil, fmt.Errorf("Invalid release format: %s", releaseId)
	}
	result := &ReleaseId{}
	result.Type = split[0]
	result.Name = strings.Join(split[1:len(split)-1], "-")

	version := split[len(split)-1]
	if version == "latest" || version == "@" || version == "v@" {
		result.Version = "latest"
	} else if strings.HasPrefix(version, "v") {
		result.Version = version[1:]
	} else {
		return nil, fmt.Errorf("Invalid version string in release ID '%s': %s", releaseId, version)
	}

	if err := result.Validate(); err != nil {
		return nil, fmt.Errorf("Invalid release ID '%s': %s", releaseId, err.Error())
	}
	return result, nil
}

func (r *ReleaseId) Validate() error {
	return validateVersion(r.Version)
}

func validateVersion(version string) error {
	if version == "latest" {
		return nil
	}
	re := regexp.MustCompile(`^[0-9]+(\.[0-9]+)*(\.@)?$`)
	matches := re.Match([]byte(version))
	if !matches {
		return fmt.Errorf("Invalid version format: %s", version)
	}
	return nil
}

func (r *ReleaseId) ToString() string {
	version := r.Version
	if version != "latest" {
		version = "v" + version
	}
	return r.Type + "-" + r.Name + "-" + version
}

func (r *ReleaseId) NeedsResolving() bool {
	return r.Version == "latest" || strings.HasSuffix(r.Version, ".@")
}
