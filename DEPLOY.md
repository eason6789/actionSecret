# Deploy: 行动代号 (Action Secret / Codenames)

## Server Path (Frontend)

`/var/www/landingPage/codenames-game/src/`

## Domain (Frontend)

https://tuteng3.site/game.html

## Files to Deploy

| File | Purpose |
|------|---------|
| `src/game.html` | Main game page |
| `src/index.html` | Entry/landing page |
| `src/words.js` | Built-in word bank (1200+ words) |
| `src/debug.html` | Debug/testing page |
| `src/test.html` | Test page |

## Backend Service

- **Binary**: `codenames-backend` (Go, compiled)
- **Port**: 8082
- **Management**: systemd
- **Source**: `backend/main.go`

## Deploy Steps

### Frontend

1. **SFTP upload src/ files to server**

```bash
cd /Users/easonlv/progs/actionSecret
scp -r src/* root@<server_ip>:/var/www/landingPage/codenames-game/src/
```

2. **Verify**

```bash
ssh root@<server_ip> "ls -la /var/www/landingPage/codenames-game/src/game.html"
```

### Backend (AI word generator)

3. **Cross-compile Go binary and upload**

```bash
cd /Users/easonlv/progs/actionSecret/backend
GOOS=linux GOARCH=amd64 go build -o codenames-backend .
scp codenames-backend root@<server_ip>:/usr/local/bin/
```

4. **Restart backend service**

```bash
ssh root@<server_ip> "systemctl restart codenames-backend"
```

5. **Verify backend**

```bash
ssh root@<server_ip> "systemctl status codenames-backend"
curl https://tuteng3.site/api/generate-words -X POST -d '{"topic":"水果"}'
```

6. **Reload Nginx (if config changed)**

```bash
ssh root@<server_ip> "nginx -t && systemctl reload nginx"
```

7. **Browser verification**: Open https://tuteng3.site/game.html

## Nginx Config Reference

```nginx
# Frontend -- serve game.html at /game.html path
location = /game.html {
    alias /var/www/landingPage/codenames-game/src/game.html;
}

# Serve other src files
location /codenames-game/ {
    alias /var/www/landingPage/codenames-game/src/;
}

# AI word generation API proxy
location /api/generate-words {
    proxy_pass http://127.0.0.1:8082;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_read_timeout 120s;
}
```

## Backend Environment Variables

Set in the systemd unit file (`/etc/systemd/system/codenames-backend.service`):

| Variable | Description |
|----------|-------------|
| `MINIMAX_API_KEY` | MiniMax API key for word generation |
| `PORT` | Listen port (default: 8082) |

## Dependencies

- Nginx
- Go runtime (backend binary is statically linked, no runtime needed on server)
- MiniMax API account (for AI word generation)
