-- roles.sql
-- Seeds the roles table with default system roles.
-- This should run AFTER the migration creates the roles table,
-- but BEFORE user_roles.sql populates user-role assignments.

-- Insert default roles with hierarchy (idempotent - won't duplicate on re-run)
INSERT INTO aion_api.roles (name, description, is_active)
VALUES
    ('owner', 'System owner with highest privileges', true),
    ('admin', 'Administrator with full system access', true),
    ('user', 'Default user role with basic access', true),
    ('blocked', 'Blocked user with no access', true)
ON CONFLICT (name) DO NOTHING;

