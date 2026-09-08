package adr

import (
	"reflect"
	"testing"
)

func TestParseOptions(t *testing.T) {
	for _, test := range []struct {
		name     string
		args     []string
		wantOpts []string
		wantRest []string
	}{
		{
			name:     "a value in the next argument",
			args:     []string{"-s", "12", "Use", "PostgreSQL"},
			wantOpts: []string{"s=12"},
			wantRest: []string{"Use", "PostgreSQL"},
		},
		{
			name:     "a value attached to its option",
			args:     []string{"-s12", "Use", "PostgreSQL"},
			wantOpts: []string{"s=12"},
			wantRest: []string{"Use", "PostgreSQL"},
		},
		{
			name:     "an option given more than once",
			args:     []string{"-s", "3", "-s", "4", "Title"},
			wantOpts: []string{"s=3", "s=4"},
			wantRest: []string{"Title"},
		},
		{
			name:     "a double hyphen ends the options",
			args:     []string{"--", "-Bar-"},
			wantRest: []string{"-Bar-"},
		},
		{
			name:     "an operand ends the options",
			args:     []string{"Title", "-s", "3"},
			wantRest: []string{"Title", "-s", "3"},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			var opts []string
			rest, err := parseOptions(test.args, "sl", func(option byte, value string) error {
				opts = append(opts, string(option)+"="+value)
				return nil
			})
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(opts, test.wantOpts) {
				t.Errorf("options = %q, want %q", opts, test.wantOpts)
			}
			if !reflect.DeepEqual(rest, test.wantRest) {
				t.Errorf("operands = %q, want %q", rest, test.wantRest)
			}
		})
	}
}

func TestParseOptionsRejectsBadOptions(t *testing.T) {
	for _, args := range [][]string{{"-x", "1"}, {"-s"}} {
		if _, err := parseOptions(args, "sl", func(byte, string) error { return nil }); err == nil {
			t.Errorf("parseOptions(%q) succeeded, want an error", args)
		}
	}
}
