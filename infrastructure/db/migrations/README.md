# Database Migrations

This directory contains migration files managed by [golang-migrate](https://github.com/golang-migrate/migrate).

## Structure

```
migrations/
├── 000001_initial_schema.up.sql       # Base schema and utility functions
├── 000001_initial_schema.down.sql
├── 000002_users_and_roles.up.sql      # Users, roles, user_roles
├── 000002_users_and_roles.down.sql
├── 000003_categories_tags_days.up.sql # Categories, tags, days
├── 000003_categories_tags_days.down.sql
├── 000004_diaries_and_day_data.up.sql # Diaries and day_* tables
├── 000004_diaries_and_day_data.down.sql
├── 000005_records.up.sql              # Records with full-text search
├── 000005_records.down.sql
├── 000006_chat_history.up.sql         # Chat history for AI
├── 000006_chat_history.down.sql
└── README.md
```

## Naming Convention

- `{version}_{descriptive_name}.up.sql` - Apply migration
- `{version}_{descriptive_name}.down.sql` - Rollback migration

Version uses 6 digits: `000001`, `000002`, etc.

## Commands

### Apply all migrations
```bash
make migrate-up
```

### Apply migrations in dev environment (via Docker)
```bash
make migrate-dev-up
```

### Rollback the last migration
```bash
make migrate-down
```

### Create a new migration
```bash
make migrate-new
# You will be prompted for the migration name
```

### Force a specific version (useful in case of errors)
```bash
make migrate-force VERSION=6
```

## Environment Variables

```bash
# Connection URL format
MIGRATION_DB=postgres://user:password@host:port/database?sslmode=disable

# Example for local development
export MIGRATION_DB="postgres://aion:aion123@localhost:5432/aionapi?sslmode=disable"
```
