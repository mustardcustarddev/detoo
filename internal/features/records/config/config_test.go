package config

import (
	"context"
	"testing"

	"github.com/mustardcustard/decision-records/internal/adapters/memlog"
)

// stubTemplates reports the project's own template, or none.
type stubTemplates struct{ file string }

func (s stubTemplates) ProjectFile(context.Context, string) (string, error) { return s.file, nil }

func TestHandleReportsTheLogDirectoryAndTemplate(t *testing.T) {
	log := memlog.New()
	if err := log.SetDir(context.Background(), "architecture-log"); err != nil {
		t.Fatal(err)
	}

	result, err := Handler{
		Log:       log,
		Templates: stubTemplates{file: "architecture-log/templates/template.md"},
	}.Handle(context.Background(), Query{})
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}

	if want := "architecture-log"; result.Dir != want {
		t.Errorf("Dir = %q, want %q", result.Dir, want)
	}
	if want := "architecture-log/templates/template.md"; result.Template != want {
		t.Errorf("Template = %q, want %q", result.Template, want)
	}
}

// A project with no template of its own reports none, rather than the built-in
// one, which has no path to report.
func TestHandleReportsNoTemplateWhenTheProjectHasNone(t *testing.T) {
	result, err := Handler{Log: memlog.New(), Templates: stubTemplates{}}.Handle(context.Background(), Query{})
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}
	if result.Template != "" {
		t.Errorf("Template = %q, want an empty string", result.Template)
	}
	if want := "docs/adr"; result.Dir != want {
		t.Errorf("Dir = %q, want %q", result.Dir, want)
	}
}
