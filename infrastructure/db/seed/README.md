# Database Seed Data

This directory contains seed data for populating the database with test/development data.

> Note: all seeded users share the same password, configurable via `DEV_PASSWORD` (default: `testpassword123`).

## 📋 Overview

The seed structure follows a hierarchical relationship:
```
Roles → Admin User → Users → User Roles → Categories → Tags → Records
```

## 🚀 Quick Start

```bash
# Seed 10 users (default) with full data
make seed-all

# Seed 1 user
make seed-all N=1

# Seed 100 users
make seed-all N=100

# Clean and repopulate with 50 users
make populate N=50
```

## 🗂️ Seed Files

| File | Description |
|------|-------------|
| `roles.sql` | System roles (owner, admin, user, blocked) |
| `admin_user.sql` | Admin user "aion" with owner/admin/user roles |
| `user_generate.sql` | Generates N users via pgcrypto |
| `user_roles.sql` | Assigns 'user' role to all users |
| `category_generate.sql` | 8 categories per user (saude_fisica, meditacao, etc.) |
| `tags_generate.sql` | 45+ tags per user (organized by category) |
| `records_generate.sql` | Records for each user (configurable days) |

## 🔧 Make Targets

### Core Seeds
```bash
make seed-all N=10       # Full seed (roles + admin + N users + data)
make seed-user1-all      # Alias for seed-all N=1
make populate N=100      # Clean all tables and repopulate
```

### Individual Seeds
```bash
make seed-roles          # Seed system roles
make seed-admin          # Seed admin user (aion)
make seed-users N=10     # Seed N users
make seed-user-roles     # Assign roles to users
make seed-categories N=10 # Seed categories for N users
make seed-tags N=10      # Seed tags for N users
make seed-records N=10   # Seed records for N users
```

### Clean Targets (Dev Only)
```bash
make seed-clean-all      # Truncate all seeded tables
make seed-clean-users    # Truncate users only
make seed-clean-records  # Truncate records only
```

### API-based Seeding
```bash
make seed-caller N=5     # Seed via API (requires running server)
```

## 📊 Data Generated Per User

| Entity | Count | Naming Pattern |
|--------|-------|----------------|
| Categories | 8 | `{category_name}_{user_id}` (e.g., `saude_fisica_1`) |
| Tags | 45+ | `{tag_name}_{user_id}` (e.g., `Running_1`) |
| Records | ~270 | 6 records per tag × 7 days |

### Categories
- saude_fisica (dumbbell icon)
- meditacao (spa icon)
- saude_mental (brain icon)
- estudo_trabalho (briefcase icon)
- idiomas (globe icon)
- pessoal (user icon)
- trabalho_de_casa (home icon)
- outros (ellipsis-h icon)

## 🔐 Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `N` | 10 | Number of users to generate |
| `DEV_PASSWORD` | testpassword123 | Password for all seeded users |
| `SEED_DAYS` | 7 | Number of days of records to generate |

## 📝 Notes

- All seeds are **idempotent** (safe to run multiple times)
- Seeds use `ON CONFLICT DO NOTHING` to skip duplicates
- Password hashing is done server-side via pgcrypto
- All timestamps use UTC timezone
- Records are spread over configurable days

---

**Last Updated:** 2026-01-18
