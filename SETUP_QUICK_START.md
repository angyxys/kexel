# 🚀 Kexel Quick Start Guide

Get Kexel running in 5 minutes with Docker Compose.

## ⚡ Fastest Setup (Recommended)

### Prerequisites
- Docker Desktop installed
- Git

### Steps

```bash
# 1. Clone repository
git clone https://github.com/angyxys/kexel.git
cd kexel

# 2. Copy dev environment
cp .env.dev .env

# 3. Start all services
docker-compose -f docker-compose.dev.yml up

# 4. Wait for services to start (~30 seconds)
# You'll see: "listening on http://0.0.0.0:8080" and "Local: http://localhost:5173"

# 5. Open in browser
# Frontend:  http://localhost:5173
# Backend:   http://localhost:8080
# Database:  http://localhost:8081 (Adminer)

# 6. Test backend health
curl http://localhost:8080/health

# 7. Create first user
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "email": "admin@example.com",
    "password": "password123"
  }'

# 8. Login and get token
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "password123"
  }'
```

### Access Kexel
- 🎨 **Admin Panel:** http://localhost:5173
- 🔌 **API:** http://localhost:8080
- 🗄️ **Database UI:** http://localhost:8081

### Stop Services
```bash
docker-compose -f docker-compose.dev.yml down

# Keep database (volumes persist)
docker-compose -f docker-compose.dev.yml down

# Reset everything (delete data)
docker-compose -f docker-compose.dev.yml down -v
```

---

## 🛠️ Local Development (Without Docker)

### Backend Setup

**Install Go 1.26.4+**
```bash
go version  # Verify installation
```

**Setup database**
```bash
# macOS with Homebrew
brew install postgresql@16

# Linux
sudo apt-get install postgresql postgresql-contrib

# Windows
# Download from https://www.postgresql.org/download/windows/

# Create database
createdb -U postgres kexel_db
```

**Install dependencies**
```bash
go mod download
go install github.com/cosmtrek/air@latest  # Hot reload
```

**Run backend**
```bash
cp .env.dev .env
# Edit .env with your database credentials
air -c .air.toml
```

### Frontend Setup

**Install Node 20+**
```bash
node --version  # Verify installation
npm install -g pnpm  # Use pnpm instead of npm
```

**Start frontend**
```bash
cd frontend
pnpm install
pnpm dev
```

Open http://localhost:5173

---

## 📋 What's Included

### Backend ✅
- ✅ JWT authentication with refresh tokens
- ✅ Role-based access control (owner, mod, vip, user)
- ✅ Player management endpoints
- ✅ VRChat integration endpoints
- ✅ PostgreSQL database
- ✅ Error handling & validation
- ✅ CORS support

### Frontend ✅
- ✅ React 19 + TypeScript
- ✅ Tailwind CSS styling
- ✅ React Router navigation
- ✅ Login/Register pages
- ✅ Player management dashboard
- ✅ Role & ban management
- ✅ Form validation with Zod
- ✅ Global state with Zustand
- ✅ Automatic token refresh

### Infrastructure ✅
- ✅ Docker & Docker Compose setup
- ✅ Hot reload for development (Air + Vite)
- ✅ Database UI (Adminer)
- ✅ Development environment (.env.dev)
- ✅ Caddy + Cloudflare HTTPS (production)

---

## 📚 Project Structure

```
kexel/
├── cmd/api/                    # Backend entry point
├── internal/
│   ├── database/               # DB models & connection
│   ├── handler/                # HTTP handlers
│   ├── service/                # Business logic
│   ├── repository/             # Data access
│   ├── config/                 # Configuration
│   └── middleware/             # Auth & middleware
├── frontend/src/
│   ├── pages/                  # Login, Register, Dashboard
│   ├── components/             # Reusable components
│   ├── api/                    # API clients
│   ├── store/                  # Zustand state
│   ├── types/                  # TypeScript types
│   └── schemas/                # Zod validation
├── gateway/                    # Nginx/Caddy config
├── docker-compose.dev.yml      # Dev environment
├── docker-compose.yml          # Production (Caddy)
├── Dockerfile.dev              # Go dev image
├── frontend/Dockerfile.dev     # Node dev image
├── .air.toml                   # Air hot reload config
├── .env.dev                    # Dev environment variables
└── DEVELOPMENT_ROADMAP.md      # Feature roadmap
```

