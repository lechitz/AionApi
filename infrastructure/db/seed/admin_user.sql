-- admin_user.sql
-- Seeds an admin user 'aion' with role 'admin'.
-- This should run AFTER roles.sql.
-- Password is 'testpassword123' (bcrypt hash provided via variable).

\if :{?user_seed_password_hash}
  SELECT CASE
    WHEN :'user_seed_password_hash' = '' THEN 'true'
    ELSE 'false'
  END AS user_seed_hash_empty;
  \gset
  \if :user_seed_hash_empty
    \echo 'ERROR: user_seed_password_hash is empty. Provide a bcrypt hash (USER_TOKEN_TEST).'
    \quit 1
  \endif
\else
  \echo 'ERROR: user_seed_password_hash not provided. Export USER_TOKEN_TEST or pass -v user_seed_password_hash=<hash>.'
  \quit 1
\endif

-- Create admin user 'aion' (with onboarding already completed)
INSERT INTO aion_api.users (name, username, password, email, onboarding_completed)
VALUES ('Aion Owner', 'aion', :'user_seed_password_hash', 'aion@aion.com', TRUE)
ON CONFLICT (username) DO NOTHING;

-- Assign 'owner' role to 'aion' user (highest privilege)
INSERT INTO aion_api.user_roles (user_id, role_id, assigned_at)
SELECT u.user_id, r.role_id, NOW()
FROM aion_api.users u
CROSS JOIN aion_api.roles r
WHERE u.username = 'aion'
  AND r.name = 'owner'
  AND r.is_active = true
  AND NOT EXISTS (
    SELECT 1 FROM aion_api.user_roles ur
    WHERE ur.user_id = u.user_id AND ur.role_id = r.role_id
  );

-- Assign 'admin' role to 'aion' (owner has all roles)
INSERT INTO aion_api.user_roles (user_id, role_id, assigned_at)
SELECT u.user_id, r.role_id, NOW()
FROM aion_api.users u
CROSS JOIN aion_api.roles r
WHERE u.username = 'aion'
  AND r.name = 'admin'
  AND r.is_active = true
  AND NOT EXISTS (
    SELECT 1 FROM aion_api.user_roles ur
    WHERE ur.user_id = u.user_id AND ur.role_id = r.role_id
  );

-- Also assign 'user' role to 'aion' (complete access)
INSERT INTO aion_api.user_roles (user_id, role_id, assigned_at)
SELECT u.user_id, r.role_id, NOW()
FROM aion_api.users u
CROSS JOIN aion_api.roles r
WHERE u.username = 'aion'
  AND r.name = 'user'
  AND r.is_active = true
  AND NOT EXISTS (
    SELECT 1 FROM aion_api.user_roles ur
    WHERE ur.user_id = u.user_id AND ur.role_id = r.role_id
  );
