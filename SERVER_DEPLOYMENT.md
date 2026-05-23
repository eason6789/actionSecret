# Action Secret (行动代号 / Codenames) - 服务器部署指南

## 访问地址

- **游戏页面**: https://tuteng3.site/game.html

## 服务器配置

- **项目路径**: `/var/www/codenames-game/`
- **游戏文件**: `game.html` (主要), `index.html`, `words.js`
- **Nginx 配置**: `/etc/nginx/conf.d/tuteng3.site.conf`

## 部署步骤

### 本地开发后部署

#### 步骤 1: 上传代码

```bash
# 在本地项目目录执行
cd /Users/easonlv/Code/web/actionSecret

# 上传 src 目录下的文件
scp -r src/. root@119.29.178.222:/var/www/codenames-game/
```

#### 步骤 2: 验证部署

```bash
# SSH 到服务器验证
ssh root@119.29.178.222 "ls -la /var/www/codenames-game/game.html"
```

## Nginx 配置

已在 `/etc/nginx/conf.d/tuteng3.site.conf` 中配置：

```nginx
location = /game.html {
    alias /var/www/codenames-game/game.html;
}
```

## 维护命令

### 重启 Nginx

```bash
ssh root@119.29.178.222 "nginx -s reload"
```

### 完整更新流程

```bash
# 1. 本地上传代码
cd /Users/easonlv/Code/web/actionSecret
scp -r src/. root@119.29.178.222:/var/www/codenames-game/

# 2. SSH 到服务器验证
ssh root@119.29.178.222 "ls -la /var/www/codenames-game/game.html"

# 3. 重启 Nginx（如需要）
ssh root@119.29.178.222 "nginx -s reload"
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