// Package output renders application results.
//
// Rendering is a presentation concern, so it lives on the CLI side of the
// boundary.  A slice returns a Result; this package decides whether the user
// sees aligned text, JSON, or nothing at all.  Adding a String method to a
// Result would move that decision into the application and fix it for every
// client at once.
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"text/tabwriter"
)

type Kind string

const (
	Human Kind = "human"
	JSON  Kind = "json"
)

// Format is a pflag.Value, so --output is parsed and validated by the flag
// package rather than checked inside a RunE.
type Format struct{ Kind Kind }

func (f *Format) String() string { return string(f.Kind) }
func (f *Format) Type() string   { return "format" }

func (f *Format) Set(value string) error {
	switch Kind(value) {
	case Human, JSON:
		f.Kind = Kind(value)
		return nil
	default:
		return fmt.Errorf("unknown output format %q: want human or json", value)
	}
}

// Write renders a result to w.
//
// It takes an io.Writer rather than reaching for os.Stdout, so that a test can
// pass a buffer and a command can pass cmd.OutOrStdout().
func (f *Format) Write(w io.Writer, result any) error {
	if f.Kind == JSON {
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		return encoder.Encode(result)
	}
	return writeHuman(w, result)
}

// writeHuman lays a result out as aligned "Field: value" lines.
//
// It reads the fields reflectively rather than asking the result to describe
// itself.  That is deliberate: a Rows or String method on create.Result would
// put the layout inside the application, where an HTTP client would inherit it
// and a second CLI format could not override it.
func writeHuman(w io.Writer, result any) error {
	value := reflect.Indirect(reflect.ValueOf(result))
	if value.Kind() != reflect.Struct {
		_, err := fmt.Fprintln(w, result)
		return err
	}

	tab := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	for i := range value.NumField() {
		field := value.Type().Field(i)
		if !field.IsExported() {
			continue
		}
		if _, err := fmt.Fprintf(tab, "%s:\t%v\n", field.Name, value.Field(i).Interface()); err != nil {
			return err
		}
	}
	return tab.Flush()
}
