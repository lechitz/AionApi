# Database Seed Data

This directory contains seed data for populating the database with test/development data.

> Note: all seeded users share the same bcrypt password hash, provided at runtime via env var `USER_TOKEN_TEST`. The value is not stored in the repo; export it before running any seed command.

## 📋 Overview

The seed structure follows a hierarchical relationship:
```
Users → Categories → Tags → Records
```

### 🔧 Seed Helper Tool

This project includes a Go-based CLI tool (`cmd/seed-helper`) that generates:
- Bcrypt hashes for user passwords
- JWT tokens for API testing
- Complete `.env.local` configuration files

**Why Go instead of Python?**
- Consistent with project tech stack (no external dependencies)
- Uses the same JWT/bcrypt libraries as the main application
- Ensures token compatibility with the running API
- No need to install Python packages in a Go project

**Deprecated:** `generate_users.py` and `generate_records.py` are no longer used. Use `user_generate.sql` (pgcrypto) for mass user generation instead.

## 🗂️ Seed Files

### 1. `user.sql`
Creates test users in the system.

**Users:**
- User ID 1: Primary test user (focus of category/tag/record seeds)
- User ID 2-5: Additional test users
- Password for all seeded users is injected via env var `USER_TOKEN_TEST` (bcrypt hash, not stored in repo)

### 2. `category.sql`
Creates tag categories for organizing tags.

