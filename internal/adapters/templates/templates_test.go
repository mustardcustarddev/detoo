package templates

import (
	"context"
	"maps"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/mustardcustard/detoo/internal/domain/record"
)

var kind = record.Kind{Name: "adr"}

// builtin is a tool's files, holding only its template for a new record.
var builtin = fstest.MapFS{builtinRecord: {Data: []byte("built-in\n")}}

func write(t *testing.T, file, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// TestRecordTemplateVariable covers the TEMPLATE variable, which names a file
// either relative to the working directory or by an absolute path.
func TestRecordTemplateVariable(t *testing.T) {
	root := t.TempDir()
	workdir := filepath.Join(root, "project")
	write(t, filepath.Join(workdir, "custom.md"), "relative\n")
	elsewhere := filepath.Join(root, "elsewhere", "template.md")
	write(t, elsewhere, "absolute\n")

	for _, test := range []struct {
		name, variable, want string
	}{
		{"unset", "", "built-in\n"},
		{"relative to the working directory", "custom.md", "relative\n"},
		{"absolute", elsewhere, "absolute\n"},
	} {
		t.Run(test.name, func(t *testing.T) {
			getenv := func(name string) string {
				if name == "ADR_TEMPLATE" {
					return test.variable
				}
				return ""
			}
			got, err := New(workdir, getenv, kind, builtin).Record(context.Background(), "docs/adr")
			if err != nil {
				t.Fatal(err)
			}
			if got != test.want {
				t.Errorf("Record() = %q, want %q", got, test.want)
			}
		})
	}
}

// TestDirectory covers the files a tool's new records are created with in the
// directory beside them, which a tool need not have.
func TestDirectory(t *testing.T) {
	files := fstest.MapFS{
		builtinRecord:                            {Data: []byte("built-in\n")},
		builtinDirectory + "/problem/README.md":  {Data: []byte("problem\n")},
		builtinDirectory + "/solution/README.md": {Data: []byte("solution\n")},
	}
	got, err := New(t.TempDir(), func(string) string { return "" }, kind, files).Directory(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]string{"problem/README.md": "problem\n", "solution/README.md": "solution\n"}
	if !maps.Equal(got, want) {
		t.Errorf("Directory() = %v, want %v", got, want)
	}

	got, err = New(t.TempDir(), func(string) string { return "" }, kind, builtin).Directory(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Errorf("Directory() of a tool with none = %v, want nothing", got)
	}
}
