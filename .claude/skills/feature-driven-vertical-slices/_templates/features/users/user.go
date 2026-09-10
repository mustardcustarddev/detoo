// Package users holds the concepts every "users" slice exchanges.
//
// Only what more than one slice genuinely shares belongs here.  Use-case
// specific input, output, rules and errors stay in the slice that owns them —
// see ./create.  Resist the pull to grow this into a service layer: it has no
// behaviour and no dependencies, and should keep neither.
package users

import "time"

// User is the domain concept, not a CLI or database representation of one.
type User struct {
	ID        string
	Name      string
	Email     string
	Role      string
	CreatedAt time.Time
}

// Roles a user may hold.
const (
	RoleMember = "member"
	RoleAdmin  = "admin"
)
