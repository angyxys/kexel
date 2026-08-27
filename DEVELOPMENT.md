# Development Guide

This guide explains how to set up Kexel for local development.

## Quick Start with Docker Compose

### Prerequisites
- Docker & Docker Compose
- Git

### Setup

1. **Clone the repository:**
```bash
git clone https://github.com/angyxys/kexel.git
cd kexel
```

2. **Create environment file:**
```bash
cp .env.dev .env
```

3. **Start all services:**
```bash
docker-compose -f docker-compose.dev.yml up
```

4. **Access services:**
   - Backend API: http://localhost:8080
   - Frontend: http://localhost:5173
   - Database Admin (Adminer): http://localhost:8081
   - API Docs (if enabled): http://localhost:8080/swagger

5. **Verify everything works:**
```bash
# Check backend health
curl http://localhost:8080/health

# Register a new user
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

---

## Local Development (Without Docker)

### Backend Setup

**Prerequisites:**
- Go 1.26.4+
- PostgreSQL 16+

**Steps:**

1. **Install dependencies:**
```bash
go mod download
```

2. **Install air (hot reload):**
```bash
go install github.com/cosmtrek/air@latest
```

3. **Setup database:**
```bash
createdb -U postgres kexel_db
```

4. **Create .env file:**
```bash
cp .env.example .env
# Edit .env with your PostgreSQL credentials
```

5. **Run with hot reload:**
```bash
air -c .air.toml
```

The backend will restart automatically when files change.

### Frontend Setup

**Prerequisites:**
- Node.js 20+
- pnpm (or npm)

**Steps:**

1. **Navigate to frontend:**
```bash
cd frontend
```

2. **Install dependencies:**
```bash
pnpm install
```

3. **Start dev server:**
```bash
pnpm dev
```

Frontend will be available at http://localhost:5173 with hot reload.

---

## Development Workflow

### Making Changes

**Backend Changes:**
```bash
# Edit files in internal/
# Air will automatically rebuild and restart
# Changes visible immediately at http://localhost:8080
```

**Frontend Changes:**
```bash
# Edit files in frontend/src/
# Vite will hot-reload automatically
# Changes visible immediately in browser
```

### Running Tests

**Backend:**
```bash
go test ./...
```

**Frontend:**
```bash
cd frontend
pnpm test
```

### Database Migrations

Create migration files in `internal/database/migrations/`:

```go
// Example: internal/database/migrations/migration_name.go
package migrations

import "gorm.io/gorm"

func MigrationName(db *gorm.DB) error {
    // Your migration logic
    return nil
}
```

Then run in main.go:
```go
// Auto-migrate models
db.AutoMigrate(&models.Player{}, &models.User{}, &models.RefreshToken{})
```

---

## Database Management

### Access Database with Adminer
1. Open http://localhost:8081
2. Login with:
   - Server: `postgres`
   - Username: `postgres`
   - Password: `postgres`
   - Database: `kexel_db`

### Direct PostgreSQL Access
```bash
psql -U postgres -d kexel_db -h localhost

# Useful commands:
# \dt - List tables
# \d <table_name> - Describe table
# \q - Quit
```

### Reset Database
```bash
# Drop and recreate database
dropdb -U postgres kexel_db
createdb -U postgres kexel_db

# If using Docker:
docker-compose -f docker-compose.dev.yml down -v
docker-compose -f docker-compose.dev.yml up
```

---

## Common Development Tasks

### Add New Endpoint

1. **Create handler method in `internal/handler/`:**
```go
func (h *PlayerHandler) NewEndpoint(c *gin.Context) {
    // Implementation
}
```

2. **Register route in `cmd/api/main.go`:**
```go
webRoute.POST("/new-endpoint", playerHandl.NewEndpoint)
```

3. **Test with curl:**
```bash
curl -X POST http://localhost:8080/web/new-endpoint \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{...}'
```

### Add New Frontend Component

1. **Create component in `frontend/src/components/`:**
```tsx
export function MyComponent() {
  return <div>Component</div>
}
```

2. **Use in pages or other components:**
```tsx
import { MyComponent } from '../components/MyComponent'

export function MyPage() {
  return <MyComponent />
}
```

### Add API Endpoint Client

1. **Add in `frontend/src/api/`:**
```typescript
export const myApi = {
  getEndpoint: async () => {
    const response = await axiosInstance.get('/endpoint')
    return response.data
  }
}
```

2. **Use in component:**
```typescript
import { myApi } from '../api'

const data = await myApi.getEndpoint()
```

---

## Environment Variables

### Backend
- `DATABASE_DSN` - PostgreSQL connection string
- `JWT_SECRET` - JWT signing secret (use strong value!)
- `MAP_SECRET` - VRChat map authentication secret
- `LOG_LEVEL` - Logging level (debug, info, warn, error)

### Frontend
- `VITE_API_URL` - Backend API URL (default: http://localhost:8080)

---

## Troubleshooting

### Port Already in Use
```bash
# Find process using port
lsof -i :8080

# Kill process
kill -9 <PID>

# Or change port in docker-compose.dev.yml
```

### Database Connection Error
```bash
# Verify PostgreSQL is running
docker-compose -f docker-compose.dev.yml ps

# Check logs
docker-compose -f docker-compose.dev.yml logs postgres

# Reset everything
docker-compose -f docker-compose.dev.yml down -v
docker-compose -f docker-compose.dev.yml up
```

### Hot Reload Not Working

**Backend:**
```bash
# Reinstall air
go install github.com/cosmtrek/air@latest
```

**Frontend:**
```bash
# Clear cache and reinstall
cd frontend
rm -rf node_modules .pnpm-store
pnpm install
pnpm dev
```

### Frontend Can't Connect to Backend
```bash
# Verify backend is running
curl http://localhost:8080/health

# Check VITE_API_URL in .env
# In docker-compose.dev.yml, it should be: http://backend:8080 (internally)
```

---

## Performance Tips

### Reduce Build Time
- Use `.air.toml` to exclude unnecessary directories
- Use pnpm instead of npm for frontend

### Hot Reload Optimization
- Air watches files efficiently but can be configured in `.air.toml`
- Vite HMR works best with local changes to modules

### Database Performance
- Use indexes on frequently queried fields
- Monitor slow queries in PostgreSQL logs

---

## Next Steps

1. Read [DEVELOPMENT_ROADMAP.md](./DEVELOPMENT_ROADMAP.md) for feature implementation plan
2. Check [CONTRIBUTING.md](./CONTRIBUTING.md) for code style guidelines
3. Join the community Discord for questions

Happy coding! 🚀
