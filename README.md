# Kexel

![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat&logo=go)
![Gin Framework](https://img.shields.io/badge/Gin-Framework-00ADD8?style=flat)
![GORM](https://img.shields.io/badge/GORM-ORM-e72528?style=flat)
![Docker](https://img.shields.io/badge/Docker-Supported-2496ED?style=flat&logo=docker)
![License](https://img.shields.io/badge/License-MIT-blue.svg)

**Kexel** is a lightweight, self-hosted role management and moderation system for VRChat worlds. It allows world creators to assign roles (Owner, Moderator, VIP) and manage bans in real-time without needing to update or re-upload the VRChat world instance.

## Overview

VRChat's `VRCStringDownloader` is limited to HTTP GET requests. Kexel bypasses this limitation by providing a secure, high-performance REST API built in Go. It acts as the backend bridge between your VRChat UdonSharp scripts and a persistent database, enabling real-time role and moderation management.

The backend is built with:
- **Go 1.26.4** - High-performance concurrent server
- **Gin Framework** - Fast HTTP routing and middleware
- **GORM** - Type-safe database ORM
- **PostgreSQL** - Production-ready database

## Features

- **Real-Time Role Management:** Assign and modify player roles instantly (Owner, Moderator, VIP, User).
- **Player Bans:** Track and manage banned players in real-time.
- **High Performance:** Built with Go and Gin for millisecond-level response times.
- **Database-Backed:** Persistent storage using PostgreSQL with GORM ORM.
- **Role-Based Access Control:** Different permission levels for administrative actions.
- **Environment Configuration:** Dynamic `.env` variable expansion for flexible deployments.
- **Docker Support:** Pre-configured Docker and Docker Compose setup for easy deployment.

## Architecture

```
VRChat World (UdonSharp)
        ↓
    HTTP GET Requests
        ↓
   Kexel Backend API
   (Go + Gin + GORM)
        ↓
   PostgreSQL Database
```

### Project Structure

```
kexel/
├── cmd/
│   └── api/
│       └── main.go           # Application entry point
├── internal/
│   ├── config/               # Configuration management
│   ├── database/             # Database setup and models
│   │   └── models/           # Data models (Player, Role, etc.)
│   ├── repository/           # Data access layer
│   ├── service/              # Business logic layer
│   ├── handler/              # HTTP request handlers
│   └── util/                 # Utility functions
├── gateway/                  # Nginx reverse proxy configuration
├── frontend/                 # Admin frontend (React)
├── Dockerfile                # Docker build configuration
├── docker-compose.yml        # Multi-container orchestration
└── .env.example              # Environment variables template
```

## Prerequisites

### For Local Development
- [Go](https://golang.org/dl/) 1.26.4 or higher
- [PostgreSQL](https://www.postgresql.org/download/) 16 or higher

### For Docker Deployment
- [Docker](https://www.docker.com/) 20.10+
- [Docker Compose](https://docs.docker.com/compose/) 2.0+

## Quick Start

### Option 1: Local Development

1. **Clone the repository:**
   ```bash
   git clone https://github.com/angyxys/kexel.git
   cd kexel
   ```

2. **Download dependencies:**
   ```bash
   go mod download
   ```

3. **Configure environment variables:**
   
   Create a `.env` file in the root directory:
   ```env
   DB_HOST=localhost
   POSTGRES_USER=postgres
   POSTGRES_PASSWORD=your_secure_password
   POSTGRES_DB=kexel_db
   DB_PORT=5432

   DATABASE_DSN="host=$DB_HOST user=$POSTGRES_USER password=$POSTGRES_PASSWORD dbname=$POSTGRES_DB port=$DB_PORT sslmode=disable"

   JWT_SECRET=your_jwt_secret_key
   ```

4. **Initialize the database:**
   ```bash
   # Ensure PostgreSQL is running and create the database
   createdb -U postgres kexel_db
   ```

5. **Run the application:**
   ```bash
   go run cmd/api/main.go
   ```
   
   The API will start on `http://localhost:8080`

### Option 2: Docker Compose with HTTPS (Recommended)

1. **Clone the repository:**
   ```bash
   git clone https://github.com/angyxys/kexel.git
   cd kexel
   ```

2. **Configure environment variables:**
   
   Create a `.env` file with Cloudflare DNS setup:
   ```env
   # Database
   POSTGRES_USER=postgres
   POSTGRES_PASSWORD=your_secure_password
   POSTGRES_DB=kexel_db
   DATABASE_DSN="host=postgres user=postgres password=your_secure_password dbname=kexel_db port=5432 sslmode=disable"
   JWT_SECRET=your_jwt_secret_key

   # Caddy & HTTPS (with Cloudflare DNS validation)
   DOMAIN=api.example.com
   CLOUDFLARE_API_TOKEN=your_cloudflare_api_token
   ACME_EMAIL=your-email@example.com
   ```

   For detailed setup instructions, see [Caddy Setup Guide](gateway/CADDY_SETUP.md)

3. **Start all services:**
   ```bash
   docker-compose up -d
   ```

   Services started:
   - PostgreSQL on port 5432
   - Backend API on port 8080 (internal)
   - Frontend on port 80 (internal)
   - **Caddy gateway** on ports 80 & 443 with automatic HTTPS

4. **Access the application:**
   - Backend API: `https://api.example.com/api`
   - Frontend: `https://api.example.com`
   - Automatic HTTP → HTTPS redirect

### Option 2b: Docker Compose without HTTPS (Development)

If you want to use Docker Compose without HTTPS (local development):

1. Edit `.env` and remove the Caddy variables (or leave defaults)
2. Update `gateway/Caddyfile` to use `http://` instead of `https://`
3. Run: `docker-compose up -d`
4. Access on `http://localhost`

## API Endpoints

### Public Endpoints (VRChat Access)

All endpoints are served under `/vrc` and return role/player data for in-game use.

#### List Endpoints
- `GET /vrc/list/vip` - Get all VIP players
- `GET /vrc/list/banned` - Get all banned players
- `GET /vrc/list/moderator` - Get all moderators
- `GET /vrc/list/owner` - Get all owners

**Response Format:**
```json
[
  {
    "vrchat_id": "usr_12345",
    "role": ["vip", "user"],
    "is_banned": false
  }
]
```

### Admin Endpoints (Web)

All endpoints under `/web` require Bearer token authentication.

#### Authentication
Include the authorization header:
```
Authorization: Bearer YOUR_JWT_SECRET
```

#### Endpoints
- `POST /web/player` - Create or update a player

**Request Body:**
```json
{
  "vrchat_id": "usr_12345",
  "roles": ["moderator", "user"],
  "is_banned": false
}
```

**Response:**
```json
{
  "message": "user usr_12345 created successfully",
  "status": 201
}
```

## Data Models

### Player
| Field | Type | Description |
|-------|------|-------------|
| `vrchat_id` | string | Primary key - VRChat user ID |
| `role` | []string | Array of roles: `owner`, `mod`, `vip`, `user` |
| `is_banned` | bool | Ban status (default: false) |

### Roles
- `owner` - World owner with full permissions
- `mod` - Moderator with ban/unban permissions
- `vip` - VIP player with special status
- `user` - Default player role

## Configuration

All configuration is managed through environment variables in `.env`:

| Variable | Required | Description | Example |
|----------|----------|-------------|---------|
| `DATABASE_DSN` | Yes | PostgreSQL connection string | `host=postgres user=postgres password=pass dbname=kexel_db sslmode=disable` |
| `JWT_SECRET` | Yes | Secret key for Bearer token authentication | `your_jwt_secret_123` |
| `POSTGRES_USER` | Yes (Docker) | PostgreSQL username | `postgres` |
| `POSTGRES_PASSWORD` | Yes (Docker) | PostgreSQL password | `secure_password` |
| `POSTGRES_DB` | Yes (Docker) | PostgreSQL database name | `kexel_db` |

## Security Considerations

### Important
Due to VRChat's architecture, Udon scripts are compiled into world files (`.vrcw`). This means a determined user could potentially extract authentication tokens.

### Best Practices
1. **Always use HTTPS** when serving the API over the internet.
2. **Rotate credentials regularly** - change `JWT_SECRET` periodically.
3. **Use strong database passwords** - especially in production.
4. **Keep backups** - regularly backup your PostgreSQL database.
5. **Monitor access logs** - review gateway/reverse proxy logs for suspicious activity.
6. **Implement rate limiting** - consider adding rate limiting in Nginx or a WAF.

## Development

### Running Tests
```bash
go test ./...
```

### Building for Production
```bash
go build -ldflags="-s -w" -o kexel ./cmd/api/main.go
```

### Docker Build
```bash
docker build -t kexel:latest .
```

## Deployment

### Local Server
See Quick Start Option 1 above.

### Cloud Deployment
The application is designed to work with services like:
- **Fly.io** - Docker-based deployment
- **Render.com** - PostgreSQL + Docker hosting
- **AWS/GCP/Azure** - Manual or Kubernetes deployment

For production deployments, ensure:
- PostgreSQL is configured with proper backups
- API is served over HTTPS with a valid SSL certificate
- Environment variables are securely managed (not hardcoded)
- Reverse proxy (Nginx/HAProxy) is configured for rate limiting

## Troubleshooting

### Database Connection Error
```
error on database: failed to connect to `host=...`
```
**Solution:** Verify PostgreSQL is running and credentials in `.env` are correct.

### Port Already in Use
```
Error binding to port :8080
```
**Solution:** Change the port in code or stop the service using that port.

### Docker Container Exits
```bash
docker-compose logs backend
```
Check logs for connection or configuration issues.

## Contributing

Contributions are welcome! Please:
1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Hosting Options

### Self-Hosted
Deploy Kexel on your own infrastructure using Docker. See [Deployment](#deployment) section for details.

### Kexel Cloud
Don't want to manage servers? Try **[Kexel Cloud](https://kexel.cloud)** - a fully managed, auto-hosted platform for Kexel.

- Zero infrastructure management
- Automatic backups and updates
- Global CDN for low latency
- 99.9% uptime SLA
- Easy migration from self-hosted

[Get started with Kexel Cloud →](https://kexel.cloud)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support & Donate

### Report Issues
For bugs, questions, or suggestions, please open an issue on the [GitHub repository](https://github.com/angyxys/kexel/issues).

### Support the Project
If Kexel has been helpful for your VRChat world, consider supporting the development:

[![Ko-fi](https://img.shields.io/badge/Ko--fi-Support%20angyxys-FF5E5B?style=flat&logo=ko-fi)](https://ko-fi.com/angyxys)

Your support helps me:
- Maintain and update Kexel
- Add new features
- Provide better support
- Keep the project alive

---

**Made with ❤️ for VRChat creators**
