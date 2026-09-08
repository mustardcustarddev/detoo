package adr

import (
	"os"
	"path/filepath"
	"testing"
)

func TestADRDirDefaultsWhenNothingIsFound(t *testing.T) {
	e := testEnv(t)
	dir, err := e.ADRDir()
	if err != nil {
		t.Fatal(err)
	}
	if dir != DefaultDir {
		t.Errorf("ADRDir() = %q, want %q", dir, DefaultDir)
	}
}

func TestADRDirFindsAMarkerFile(t *testing.T) {
	e := testEnv(t)
	if err := os.WriteFile(filepath.Join(e.Dir, DirMarkerFile), []byte("architecture-log\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	dir, err := e.ADRDir()
	if err != nil {
		t.Fatal(err)
	}
	if want := "architecture-log"; dir != want {
		t.Errorf("ADRDir() = %q, want %q", dir, want)
	}
}

func TestADRDirSearchesAncestors(t *testing.T) {
	for _, test := range []struct {
		name  string
		setup func(root string)
		want  string
	}{
		{
			name:  "default directory in an ancestor",
			setup: func(root string) { mkdirAll(t, filepath.Join(root, "doc", "adr")) },
			want:  "../../doc/adr",
		},
		{
			name:  "marker file in an ancestor",
			setup: func(root string) { write(t, filepath.Join(root, DirMarkerFile), "architecture-log\n") },
			want:  "../../architecture-log",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			test.setup(root)
			deep := filepath.Join(root, "src", "main")
			mkdirAll(t, deep)
			e := testEnv(t)
			e.Dir = deep

			dir, err := e.ADRDir()
			if err != nil {
				t.Fatal(err)
			}
			if dir != test.want {
				t.Errorf("ADRDir() = %q, want %q", dir, test.want)
			}
		})
	}
}

// TestADRDirPrefersTheNearestMarker checks that a marker file closer to the
// working directory wins over one further up.
func TestADRDirPrefersTheNearestMarker(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, DirMarkerFile), "outer\n")
	inner := filepath.Join(root, "component")
	mkdirAll(t, inner)
	write(t, filepath.Join(inner, DirMarkerFile), "inner\n")

	e := testEnv(t)
	e.Dir = inner
	dir, err := e.ADRDir()
	if err != nil {
		t.Fatal(err)
	}
	if want := "inner"; dir != want {
		t.Errorf("ADRDir() = %q, want %q", dir, want)
	}
}

func mkdirAll(t *testing.T, dir string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
}

func write(t *testing.T, file, content string) {
	t.Helper()
	if err := os.WriteFile(file, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
