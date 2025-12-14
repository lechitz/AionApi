# API Seed Caller

CLI para gerar tráfego de sucesso via endpoints da API em vez de acessar o banco diretamente. Útil para observar métricas e tracing ao longo do fluxo autenticado (login + mutações GraphQL).

## Uso rápido

```bash
go run ./cmd/api-seed-caller
```

Por padrão, o cliente chama:
- Host: `http://localhost:5001`
- Contexto: `/aion`
- API root: `/api/v1`
- GraphQL: `/graphql`
- Usuário/senha: `user1` / `testpassword123`
- Log de sucesso: `infrastructure/db/seed/api_success.log`
- Multiusuário: exporte `API_CALLER_COUNT` (ex.: 50) e `API_CALLER_USER_PREFIX` (ex.: `user`) para gerar usuários `user1..userN` e semear para cada um.

## Variáveis de ambiente

| Variável | Descrição |
| --- | --- |
| `API_CALLER_HOST` | Host base da API (default: `http://localhost:5001`) |
| `API_CALLER_CONTEXT` | Context path (default: `/aion`) |
| `API_CALLER_ROOT` | API root (default: `/api/v1`) |
| `API_CALLER_GRAPHQL` | Caminho GraphQL (default: `/graphql`) |
| `API_CALLER_USER` | Usuário para login (default: `user1`) |
| `API_CALLER_PASS` / `API_CALLER_PASSWORD` | Senha para login (default: `testpassword123`) |
| `API_CALLER_SUCCESS_LOG` | Arquivo onde são registrados apenas casos de sucesso (default: `infrastructure/db/seed/api_success.log`) |
| `API_CALLER_AUTO_CREATE` | Se `true`, cria um usuário novo quando o login falhar (default: `false`) |
| `API_CALLER_CLEAN` | Se `true`, executa soft delete de records antes de criar novos (ou apenas limpa, se `ONLY_CLEAN`) (default: `false`) |
| `API_CALLER_ONLY_CLEAN` | Se `true` e `API_CALLER_CLEAN=true`, só limpa: records, tags, categorias e o próprio usuário (soft delete) e encerra (default: `false`) |
| `API_CALLER_COUNT` | Quantidade de usuários a semear em sequência (default: `1`) |
| `API_CALLER_USER_PREFIX` | Prefixo para construir usernames quando `API_CALLER_COUNT>1` (default: `user`, ex.: `user37`) |
| `API_CALLER_RUN_ID` | Identificador do run (usado em nomes/títulos) — default: vazio; se vazio não adiciona sufixo aos nomes (defina para algo único se quiser evitar colisão) |
| `API_CALLER_DEBUG` | Se `true`, loga payloads GraphQL para depuração |
| `API_CALLER_TIMEOUT` | Timeout HTTP (ex.: `15s`, default interno `10s`) |

### Multiusuário

```bash
# 50 usuários user1..user50 com defaults (prefixo=user, senha padrão ou definida por env), criando usuário se faltar
make seed-caller N=50

# Ou via Go
API_CALLER_COUNT=50 API_CALLER_AUTO_CREATE=true go run ./cmd/api-seed-caller
```

## O que o comando faz

1. Autentica via `/auth/login` e obtém o token.
2. (Opcional) Se `API_CALLER_AUTO_CREATE=true` e o login falhar, cria o usuário pelo endpoint `/user/create` e tenta novamente; em modo estrito falha cedo para expor problemas de credencial.
3. Gera um token por usuário (ex.: `user1..userN`) e semeia categorias, tags e registros via GraphQL usando nomes com `run_id` específico de cada usuário para evitar conflitos e reusa entidades existentes quando já presentes.
4. (Opcional) Se `API_CALLER_CLEAN=true`, executa `softDeleteAllRecords` antes de recriar registros (por usuário). Se também definir `API_CALLER_ONLY_CLEAN=true`, só executa a limpeza (records, tags, categorias e soft delete do usuário) e encerra.
5. Registra somente as execuções bem-sucedidas no arquivo de log configurado e valida a presença das entidades criadas.
