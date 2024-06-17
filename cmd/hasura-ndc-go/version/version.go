// Package version implements cli handling.
package version

import (
	"runtime/debug"
)

// DevVersion is the version string for development versions.
const DevVersion = "v0"

// BuildVersion is the version string with which CLI is built. Set during
// the build time.
var BuildVersion = ""

func init() {
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

	for _, s := range bi.Settings {
		if s.Key == "vcs.revision" && s.Value != "" {
			BuildVersion = s.Value
			return
		}
	}
}
