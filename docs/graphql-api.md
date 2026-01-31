# GraphQL API Reference

A API GraphQL do Aion oferece uma interface flexível e tipada para interagir com todos os recursos da plataforma.

## 🎮 GraphQL Playground

Acesse o playground interativo para testar queries e mutations:

- Básico: `http://localhost:5001/aion/api/v1/graphql/playground`
- Com exemplos prontos: `http://localhost:5001/aion/api/v1/graphql/playground/examples`

O Playground oferece:
- ✅ Autocompletar com schema introspection
- 📖 Documentação inline de todos os tipos
- 🔍 Validação de queries em tempo real
- 📝 Histórico de queries executadas

## 🔐 Autenticação

Todas as queries e mutations (exceto login) requerem autenticação via Bearer token:

```http
Authorization: Bearer YOUR_JWT_TOKEN
```

### Obter Token

```graphql
mutation Login {
  login(username: "user", password: "password") {
    token
  }
}
```

---

## 📦 Records (Registros)

### Queries

#### `recordById` - Buscar por ID
```graphql
query RecordById($id: ID!) {
  recordById(id: $id) {
    id
    userId
    categoryId
    tagId
    title
    description
    eventTime
    recordedAt
    durationSeconds
    value
    source
    timezone
    status
    createdAt
    updatedAt
  }
}
```

**Variables:**
```json
{
  "id": "1"
}
```

---

#### `records` - Listar todos (com paginação)
```graphql
query Records($limit: Int, $afterEventTime: String, $afterId: ID) {
  records(limit: $limit, afterEventTime: $afterEventTime, afterId: $afterId) {
    id
    title
    categoryId
    tagId
    eventTime
    durationSeconds
    value
    status
  }
}
```

**Variables:**
```json
{
  "limit": 50,
  "afterEventTime": null,
  "afterId": null
}
```

---

#### `recordsByTag` - Listar por tag
```graphql
query RecordsByTag($tagId: ID!, $limit: Int) {
  recordsByTag(tagId: $tagId, limit: $limit) {
    id
    title
    tagId
    eventTime
    durationSeconds
    value
  }
}
```

**Variables:**
```json
{
  "tagId": "1",
  "limit": 20
}
```

---

#### `recordsByDay` - Listar por dia
```graphql
query RecordsByDay($date: String!) {
  recordsByDay(date: $date) {
    id
    title
    eventTime
    durationSeconds
    value
    status
  }
}
```

**Variables:**
```json
{
  "date": "2025-01-06"
}
```

---

#### `recordsUntil` - Listar até data
```graphql
query RecordsUntil($until: String!, $limit: Int) {
  recordsUntil(until: $until, limit: $limit) {
    id
    title
    eventTime
    durationSeconds
    value
  }
}
```

**Variables:**
```json
{
  "until": "2025-01-10T23:59:59Z",
  "limit": 50
}
```

---

#### `recordsBetween` - Listar entre datas
```graphql
query RecordsBetween($startDate: String!, $endDate: String!, $limit: Int) {
  recordsBetween(startDate: $startDate, endDate: $endDate, limit: $limit) {
    id
    title
    eventTime
    durationSeconds
    value
    status
  }
}
```

**Variables:**
```json
{
  "startDate": "2025-01-05T00:00:00Z",
  "endDate": "2025-01-09T23:59:59Z",
  "limit": 100
}
```

---

### Mutations

#### `createRecord` - Criar registro
```graphql
mutation CreateRecord($input: CreateRecordInput!) {
  createRecord(input: $input) {
    id
    title
    categoryId
    tagId
    eventTime
    durationSeconds
    value
    status
    createdAt
  }
}
```

**Variables:**
```json
{
  "input": {
    "title": "Treino de corrida",
    "description": "Corrida leve no parque",
    "categoryId": "1",
    "tagId": "4",
    "eventTime": "2025-01-14T06:00:00Z",
    "recordedAt": "2025-01-14T07:30:00Z",
    "durationSeconds": 2400,
    "value": 5.5,
    "source": "mobile_app",
    "timezone": "America/Sao_Paulo",
    "status": "completed"
  }
}
```

---

#### `updateRecord` - Atualizar registro
```graphql
mutation UpdateRecord($input: UpdateRecordInput!) {
  updateRecord(input: $input) {
    id
    title
    description
    categoryId
    tagId
    eventTime
    durationSeconds
    value
    status
    updatedAt
  }
}
```

**Variables (atualização parcial):**
```json
{
  "input": {
    "id": "1",
    "title": "Treino atualizado",
    "durationSeconds": 3000,
    "value": 6.0
  }
}
```

---

#### `softDeleteRecord` - Deletar registro
```graphql
mutation SoftDeleteRecord($input: DeleteRecordInput!) {
  softDeleteRecord(input: $input)
}
```

**Variables:**
```json
{
  "input": {
    "id": "1"
  }
}
```

---

#### `softDeleteAllRecords` - Deletar todos os registros
```graphql
mutation SoftDeleteAllRecords {
  softDeleteAllRecords
}
```

**Sem variables** - deleta todos os records do usuário autenticado.

---

## 🏷️ Tags

### Queries

#### `tags` - Listar todas
```graphql
query Tags {
  tags {
    id
    userId
    name
    categoryId
    description
    createdAt
    updatedAt
  }
}
```

---

#### `tagById` - Buscar por ID
```graphql
query TagById($id: ID!) {
  tagById(id: $id) {
    id
    userId
    name
    categoryId
    description
    createdAt
    updatedAt
  }
}
```

