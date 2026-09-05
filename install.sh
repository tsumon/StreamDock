#!/usr/bin/env bash
set -euo pipefail

# StreamDock — Emby reverse proxy management panel
# Interactive installer / updater / uninstaller
# Usage: bash <(curl -sL https://raw.githubusercontent.com/tsumon/StreamDock/main/install.sh)

REPO="tsumon/StreamDock"
INSTALL_DIR="/usr/local/bin"
DATA_DIR="/opt/streamdock"
SERVICE_FILE="/etc/systemd/system/streamdock.service"
BIN_NAME="streamdock"

# ─── Colors ───
RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[0;33m'; CYAN='\033[0;36m'; BOLD='\033[1m'; NC='\033[0m'

info()  { echo -e "${CYAN}[INFO]${NC} $*"; }
ok()    { echo -e "${GREEN}[OK]${NC} $*"; }
warn()  { echo -e "${YELLOW}[WARN]${NC} $*"; }
fail()  { echo -e "${RED}[ERROR]${NC} $*"; exit 1; }

# ─── Detect platform ───
detect_platform() {
    local os arch suffix
    os=$(uname -s | tr '[:upper:]' '[:lower:]')
    arch=$(uname -m)

    case "$os" in
        linux)  os="linux" ;;
        darwin) os="darwin" ;;
        *)      fail "不支持的操作系统: $os" ;;
    esac

    case "$arch" in
        x86_64|amd64)   arch="amd64" ;;
        aarch64|arm64)  arch="arm64" ;;
        *)              fail "不支持的架构: $arch" ;;
    esac

    suffix="${os}-${arch}"
    echo "$suffix"
}

# ─── Get latest version tag ───
get_latest_version() {
    curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" 2>/dev/null \
        | grep '"tag_name"' | head -1 | sed 's/.*"tag_name": *"//;s/".*//'
}

# ─── Get current installed version ───
get_current_version() {
    if command -v "$BIN_NAME" &>/dev/null; then
        echo "已安装"
    else
        echo ""
    fi
}

# ─── Install / Update ───
do_install() {
    local suffix version url

    info "检测平台..."
    suffix=$(detect_platform)
    ok "平台: $suffix"

    info "获取最新版本..."
    version=$(get_latest_version)
    if [ -z "$version" ]; then
        fail "当前仓库还没有可用的 GitHub Release。请先从 Releases 页面下载，或改用 Docker / 源码构建。"
    fi
    ok "最新版本: $version"

    url="https://github.com/${REPO}/releases/download/${version}/${BIN_NAME}-${suffix}"
    info "下载 $url ..."
    curl -fSL -o "/tmp/${BIN_NAME}-${suffix}" "$url" || fail "下载失败"

    local sums_url="https://github.com/${REPO}/releases/download/${version}/sha256sums.txt"
    if curl -fsSL -o /tmp/streamdock-sha256sums.txt "$sums_url"; then
        info "校验 sha256sums.txt ..."
        if ! grep -E "[[:space:]]${BIN_NAME}-${suffix}$" /tmp/streamdock-sha256sums.txt > /tmp/streamdock-sha256.one; then
            fail "sha256sums.txt 里没有 ${BIN_NAME}-${suffix}"
        fi
        (cd /tmp && sha256sum -c streamdock-sha256.one) || fail "校验和不匹配"
        ok "校验通过"
    else
        warn "这个 Release 没有 sha256sums.txt，跳过校验"
    fi

    chmod +x "/tmp/${BIN_NAME}-${suffix}"

    info "安装到 ${INSTALL_DIR}/${BIN_NAME} ..."
    sudo mv "/tmp/${BIN_NAME}-${suffix}" "${INSTALL_DIR}/${BIN_NAME}"
    ok "二进制已安装"

    # Create data directory and system user
    if [ ! -d "$DATA_DIR" ]; then
        sudo mkdir -p "$DATA_DIR"
        ok "数据目录已创建: $DATA_DIR"
    fi
    if ! id streamdock >/dev/null 2>&1; then
        sudo useradd --system --home-dir "$DATA_DIR" --shell /usr/sbin/nologin streamdock \
            || sudo useradd --system --home-dir "$DATA_DIR" --shell /bin/false streamdock
        ok "已创建系统用户 streamdock"
    fi
    sudo chown -R streamdock:streamdock "$DATA_DIR"

    # Generate JWT secret if not exists
    local env_file="${DATA_DIR}/.env"
    if [ ! -f "$env_file" ]; then
        local secret setup_token
        secret=$(openssl rand -hex 32)
        setup_token=$(openssl rand -hex 32)
        sudo bash -c "cat > $env_file" <<ENVEOF
JWT_SECRET=${secret}
SETUP_TOKEN=${setup_token}
PORT=9090
DB_PATH=${DATA_DIR}/streamdock.db
PANEL_BIND_ADDR=127.0.0.1
ENVEOF
        sudo chmod 600 "$env_file"
        sudo chown streamdock:streamdock "$env_file"
        ok "配置文件已生成: $env_file"
        echo -e "  首次初始化令牌 SETUP_TOKEN=${BOLD}${setup_token}${NC}"
        echo -e "  面板默认只监听 127.0.0.1；需要公网访问时改 ${DATA_DIR}/.env 中的 PANEL_BIND_ADDR 并设 ALLOW_INSECURE_HTTP=true"
    else
        info "配置文件已存在，跳过: $env_file"
    fi

    # Create systemd service
    if [ -d /run/systemd/system ]; then
        info "配置 systemd 服务..."
        sudo bash -c "cat > $SERVICE_FILE" <<SVCEOF
[Unit]
Description=StreamDock — Emby reverse proxy management panel
After=network.target

[Service]
Type=simple
User=streamdock
Group=streamdock
EnvironmentFile=${DATA_DIR}/.env
ExecStart=${INSTALL_DIR}/${BIN_NAME}
WorkingDirectory=${DATA_DIR}
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
SVCEOF
        sudo systemctl daemon-reload
        sudo systemctl enable streamdock
        ok "systemd 服务已配置"

        echo ""
        read -rp "$(echo -e "${CYAN}是否立即启动 StreamDock？[Y/n]:${NC} ")" start_now
        if [[ "$start_now" != "n" && "$start_now" != "N" ]]; then
            sudo systemctl restart streamdock
            ok "StreamDock 已启动"
        fi
    else
        warn "未检测到 systemd，跳过服务配置"
        echo -e "  手动启动: ${BOLD}source ${DATA_DIR}/.env && ${INSTALL_DIR}/${BIN_NAME}${NC}"
    fi

    echo ""
    echo -e "${GREEN}════════════════════════════════════════${NC}"
    echo -e "${GREEN}  StreamDock $version 安装完成${NC}"
    echo -e "${GREEN}════════════════════════════════════════${NC}"
    echo -e "  面板地址:  ${BOLD}http://$(hostname -I 2>/dev/null | awk '{print $1}' || echo 'localhost'):9090${NC}"
    echo -e "  配置文件:  ${DATA_DIR}/.env"
    echo -e "  数据目录:  ${DATA_DIR}"
    echo -e "  服务管理:  systemctl {start|stop|restart|status} streamdock"
    echo ""
}

