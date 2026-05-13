# 行动代号 (CodeNames)

一款基于网页的猜词派对游戏，支持双人对战，内置 AI 智能词库生成。

## 游戏玩法

- **猜词者**：根据出题者给出的线索，猜测场上哪些词语属于己方队伍
- **出题者**：查看"密码本"，用一个词语和数字作为线索，引导队友猜词
- 红蓝双方轮流进行，率先猜中己方所有词语的队伍获胜
- 小心卧底（黑色卡牌），一旦翻开游戏立即结束

## 功能特性

- 🎮 红蓝双人对战，经典猜词玩法
- 🤖 AI 智能词库生成（基于 MiniMax M2.7 大模型）
- 📋 15+ 预设主题词库（旅行、美食、科技、亲密关系、两性情感等）
- ✏️ 支持自定义主题，AI 实时生成相关词语
- 🎨 5 种主题配色（深色、日落、森林、粉色、海洋）
- 📱 响应式设计，手机端适配

## 项目结构

```
├── src/
│   ├── game.html          # 主游戏页面
│   ├── words.js           # 内置词库（1200+ 词语）
│   ├── index.html         # 入口页面
│   ├── debug.html         # 调试页面
│   └── test.html          # 测试页面
├── backend/
│   ├── main.go            # Go 后端源码
│   └── go.mod             # Go 模块定义
├── assets/                # 静态资源
└── .gitignore
```

## 技术架构

### 前端
- 纯原生 HTML/CSS/JS，无框架依赖
- 数据持久化使用浏览器 localStorage
- 响应式网格布局，适配手机端

### 后端 (AI 词库生成)
- Go 语言编写，纯标准库实现
- 调用 MiniMax API 生成主题相关中文词语
- 滑动窗口限流（单 IP 每 60 秒最多 5 次）
- CORS 支持

## 本地运行

### 前端
直接用浏览器打开 `src/game.html` 即可游玩（除 AI 词库外功能均可用）。

### 后端（AI 词库生成服务）

```bash
cd backend

# 设置 MiniMax API Key
export MINIMAX_API_KEY="your-api-key"

# 运行
go run main.go

# 或编译
go build -o codenames-backend .
./codenames-backend
```

服务默认监听 `:8082` 端口，可通过 `PORT` 环境变量修改。

### API 接口

**POST /api/generate-words**

请求：
```json
{ "topic": "水果" }
```

响应：
```json
{
  "words": ["苹果", "香蕉", "橙子", ...]
}
```

限流：单 IP 每 60 秒最多 5 次请求，超出返回 HTTP 429。

## 部署

前端通过 Nginx 提供静态文件服务，后端通过 Nginx 反向代理：

```nginx
# AI 词库 API 代理
location /api/generate-words {
    proxy_pass http://127.0.0.1:8082;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_read_timeout 120s;
}
```

后端使用 systemd 管理进程。

## License

MIT
