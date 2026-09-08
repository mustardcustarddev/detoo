package adr

const listHelp = `usage: adr list

Lists the architecture decision records
`

func runList(e *Env, args []string) error {
	records, err := e.List()
	if err != nil {
		return err
	}
	for _, record := range records {
		e.print(record)
	}
	return nil
}
