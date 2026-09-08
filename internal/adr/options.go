package adr

import (
	"fmt"
	"strings"
)

// parseOptions parses the options at the front of args and returns the operands
// that follow.  Every option in spec takes a value, given either attached to the
// option or as the next argument.  The first operand, or "--", ends the options.
func parseOptions(args []string, spec string, set func(option byte, value string) error) ([]string, error) {
	i := 0
	for ; i < len(args); i++ {
		arg := args[i]
		if arg == "--" {
			i++
			break
		}
		if len(arg) < 2 || !strings.HasPrefix(arg, "-") {
			break
		}
		option := arg[1]
		if !strings.ContainsRune(spec, rune(option)) {
			return nil, fmt.Errorf("unknown option -%c", option)
		}
		value := arg[2:]
		if value == "" {
			if i++; i >= len(args) {
				return nil, fmt.Errorf("option -%c requires an argument", option)
			}
			value = args[i]
		}
		if err := set(option, value); err != nil {
			return nil, err
		}
	}
	return args[i:], nil
}
