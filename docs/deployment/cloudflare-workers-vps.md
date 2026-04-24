# Cloudflare Workers 前端 + VPS Docker Compose 后端部署

本项目现在支持前后端分离部署：

- 前端：`web/` 构建后发布到 Cloudflare Workers
- 后端：通过 `docker compose -f docker-compose.vps.yml` 部署到 VPS

## 1. 前端部署到 Cloudflare Workers

Cloudflare Workers 静态资源配置位于 `web/wrangler.toml`，其中：

- `directory = "./dist"` 指向 Vite 构建产物
- `not_found_handling = "single-page-application"` 让 SPA 路由回退到 `index.html`

### 构建前端

```bash
cd web
bun install
VITE_REACT_APP_SERVER_URL=https://api.example.com bun run build
```

### 发布前端

```bash
cd web
bunx wrangler deploy
```

> `VITE_REACT_APP_SERVER_URL` 应设置为你的后端公开 API 地址，例如 `https://api.example.com`。

## 2. 后端部署到 VPS

### 准备环境变量

```bash
cp .env.backend.example .env.backend
```

至少修改以下变量：

- `FRONTEND_BASE_URL=https://app.example.com`
- `CORS_ALLOW_ORIGINS=https://app.example.com`
- `SESSION_SECRET`
- `POSTGRES_PASSWORD`
- `REDIS_PASSWORD`

### 启动后端栈

```bash
docker compose --env-file .env.backend -f docker-compose.vps.yml up -d --build
```

## 3. 分离部署时的关键变量

### `FRONTEND_BASE_URL`

当该变量存在时：

- 后端不再负责网页路由
- 非 API 路由会重定向到前端站点
- 支付/登录等用户界面回跳地址会优先指向前端域名
- Passkey 默认来源会使用前端公开域名

### `CORS_ALLOW_ORIGINS`

跨域前端访问后端 API 时必须设置，例如：

```env
CORS_ALLOW_ORIGINS=https://app.example.com
```

多个源可用逗号分隔。

### Cookie 设置

当前端和后端处于不同源时，建议：

```env
SESSION_COOKIE_SECURE=true
SESSION_COOKIE_SAME_SITE=none
```

否则浏览器不会在跨站请求中携带登录态 Cookie。

## 4. 本地一体化运行仍然可用

- `Dockerfile` 仍会构建并打包前端静态产物，适合单容器部署
- `Dockerfile.backend` 只构建后端，适合 VPS 分离部署

## 5. 推荐域名规划

- 前端：`https://app.example.com`
- 后端：`https://api.example.com`

前端构建变量示例：

```bash
VITE_REACT_APP_SERVER_URL=https://api.example.com
```
