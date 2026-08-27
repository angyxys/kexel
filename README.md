# Kexel

![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat&logo=go)
![Gin Framework](https://img.shields.io/badge/Gin-Framework-00ADD8?style=flat)
![GORM](https://img.shields.io/badge/GORM-ORM-e72528?style=flat)
![License](https://img.shields.io/badge/License-MIT-blue.svg)

Kexel is a lightweight, self-hosted role management and moderation system for VRChat worlds. It allows world creators to assign roles (Owner, Moderator, VIP) and manage bans in real-time without needing to update or re-upload the VRChat world instance.

## Overview

VRChat's `VRCStringDownloader` is limited to HTTP GET requests. Kexel bypasses this limitation by providing a secure, read-and-execute REST API built in Go. It acts as the backend bridge between your VRChat UdonSharp scripts and a persistent database, ensuring high performance, low memory usage, and cross-platform compatibility.

## Features

*   **Real-Time Moderation:** Apply bans or grant VIP status instantly.
*   **Role Management:** Granular control over permissions (Owner, Moderator, VIP, Default).
*   **High Performance:** Built with Go and Gin for sub-millisecond response times.
*   **Database Agnostic:** Uses GORM, supporting both PostgreSQL (Production) and SQLite3 (Development).
*   **Map-Level Security:** Secures in-game requests using static Map Secret Tokens to prevent unauthorized API execution.
*   **Environment Composition:** Supports dynamic `.env` variable expansion for flexible deployments.

## Architecture

1.  **Backend API:** Go + Gin + GORM.
2.  **Database:** PostgreSQL (Recommended for Fly.io/Render) or SQLite3.
3.  **VRChat Client:** UdonSharp using `VRCStringDownloader` via HTTPS.

## Prerequisites

*   [Go](https://golang.org/dl/) 1.21 or higher.
*   PostgreSQL database (if running in production).
*   A domain with SSL/HTTPS (Required by VRChat).

## Installation

1. Clone the repository:
   ```bash
   git clone [https://github.com/angyxys/kexel.git](https://github.com/angyxys/kexel.git)
   cd kexel

```

2. Download dependencies:
```bash
go mod tidy

```


3. Create a `.env` file in the root directory. Kexel supports composite variables for easier configuration:
```env
DB_HOST=localhost
POSTGRES_USER=postgres
POSTGRES_PASSWORD=mi_password_seguro
POSTGRES_DB=kexel_db
DB_PORT=5432

DATABASE_DSN="host=$DB_HOST user=$POSTGRES_USER password=$POSTGRES_PASSWORD dbname=$POSTGRES_DB port=$DB_PORT sslmode=disable"

JWT_SECRET=super_secreto_kexel
```


4. Run the application:
```bash
go run cmd/server/main.go
```

## API Endpoints Overview

### VRChat Endpoints (Requires `?secret=` query parameter)

* `GET /api/vrc/player/:id` - Fetches the current role and ban status of a player.
* `GET /api/vrc/action` - Executes a moderation command directly from the game (requires valid target and action type).

### Web Administration Endpoints (Requires OAuth2/JWT)

* `POST /api/web/player` - Creates or updates a player's role.
* `POST /api/web/ban` - Toggles the ban status of a player.

## Security Considerations

Due to VRChat's architecture, Udon scripts are compiled into the `.vrcw` file and executed on the client-side. This means a malicious user could potentially extract your `MAP_SECRET` via world-ripping tools.

To mitigate this:

1. Always serve the Kexel API over **HTTPS**.
2. Rotate your `MAP_SECRET` in the `.env` file and Unity Inspector if you suspect your world has been compromised.
3. Use the planned Web Panel to handle administrative tasks safely outside the game client.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
