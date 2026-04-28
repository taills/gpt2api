# gpt2api 部署指南

单文件部署:一个可执行 + 一个配置文件,无需 MySQL / Redis / Nginx。

- 数据库:SQLite(自动创建,无需安装任何服务)
- 前端:编译期通过 `go:embed` 内嵌到二进制,无需单独部署静态资源
- 限流/调度锁:进程内实现,无需 Redis

---

## 快速开始

### 0. 宿主机依赖

| 工具 | 版本要求 |
|------|---------|
| Go   | 1.22+   |
| Node | 18+ / 20 LTS |
| Docker + compose v2 | 可选,仅容器部署需要 |

> 建议加速:
> ```bash
> go env -w GOPROXY=https://goproxy.cn,direct
> npm config set registry https://registry.npmmirror.com
> ```

### 1. 构建单文件可执行

**Linux / macOS / WSL:**
```bash
bash deploy/build-local.sh
```

**Windows PowerShell:**
```powershell
powershell -NoProfile -File deploy\build-local.ps1
```

产物:
```
deploy/bin/gpt2api    # Linux/amd64 单文件(已内嵌前端)
```

### 2a. 直接在宿主机运行

```bash
cp configs/config.example.yaml configs/config.yaml
# 修改 jwt.secret / crypto.aes_key
./deploy/bin/gpt2api -c configs/config.yaml
```

打开 http://localhost:8080 即可访问。  
**第一个注册的用户自动成为 admin。**

### 2b. Docker 部署

```bash
cd deploy
cp .env.example .env        # 修改 JWT_SECRET / CRYPTO_AES_KEY
docker compose build server
docker compose up -d
docker compose logs -f server
```

暴露端口:

| 服务   | 端口   | 说明                     |
|--------|--------|--------------------------|
| server | `8080` | OpenAI 兼容网关 + 管理后台 |

---

## 日常更新

| 场景 | 命令 |
|------|------|
| 代码/前端有变更 | `bash deploy/build-local.sh` → `docker compose build server && docker compose up -d server` |
| 仅改了后端代码  | `bash deploy/build-local.sh --no-web` → `docker compose build server && docker compose up -d server` |
| 仅改了 `.env`   | `docker compose up -d` |
| 快速重启        | `docker compose restart server` |

---

## 数据持久化

Docker compose 挂载了 `data` 卷到 `/app/data`,SQLite 文件存于其中。

冷备份:
```bash
docker compose exec server cp /app/data/gpt2api.db /app/data/gpt2api.db.bak
docker compose cp server:/app/data/gpt2api.db ./gpt2api.db.bak
```

---

## 安全红线

生产环境必须在 `.env` 中覆盖:

- `JWT_SECRET`:至少 32 字符随机串
- `CRYPTO_AES_KEY`:严格 64 位 hex(32 字节 AES-256 key)
