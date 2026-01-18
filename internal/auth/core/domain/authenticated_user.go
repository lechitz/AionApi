package domain

// AuthenticatedUser is a minimal user view returned by the auth context.
//
// It intentionally does NOT depend on the /user domain entity to avoid
// bounded-context leakage. It also carries the roles snapshot used for
// token issuance.
type AuthenticatedUser struct {
	ID       uint64
	Name     string
	Username string
	Email    string
	Roles    []string
}
