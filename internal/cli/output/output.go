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

	"github.com/mustardcustard/decision-records/internal/domain/adr"
)

// Path writes the name of a record, alone on a line, so that the command can be
// used in a script.
func Path(w io.Writer, file string) {
	fmt.Fprintln(w, file)
}

// Records writes a log, one record a line.
func Records(w io.Writer, records []adr.Record) {
	for _, record := range records {
		fmt.Fprintln(w, record.File)
	}
}

// Config writes the configuration as shell-style assignments, so that it can be
// read by eval as well as by a person.
func Config(w io.Writer, dir, template string) {
	fmt.Fprintf(w, "adr_dir=%q\n", dir)
	fmt.Fprintf(w, "adr_template=%q\n", template)
}

// Document writes a generated report, which already ends in a newline.
func Document(w io.Writer, text string) error {
	_, err := io.WriteString(w, text)
	return err
}
