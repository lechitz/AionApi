// Package domain contains role hierarchy and authorization logic.
package domain

func roleLevel(role string) (int, bool) {
	switch role {
	case RoleOwner:
		return 100, true
	case RoleAdmin:
		return 50, true
	case RoleUser:
		return 10, true
	case RoleBlocked:
		return 0, true
	default:
		return 0, false
	}
}

// CanManageRole checks if actorRole can manage targetRole based on hierarchy.
// Rule: An actor can only manage roles with LOWER privilege level.
func CanManageRole(actorRole, targetRole string) bool {
	actorLevel, actorExists := roleLevel(actorRole)
	targetLevel, targetExists := roleLevel(targetRole)

	if !actorExists || !targetExists {
		return false // Invalid roles
	}

	// Actor must have HIGHER privilege than target
	return actorLevel > targetLevel
}

// HasRole checks if a user has a specific role.
func HasRole(userRoles []string, role string) bool {
	for _, r := range userRoles {
		if r == role {
			return true
		}
	}
	return false
}

// GetHighestRole returns the highest privilege role from a list of roles.
func GetHighestRole(roles []string) string {
	highestRole := RoleBlocked
	highestLevel := 0

	for _, role := range roles {
		if level, exists := roleLevel(role); exists && level > highestLevel {
			highestLevel = level
			highestRole = role
		}
	}

	return highestRole
}
