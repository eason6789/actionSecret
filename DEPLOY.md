<div align="center">
  <a href="#english">English</a> | <a href="#chinese">中文</a>
</div>

---

<h2 id="english">Deploy: CodeNames (Action Secret)</h2>

### Server Path (Frontend)

`/var/www/codenames-game/`

### Domain (Frontend)

https://tuteng3.site/game.html

### Files to Deploy

| File | Purpose |
|------|---------|
| `src/game.html` | Main game page |
| `src/index.html` | Entry/landing page |
| `src/words.js` | Built-in word bank (1200+ words) |
| `src/debug.html` | Debug/testing page |
| `src/test.html` | Test page |

### Backend Service

- **Binary**: `codenames-backend` (Go, compiled)
- **Port**: 8082
- **Management**: systemd
- **Source**: `backend/main.go`

### Deploy Steps

#### Frontend

1. **SFTP upload src/ files to server**

```bash
cd /Users/easonlv/progs/actionSecret
scp -r src/* root@<server_ip>:/var/www/codenames-game/
```

2. **Verify**

```bash
ssh root@<server_ip> "ls -la /var/www/codenames-game/game.html"
```

#### Backend (AI word generator)

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

### Nginx Config Reference

```nginx
# Frontend -- serve game.html at /game.html path
location = /game.html {
    alias /var/www/codenames-game/game.html;
}

# Serve other src files
location /codenames-game/ {
    alias /var/www/codenames-game/;
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

### Backend Environment Variables

Set in the systemd unit file (`/etc/systemd/system/codenames-backend.service`):

| Variable | Description |
|----------|-------------|
| `DEEPSEEK_API_KEY` | DeepSeek API key for word generation |
| `PORT` | Listen port (default: 8082) |

### Dependencies

- Nginx
- Go runtime (backend binary is statically linked, no runtime needed on server)
- DeepSeek API account (for AI word generation)

---

<h2 id="chinese">部署：行动代号 (CodeNames)</h2>

### 服务器路径（前端）

`/var/www/codenames-game/`

### 域名（前端）

https://tuteng3.site/game.html

### 需要部署的文件

| 文件 | 说明 |
|------|------|
| `src/game.html` | 主游戏页面 |
| `src/index.html` | 入口页面 |
| `src/words.js` | 内置词库（1200+ 词语） |
| `src/debug.html` | 调试页面 |
| `src/test.html` | 测试页面 |

### 后端服务

- **二进制文件**: `codenames-backend`（Go，已编译）
- **端口**: 8082
- **管理方式**: systemd
- **源码**: `backend/main.go`

### 部署步骤

#### 前端

1. **SFTP 上传 src/ 文件到服务器**

```bash
cd /Users/easonlv/progs/actionSecret
scp -r src/* root@<server_ip>:/var/www/codenames-game/
```

2. **验证**

```bash
ssh root@<server_ip> "ls -la /var/www/codenames-game/game.html"
```

#### 后端（AI 词库生成）

3. **交叉编译 Go 二进制并上传**

```bash
cd /Users/easonlv/progs/actionSecret/backend
GOOS=linux GOARCH=amd64 go build -o codenames-backend .
scp codenames-backend root@<server_ip>:/usr/local/bin/
```

4. **重启后端服务**

```bash
ssh root@<server_ip> "systemctl restart codenames-backend"
```

5. **验证后端**

```bash
ssh root@<server_ip> "systemctl status codenames-backend"
curl https://tuteng3.site/api/generate-words -X POST -d '{"topic":"水果"}'
```

6. **重载 Nginx（如配置有变更）**

```bash
ssh root@<server_ip> "nginx -t && systemctl reload nginx"
```

7. **浏览器验证**：打开 https://tuteng3.site/game.html

### Nginx 配置参考

```nginx
# 前端 -- 在 /game.html 路径提供游戏页面
location = /game.html {
    alias /var/www/codenames-game/game.html;
}

# 提供其他 src 文件
location /codenames-game/ {
    alias /var/www/codenames-game/;
}

# AI 词库生成 API 代理
location /api/generate-words {
    proxy_pass http://127.0.0.1:8082;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_read_timeout 120s;
}
```

### 后端环境变量

在 systemd unit 文件中配置（`/etc/systemd/system/codenames-backend.service`）：

| 变量 | 说明 |
|------|------|
| `DEEPSEEK_API_KEY` | DeepSeek API 密钥，用于词库生成 |
| `PORT` | 监听端口（默认: 8082） |

### 依赖

- Nginx
- Go 运行时（后端二进制为静态链接，服务器端无需 Go 运行时）
- DeepSeek API 账号（用于 AI 词库生成）
