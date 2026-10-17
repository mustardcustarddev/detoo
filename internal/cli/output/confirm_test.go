package output

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

// Enter takes the answer chosen to begin with; y and n answer directly.
func TestAskTakesTheDefaultOnEnter(t *testing.T) {
	for _, test := range []struct {
		keys      string
		yes, want bool
	}{
		{"\r", false, false},
		{"\r", true, true},
		{"y", false, true},
		{"n", true, false},
	} {
		done := make(chan struct{})
		var got bool
		var err error
		var out bytes.Buffer
		go func() {
			defer close(done)
			got, err = Ask(strings.NewReader(test.keys), &out, Terminal{Width: 80}, Question{Title: "Renumber?", Yes: test.yes})
		}()
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatalf("keys %q: Ask did not return", test.keys)
		}
		if err != nil || got != test.want {
			t.Errorf("keys %q, default %v = %v, %v, want %v", test.keys, test.yes, got, err, test.want)
		}
	}
}
