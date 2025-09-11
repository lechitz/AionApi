# cmd/aion-api (Entrypoint)

Responsabilidade
- Orquestrar ciclo de vida do processo: logger/tracer/metrics, carregar config, inicializar dependências, construir servidor(es) e gerenciar shutdown gracioso.

Como
- `main.go`:
  - Cria logger e keygen.
  - `loadConfig` → `config.New(...).Load` + `Validate`.
  - Inicializa OTel (métricas e traços).
  - `bootstrap.InitializeDependencies` para montar `AppDependencies` do domínio.
  - `serverHTTP.ComposeHandler(cfg, deps, log)` para compor o handler HTTP único (REST + GraphQL montado).
  - `serverHTTP.Build(appCtx, FromHTTP(cfg, handler), log)` cria `*http.Server`.
  - Sobe goroutines dos servidores e aguarda sinais de SO para shutdown.

Split opcional
- Existe bloco comentado que demonstra como criar dois servidores (REST + GraphQL) usando `internal/platform/server/graphql`.

Boas práticas
- Não criar rotas aqui.
- Não manipular dependências de domínio além de passá-las para a Plataforma.
