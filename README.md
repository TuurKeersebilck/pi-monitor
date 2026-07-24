# Pi Dashboard

A self-hosted homelab dashboard for Raspberry Pi. Shows live system stats, Docker container status, Pi-hole metrics, and Immich stats — all in one clean interface.

![Go](https://img.shields.io/badge/Go-1.24-blue) ![Svelte](https://img.shields.io/badge/Svelte-5-orange) ![Docker](https://img.shields.io/badge/Docker-ready-2496ED) ![License](https://img.shields.io/badge/license-MIT-green)

## Features

- **System overview** — CPU, RAM, disk, temperature, and network — live gauges or live sparkline graphs with rolling averages
- **Docker containers** — live status for all running containers, grouped by label
- **Pi-hole** — query stats, blocked percentage, and domains blocked
- **Immich** — photo/video count and storage usage
- **Service bookmarks** — customizable shortcut cards to your self-hosted services
- **Dark/light mode** — persisted per browser
- **Custom background** — set a URL, hex color, or upload an image
- **Live updates** — WebSocket connection, no page refresh needed

## Stack

| Layer | Tech |
|---|---|
| Backend | Go — WebSocket hub, REST API, Docker & Pi-hole & Immich clients |
| Frontend | Svelte 5, plain CSS, Chart.js for live sparklines |
| Deployment | Single Docker image, multi-stage build (~50MB) |

## Getting started

### 1. Clone

```bash
git clone https://github.com/TuurKeersebilck/pi-monitor.git
cd pi-monitor
```

### 2. Configure

```bash
cp .env.example .env
```

Edit `.env`:

```env
PIHOLE_URL=http://pihole:80        # URL to your Pi-hole instance
PIHOLE_APP_PASSWORD=               # Pi-hole app password (Settings → API)

IMMICH_URL=http://immich_server:2283   # Optional — leave blank to disable
IMMICH_API_KEY=                        # Immich API key (Account Settings → API Keys)

STATIC_DIR=/app/frontend/dist
LISTEN_ADDR=:8080
```

Immich integration is optional — if `IMMICH_URL` or `IMMICH_API_KEY` are not set, the widget simply won't appear.

### 3. Run

```bash
docker compose up -d
```

The dashboard will be available on port `8080`. If you use a reverse proxy (Nginx Proxy Manager, Traefik, Caddy), point it there.

> **Note:** The `group_add` GID in `docker-compose.yml` is the Docker socket group GID on Raspberry Pi OS (`991`). Check yours with `stat -c '%g' /var/run/docker.sock` and update if different.

## Docker container grouping

You can group containers into named sections in the dashboard by adding a label to your other containers (Jellyfin, Pi-hole, etc.):

```yaml
labels:
  - dashboard.group=Media
```

Without the label, all containers appear under a single "Containers" section.

## Building from source

```bash
# Backend
cd backend && go build ./cmd/server

# Frontend
cd frontend && npm install && npm run build
```

Or build the Docker image:

```bash
docker build -t pi-dashboard .
```

## Environment variables

| Variable | Default | Description |
|---|---|---|
| `PIHOLE_URL` | `http://pihole:80` | Pi-hole base URL |
| `PIHOLE_APP_PASSWORD` | _(empty)_ | Pi-hole app password |
| `IMMICH_URL` | _(empty)_ | Immich base URL (optional) |
| `IMMICH_API_KEY` | _(empty)_ | Immich API key (optional) |
| `LISTEN_ADDR` | `:8080` | Address the server listens on |
| `STATIC_DIR` | `/app/frontend/dist` | Path to built frontend |
| `DATA_DIR` | `/app/data` | Persistent data directory |

## License

MIT — see [LICENSE](LICENSE)
