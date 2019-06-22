package version

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ErrInvalidVersion is an error that can be returned when trying to parse a string that does not represents a version.
var ErrInvalidVersion = errors.New("invalid version")

// version is a type that represents a version (mostly) following the Semantic Versioning specification 2.0.0.
// See: https://semver.org
type Version struct {
	Major, Minor, Patch int
}

// Parse returns a version type from a string that represents a version according to the .
func Parse(v string) (Version, error) {
	if len(v) < 5 { // Minimum size must be 5 characters: 0.0.0
		return Version{}, ErrInvalidVersion
	}

	if v[0] == 'v' {
		v = v[1:]
	}

	parts := strings.Split(v, ".")
	if len(parts) != 3 { // Must have all its part explicit: major, minor and patch
		return Version{}, ErrInvalidVersion
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return Version{}, err
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return Version{}, err
	}

	patch, err := strconv.Atoi(parts[3])
	if err != nil {
		return Version{}, err
	}

	return Version{
		Major: major,
		Minor: minor,
		Patch: patch,
	}, nil
}

// String returns a string representation of the version provided
func String(v Version) string {
	return fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// Compare returns:
//
// 0 if the two versions are equal
//
// +1 if v1 > v2
//
// -1 if v2 > v1
func Compare(v1, v2 Version) int {
	if v1.Major != v2.Major {
		return v1.Major - v2.Major
	}
	if v1.Minor != v2.Minor {
		return v1.Minor - v2.Minor
	}
	if v1.Patch != v2.Patch {
		return v1.Patch - v2.Patch
	}
	return 0
}

// Newer returns true if the first version (v1) is newer than the second (v2)
func Newer(v1, v2 Version) bool {
	return Compare(v1, v2) > 0
}

// Older returns true if the first version (v1) is older than the second (v2)
func Older(v1, v2 Version) bool {
	return Compare(v1, v2) < 0
}

// Equal returns true if the first version (v1) is equal to the second (v2)
func Equal(v1, v2 Version) bool {
	return Compare(v1, v2) == 0
}