**User 1 Categories (8 total):**
1. **Learning** - Books, courses, studies (#F8B400, book icon)
2. **Fitness** - Exercise, sports, physical activities (#E94F37, dumbbell icon)
3. **Mindfulness** - Meditation, breathing, mental health (#9C27B0, spa icon)
4. **Career** - Work projects, coding, professional development (#1976D2, briefcase icon)
5. **Social** - Friends, family, social connections (#FF6F00, users icon)
6. **Creative** - Art, music, writing, hobbies (#00ACC1, palette icon)
7. **Nutrition** - Meal planning, cooking, healthy eating (#388E3C, utensils icon)
8. **Rest** - Sleep, relaxation, recovery (#5E35B1, moon icon)

### 3. `tags.sql`
Creates specific tags within each category (2 tags per category = 16 total).

**User 1 Tags (16 total):**
- **Learning:** Reading, Online Course
- **Fitness:** Running, Weight Training
- **Mindfulness:** Meditation, Breathing
- **Career:** Coding, Meeting
- **Social:** Family Time, Friends
- **Creative:** Writing, Music
- **Nutrition:** Meal Prep, Hydration
- **Rest:** Sleep, Stretching

### 4. `records.sql`
Creates habit tracking records (128 total for User 1).

**Distribution:**
- **Total records:** 128
- **Days covered:** 16 (2025-01-01 to 2025-01-16)
- **Records per day:** 8 (one from each category)
- **User focused:** User ID 1

**Record Types:**
- Mixed duration activities (5 minutes to 4 hours)
- Various value types (distance, pages, minutes, liters)
- Realistic timestamps throughout each day
- All marked as "completed" status

## 🔄 Seed Order

Seeds must be applied in this order due to foreign key constraints:

```bash
1. user.sql          # Creates users first
2. category.sql      # Creates categories (depends on users)
3. tags.sql          # Creates tags (depends on categories)
4. records.sql       # Creates records (depends on tags)
```

## 🚀 How to Apply Seeds

### Quick Start (Recommended)

The easiest way to seed your database is using the `seed-helper` tool:

```bash
# 1. Setup seed environment (generates .env.local with all credentials)
make seed-setup

# 2. Seed the database
make seed-quick
```

That's it! The setup will ask how many users you want to generate, and create all necessary credentials automatically.

**How it works:**
1. `seed-setup` builds the Go-based `seed-helper` binary
2. It generates a `.env.local` file with properly quoted values (shell-safe)
3. Bcrypt hashes are wrapped in single quotes to prevent shell variable expansion (`$2a$` → `'$2a$...'`)
4. The SECRET_KEY is read from your `infrastructure/docker/environments/dev/.env.dev` file
5. JWT tokens are generated using the same libraries as your main application
6. `seed-quick` loads these values and populates the database

### Advanced: Custom Configuration

#### Generate specific number of users
```bash
# Generate .env.local for 100 users
make seed-helper
./bin/seed-helper generate-env 100

# Then seed
make seed-all-local
```

#### Generate only a JWT token for API testing
```bash
make seed-helper
./bin/seed-helper generate-token 1
```

#### Generate only a bcrypt hash
```bash
make seed-helper
./bin/seed-helper generate-bcrypt mypassword
```

### Manual Method (Advanced Users)

If you prefer to manage credentials yourself:

```bash
# Export the bcrypt hash
export USER_TOKEN_TEST='$2a$10$...'

# Seed all tables
make seed-all
```

### Docker Direct Method

```bash
USER_TOKEN_TEST='$2a$10$...' docker exec -i postgres-dev psql -U aion -d aionapi -v user_seed_password_hash="$USER_TOKEN_TEST" < infrastructure/db/seed/user.sql
docker exec -i postgres-dev psql -U aion -d aionapi < infrastructure/db/seed/category.sql
docker exec -i postgres-dev psql -U aion -d aionapi < infrastructure/db/seed/tags.sql
docker exec -i postgres-dev psql -U aion -d aionapi < infrastructure/db/seed/records.sql
```

### Understanding the .env.local file

The `seed-helper generate-env` command creates a `.env.local` file with:

- `USER_TOKEN_TEST`: bcrypt hash for all seeded users (wrapped in single quotes)
- `SEED_USER_COUNT`: number of users to generate
- `DEV_PASSWORD`: plaintext password (for your reference only, wrapped in single quotes)
- `SECRET_KEY`: your JWT secret (read from .env.dev, wrapped in single quotes)
- `JWT_TOKEN`: a valid JWT token for user 1 (for API testing, wrapped in single quotes)

**Important:** 
- `.env.local` is gitignored. Never commit it!
- All values are wrapped in single quotes to prevent shell variable expansion issues
- Bcrypt hashes start with `$2a$` which would be interpreted as shell variables without quoting

To use it locally, copy the example and fill the bcrypt hash (never commit `.env`):

```bash
cp infrastructure/db/seed/.env.example infrastructure/db/seed/.env
# edit infrastructure/db/seed/.env and replace <bcrypt-hash-placeholder>
```

Alternatively, you can run the Makefile target that loads the `.env` if present:

```bash
# Loads infrastructure/db/seed/.env (if present) and runs the user seed
make seed-users-local
# or: USER_TOKEN_TEST=<bcrypt-hash> make seed-users-local
```

### Quick ways to generate a bcrypt hash locally

Use the Go-based `seed-helper` tool (recommended):

```bash
make seed-helper
./bin/seed-helper generate-bcrypt testpassword123
```

This ensures compatibility with the main application's bcrypt implementation.

### In-database bcrypt hashing (no Python required)

To avoid adding language runtime dependencies, we support generating users directly in Postgres using `pgcrypto`. This hashes the password server-side using `crypt()` and `gen_salt()`.

Usage (example):

```bash
# generate 60 users using plain password 'testpassword123' (hashed inside Postgres)
SEED_USER_COUNT=60 DEV_PASSWORD=testpassword123 make seed-users-local
```

This uses the SQL file `infrastructure/db/seed/user_generate.sql` which requires the Postgres `pgcrypto` extension (the script will create it if missing).

If you prefer to provide an explicit bcrypt hash instead, keep using `USER_TOKEN_TEST` in `.env` or export it and run `make seed-users-local` or `make seed-users`.

### Seeding variable-sized user sets

The example `.env` includes `SEED_USER_COUNT` which your seed scripts can read (if implemented) to generate a variable number of users (e.g., 100, 60, 43). If you prefer, write a small wrapper script (Bash/Python) that reads `SEED_USER_COUNT` and generates a `user.sql` with the requested number of INSERTs before running the seeds.

Example workflow to generate 60 users then seed:

```bash
# set desired count in env file or export variable
export SEED_USER_COUNT=60
make seed-users-local
```

If you want, I can help add a small Python script `generate_users.py` that creates `user.sql` with N users reading `SEED_USER_COUNT` from the environment. This makes end-to-end seeding easy and reproducible.

## 🧪 Testing Queries

After seeding, test with these GraphQL queries:

### Get All Categories for User 1
```graphql
query {
  categories {
    id
    name
    description
    colorHex
    icon
  }
}
```

### Get Tags by Category
```graphql
query {
  tagsByCategoryId(categoryId: "1") {
    id
    name
    description
  }
}
```

### Get Records for a Specific Day
```graphql
query {
  recordsByDay(date: "2025-01-01") {
    id
    title
    eventTime
    categoryId
    tagId
    value
    durationSeconds
  }
}
```

### Get Records by Tag
```graphql
query {
  recordsByTag(tagId: "1", limit: 20) {
    id
    title
    eventTime
    value
  }
}
```

### Get Records in Date Range
```graphql
query {
  recordsBetween(
    startDate: "2025-01-01T00:00:00Z"
    endDate: "2025-01-16T23:59:59Z"
    limit: 128
  ) {
    id
    title
    eventTime
    categoryId
    tagId
  }
}
```

## 📊 Data Statistics (User 1)

| Entity | Count | Distribution |
|--------|-------|-------------|
| Categories | 8 | Evenly distributed across life aspects |
| Tags | 16 | 2 per category |
| Records | 128 | 8 per day × 16 days |

### Records by Category (approximate)
- Learning: 16 records
- Fitness: 16 records
- Mindfulness: 16 records
- Career: 16 records
- Social: 16 records
- Creative: 16 records
- Nutrition: 16 records
- Rest: 16 records

## 🔧 Regenerating Records

The `records.sql` file is generated using `generate_records.py`:

```bash
cd infrastructure/db/seed
python3 generate_records.py > records.sql
```

This script creates realistic, varied records with:
- Random but sensible durations
- Appropriate values for each activity type
- Spread throughout the day
- Realistic descriptions

## 🎯 Design Principles

1. **Realistic Data**: Activities mirror real-world habit tracking
2. **Variety**: Mix of different activity types, durations, and times
3. **Consistency**: Regular daily tracking pattern
4. **Completeness**: All categories and tags are represented
5. **Testability**: Enough data to test all query types and filters

## 📝 Notes

- All timestamps use 'America/Sao_Paulo' timezone
- All records marked as 'completed' status
- Source varies between 'mobile_app' and 'web'
- Values represent different metrics (minutes, km, liters, etc.)
- User 1 is the primary test user with full dataset
- Other users (2-5) have minimal data for multi-user testing

## 🔄 Updating Seeds

When adding new seeds:
1. Maintain foreign key order (users → categories → tags → records)
2. Use consistent timestamp format: `YYYY-MM-DD HH:MM:SS`
3. Include descriptive comments
4. Test seeds in isolated database first
5. Update this README with new counts/structure

---

**Last Updated:** 2025-01-14  
**Total Records (User 1):** 128  
**Date Range:** 2025-01-01 to 2025-01-16
