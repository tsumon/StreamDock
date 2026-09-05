<p align="center">
  <img src="./assets/readme/hero.svg" width="100%" alt="StreamDock：单二进制 Emby 反代管理面板。明文导入导出站点，独立流量页。">
</p>

<p align="center">
  <a href="https://go.dev"><img src="https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go&logoColor=white" alt="Go 1.26+"></a>
  <a href="https://pkg.go.dev/modernc.org/sqlite"><img src="https://img.shields.io/badge/SQLite-embedded-003B57?logo=sqlite&logoColor=white" alt="Embedded SQLite"></a>
  <a href="https://github.com/tsumon/StreamDock/releases/latest"><img src="https://img.shields.io/github/v/release/tsumon/StreamDock?label=Release" alt="Latest Release"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="MIT License"></a>
</p>

**Emby 反向代理管理面板。** 单二进制：浏览器加站点，每个站点一个监听口；API 走主回源，播放可另走 playback。

仓库 [`tsumon/StreamDock`](https://github.com/tsumon/StreamDock)

---

## 面板长什么样

对角切面：左上暗色 · 右下亮色。同一布局，两套墨色。顶栏可切 **亮 / 暗 / 系统**。

<p align="center">
  <img src="docs/theme-split.png" width="100%" alt="StreamDock 仪表盘对角深浅对比：运行灯、主回源、播放分流">
</p>

<p align="center">
  <img src="docs/theme-split-sites.png" width="100%" alt="StreamDock 站点管理对角深浅对比：监听端口与明文导入导出">
</p>

---

## 30 秒验证

```bash
cp .env.example .env
# 至少填 JWT_SECRET（openssl rand -hex 32）
docker compose up -d
curl -sS -o /dev/null -w "%{http_code}\n" http://127.0.0.1:9090
```

应打印 `200`。空库第一次创建管理员用启动日志里的 `SETUP_TOKEN`。

加一个站点后：

1. `curl -sS -o /dev/null -w "%{http_code}\n" http://127.0.0.1:<监听端口>/`
2. 打开仪表盘，该节点右侧运行灯应为绿点

站点端口默认映射 `8001-8010`，不够就改 compose。

---

## 先跑起来（Compose 优先）

<p align="center">
  <img src="./assets/readme/flow.svg" width="100%" alt="Compose → 打开面板 → 添加站点 → curl 验证监听口">
</p>

仓库根目录 [`docker-compose.yml`](docker-compose.yml)：

```bash
cp .env.example .env
docker compose up -d
```

面板：`http://127.0.0.1:9090`。镜像 `ghcr.io/tsumon/streamdock`。容器里面板绑 `0.0.0.0`（容器网络）；宿主机只映射本机端口。

裸二进制默认只听 `127.0.0.1:9090`。

---

## Emby 反代

每个站点：

| 字段 | 作用 |
|------|------|
| 监听端口 | Emby 客户端连这个口，不是 9090 |
| 主回源 `target_url` | API / WebSocket |
| 播放回源 `playback_target_url` | `/Videos/`、`/emby/Videos/` 等播放路径 |
| `stream_hosts` | 额外播放节点 |
| UA 模式 | Infuse / Web / 客户端 |

管理面在 `:9090`。Infuse、Emby Theater 等客户端指向 **站点监听口**。

---

## 它能做什么

- **站点反代**：一站点一口；播放可分流
- **明文备份**：导出 / 导入 JSON（`streamdock-v1`）；只新建。端口冲突返回 `skipped_items`（名字、端口、原因）
- **独立流量页**：按站点看入站 / 出站
- **故障诊断**：回源、证书、本地监听
- **亮 / 暗 / 系统**：同一套 CSS 变量；无存储 = 跟随系统

---

## 和上游差在哪

fork 自 [snnabb/Meridian](https://github.com/snnabb/Meridian)。相对当前上游（v1.12.x），不要把「分离推流」当成独有功能——上游站点模型里已经有 `playback_target_url` / `playback_mode` / `stream_hosts`。

| 能力 | 本仓 | 上游 |
|------|------|------|
| 明文站点 JSON 导出 / 导入 | `GET /api/sites/export`、`POST /api/sites/import` | 加密全量备份 |
| 独立流量页 `#traffic` | 有 | 并进仪表盘 |
| Topology Worksheet 亮 / 暗 / 系统 | 有 | — |
| 分离推流 | 有，实现更窄 | 字段更全 |
| 单文件 Go + 原生前端 | 有意保持 | 已拆到 `cmd/meridian/` |

备份不含用户、流量或 JWT。旧文件 `meridian-v1` 仍能导入。

---

## 其它安装方式

### 二进制

从 [Releases](https://github.com/tsumon/StreamDock/releases/latest) 下载 `streamdock-*`（含 `darwin-amd64`）。Windows 另有 `streamdock-windows-amd64.zip`。校验用同级 `sha256sums.txt`。

```bash
chmod +x streamdock-linux-amd64
JWT_SECRET=$(openssl rand -hex 32) SETUP_TOKEN=$(openssl rand -hex 32) ./streamdock-linux-amd64
```

### 安装脚本

```bash
bash <(curl -sL https://raw.githubusercontent.com/tsumon/StreamDock/main/install.sh)
```

从 GitHub Releases 拉二进制并校验 checksum。systemd unit：`streamdock.service`，数据目录 `/opt/streamdock`。

### 源码

```bash
git clone https://github.com/tsumon/StreamDock.git
cd StreamDock
go build -o streamdock .
JWT_SECRET=$(openssl rand -hex 32) SETUP_TOKEN=$(openssl rand -hex 32) ./streamdock
```

---

## 配置

```bash
./streamdock                          # 127.0.0.1:9090，数据库 streamdock.db
./streamdock --port 8080
./streamdock --db /data/streamdock.db
./streamdock --version
```

| 变量 | 默认 | 说明 |
|------|------|------|
| `PORT` | `9090` | 管理面板端口 |
| `DB_PATH` | `streamdock.db` | SQLite。未指定且只有旧的 `meridian.db` 时打开旧文件一轮 |
| `JWT_SECRET` | 进程内随机 | 至少 32 字节。不设则重启后会话作废 |
| `SETUP_TOKEN` | 空库时生成 | 创建第一个管理员 |
| `PANEL_BIND_ADDR` | `127.0.0.1` | 非回环必须再设 `ALLOW_INSECURE_HTTP=true`，或前面加 HTTPS 反代 |
| `TRUSTED_PROXY_CIDRS` | 空 | 逗号分隔 CIDR。空 = 登录限流只信直连 peer；命中后才认 `X-Real-IP` / `X-Forwarded-For` |

会话在 HttpOnly Cookie `streamdock_session`。脚本仍可用 `Authorization: Bearer`。

---

## 备份

站点管理 → 导出配置 → `streamdock_backup_YYYY-MM-DD.json`。导入只新建；跳过项会列出端口、名字、原因。

也可以备份 SQLite：`streamdock.db` 及其 WAL/SHM。旧部署备份 `meridian.db`。恢复时停进程，用原来的 `JWT_SECRET` 启动。

---

## 边界

- 只有一个管理员，没有角色。
- 管理面自己不提供 HTTPS。
- JSON 备份是明文，里面有回源 URL。
- 流量每 60 秒刷进 SQLite，异常退出可能丢掉最近一分钟。
- 这不是上游 Meridian 的功能全集。

安全报告发到本仓库 [Security Advisory](https://github.com/tsumon/StreamDock/security/advisories/new)，不要发到上游。见 [SECURITY.md](SECURITY.md)、[CONTRIBUTING.md](CONTRIBUTING.md)。

---

## 上游

- [Meridian](https://github.com/snnabb/Meridian) — MIT，原版 Emby 反代管理面板

## License

MIT
