package adr

import "regexp"

const upgradeHelp = `usage: adr upgrade-repository

Upgrades the ADR documents to the latest format.
`

// ukDate matches the day/month/year date format that early versions of the tool
// wrote, which ISO 8601 dates replaced.
var ukDate = regexp.MustCompile(`Date: ([0-9]{2})/([0-9]{2})/([0-9]{4})`)

func runUpgrade(e *Env, args []string) error {
	records, err := e.List()
	if err != nil {
		return err
	}
	for _, record := range records {
		lines, err := readLines(e.path(record))
		if err != nil {
			return err
		}
		for i, line := range lines {
			if match := ukDate.FindStringSubmatchIndex(line); match != nil {
				lines[i] = line[:match[0]] + string(ukDate.ExpandString(nil, "Date: $3-$2-$1", line, match)) + line[match[1]:]
			}
		}
		if err := writeLines(e.path(record), lines); err != nil {
			return err
		}
	}
	return nil
}
