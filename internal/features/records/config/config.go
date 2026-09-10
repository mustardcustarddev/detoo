// Package config reports the configuration the tool is running under.
package config

import "context"

// Query asks for the configuration in use.
type Query struct{}

// Result is where this project keeps its decision records, and which template
// it creates them from.
type Result struct {
	// Dir is the decision record directory, relative to the working directory.
	Dir string

	// Template is the project's own record template, or "" if it has none and
	// the built-in template is used.
	Template string
}

// Store is the port this slice needs to find the log.
type Store interface {
	Dir(ctx context.Context) (string, error)
}

// Templates is the port this slice needs to find the project's template.
type Templates interface {
	ProjectFile(ctx context.Context, dir string) (string, error)
}

// Handler reports the configuration.
type Handler struct {
	Log       Store
	Templates Templates
}

func (h Handler) Handle(ctx context.Context, _ Query) (Result, error) {
	dir, err := h.Log.Dir(ctx)
	if err != nil {
		return Result{}, err
	}
	template, err := h.Templates.ProjectFile(ctx, dir)
	if err != nil {
		return Result{}, err
	}
	return Result{Dir: dir, Template: template}, nil
}