**Variables:**
```json
{
  "id": "1"
}
```

---

#### `tagByName` - Buscar por nome
```graphql
query TagByName($name: String!) {
  tagByName(name: $name) {
    id
    userId
    name
    categoryId
    description
    createdAt
    updatedAt
  }
}
```

**Variables:**
```json
{
  "name": "Golang"
}
```

---

#### `tagsByCategoryId` - Listar por categoria
```graphql
query TagsByCategoryId($categoryId: ID!) {
  tagsByCategoryId(categoryId: $categoryId) {
    id
    name
    categoryId
    description
    createdAt
  }
}
```

**Variables:**
```json
{
  "categoryId": "1"
}
```

---

### Mutations

#### `createTag` - Criar tag
```graphql
mutation CreateTag($input: CreateTagInput!) {
  createTag(input: $input) {
    id
    name
    categoryId
    description
    createdAt
  }
}
```

**Variables:**
```json
{
  "input": {
    "name": "Yoga",
    "categoryId": "1",
    "description": "Yoga practice sessions"
  }
}
```

---

## 📂 Categories

### Queries

#### `categories` - Listar todas
```graphql
query Categories {
  categories {
    id
    userId
    name
    description
    colorHex
    icon
  }
}
```

---

#### `categoryById` - Buscar por ID
```graphql
query CategoryById($id: ID!) {
  categoryById(id: $id) {
    id
    userId
    name
    description
    colorHex
    icon
  }
}
```

**Variables:**
```json
{
  "id": "1"
}
```

---

#### `categoryByName` - Buscar por nome
```graphql
query CategoryByName($name: String!) {
  categoryByName(name: $name) {
    id
    userId
    name
    description
    colorHex
    icon
  }
}
```

**Variables:**
```json
{
  "name": "saude_fisica"
}
```

---

### Mutations

#### `createCategory` - Criar categoria
```graphql
mutation CreateCategory($input: CreateCategoryInput!) {
  createCategory(input: $input) {
    id
    name
    description
    colorHex
    icon
  }
}
```

**Variables:**
```json
{
  "input": {
    "name": "hobbies",
    "description": "Atividades recreativas",
    "colorHex": "#FF5722",
    "icon": "🎨"
  }
}
```

---

#### `updateCategory` - Atualizar categoria
```graphql
mutation UpdateCategory($input: UpdateCategoryInput!) {
  updateCategory(input: $input) {
    id
    name
    description
    colorHex
    icon
  }
}
```

**Variables:**
```json
{
  "input": {
    "id": "1",
    "name": "fitness",
    "colorHex": "#00BCD4"
  }
}
```

---

#### `softDeleteCategory` - Deletar categoria
```graphql
mutation SoftDeleteCategory($input: DeleteCategoryInput!) {
  softDeleteCategory(input: $input)
}
```

**Variables:**
```json
{
  "input": {
    "id": "1"
  }
}
```

---

## 📊 Tipos e Inputs

### Record
```graphql
type Record {
  id: ID!
  userId: ID!
  categoryId: ID!
  tagId: ID!
  title: String!
  description: String
  eventTime: String!
  recordedAt: String
  durationSeconds: Int
  value: Float
  source: String
  timezone: String
  status: String
  createdAt: String!
  updatedAt: String!
}
```

### Tag
```graphql
type Tag {
  id: ID!
  userId: ID!
  name: String!
  categoryId: ID!
  description: String
  createdAt: String!
  updatedAt: String!
}
```

### Category
```graphql
type Category {
  id: ID!
  userId: ID!
  name: String!
  description: String
  colorHex: String
  icon: String # single emoji
}
```

---

## 🔍 Dicas de Uso

### 1. Campos Seletivos
Peça apenas os campos que você precisa:

```graphql
query Records {
  records(limit: 10) {
    id
    title
    eventTime
  }
}
```

### 2. Aliases
Use aliases para queries múltiplas:

```graphql
query MultipleQueries {
  running: recordsByTag(tagId: "4", limit: 5) {
    id
    title
  }
  gym: recordsByTag(tagId: "7", limit: 5) {
    id
    title
  }
}
```

### 3. Fragments
Reutilize campos com fragments:

```graphql
fragment RecordFields on Record {
  id
  title
  eventTime
  durationSeconds
  value
}

query {
  recordById(id: "1") {
    ...RecordFields
  }
}
```

### 4. Variáveis
Sempre use variáveis ao invés de valores fixos:

```graphql
# ✅ Bom
query RecordById($id: ID!) {
  recordById(id: $id) { ... }
}

# ❌ Evite
query {
  recordById(id: "1") { ... }
}
```

---

## 🐛 Tratamento de Erros

Respostas de erro seguem a especificação GraphQL:

```json
{
  "errors": [
    {
      "message": "record not found",
      "path": ["recordById"],
      "extensions": {
        "code": "NOT_FOUND"
      }
    }
  ],
  "data": {
    "recordById": null
  }
}
```

---

## 📚 Recursos Adicionais

- [GraphQL Spec](https://spec.graphql.org/)
- [gqlgen Documentation](https://gqlgen.com/)
- [Best Practices](https://graphql.org/learn/best-practices/)

---

**URL do Playground**: [http://localhost:5001/aion/api/v1/graphql/playground](http://localhost:5001/aion/api/v1/graphql/playground)

**Endpoint GraphQL**: `http://localhost:5001/aion/api/v1/graphql`
