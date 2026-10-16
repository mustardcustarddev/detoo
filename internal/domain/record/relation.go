package record

// Relation is how a later record bears on an earlier one.
//
// There are only two, so that every log describes its history in the same
// terms.  What each is called belongs to the Kind of log: an ADR supersedes the
// record it replaces, and an RFC obsoletes it.
type Relation int

const (
	// Replaces retires the earlier record, which loses its status.
	Replaces Relation = iota + 1

	// Updates makes substantive changes to the earlier record, which keeps its
	// status.
	Updates
)

// Valid reports whether r is one of the relations above.  The zero Relation is
// not, so that a request which names none is recognised as such.
func (r Relation) Valid() bool { return r == Replaces || r == Updates }

// Retires reports whether the earlier record loses its status.
func (r Relation) Retires() bool { return r == Replaces }
