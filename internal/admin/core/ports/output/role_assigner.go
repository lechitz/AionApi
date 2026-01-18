// Package output defines interfaces for admin operations that other contexts can use.
package output

import "context"

// RoleAssigner is the interface for assigning roles to users.
// This interface allows other bounded contexts (like /user) to delegate
// role assignment to the /admin context without knowing implementation details.
type RoleAssigner interface {
	// AssignDefaultRole assigns the default 'user' role to a newly created user.
	// This is called by UserRepository.Create() to delegate role management to admin context.
	AssignDefaultRole(ctx context.Context, userID uint64) error
}
