// Package input defines commands for admin operations.
package input

// PromoteToAdminCommand promotes a user to admin role.
type PromoteToAdminCommand struct {
	UserID      uint64   // Target user to promote
	ActorUserID uint64   // Who is performing the action
	ActorRoles  []string // Roles of the actor (for authorization)
}

// DemoteFromAdminCommand removes admin role from a user.
type DemoteFromAdminCommand struct {
	UserID      uint64   // Target user to demote
	ActorUserID uint64   // Who is performing the action
	ActorRoles  []string // Roles of the actor (for authorization)
}

// BlockUserCommand blocks a user.
type BlockUserCommand struct {
	UserID      uint64   // Target user to block
	ActorUserID uint64   // Who is performing the action
	ActorRoles  []string // Roles of the actor (for authorization)
}

// UnblockUserCommand unblocks a user.
type UnblockUserCommand struct {
	UserID      uint64   // Target user to unblock
	ActorUserID uint64   // Who is performing the action
	ActorRoles  []string // Roles of the actor (for authorization)
}