# ─── Uninstall ───
do_uninstall() {
    echo ""
    warn "即将卸载 StreamDock，以下内容将被移除："
    echo "  - ${INSTALL_DIR}/${BIN_NAME}"
    echo "  - ${SERVICE_FILE}"
    echo ""

    read -rp "$(echo -e "${RED}是否同时删除数据目录 ${DATA_DIR}？（含数据库和配置）[y/N]:${NC} ")" remove_data
    echo ""
    read -rp "$(echo -e "${RED}确认卸载？[y/N]:${NC} ")" confirm
    if [[ "$confirm" != "y" && "$confirm" != "Y" ]]; then
        info "已取消"
        exit 0
    fi

    # Stop service
    if [ -f "$SERVICE_FILE" ]; then
        sudo systemctl stop streamdock 2>/dev/null || true
        sudo systemctl disable streamdock 2>/dev/null || true
        sudo rm -f "$SERVICE_FILE"
        sudo systemctl daemon-reload
        ok "systemd 服务已移除"
    fi

    # Remove binary
    sudo rm -f "${INSTALL_DIR}/${BIN_NAME}"
    ok "二进制已移除"

    # Remove data
    if [[ "$remove_data" == "y" || "$remove_data" == "Y" ]]; then
        sudo rm -rf "$DATA_DIR"
        ok "数据目录已移除"
    else
        info "数据目录已保留: $DATA_DIR"
    fi

    echo ""
    ok "StreamDock 已卸载"
}

# ─── Main menu ───
main() {
    echo ""
    echo -e "${BOLD}╔══════════════════════════════════════╗${NC}"
    echo -e "${BOLD}║     StreamDock 安装管理工具           ║${NC}"
    echo -e "${BOLD}║     Emby reverse proxy panel         ║${NC}"
    echo -e "${BOLD}╚══════════════════════════════════════╝${NC}"
    echo ""

    local current
    current=$(get_current_version)
    if [ -n "$current" ]; then
        echo -e "  当前状态: ${GREEN}${current}${NC}"
    else
        echo -e "  当前状态: ${YELLOW}未安装${NC}"
    fi
    echo ""
    echo "  1) 安装 / 更新"
    echo "  2) 卸载"
    echo "  0) 退出"
    echo ""

    read -rp "请选择 [0-2]: " choice
    case "$choice" in
        1) do_install ;;
        2) do_uninstall ;;
        0) exit 0 ;;
        *) fail "无效选项" ;;
    esac
}

# Allow direct action via argument: install.sh install / uninstall
case "${1:-}" in
    install|update) do_install ;;
    uninstall|remove) do_uninstall ;;
    *) main ;;
esac
