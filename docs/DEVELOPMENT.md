# Development Guide

## Development Environment Setup

### Prerequisites

- Go 1.21+
- Node.js 18+
- `air` for Go hot reload: `go install github.com/air-verse/air@latest`

### Development Commands

#### Full Stack Development (Recommended)

Start both frontend and backend with hot reload:

```bash
make dev
```

This will:
- Start Go backend on port 8080 with hot reload (using `air`)
- Start React frontend on port 5173 with HMR (using Vite)
- Proxy API requests from frontend to backend automatically

Access the app at: **http://localhost:5173**

#### Frontend Only

```bash
make dev-webui
```

Runs Vite dev server with hot module replacement on port 5173.

#### Backend Only

```bash
make dev-go
```

Runs Go backend with hot reload on port 8080 (using `air`).

### Production Build

```bash
make build
```

Builds both frontend and backend for production.

### How Hot Reload Works

**Frontend (Vite/React):**
- Any change to `.tsx`, `.ts`, `.css` files triggers instant HMR
- No page refresh needed in most cases

**Backend (Air/Go):**
- Any change to `.go` files triggers automatic rebuild and restart
- Server restarts in ~1-2 seconds
- Database and state persist across restarts

**API Proxying:**
- During development, Vite proxies `/api/*` requests to `localhost:8080`
- In production, Go serves both the built frontend and API endpoints

### Directory Structure

```
clai/
├── cmd/              # CLI commands
├── internal/         # Internal packages
│   ├── ai/          # AI integration
│   ├── storage/     # SQLite persistence
│   └── webui/       # Web server & API
├── webui/           # React frontend
│   ├── src/         # React components
│   └── dist/        # Built frontend (embedded in Go binary)
├── .air.toml        # Air configuration
└── Makefile         # Build commands
```

### Troubleshooting

**"air: command not found"**
```bash
go install github.com/air-verse/air@latest
```

**Port already in use**
```bash
# Kill existing processes
pkill -f "clai webui"
# Or change port in .air.toml
```

**Frontend can't reach backend**
- Ensure backend is running on port 8080
- Check `webui/vite.config.js` proxy settings
- Verify `/api` requests are being proxied

### Database Location

Development database: `~/.local/share/clai/clai.db`

To reset database:
```bash
rm -rf ~/.local/share/clai/clai.db*
```

