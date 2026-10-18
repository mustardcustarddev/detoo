package app

import (
	"github.com/mustardcustard/detoo/internal/domain/record"
	"github.com/mustardcustard/detoo/internal/features/agents/skills"
	"github.com/mustardcustard/detoo/internal/features/health/doctor"
)

// Detoo groups the handlers of the detoo command, which looks across every
// log in a project rather than keeping one, and installs the agent skills
// that drive the tools.  Like Application, it holds wiring only.
type Detoo struct {
	Health Health
	Skills Skills
}

// Health are the capabilities that check the logs are where the tools look,
// and put them there.
type Health struct {
	Doctor doctor.Handler
	Fix    doctor.Fixer
}

// Skills are the capabilities that install agent skills, update them and
// remove them.
type Skills struct {
	Install skills.Installer
	Update  skills.Updater
	Remove  skills.Remover
}

// NewDetoo wires detoo for logs of kinds, in the project given, which git
// moves files in when a fix asks it to, and whose history orders documents
// that share a number; and for skills fetched by fetcher, installed in the
// project or for the user.
func NewDetoo(kinds []record.Kind, project doctor.Workspace, git doctor.Git, fetcher skills.Fetcher, installations skills.Installations) *Detoo {
	return &Detoo{
		Health: Health{
			Doctor: doctor.Handler{Kinds: kinds, Project: project, History: git},
			Fix:    doctor.Fixer{Kinds: kinds, Project: project, Git: git},
		},
		Skills: Skills{
			Install: skills.Installer{Fetcher: fetcher, Installations: installations},
			Update:  skills.Updater{Fetcher: fetcher, Installations: installations},
			Remove:  skills.Remover{Installations: installations},
		},
	}
}
