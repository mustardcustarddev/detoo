package list

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/mustardcustard/decision-records/internal/adapters/memlog"
	"github.com/mustardcustard/decision-records/internal/domain/record"
)

func TestHandleReturnsTheRecordsInOrder(t *testing.T) {
	log := memlog.WithRecords(map[string]string{
		"docs/adr/0002-second.md": "# 2. Second\n",
		"docs/adr/0001-first.md":  "# 1. First\n",
		"docs/adr/README.md":      "not a record\n",
	})

	result, err := Handler{Log: log}.Handle(context.Background(), Query{})
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}

	want := []record.Record{
		{Number: 1, File: "docs/adr/0001-first.md"},
		{Number: 2, File: "docs/adr/0002-second.md"},
	}
	if !reflect.DeepEqual(result.Records, want) {
		t.Errorf("Records = %v, want %v", result.Records, want)
	}
}

// A project with no log is a condition the caller recognises, not an opaque
// filesystem error, so the store's *record.MissingDirError reaches it intact.
func TestHandleReportsAMissingLog(t *testing.T) {
	_, err := Handler{Log: memlog.New()}.Handle(context.Background(), Query{})

	var missing *record.MissingDirError
	if !errors.As(err, &missing) {
		t.Fatalf("err = %v, want an *record.MissingDirError", err)
	}
	if missing.Dir != "docs/adr" {
		t.Errorf("missing directory = %q, want docs/adr", missing.Dir)
	}
}
