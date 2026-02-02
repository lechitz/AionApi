# 🔥 Hot Reload Development Setup

## Overview

All three projects now support **hot reload** when running via `make dev`:

- **Go (aion-api)**: Air auto-recompiles on `.go` changes (~3-5s)
- **Python (aion-chat)**: Uvicorn reloads on `.py` changes (~1-2s)
- **TypeScript (dashboard)**: Vite HMR on `.tsx` changes (<1s)

## Quick Start

```bash
# Start all services with hot reload
cd AionApi
make dev

# In another terminal, edit code:
vim AionApi/internal/adapter/primary/graphql/resolver/category_resolver.go
vim aion-chat/src/application/services/chat_service.py
vim aionapi-dashboard/src/features/chat/components/ChatMessage.tsx

# Save files (Ctrl+S) and see changes automatically!
```

## How It Works

### Dashboard (TypeScript + Vite)

**What changed:**
- Added volumes in `docker-compose-dev.yaml` to mount source code
- Added named volume for `node_modules` (prevents overwriting)
- Vite HMR already configured in Dockerfile

**Hot reload:**
- Edit `.tsx`, `.ts`, `.css` files
- Save → Vite detects change → Browser updates via WebSocket (<1s)
- No page refresh needed!

**Logs:**
```bash
make logs-dashboard
# Look for: "hmr update /src/features/..."
```

---

### Python (aion-chat + Uvicorn)

**What changed:**
- Modified `entrypoint.sh` to detect `DEV_MODE` environment variable
- Added `DEV_MODE=true` in docker-compose
- Added volume to mount `/app/src`
- Uvicorn runs with `--reload` flag in dev mode

**Hot reload:**
- Edit `.py` files
- Save → Uvicorn detects change → Reloads app (~1-2s)
- No restart needed!

**Logs:**
```bash
make logs-chat
# Look for: "Reloading..." and "Application startup complete"
```

---

### Go (aion-api + Air)

**What changed:**
- Created `Dockerfile.dev` with Air hot reload tool
- Updated docker-compose to use `Dockerfile.dev` instead of production Dockerfile
- Added volumes to mount source code + Go caches
- Air watches for `.go` file changes and recompiles automatically

**Hot reload:**
- Edit `.go` files
- Save → Air detects change → Recompiles binary (~3-5s)
- No manual rebuild needed!

**Logs:**
```bash
make logs-api
# Look for: "building..." and "running..."
```

---

## Technical Details

### Files Changed

**1. docker-compose-dev.yaml**
- Added volumes for all 3 services
- Added named volumes for dependencies (node_modules, go cache)
- Added `DEV_MODE=true` for aion-chat
- Changed aion-api to use `Dockerfile.dev`

**2. Dockerfile.dev (aion-api)**
- New file for Go development
- Uses `golang:1.25` (not Alpine) to have build tools
- Installs Air for hot reload
- Does NOT copy source code (comes from volume)

**3. entrypoint.sh (aion-chat)**
- Added logic to detect `DEV_MODE` environment variable
- Runs `uvicorn --reload` when DEV_MODE=true
- Runs normal `uvicorn` in production

**4. Makefile**
- Updated help text to mention hot reload

### Volume Strategy

**Why named volumes for dependencies?**

```yaml
volumes:
  # ✅ Correct: Preserve installed dependencies
  - aionapi-dashboard-node-modules:/app/node_modules
  - go-mod-cache:/go/pkg/mod
  
  # ❌ Wrong: Would delete node_modules!
  - ../aionapi-dashboard:/app  # This overwrites EVERYTHING
```

By using named volumes for `node_modules` and Go caches, we:
1. Mount source code for hot reload
2. Preserve installed dependencies
3. Get faster rebuilds (Go cache reused)

### Performance

| Before (rebuild) | After (hot reload) | Speedup |
|-----------------|-------------------|---------|
| Go: ~60s | Go: ~3-5s | **12-20x** |
| Python: ~30s | Python: ~1-2s | **15-30x** |
| TypeScript: ~40s | TypeScript: <1s | **40x+** |

---

