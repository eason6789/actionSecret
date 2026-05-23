# Action Secret (行动代号 / Codenames) - 服务器部署指南

## 访问地址

- **游戏页面**: https://your-domain.com/game.html

## 服务器配置

- **项目路径**: `<YOUR_PROJECT_PATH>/codenames-game/src/`
- **游戏文件**: `game.html` (主要), `index.html`, `words.js`
- **Nginx 配置**: `<YOUR_NGINX_CONF_PATH>`

## 部署步骤

### 本地开发后部署

#### 步骤 1: 上传代码

```bash
# 在本地项目目录执行
cd /path/to/your/project

# 上传 src 目录下的文件到服务器
scp -r src/. user@<YOUR_SERVER_IP>:<YOUR_PROJECT_PATH>/src/
```

#### 步骤 2: 验证部署

```bash
# SSH 到服务器验证
ssh user@<YOUR_SERVER_IP> "ls -la <YOUR_PROJECT_PATH>/src/game.html"
```

## Nginx 配置

在 Nginx 配置文件中配置：

```nginx
location = /game.html {
    alias <YOUR_PROJECT_PATH>/src/game.html;
}
```

## 维护命令

### 重启 Nginx

```bash
ssh user@<YOUR_SERVER_IP> "nginx -s reload"
```

### 完整更新流程

```bash
# 1. 本地上传代码
cd /path/to/your/project
scp -r src/. user@<YOUR_SERVER_IP>:<YOUR_PROJECT_PATH>/src/

# 2. SSH 到服务器验证
ssh user@<YOUR_SERVER_IP> "ls -la <YOUR_PROJECT_PATH>/src/game.html"

# 3. 重启 Nginx（如需要）
ssh user@<YOUR_SERVER_IP> "nginx -s reload"
```

## 项目结构

```
src/
├── game.html          # 主要游戏页面
├── index.html         # 备用/测试页面
├── words.js           # 游戏词语库
├── debug.html         # 调试页面
└── test.html          # 测试页面
```

---

*最后更新: 2026-04-25*