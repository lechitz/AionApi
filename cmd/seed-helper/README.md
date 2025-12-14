# Seed Helper Tool

A Go-based CLI utility for generating seed data credentials and authentication tokens for the AionAPI project.

## Purpose

This tool helps developers quickly set up a local development database with test data by:
- Generating bcrypt password hashes compatible with the database
- Creating JWT tokens that work with the authentication system
- Producing complete `.env.local` configuration files

## Why Go and not Python?

1. **Consistency**: Uses the same libraries as the main application (golang-jwt, bcrypt)
2. **Compatibility**: Tokens generated are guaranteed to work with the API
3. **No external dependencies**: No need to install Python/pip in a Go project
4. **Maintainability**: Single tech stack for the entire project

## Commands

### `generate-env` - Complete Setup (Recommended)

Generate a complete `.env.local` file with all seed variables:

```bash
./bin/seed-helper generate-env [userCount] [secretKey] [password]
```

**Parameters:**
- `userCount` (optional): Number of users to generate (default: 10)
- `secretKey` (optional): JWT secret key (default: read from .env.dev)
- `password` (optional): Password for all seeded users (default: testpassword123)

**Example:**
```bash
# Generate config for 100 users
./bin/seed-helper generate-env 100

# Generate with custom password
./bin/seed-helper generate-env 50 my-secret-key mypassword
```

**Output:**
- Creates `infrastructure/db/seed/.env.local`
- Contains: USER_TOKEN_TEST, SEED_USER_COUNT, DEV_PASSWORD, SECRET_KEY, JWT_TOKEN

### `generate-token` - JWT Token Only

Generate a JWT token for API testing:

```bash
./bin/seed-helper generate-token [userID] [secretKey]
```

**Parameters:**
- `userID` (optional): User ID for the token (default: 1)
- `secretKey` (optional): JWT secret key (default: read from .env.dev)

**Example:**
```bash
# Generate token for user 1
./bin/seed-helper generate-token 1

# Generate for specific user
./bin/seed-helper generate-token 42 my-secret-key
```

**Use case:**
- Testing API endpoints with curl/Postman
- Debugging authentication issues
- Creating tokens for different users

### `generate-bcrypt` - Password Hash Only

Generate a bcrypt hash for a password:

```bash
./bin/seed-helper generate-bcrypt [password]
```

**Parameters:**
- `password` (optional): Password to hash (default: testpassword123)

**Example:**
```bash
# Hash default password
./bin/seed-helper generate-bcrypt

# Hash custom password
./bin/seed-helper generate-bcrypt supersecret
```

**Use case:**
- Manual database seeding
- Testing password validation
- Creating hashes for specific test scenarios

## Usage with Make

The easiest way to use this tool is through Make targets:

```bash
# Build the tool
make seed-helper

# Interactive setup (builds + generates .env.local)
make seed-setup

# Quick seed (uses .env.local)
make seed-quick
```

## Workflow

### First Time Setup

```bash
# 1. Setup seed environment (interactive)
make seed-setup
# Enter number of users when prompted (e.g., 100)

# 2. Seed the database
make seed-quick

# 3. (Optional) Get the JWT token for API testing
cat infrastructure/db/seed/.env.local | grep JWT_TOKEN
```

### Subsequent Seeds

```bash
# Clean and reseed
make seed-clean-all
make seed-quick
```

### Custom Scenarios

```bash
# Build tool
make seed-helper

# Generate 1000 users with specific password
./bin/seed-helper generate-env 1000 "" mypassword

# Seed with custom config
make seed-all-local
```

## Integration with Database

The generated values are used by SQL scripts:

- **user_generate.sql**: Uses `DEV_PASSWORD` via `user_seed_password_plain` variable to hash passwords server-side with pgcrypto
- **user.sql**: Uses `USER_TOKEN_TEST` for pre-hashed passwords
- **Make targets**: Source `.env.local` to pass variables to psql

## Security Notes

1. **Never commit `.env.local`**: It's gitignored for a reason
2. **Development only**: These tools are for local development, not production
3. **Bcrypt cost**: Uses cost factor of 10 (fast for testing, secure enough)
4. **Token expiration**: JWT tokens expire in 24 hours by default

## Troubleshooting

### "SECRET_KEY not found" warning

The tool tries to read SECRET_KEY from:
1. Environment variable `SECRET_KEY`
2. File `infrastructure/docker/environments/dev/.env.dev`

If not found, it uses a placeholder. Provide it as an argument:
```bash
./bin/seed-helper generate-env 100 your-actual-secret-key
```

### File permission error

The `.env.local` file is created with mode 0600 (read/write owner only). If you get permission errors:
```bash
rm infrastructure/db/seed/.env.local
make seed-setup
```

### Token doesn't work with API

Ensure the SECRET_KEY used to generate the token matches the one in your running API:
```bash
# Check API secret
grep SECRET_KEY infrastructure/docker/environments/dev/.env.dev

# Generate token with same secret
./bin/seed-helper generate-token 1 <secret-from-above>
```

## Related Files

- `infrastructure/db/seed/user_generate.sql` - Generates N users with pgcrypto
- `infrastructure/db/seed/.env.example` - Example configuration file
- `makefiles/seed.mk` - Make targets for seeding
- `internal/adapter/secondary/token/jwt_impl.go` - JWT implementation used by API

## Future Enhancements

Potential additions:
- Generate test data for categories/tags/records
- Support for multiple authentication schemes
- Export tokens in different formats (curl, Postman collection)
- Database connectivity for direct seeding