## Troubleshooting

### Go: Air not recompiling

**Check logs:**
```bash
make logs-api | grep -i "error\|building"
```

**Common issues:**
- Syntax error in `.go` file → Air stops watching
- Wrong file pattern in `.air.toml` → Not watching file
- Permission issues → Volume mounted as root

**Solution:**
```bash
# Rebuild with fresh image
make rebuild-api
```

---

### Python: Uvicorn not reloading

**Check if DEV_MODE is enabled:**
```bash
docker exec aion-chat-dev env | grep DEV_MODE
# Should show: DEV_MODE=true
```

**Check logs:**
```bash
make logs-chat | head -20
# Look for: "🔥 DEV MODE: Hot reload enabled"
```

**If not showing:**
```bash
# Rebuild container
make rebuild-chat
```

---

### TypeScript: Vite HMR not working

**Check volumes are mounted:**
```bash
docker exec aionapi-dashboard-dev ls -la /app/src
# Should show your files
```

**Check node_modules:**
```bash
docker exec aionapi-dashboard-dev ls -la /app/node_modules | head -10
# Should show packages, not empty
```

**Clear browser cache:**
- Open DevTools (F12)
- Right-click Refresh → Hard Reload (Ctrl+Shift+R)

**Rebuild if needed:**
```bash
make rebuild-dashboard
```

---

### Permission Issues

**Symptom:**
```
Permission denied: '/app/src/...'
```

**Cause:**
Files in volume mounted with different user ID

**Solution:**
```bash
# Check ownership
ls -la aion-chat/src/

# Should be your user (1000:1000 typically)
# If root owns it:
sudo chown -R $USER:$USER aion-chat/src/
```

---

## Production vs Development

### Production (CI/CD, Docker Hub)

**Uses:**
- `infrastructure/docker/Dockerfile` (multi-stage, optimized)
- No volumes
- Compiled binaries
- Minimal Alpine images

**Build:**
```bash
docker build -f infrastructure/docker/Dockerfile -t aion-api:prod .
```

### Development (local)

**Uses:**
- `infrastructure/docker/Dockerfile.dev` (Go only)
- Volumes for hot reload
- Dev tools (Air, Uvicorn --reload, Vite)
- Full golang/python images

**Build:**
```bash
make dev  # Uses docker-compose-dev.yaml
```

**No conflicts!** Both setups coexist peacefully.

---

## Tips for Fast Development

### 1. Keep `make dev` running

Don't stop/restart unless you:
- Changed Dockerfile
- Changed docker-compose.yaml
- Installed new dependencies

### 2. Use multiple terminals

```bash
# Terminal 1: Logs dashboard
make logs-dashboard -f

# Terminal 2: Logs chat
make logs-chat -f

# Terminal 3: Logs API
make logs-api -f

# Terminal 4: Edit code
vim aion-chat/src/...
```

### 3. Combine with editor integration

**VS Code:**
- Save on focus lost: `"files.autoSave": "onFocusChange"`
- Switch to browser → File auto-saves → HMR triggers

**Vim:**
- `:set autowrite` → Saves on `:make`, `:next`

### 4. Use `make dev-fast` when restarting

```bash
# First time (builds images):
make dev

# Later, if just restarting containers:
make dev-down
make dev-fast  # Skips rebuild (faster)
```

---

## Rollback

If hot reload causes issues, revert to old behavior:

```bash
# 1. Checkout previous docker-compose
git checkout HEAD~1 -- infrastructure/docker/environments/dev/docker-compose-dev.yaml

# 2. Checkout previous entrypoint
git checkout HEAD~1 -- ../aion-chat/infrastructure/docker/scripts/entrypoint.sh

# 3. Remove Dockerfile.dev
rm infrastructure/docker/Dockerfile.dev

# 4. Rebuild
make clean-dev
make dev
```

---

## Summary

✅ **15-40x faster** development cycle  
✅ **No manual rebuilds** for code changes  
✅ **Backward compatible** with production builds  
✅ **Easy to disable** (just revert docker-compose changes)

**Enjoy your blazing fast development! 🚀**
