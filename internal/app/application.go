// Package app is the composition root: it constructs every slice handler and
// gives each one the adapters it depends on.
//
// This is the one package that knows the whole application, and it grows a
// field per use case.  That is what a composition root is for — and it is also
// how such a type turns into a god object, so the rule here is that Application
// holds wiring only.  It has no methods, and must not acquire any: the moment
// behaviour hangs off it, a caller that wants one capability starts depending
// on all of them.
package app

import (
	"github.com/mustardcustard/decision-records/internal/domain/record"
	"github.com/mustardcustard/decision-records/internal/features/records/config"
	"github.com/mustardcustard/decision-records/internal/features/records/create"
	"github.com/mustardcustard/decision-records/internal/features/records/initialise"
	"github.com/mustardcustard/decision-records/internal/features/records/link"
	"github.com/mustardcustard/decision-records/internal/features/records/list"
	"github.com/mustardcustard/decision-records/internal/features/reports/graph"
	"github.com/mustardcustard/decision-records/internal/features/reports/toc"
)

// Application groups the slice handlers by capability, so that a caller reads
// as application.Records.Create.Handle(ctx, cmd).
type Application struct {
	Records Records
	Reports Reports
}

// Records are the capabilities that read and write the log.
type Records struct {
	Initialise initialise.Handler
	Create     create.Handler
	Link       link.Handler
	List       list.Handler
	Config     config.Handler
}

// Reports are the capabilities that describe the log as a document.
type Reports struct {
	TOC   toc.Handler
	Graph graph.Handler
}

// Log is everything the slices ask of a record store, gathered so that New has
// a signature.  No slice refers to it: each declares the two or three methods
// it uses, and the store adapter satisfies all of them at once without naming
// any.
type Log interface {
	initialise.Store
	create.Store
	link.Store
	list.Store
	config.Store
	toc.Store
	graph.Store
}

// Templates is likewise the union of what the slices ask of a template source.
type Templates interface {
	create.Templates
	initialise.Templates
	config.Templates
}

// Clock is the union of what the slices ask of a clock.
type Clock interface {
	create.Clock
	initialise.Clock
}

// New wires the application for logs of kind.  Its other parameters are the
// adapters; below this line nothing knows that records are Markdown files on a
// disk.
func New(kind record.Kind, log Log, templates Templates, clock Clock, editor create.Editor) *Application {
	return &Application{
		Records: Records{
			Initialise: initialise.Handler{Kind: kind, Log: log, Templates: templates, Clock: clock},
			Create:     create.Handler{Kind: kind, Log: log, Templates: templates, Clock: clock, Editor: editor},
			Link:       link.Handler{Kind: kind, Log: log},
			List:       list.Handler{Log: log},
			Config:     config.Handler{Log: log, Templates: templates},
		},
		Reports: Reports{
			TOC:   toc.Handler{Kind: kind, Log: log},
			Graph: graph.Handler{Kind: kind, Log: log},
		},
	}
}
