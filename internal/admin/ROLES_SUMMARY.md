# Resumo: Roles Architecture - Separação de Responsabilidades

## ✅ Estado Final

### User Domain (`internal/user`)
```go
type User struct {
    Roles []string  // READ-ONLY: populated from user_roles table
    // ... other fields
}
```

**Responsabilidades:**
- ✅ **LER** roles (para autenticação/autorização)
- ✅ Criar usuário com role padrão "user"
- ❌ **NÃO pode modificar** roles

### Admin Domain (`internal/admin`)
```go
type AdminService interface {
    UpdateUserRoles(ctx, cmd) (User, error)  // WRITE roles
}
```

**Responsabilidades:**
- ✅ **ESCREVER** roles (add, remove, update)
- ✅ Validar regras de negócio (não bloquear admin, etc)
- ✅ Gerenciar tabela `user_roles`

---

## 📊 Database Schema

```sql
-- Tabela users SEM coluna roles
CREATE TABLE aion_api.users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR NOT NULL,
    -- NO roles column!
);

-- Tabela de roles disponíveis
CREATE TABLE aion_api.roles (
    role_id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,  -- 'user', 'admin', 'blocked'
    is_active BOOLEAN DEFAULT TRUE
);

-- Junction table (many-to-many)
CREATE TABLE aion_api.user_roles (
    user_role_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id),
    role_id INTEGER REFERENCES roles(role_id),
    assigned_by INTEGER REFERENCES users(user_id),  -- audit: quem atribuiu
    assigned_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, role_id)
);
```

---

## 🔄 Fluxos

### 1. Criar Usuário (sempre role "user")
```
POST /user/create
{
  "username": "joao",
  "email": "joao@example.com"
}

↓ UserRepository.Create()
  1. INSERT INTO users (...)
  2. INSERT INTO user_roles (user_id, role_id)  -- role = 'user'
  3. SELECT user WITH roles via JOIN
  
→ User{Roles: ["user"]}
```

### 2. Promover para Admin
```
PUT /admin/users/{id}/roles
Authorization: Bearer <admin_token>
{
  "roles": ["user", "admin"]
}

↓ AdminRepository.UpdateRoles()
  1. DELETE FROM user_roles WHERE user_id = ?
  2. INSERT INTO user_roles (user_id, role_id) VALUES (?, role_admin)
  3. INSERT INTO user_roles (user_id, role_id) VALUES (?, role_user)
  
→ User{Roles: ["user", "admin"]}
```

### 3. Ler Usuário (busca roles)
```
GET /user/{id}

↓ UserRepository.GetByID()
  1. SELECT * FROM users WHERE user_id = ?
  2. SELECT r.name FROM user_roles ur 
     JOIN roles r ON ur.role_id = r.role_id 
     WHERE ur.user_id = ?
  3. user.Roles = [query results]
  
→ User{Roles: ["user", "admin"]}
```

---

## 🎯 Seed Strategy (make seed-caller n=3)

```bash
make seed-caller n=3
```

**Execução:**
1. `seed-roles` → Cria tabela `roles` + roles padrão
2. `seed-admin` → Cria user `aion` com role `admin` (via SQL)
3. `api-seed-caller` → Cria 3 users via API (user1, user2, user3) com role `user`

**Resultado:**
```
aion   → ["admin", "user"]  (via SQL seed)
user1  → ["user"]           (via API /user/create)
user2  → ["user"]           (via API /user/create)
user3  → ["user"]           (via API /user/create)
```

---

## 🚫 Regras de Negócio

1. ✅ **Default role**: Todo novo usuário recebe role "user"
2. ✅ **Não pode bloquear admin**: Validação no `AdminService`
3. ✅ **Apenas admin pode modificar roles**: Middleware de autenticação
4. ✅ **Roles são normalizadas**: Tabela separada, não CSV em coluna

---

## 📝 Pontos de Atenção

### ⚠️ User.Roles é Read-Only
```go
// ❌ ERRADO: modificar roles no user context
user.Roles = append(user.Roles, "admin")
userRepository.Update(ctx, user)  // Roles serão IGNORADAS!

// ✅ CORRETO: usar admin context
adminService.UpdateUserRoles(ctx, cmd)
```

### ⚠️ Create User ignora domain.User.Roles
```go
// No UserRepository.Create():
user := domain.User{
    Username: "joao",
    Roles: []string{"admin"},  // ← IGNORADO!
}
repo.Create(ctx, user)
// Resultado: user criado com role "user" (hardcoded)
```

### ⚠️ Mapper não preenche Roles
```go
// mapper.UserFromDB() retorna:
User{
    Username: "joao",
    Roles: nil,  // ← Repository precisa popular!
}

// Repository faz:
user := mapper.UserFromDB(userDB)
user.Roles = getUserRolesFromJunction(ctx, userID)  // ← Busca separada
return user
```

---

## 🎓 Lições Aprendidas

1. **Separação de Responsabilidades**: `/user` lê, `/admin` escreve
2. **Read-Only Domain Field**: User.Roles existe mas não é modificável no user context
3. **Transações são críticas**: Create user + assign role devem ser atômicos
4. **Seed order matters**: roles → users → user_roles
5. **Mapper é "dumb"**: Apenas converte struct, não busca relacionamentos

---

## 📚 Referências

- `internal/admin/ROLES_ARCHITECTURE.md` - Arquitetura detalhada
- `infrastructure/db/migrations/01a_roles.sql` - Schema roles
- `infrastructure/db/migrations/01c_user_roles.sql` - Schema junction table
- `infrastructure/db/seed/admin_user.sql` - Seed admin user

