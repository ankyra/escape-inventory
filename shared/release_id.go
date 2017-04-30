package shared

import (
	"errors"
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
		return nil, errors.New("Invalid release format: " + releaseId)
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
		return nil, errors.New("Invalid version string in release ID '" + releaseId + "': " + version)
	}
	return result, nil
}

func (r *ReleaseId) ToString() string {
    version := r.Version
    if version != "latest" {
        version = "v" + version
    }
    return r.Type + "-" + r.Name + "-" + version
}
