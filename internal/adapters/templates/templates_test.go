package templates

import (
	"context"
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
