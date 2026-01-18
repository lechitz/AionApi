# Roles Management - Architecture

## 📋 Overview

O sistema de roles foi **separado** do bounded context `/user` para o bounded context `/admin`.

### 🔑 Design Decision: Read-Only Roles in User Domain

```go
// User domain (internal/user/core/domain/user_domain.go)
type User struct {
    // ... other fields
    Roles []string  // READ-ONLY: populated by repository from user_roles table
}
```

**Rationale:**
- ✅ User context pode **LER** roles (para autenticação, autorização, JWT)
- ✅ Admin context pode **ESCREVER** roles (via AdminService.UpdateUserRoles)
- ✅ Segue padrão CQRS leve (Command/Query Responsibility Segregation)
- ✅ Mantém User domain "aware" das suas permissões sem gerenciá-las

**Alternative Considered:**
- ❌ Remover Roles do User domain completamente
  - Problema: Toda leitura precisaria de cross-context call
  - Problema: JWT precisaria buscar roles de outro lugar
  - Problema: Middleware de autorização ficaria mais complexo

## 🏗️ Arquitetura

```
/user     -> READ roles (via joins em user_roles)
           -> Cria usuários com role "user" padrão
           
/admin    -> WRITE roles (INSERT/UPDATE/DELETE em user_roles)
           -> Gerencia permissões (promover admin, bloquear, etc)
```

## 🔐 Roles Disponíveis

- `user` (padrão) - Usuário comum
- `admin` - Administrador com acesso total
- `blocked` - Usuário bloqueado sem acesso

## 📊 Fluxo de Criação de Usuário

### 1. Via API `/user/create`
```
POST /user/create
{
  "username": "joao",
  "email": "joao@example.com",
  "password": "senha123"
}

→ User criado com role "user" (hardcoded no repository)
```

### 2. Via Seed SQL
```sql
-- Cria usuário na tabela users
INSERT INTO aion_api.users (...) VALUES (...);

-- Atribui role na tabela user_roles
INSERT INTO aion_api.user_roles (user_id, role_id, assigned_at)
SELECT u.user_id, r.role_id, NOW()
FROM aion_api.users u
CROSS JOIN aion_api.roles r
WHERE u.username = 'aion' AND r.name = 'admin';
```

## 🎯 Seed Strategy

### make seed-caller n=3

Executa na seguinte ordem:

1. **seed-roles** - Cria tabela `roles` com roles padrão
2. **seed-admin** - Cria usuário `aion` com role `admin` via SQL
3. **API Caller** - Cria `n` usuários via API (user1, user2, user3) todos com role `user`

**Resultado:**
- `aion` → admin + user (via seed SQL)
- `user1` → user (via API)
- `user2` → user (via API)
- `user3` → user (via API)

## 🔄 Promovendo Usuário para Admin

Use o endpoint `/admin`:

```
PUT /admin/users/{user_id}/roles
Authorization: Bearer <admin_token>

{
  "roles": ["admin", "user"]
}
```

## 🗄️ Database Schema

```sql
-- Roles fixas no sistema
aion_api.roles (role_id, name, description, is_active)

-- Junction table (many-to-many)
aion_api.user_roles (user_role_id, user_id, role_id, assigned_by, assigned_at)

-- Users sem coluna roles
aion_api.users (user_id, username, email, ...) -- NO roles column!
```

## ⚠️ Regras de Negócio

1. **Todo usuário criado** → role `user` (automático)
2. **Não pode bloquear admin** → validação no usecase
3. **Roles são gerenciadas apenas por admins** → middleware de autenticação

## 📝 Comandos Úteis

```bash
# Seed completo (roles + admin + users via SQL)
make seed-all

# Seed via API (n usuários regulares + 1 admin via SQL)
make seed-caller n=5

# Seed apenas roles
make seed-roles

# Seed apenas admin
make seed-admin

# Limpar roles
make seed-clean-roles
make seed-clean-user-roles
```

## 🔍 Verificação

```sql
-- Ver roles de um usuário
SELECT u.username, r.name as role
FROM aion_api.users u
JOIN aion_api.user_roles ur ON u.user_id = ur.user_id
JOIN aion_api.roles r ON ur.role_id = r.role_id
WHERE u.username = 'aion';
```

