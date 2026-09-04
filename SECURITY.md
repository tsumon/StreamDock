# 安全策略

## 报告安全问题

如果你发现了安全漏洞，**请不要在公开的 Issue 中报告**。

请通过以下方式私下报告：

- 使用 GitHub 的 [Security Advisory](https://github.com/tsumon/StreamDock/security/advisories/new) 功能
- 或发送邮件至仓库维护者（见仓库 Profile）

我们会在收到报告后尽快确认并处理。

产品名是 **StreamDock**。GitHub 仓库是 `tsumon/StreamDock`。安全报告请发到 **本仓库**，不要发到上游 `snnabb/Meridian`。

## 当前安全边界

以下是 StreamDock 当前的安全设计和已知限制，请在部署前了解：

### 认证

- 管理面板使用 HttpOnly + SameSite=Strict 会话 Cookie（`streamdock_session`）；API 仍可选用 `Authorization: Bearer`
- 密码使用 bcrypt 哈希存储
- JWT 使用 HMAC-SHA256 签名，会话有效期 72 小时
- `JWT_SECRET` 至少 32 字节；未设置时进程会生成临时密钥（重启后旧会话全部失效）
- 首次创建管理员必须提供 `SETUP_TOKEN`（至少 32 字节）。未设置时启动会生成并打印到日志
- 登录 / 首次 setup：同一客户端 1 分钟内 5 次失败后锁定 60 秒，并返回 `Retry-After`
- 新管理员密码 12–72 字节；用户名 1–64 字符

### 单用户

- 当前只支持一个管理员账户
- 不支持多用户、角色划分或权限隔离
- 任何持有有效会话的请求都可以执行所有管理操作

### 网络

- 管理面板默认监听 `127.0.0.1`。公网暴露必须显式设置 `PANEL_BIND_ADDR`，并启用 HTTPS 或 `ALLOW_INSECURE_HTTP=true`
- 管理面板本身不提供 HTTPS，需要外层反代处理 TLS 终止
- 跨源请求被拒绝（不再使用 `Access-Control-Allow-Origin: *`）
- 反代站点的监听端口仍然绑定在所有接口上（这是给 Emby 客户端用的，不是管理面）

### 数据

- SQLite 数据库文件包含管理员密码哈希和站点配置，注意文件权限
- 站点 JSON 导出是明文配置（含回源 URL），不含用户、流量或 JWT
- 没有审计日志，操作不可追溯

### 上游通信

- 与上游 Emby 服务器的通信基于配置的 URL scheme（HTTP 或 HTTPS）
- **业务隧道和 WebSocket 校验证书**；诊断探针在读取证书信息时仍使用 `InsecureSkipVerify`（不影响实际代理流量）
- 反代响应不再删除上游的 `Content-Security-Policy` / `X-Frame-Options`

## 支持的版本

当前项目处于活跃开发阶段，安全修复仅针对最新版本。