---

## 🔑 Default Credentials

**Database (Adminer)**
- Server: `postgres` (Docker) or `localhost` (Local)
- Username: `postgres`
- Password: `postgres`
- Database: `kexel_db`

**JWT Secret (Dev)**
- `dev_jwt_secret_change_in_production_123456789`

---

## 📝 Useful Commands

### Docker Compose
```bash
# Start all services
docker-compose -f docker-compose.dev.yml up

# Start in background
docker-compose -f docker-compose.dev.yml up -d

# View logs
docker-compose -f docker-compose.dev.yml logs -f

# View specific service logs
docker-compose -f docker-compose.dev.yml logs -f backend

# Stop services
docker-compose -f docker-compose.dev.yml down

# Reset database
docker-compose -f docker-compose.dev.yml down -v
docker-compose -f docker-compose.dev.yml up
```

### Backend
```bash
# Run with air (hot reload)
air -c .air.toml

# Run tests
go test ./...

# Build for production
go build -o kexel ./cmd/api/main.go

# Format code
go fmt ./...

# Lint code
golangci-lint run ./...
```

### Frontend
```bash
# Development server
pnpm dev

# Build for production
pnpm build

# Type checking
pnpm tsc

# Linting
pnpm lint

# Preview production build
pnpm preview
```

### Database
```bash
# Connect to database
psql -U postgres -d kexel_db

# List tables
\dt

# View table structure
\d player

# Execute SQL file
psql -U postgres -d kexel_db -f migrations.sql

# Exit
\q
```

### API Testing
```bash
# Health check
curl http://localhost:8080/health

# Register user
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@example.com","password":"password123"}'

# Login
TOKEN=$(curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"password123"}' | jq -r '.access_token')

# Get players with token
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/web/players

# Get VIP list (requires MAP_SECRET)
curl http://localhost:8080/vrc/list/vip?secret=dev_map_secret
```

---

## 🐛 Troubleshooting

### Port Already in Use
```bash
# Find process using port 8080
lsof -i :8080

# Kill process
kill -9 <PID>

# Or use a different port (edit docker-compose.dev.yml)
```

### Database Connection Error
```bash
# Verify PostgreSQL is running
docker-compose -f docker-compose.dev.yml ps

# Check logs
docker-compose -f docker-compose.dev.yml logs postgres

# Restart database
docker-compose -f docker-compose.dev.yml restart postgres
```

### Hot Reload Not Working
```bash
# Backend (Go)
go install github.com/cosmtrek/air@latest

# Frontend (Node)
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
echo $VITE_API_URL

# Rebuild frontend
cd frontend
pnpm build
pnpm preview
```

---

## 📖 Documentation

- **Development:** See `DEVELOPMENT.md`
- **Roadmap:** See `DEVELOPMENT_ROADMAP.md`
- **Roadmap Visual:** See `ROADMAP_VISUAL.md`
- **Production:** See `README.md` and `gateway/CADDY_SETUP.md`

---

## 🎯 Next Steps

1. ✅ **Start services**
   ```bash
   docker-compose -f docker-compose.dev.yml up
   ```

2. ✅ **Open frontend**
   - Go to http://localhost:5173

3. ✅ **Register & login**
   - Click "Register here" to create account
   - Login with your credentials

4. ✅ **Try it out**
   - Add a player
   - Edit roles
   - View players

5. ✅ **Read the roadmap**
   - Check `DEVELOPMENT_ROADMAP.md` for next features
   - Start with Phase 1: Audit Logging

---

## 🚀 Ready to Code?

**Start implementing Milestone 1.1 (Audit Logging):**

```bash
# Create feature branch
git checkout -b feat/audit-logging

# Follow tasks in DEVELOPMENT_ROADMAP.md
# Implement audit log model, repository, service, handler
# Create frontend page

# Test locally
curl http://localhost:8080/web/audit-logs \
  -H "Authorization: Bearer <token>"

# Commit and push
git add .
git commit -m "feat: implement audit logging system"
git push origin feat/audit-logging
```

---

## 📞 Need Help?

- 📚 Check `DEVELOPMENT.md` for detailed setup
- 🗺️ Review `DEVELOPMENT_ROADMAP.md` for features
- 🐛 Check troubleshooting section above
- 💬 Open an issue on GitHub

**Happy coding! 🎉**
