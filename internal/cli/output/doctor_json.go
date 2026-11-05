package output

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/mustardcustard/detoo/internal/features/health/doctor"
)

// The commands that fix a problem that does not need fixing by hand.
const (
	fixCommand      = "detoo doctor --fix"
	renumberCommand = "detoo doctor --fix --renumber"
)

// DoctorJSON writes the health of a project's logs as JSON, for a program or
// an agent to read rather than a person: what was checked, a summary, and
// each log with its problems, each saying how to fix it.
func DoctorJSON(w io.Writer, result doctor.Result) error {
	return writeJSON(w, newJSONReport(result))
}

// DoctorFixJSON writes what a fix changed, and the health of the logs after
// it, as DoctorJSON writes it.
func DoctorFixJSON(w io.Writer, fixed doctor.Fixed) error {
	report := newJSONReport(fixed.After)
	report.Fixed = newJSONFixed(fixed)
	return writeJSON(w, report)
}

func writeJSON(w io.Writer, v any) error {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

// jsonReport is the whole report.  Every path in it but Project is slash
// separated and relative to Project.
type jsonReport struct {
	Project    string      `json:"project"`
	Healthy    bool        `json:"healthy"`
	ExitStatus int         `json:"exitStatus"`
	Summary    jsonSummary `json:"summary"`
	Fixed      *jsonFixed  `json:"fixed,omitempty"`
	Logs       []jsonLog   `json:"logs"`
}

// jsonSummary counts across every log.  Fixable counts the problems a fix
// can fix, including the shared numbers that need --renumber; ByHand those
// left to fix by hand.  Renumberable counts the documents --renumber would
// renumber.
type jsonSummary struct {
	Documents    int `json:"documents"`
	Directories  int `json:"directories"`
	Problems     int `json:"problems"`
	Fixable      int `json:"fixable"`
	ByHand       int `json:"byHand"`
	Renumberable int `json:"renumberable"`
}

type jsonLog struct {
	Tool        string         `json:"tool"`
	Noun        string         `json:"noun"`
	Prefix      string         `json:"prefix,omitempty"`
	Dir         string         `json:"dir"`
	Found       bool           `json:"found"`
	Documents   int            `json:"documents"`
	Directories []jsonDir      `json:"directories"`
	Problems    []jsonProblem  `json:"problems"`
	Renumbers   []jsonRenumber `json:"renumbers,omitempty"`
}

type jsonDir struct {
	Path      string `json:"path"`
	Documents int    `json:"documents"`
}

// jsonProblem is one problem with a log.  Type says which, and which of the
// other fields it has:
//
//   - misplaced: the log is kept in directory From, rather than To, for
//     Reason;
//   - misnamed: the document From, in Dir, should be named To;
//   - shared-number: the Documents in Dir share Number;
//   - mistitled: the document File, in Dir, has a Title that does not match
//     its name.
//
// Fix is the command that fixes it, unless ByHand, when Clash says why it
// cannot be fixed by a command.  Message says the same in a sentence or two.
type jsonProblem struct {
	Type          string   `json:"type"`
	ByHand        bool     `json:"byHand"`
	Fix           string   `json:"fix,omitempty"`
	Message       string   `json:"message"`
	Reason        string   `json:"reason,omitempty"`
	Clash         string   `json:"clash,omitempty"`
	ClashWith     string   `json:"clashWith,omitempty"`
	Dir           string   `json:"dir,omitempty"`
	File          string   `json:"file,omitempty"`
	From          string   `json:"from,omitempty"`
	To            string   `json:"to,omitempty"`
	Number        int      `json:"number,omitempty"`
	Documents     []string `json:"documents,omitempty"`
	Title         string   `json:"title,omitempty"`
	NumberDiffers bool     `json:"numberDiffers,omitempty"`
	WordsDiffer   bool     `json:"wordsDiffer,omitempty"`
}

// jsonRenumber gives the document File, in Dir, the name To: either because
// it shares its number with one added before it, or to make room for one
// that does.  After and Before are where a document that shares a number
// goes in the log.
type jsonRenumber struct {
	Dir          string `json:"dir"`
	File         string `json:"file"`
	To           string `json:"to"`
	SharesNumber bool   `json:"sharesNumber"`
	After        string `json:"after,omitempty"`
	Before       string `json:"before,omitempty"`
	ByHand       bool   `json:"byHand"`
}

// jsonFixed is what a fix changed.
type jsonFixed struct {
	Problems   int            `json:"problems"`
	Changes    int            `json:"changes"`
	Git        bool           `json:"git"`
	Moved      []jsonMoved    `json:"moved"`
	Renamed    []jsonRenamed  `json:"renamed"`
	Renumbered []jsonRenamed  `json:"renumbered"`
	Removed    []string       `json:"removed"`
	Relinked   []jsonRelinked `json:"relinked"`
}

type jsonMoved struct {
	Tool string `json:"tool"`
	From string `json:"from"`
	To   string `json:"to"`
}

// jsonRenamed is a document renamed in Dir, the directory it was in before
// the fix, from From to To.
type jsonRenamed struct {
	Tool string `json:"tool"`
	Dir  string `json:"dir"`
	From string `json:"from"`
	To   string `json:"to"`
}

type jsonRelinked struct {
	Path     string `json:"path"`
	Links    int    `json:"links"`
	Retitled bool   `json:"retitled"`
}

func newJSONReport(result doctor.Result) jsonReport {
	report := jsonReport{
		Project: result.Root,
		Healthy: result.Healthy(),
		Summary: jsonSummary{
			Documents:    result.Documents(),
			Directories:  result.Directories(),
			Problems:     result.Problems(),
			Fixable:      result.Fixable(),
			ByHand:       result.ByHand(),
			Renumberable: result.Renumberable(),
		},
		Logs: make([]jsonLog, 0, len(result.Logs)),
	}
	if !report.Healthy {
		report.ExitStatus = 1
	}
	for _, l := range result.Logs {
		report.Logs = append(report.Logs, newJSONLog(l))
	}
	return report
}

func newJSONLog(l doctor.Log) jsonLog {
	log := jsonLog{
		Tool:        l.Kind.Name,
		Noun:        l.Kind.Noun,
		Prefix:      l.Kind.Prefix,
		Dir:         l.Dir,
		Found:       l.Found(),
		Documents:   l.Documents(),
		Directories: make([]jsonDir, 0, len(l.Dirs)),
		Problems:    make([]jsonProblem, 0, l.Problems()),
	}
	for _, d := range l.Dirs {
		log.Directories = append(log.Directories, jsonDir{Path: d.Path, Documents: len(d.Documents)})
	}
	for _, m := range l.Misplaced {
		log.Problems = append(log.Problems, misplacedProblem(l, m))
	}
	for _, rn := range l.Misnamed {
		log.Problems = append(log.Problems, misnamedProblem(l, rn))
	}
	for _, d := range l.Duplicates {
		log.Problems = append(log.Problems, sharedNumberProblem(l, d))
	}
	for _, m := range l.Mistitled {
		log.Problems = append(log.Problems, mistitledProblem(l, m))
	}
	for _, rn := range l.Renumbers {
		log.Renumbers = append(log.Renumbers, jsonRenumber{
			Dir: rn.Dir, File: rn.File, To: rn.To, SharesNumber: rn.Duplicate,
			After: rn.After, Before: rn.Before, ByHand: rn.Clash != doctor.NoClash,
		})
	}
	return log
}

func misplacedProblem(l doctor.Log, m doctor.Move) jsonProblem {
	p := jsonProblem{Type: "misplaced", From: m.From, To: m.To, Fix: fixCommand}
	switch m.Reason {
	case doctor.UnderDoc:
		p.Reason = "under-doc"
		p.Message = fmt.Sprintf("The %s log is kept under doc/ instead of docs/.", l.Kind.Name)
	case doctor.Former:
		p.Reason = "former-directory"
		p.Message = fmt.Sprintf("The %s log is kept in a former directory, where %s no longer looks: %s now live in %s.",
			l.Kind.Name, l.Kind.Name, l.Kind.Plural, m.To)
	}
	switch m.Clash {
	case doctor.NoClash:
		p.Message += fmt.Sprintf(" The fix moves %s to %s.", m.From, m.To)
	case doctor.TargetExists:
		p.ByHand, p.Fix, p.Clash = true, "", "target-exists"
		p.Message += fmt.Sprintf(" Needs fixing by hand: both %s and %s exist, so moving one would clash.", m.From, m.To)
	case doctor.SharedTarget:
		p.ByHand, p.Fix, p.Clash, p.ClashWith = true, "", "shared-target", m.ClashWith
		p.Message += fmt.Sprintf(" Needs fixing by hand: %s moves to %s too, so moving both would clash.", m.ClashWith, m.To)
	}
	return p
}

func misnamedProblem(l doctor.Log, rn doctor.Rename) jsonProblem {
	p := jsonProblem{Type: "misnamed", Dir: rn.Dir, From: rn.From, To: rn.To, Fix: fixCommand,
		Message: fmt.Sprintf("The file name does not begin with %s-. The fix renames it to %s, and updates the links to it.", l.Kind.Prefix, rn.To)}
	if !rn.Clash {
		return p
	}
	p.ByHand, p.Fix, p.ClashWith = true, "", rn.ClashWith
	if rn.ClashWith == rn.To {
		p.Clash = "name-taken"
		p.Message = fmt.Sprintf("The file name does not begin with %s-. Needs fixing by hand: %s is there already.", l.Kind.Prefix, rn.To)
	} else {
		p.Clash = "same-new-name"
		p.Message = fmt.Sprintf("The file name does not begin with %s-. Needs fixing by hand: %s would be renamed %s too.", l.Kind.Prefix, rn.ClashWith, rn.To)
	}
	return p
}

func sharedNumberProblem(l doctor.Log, d doctor.Duplicate) jsonProblem {
	p := jsonProblem{Type: "shared-number", Dir: d.Dir, Number: d.Number, Documents: d.Documents}
	shared := fmt.Sprintf("%d %s share the number %d.", len(d.Documents), l.Kind.Plural, d.Number)
	switch d.Clash {
	case doctor.NoClash:
		p.Fix = renumberCommand
		p.Message = shared + fmt.Sprintf(" The first git added keeps it, and %s renumbers the others, as the log's renumbers say, and updates the links to them in the project's Markdown. Other references to them, by name or number, could be affected.", renumberCommand)
	case doctor.NoHistory:
		p.ByHand, p.Clash = true, "no-history"
		p.Message = shared + " Needs fixing by hand: this is not a git repository, so which came first is not known. Keep the first, and renumber the others, such as the log's renumbers suggest, updating the title of each, the links to them, and anything else that refers to them."
	case doctor.NamesClash:
		p.ByHand, p.Clash = true, "names-clash"
		p.Message = shared + fmt.Sprintf(" Needs fixing by hand: fix the file names that clash first, and %s can renumber these.", renumberCommand)
	}
	return p
}

func mistitledProblem(l doctor.Log, m doctor.Mistitle) jsonProblem {
	var differs []string
	if m.Number {
		differs = append(differs, "number")
	}
	if m.Words {
		differs = append(differs, "words")
	}
	return jsonProblem{
		Type: "mistitled", ByHand: true, Dir: m.Dir, File: m.File, Title: m.Title,
		NumberDiffers: m.Number, WordsDiffer: m.Words,
		Message: fmt.Sprintf("The title does not match the file name: the %s differ. Needs fixing by hand: correct the title, or rename the %s and the links to it, whichever is wrong.",
			strings.Join(differs, " and "), l.Kind.Noun),
	}
}

func newJSONFixed(fixed doctor.Fixed) *jsonFixed {
	f := &jsonFixed{
		Problems:   fixed.Count(),
		Changes:    fixed.Changes(),
		Git:        fixed.Git,
		Moved:      []jsonMoved{},
		Renamed:    []jsonRenamed{},
		Renumbered: []jsonRenamed{},
		Removed:    append([]string{}, fixed.Removed...),
		Relinked:   make([]jsonRelinked, 0, len(fixed.Relinked)),
	}
	for i, l := range fixed.Logs {
		tool := fixed.Before.Logs[i].Kind.Name
		for _, m := range l.Moved {
			f.Moved = append(f.Moved, jsonMoved{Tool: tool, From: m.From, To: m.To})
		}
		for _, rn := range l.Renamed {
			f.Renamed = append(f.Renamed, jsonRenamed{Tool: tool, Dir: rn.Dir, From: rn.From, To: rn.To})
		}
		for _, rn := range l.Renumbered {
			f.Renumbered = append(f.Renumbered, jsonRenamed{Tool: tool, Dir: rn.Dir, From: rn.File, To: rn.To})
		}
	}
	for _, r := range fixed.Relinked {
		f.Relinked = append(f.Relinked, jsonRelinked{Path: r.Path, Links: r.Links, Retitled: r.Retitled})
	}
	return f
}
