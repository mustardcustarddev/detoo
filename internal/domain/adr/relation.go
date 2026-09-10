package adr

// Relation is how a later decision record bears on an earlier one.
//
// There are only two, so that every log describes its history in the same
// words.
type Relation int

const (
	// Supercedes replaces the earlier decision, which is no longer accepted.
	Supercedes Relation = iota + 1

	// Updates makes substantive changes to the earlier decision, which remains
	// in force.
	Updates
)

// relations describes each Relation.  The zero Relation has no entry, so that a
// request which names none is recognised as such.
var relations = map[Relation]struct {
	forward, reverse string
	retires          bool
}{
	Supercedes: {forward: "Supercedes", reverse: "Superceded by", retires: true},
	Updates:    {forward: "Updates", reverse: "Updated by"},
}

// Valid reports whether r is one of the relations above.
func (r Relation) Valid() bool {
	_, ok := relations[r]
	return ok
}

// Forward is how the later record describes the relation.
func (r Relation) Forward() string { return relations[r].forward }

// Reverse is how the earlier record describes it.  Each ends in " by", which is
// how the graph report tells a reverse link from the forward one it repeats.
func (r Relation) Reverse() string { return relations[r].reverse }

// Retires reports whether the earlier record loses its accepted status.
func (r Relation) Retires() bool { return relations[r].retires }
