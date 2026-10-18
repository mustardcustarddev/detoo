package tool

import (
	"runtime/debug"
	"testing"
)

// A release reports its version, without the tag's v.  Any other build
// reports the commit it was built from, marked dirty if it had changes, since
// its pseudo-version is long and says little more.
func TestVersionOf(t *testing.T) {
	const revision = "eb82bb5ed344d8c1b2f6a9e0c3d4b5a6f7e8d9c0"
	tests := []struct {
		name     string
		version  string
		settings map[string]string
		want     string
	}{
		{"release", "v1.0.0", map[string]string{"vcs.revision": revision, "vcs.modified": "false"}, "1.0.0"},
		{"release with changes", "v1.0.0+dirty", map[string]string{"vcs.revision": revision, "vcs.modified": "true"}, "eb82bb5-dirty"},
		{"commit", "v0.0.0-20261018151304-eb82bb5ed344", map[string]string{"vcs.revision": revision, "vcs.modified": "false"}, "eb82bb5"},
		{"commit after a release", "v1.0.1-0.20261018151304-eb82bb5ed344", map[string]string{"vcs.revision": revision}, "eb82bb5"},
		{"commit with changes", "v0.0.0-20261018151304-eb82bb5ed344+dirty", map[string]string{"vcs.revision": revision, "vcs.modified": "true"}, "eb82bb5-dirty"},
		{"module without version control", "v0.0.0-20261018151304-eb82bb5ed344", nil, "v0.0.0-20261018151304-eb82bb5ed344"},
		{"unknown", "(devel)", nil, "dev"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := &debug.BuildInfo{Main: debug.Module{Version: tt.version}}
			for key, value := range tt.settings {
				info.Settings = append(info.Settings, debug.BuildSetting{Key: key, Value: value})
			}
			if got := versionOf(info); got != tt.want {
				t.Errorf("versionOf(%q) = %q, want %q", tt.version, got, tt.want)
			}
		})
	}
}
