# StreamDock v1.5.0

单二进制 **Emby 反向代理** 管理面板：浏览器里加站点，每个站点一个监听口；API / WebSocket 走主回源，播放路径可以另走 playback 分流。

仓库：https://github.com/tsumon/StreamDock

## 亮 / 暗 / 跟随系统

面板是同一套 Topology Worksheet（纸面网格），不是两套皮肤。顶栏和登录页都可以切：

- **亮**
- **暗**
- **系统**（跟随 `prefers-color-scheme`）

偏好写在 `localStorage`，取值 `light` / `dark` / `system`。没有存过 = 跟随系统。

## 明文导入导出

站点管理 → 导出配置，得到 `streamdock_backup.json`（`version=streamdock-v1`）。导入只新建、不覆盖。端口冲突或字段不全时，接口返回 `skipped_items`（名字、端口、原因），面板会用 Toast 和列表写出来，不再只给一个 `skipped` 数字。备份是明文站点配置（含回源 URL），不含用户、流量或 JWT。旧文件 `meridian-v1` 仍能导入。

## 和 Meridian 差在哪

本仓是 [snnabb/Meridian](https://github.com/snnabb/Meridian) 的 fork，不是功能全集。

| | StreamDock | 上游 Meridian |
|---|---|---|
| 站点备份 | 明文 JSON 导入/导出 | 加密全量备份 |
| 流量 | 独立 `#traffic` 页 | 并进仪表盘 |
| 主题 | 亮 / 暗 / 跟随系统 | — |
| 分离推流 | 有，字段更窄 | 字段更全 |
| 结构 | 单二进制 + 内嵌原生前端 | 已拆到 `cmd/meridian/`，页面更多 |

## 本版打包

- `streamdock-linux-amd64` / `streamdock-linux-arm64`
- `streamdock-darwin-amd64` / `streamdock-darwin-arm64`
- `streamdock-windows-amd64.exe`，以及内含 `streamdock.exe` 的 `streamdock-windows-amd64.zip`
- `sha256sums.txt`
- 镜像 `ghcr.io/tsumon/streamdock:v1.5.0`

`install.sh` 会拉对应平台的二进制，并用 `sha256sums.txt` 校验（Release 里有这份文件时）。

推荐用仓库根目录的 `docker-compose.yml`。管理面默认 `127.0.0.1:9090`。挂在 nginx 后面时，把 `TRUSTED_PROXY_CIDRS` 设成反代 peer 网段，登录限流才会认 `X-Real-IP` / `X-Forwarded-For`；默认空 = 只信直连 TCP peer。

---

下面是 GitHub 自动生成的提交说明（附录）。
