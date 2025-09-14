# HTTP Router Adapters (Platform)

**Folder:** `internal/platform/server/http/router`

## What is this?

A home for **concrete router adapters** that implement the platform routing contract.
Contexts (Auth/User/etc.) never import these directly—they only depend on the port:

* **Contract:** `internal/platform/server/http/ports.Router`

## Structure

* One subfolder per engine:

    * `chi/` – current adapter using `github.com/go-chi/chi/v5`
    * (optional) `mux/`, `echo/`, etc. — add here if you want alternatives
* Each adapter exposes a `New() ports.Router`.

## Choosing the adapter

Selection happens **only** in the platform HTTP composer:

```go
// internal/platform/server/http/composer.go
import chiadapter "github.com/lechitz/AionApi/internal/platform/server/http/router/chi"

r := chiadapter.New() // returns ports.Router
```

Swap engines by importing a different subfolder and calling its `New()`.

## Why this matters

* **Decoupling:** Context registrars (`RegisterHTTP`) use only `ports.Router`.
* **Swap-friendly:** You can change the routing engine without touching domain code.
* **Consistency:** Same API for grouping, mounting (e.g., GraphQL), and defaults (404/405).

## Authoring a new adapter

1. Create `internal/platform/server/http/router/<impl>/`.
2. Implement all methods of `ports.Router`.
3. Export `New() ports.Router`.
4. Update the composer to use your `<impl>.New()`.

## Rules of use

* Context code must import **only** `internal/platform/server/http/ports`.
* Prefer `Mount()` to attach composite handlers (e.g., GraphQL) at a path.
* Use `GroupWith(middleware, fn)` to protect subtrees.
* `SetError(...)` is stored for platform recovery/error paths (frameworks don’t call it automatically).

## Testing (quick hint)

* Table-test context registrars with a **fake** `ports.Router`.
* Smoke-test concrete adapters (404/405, grouping, mount) within their subfolder.
