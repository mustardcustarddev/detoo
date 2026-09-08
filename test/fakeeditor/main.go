// Command fakeeditor stands in for a text editor in the tests: it reports how it
// was invoked instead of editing anything.
package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println(strings.Join(os.Args[1:], " "))
}
