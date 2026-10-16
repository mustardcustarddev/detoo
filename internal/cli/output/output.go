// Package output renders the results the slices return.
//
// This is where presentation lives.  A slice returns a Result describing what
// happened; how it is laid out for a terminal — a bare path a script can read,
// or shell-style key="value" pairs — is decided here, on the CLI side of the
// boundary, and nowhere else.
package output

import (
	"fmt"
	"io"

	"github.com/mustardcustard/decision-records/internal/domain/record"
)

// Path writes the name of a record, alone on a line, so that the command can be
// used in a script.
func Path(w io.Writer, file string) {
	fmt.Fprintln(w, file)
}

// Records writes a log, one record a line.
func Records(w io.Writer, records []record.Record) {
	for _, r := range records {
		fmt.Fprintln(w, r.File)
	}
}

// Config writes the configuration of the tool called name as shell-style
// assignments, so that it can be read by eval as well as by a person.  Each
// variable is named after the tool, so that the configurations of two tools can
// be evaluated together.
func Config(w io.Writer, name, dir, template string) {
	fmt.Fprintf(w, "%s_dir=%q\n", name, dir)
	fmt.Fprintf(w, "%s_template=%q\n", name, template)
}

// Document writes a generated report, which already ends in a newline.
func Document(w io.Writer, text string) error {
	_, err := io.WriteString(w, text)
	return err
}
