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

// Health are the capabilities that check the logs are where the tools look,
// and put them there.
type Health struct {
	Doctor doctor.Handler
	Fix    doctor.Fixer
}

// NewDetoo wires detoo for logs of kinds, in the project given, which git
// moves files in when a fix asks it to.
func NewDetoo(kinds []record.Kind, project doctor.Workspace, git doctor.Git) *Detoo {
	return &Detoo{
		Health: Health{
			Doctor: doctor.Handler{Kinds: kinds, Project: project},
			Fix:    doctor.Fixer{Kinds: kinds, Project: project, Git: git},
		},
	}
}
