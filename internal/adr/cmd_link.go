package adr

const linkHelp = `usage: adr link SOURCE LINK TARGET REVERSE-LINK

Creates a link between two ADRs, from SOURCE to TARGET new.
SOURCE and TARGET are both a reference (number or partial filename) to an ADR
LINK is the description of the link created in the SOURCE.
REVERSE-LINK is the description of the link created in the TARGET

E.g. to create link ADR 12 to ADR 10

    adr link 12 Amends 10 "Amended by"
`

func runLink(e *Env, args []string) error {
	if len(args) != 4 {
		return &UsageError{Usage: "adr link SOURCE LINK TARGET REVERSE-LINK"}
	}
	source, forward, target, reverse := args[0], args[1], args[2], args[3]
	if err := e.AddLink(source, forward, target); err != nil {
		return err
	}
	return e.AddLink(target, reverse, source)
}
