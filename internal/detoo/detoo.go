// Package detoo defines the detoo command: the logs it looks across, and the
// help it ships with.
//
// It holds data and nothing else.  cmd/detoo puts it together with the
// behaviour, as cmd/adr does for the adr tool.
package detoo

import (
	"embed"

	"github.com/mustardcustard/detoo/internal/adr"
	"github.com/mustardcustard/detoo/internal/domain/record"
	"github.com/mustardcustard/detoo/internal/plm"
	"github.com/mustardcustard/detoo/internal/rfc"
	"github.com/mustardcustard/detoo/internal/sop"
	"github.com/mustardcustard/detoo/internal/ust"
)

// Kinds are the logs detoo looks across, in the order it reports them.
var Kinds = []record.Kind{adr.Kind, rfc.Kind, ust.Kind, plm.Kind, sop.Kind}

// Files holds what detoo's commands print about themselves, under help/.
//
//go:embed help
var Files embed.FS
