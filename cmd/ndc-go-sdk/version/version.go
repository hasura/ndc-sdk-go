// Package version implements cli handling.
package version

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// DevVersion is the version string for development versions.
const DevVersion = "dev"

// BuildVersion is the version string with which CLI is built. Set during
// the build time.
var BuildVersion = DevVersion

// Version defines the version object.
type Version struct {
	// CLI is the version of CLI
	CalVer string
	year   int
	month  int
	date   int
	patch  *int
}

// GetCLIVersion return the CLI version string.
func (v *Version) GetCLIVersion() string {
	return v.CalVer
}

func (currentVersion *Version) ShouldUpdate(incomingVersion *Version) bool {
	if currentVersion.year == 0 {
		return false
	}
	current := time.Date(currentVersion.year, time.Month(currentVersion.month), currentVersion.date, 0, 0, 0, 0, time.UTC)
	incoming := time.Date(incomingVersion.year, time.Month(incomingVersion.month), incomingVersion.date, 0, 0, 0, 0, time.UTC)
	if current.Equal(incoming) {
		if incomingVersion.patch != nil {
			currentPatch := 0
			if currentVersion.patch != nil {
				currentPatch = *currentVersion.patch
			}
			return currentPatch < *incomingVersion.patch
		}
		return false
	}
	return current.Before(incoming)
}

// New returns a new version object with BuildVersion filled in as CLI
func New() *Version {
	dev := &Version{
		CalVer: BuildVersion,
		year:   0,
		month:  0,
		date:   0,
		patch:  nil,
	}
	v, err := ParseCalVer(BuildVersion)
	if err != nil {
		return dev
	}
	return v
}

// NewCLIVersion returns a new version object with CLI info filled in
func NewCLIVersion(cli string) (*Version, error) {
	return ParseCalVer(cli)
}

func ParseCalVer(s string) (*Version, error) {
	verErr := fmt.Errorf("unable to parse version %s. should be of the form year.month.day or year.month.day.patch. ex. 2023.11.20 or 2023.11.20.1", s)
	parts := strings.Split(s, ".")
	if len(parts) != 3 && len(parts) != 4 {
		return nil, verErr
	}
	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, verErr
	}
	month, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, verErr
	}
	day, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, verErr
	}
	var patch *int
	if len(parts) == 4 {
		p, err1 := strconv.Atoi(parts[3])
		if err1 != nil {
			return nil, verErr
		}
		patch = &p
	}
	return &Version{
		CalVer: s,
		year:   year,
		month:  month,
		date:   day,
		patch:  patch,
	}, nil

}
