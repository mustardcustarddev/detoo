// Package adrtest drives the adr command the way a user does — a sequence of
// commands in a working directory — and compares the whole transcript against a
// golden file.  These are the tests the shell implementation of adr-tools ran,
// ported along with the tool.
package adrtest

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

var update = flag.Bool("update", false, "rewrite the golden files from the transcripts produced by this run")

// Paths of the binaries the tests drive, built once by TestMain.
var (
	adrBin        string
	fakeEditorBin string
)

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(runTests(m))
}

func runTests(m *testing.M) int {
	dir, err := os.MkdirTemp("", "adr-tools-test")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	defer os.RemoveAll(dir)

	adrBin = filepath.Join(dir, "adr")
	fakeEditorBin = filepath.Join(dir, "fakeeditor")
	binaries := map[string]string{
		"github.com/mustardcustard/decision-records/cmd/adr":         adrBin,
		"github.com/mustardcustard/decision-records/test/fakeeditor": fakeEditorBin,
	}
	for pkg, bin := range binaries {
		build := exec.Command("go", "build", "-o", bin, pkg)
		build.Stdout, build.Stderr = os.Stderr, os.Stderr
		if err := build.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "building %s: %v\n", pkg, err)
			return 1
		}
	}
	return m.Run()
}
