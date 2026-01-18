-- Rollback: 000002_users_and_roles

DROP INDEX IF EXISTS aion_api.idx_user_roles_role_id;
DROP INDEX IF EXISTS aion_api.idx_user_roles_user_id;
DROP TABLE IF EXISTS aion_api.user_roles;
DROP TABLE IF EXISTS aion_api.users;
DROP TABLE IF EXISTS aion_api.roles;
