// Package version implements cli handling.
package version

import (
	"runtime/debug"
)

// DevVersion is the version string for development versions.
const DevVersion = "latest"

// BuildVersion is the version string with which CLI is built. Set during
// the build time.
var BuildVersion = ""

func init() { //nolint:gochecknoinits
	initBuildVersion()
}

func initBuildVersion() {
	if BuildVersion != "" {
		return
	}

	BuildVersion = DevVersion

	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}

	if bi.Main.Version != "" {
		BuildVersion = bi.Main.Version

		return
	}

	for _, s := range bi.Settings {
		if s.Key == "vcs.revision" && s.Value != "" {
			BuildVersion = s.Value

			return
		}
	}
}
