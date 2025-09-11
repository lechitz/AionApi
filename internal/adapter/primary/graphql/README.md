# GraphQL (Adapter Primário Central)

Responsabilidade
- Centralizar schema (gqlgen), resolvers root e o handler HTTP do GraphQL.
- Aplicar middlewares de domínio (ex.: Auth) antes de despachar resolvers.

Onde
- Schema: `internal/adapter/primary/graph/schema/*.graphqls` (root, category, tags).
- Resolvers root: `internal/adapter/primary/graph/resolver/*` (orquestram controladores por contexto).
- Handler: `internal/adapter/primary/graph/graphqlserver` (chi router interno + gqlgen + OTel + middlewares).

Single server (padrão)
- O handler retornado por `graphqlserver.NewGraphqlHandler(...)` é montado pela Plataforma HTTP em `cfg.ServerGraphql.Path` via `Router.Mount`.

Como estender
- Adicione `*.graphqls` ao diretório `schema/` e regenere com gqlgen.
- Implemente resolvers chamando controladores dos contextos (ex.: `category/adapter/primary/graphql/handler`).
- Diretivas customizadas ficam em `internal/adapter/primary/graph/middleware` (ex.: `@auth`).

Boas práticas
- Resolver/root não contém regra de negócio; delega para usecases via controladores do contexto.
- Mantenha autenticação/autorização no middleware + controladores.
