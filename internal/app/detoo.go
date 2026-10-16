package app

import (
	"github.com/mustardcustard/detoo/internal/domain/record"
	"github.com/mustardcustard/detoo/internal/features/health/doctor"
)

// Detoo groups the handlers of the detoo command, which looks across every
// log in a project rather than keeping one.  Like Application, it holds
// wiring only.
type Detoo struct {
	Health Health
}

// Health are the capabilities that check the logs are where the tools look.
type Health struct {
	Doctor doctor.Handler
}

// NewDetoo wires detoo for logs of kinds, in the project given.
func NewDetoo(kinds []record.Kind, project doctor.Project) *Detoo {
	return &Detoo{
		Health: Health{
			Doctor: doctor.Handler{Kinds: kinds, Project: project},
		},
	}
}
