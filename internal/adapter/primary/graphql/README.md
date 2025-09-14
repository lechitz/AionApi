# GraphQL (Primary Adapter – Central)

## Responsibilities

* Centralize the GraphQL **schema (gqlgen)**, **root resolvers**, and the **HTTP handler** for GraphQL.
* Apply **domain middlewares** (e.g., Auth directive) before dispatching to context-specific controllers/use cases.
* Stay **transport-only**: no business rules here—just mapping, orchestration, and delegation.

## Where things live

* **Schema**

    * Root: `internal/adapter/primary/graphql/schema/root.graphqls`
    * Modules: `internal/adapter/primary/graphql/schema/modules/*.graphqls`
      (e.g., `category.graphqls`, `tags.graphqls`)
* **Resolvers (root layer)**

    * `internal/adapter/primary/graphql/root.resolvers.go`
    * `internal/adapter/primary/graphql/category.resolvers.go`
    * `internal/adapter/primary/graphql/tags.resolvers.go`
* **Custom directives / middleware**

    * `internal/adapter/primary/graphql/directives/auth.go` (`@auth`)
* **GraphQL server (handler)**

    * `internal/adapter/primary/graphql/server.go`
* **gqlgen config & generated code**

    * `internal/adapter/primary/graphql/gqlgen.yml`
    * `internal/adapter/primary/graphql/generated.go`
    * `internal/adapter/primary/graphql/model/models_gen.go`
* **Resolver wiring**

    * `internal/adapter/primary/graphql/resolver.go`
      Wires the GraphQL layer to context controllers, e.g.:

        * Category controller: `internal/category/adapter/primary/graphql/controller`

## Runtime shape

* A single GraphQL HTTP handler is exposed from the package (see `server.go`) and mounted by your HTTP platform server.
* The **root resolvers** call the **Category/Tag controllers** (primary adapters for each bounded context), which then delegate to **use cases (input ports)**.

## Extending the schema

1. Add a `.graphqls` file under `internal/adapter/primary/graphql/schema/modules/`.
2. Run gqlgen to re-generate code.
3. Implement the new resolver(s) in this package and **delegate to the appropriate context controller** (e.g., `category/.../controller`).
4. If you need a new directive, implement it under `directives/` and reference it in `schema/root.graphqls`.

## Best practices

* **Keep resolvers thin**: mapping + tracing + calling controllers; no domain logic.
* **Auth/Authorization**: enforce via **directives** (e.g., `@auth`) and/or in controllers—not in resolvers.
* **Trace & log** at the boundary; push domain decisions into use cases.

## Codegen workflow

* Configuration lives in `gqlgen.yml`.
* Generated files:

    * Execution glue: `generated.go`
    * Go types for schema models: `model/models_gen.go`
* Re-generate whenever the schema changes (your repo’s `make verify`/codegen steps already handle this).

## Mounting

* The HTTP platform mounts the GraphQL handler returned by this package (see your server integration) on the configured path.
* No direct router logic here—**this package only provides the GraphQL handler** to be mounted by the platform server.
