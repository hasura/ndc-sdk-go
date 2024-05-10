package version

import "testing"

func TestVersion(t *testing.T) {
	BuildVersion = "test"
	initBuildVersion()
	if BuildVersion != "test" {
		t.Errorf("expected BuildVersion == test, got %s", BuildVersion)
		t.FailNow()
	}
}
