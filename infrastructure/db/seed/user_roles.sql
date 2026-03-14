-- user_roles.sql
-- Assigns the 'user' role to ALL existing users that don't have it yet.
-- Unlike *_generate.sql files, this does NOT use seed_count parameter.
-- It automatically operates on all users in the database.
--
-- Run AFTER:
--   1. roles.sql (creates the roles)
--   2. user_generate.sql (creates the users)

-- Assign 'user' role to all users that don't have any roles yet
INSERT INTO aion_api.user_roles (user_id, role_id, assigned_at)
SELECT u.user_id, r.role_id, NOW()
FROM aion_api.users u
CROSS JOIN aion_api.roles r
WHERE r.name = 'user'
  AND r.is_active = true
  AND u.deleted_at IS NULL
  AND NOT EXISTS (
    SELECT 1 FROM aion_api.user_roles ur
    WHERE ur.user_id = u.user_id AND ur.role_id = r.role_id
  );

