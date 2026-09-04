#!/bin/sh
set -eu

data_dir=/app/data
secrets_file="${data_dir}/.streamdock-secrets"
legacy_secrets_file="${data_dir}/.meridian-secrets"
streamdock_user=streamdock
streamdock_group=streamdock
streamdock_binary=/app/streamdock
default_db="${data_dir}/streamdock.db"
legacy_db="${data_dir}/meridian.db"

fail() {
    printf 'streamdock-entrypoint: %s\n' "$*" >&2
    exit 1
}

is_streamdock_command=0
if [ "$#" -eq 0 ]; then
    set -- "$streamdock_binary"
    is_streamdock_command=1
else
    case "$1" in
        ./streamdock|streamdock|/app/streamdock|./meridian|meridian|/app/meridian)
            shift
            set -- "$streamdock_binary" "$@"
            is_streamdock_command=1
            ;;
        --healthcheck|--version|-v|-p|--port|--db)
            set -- "$streamdock_binary" "$@"
            is_streamdock_command=1
            ;;
    esac
fi

if [ "$is_streamdock_command" -ne 1 ]; then
    exec "$@"
fi

if [ -z "${DB_PATH:-}" ]; then
    DB_PATH=$default_db
fi
if [ "$DB_PATH" = "$default_db" ] && [ ! -f "$default_db" ] && [ -f "$legacy_db" ]; then
    DB_PATH=$legacy_db
    printf 'streamdock-entrypoint: using legacy database %s\n' "$legacy_db" >&2
fi
export DB_PATH

generate_secret() {
    secret=$(od -An -N32 -tx1 /dev/urandom | tr -d ' \n')
    [ "${#secret}" -eq 64 ] || fail "failed to generate a secure random secret"
    printf '%s' "$secret"
}

validate_secret() {
    name=$1
    value=$2
    [ "${#value}" -ge 32 ] || fail "$name must be at least 32 bytes"
    case "$value" in
        *[[:space:]]*) fail "$name must not contain whitespace" ;;
    esac
}

if [ ! -e "$secrets_file" ] && [ -f "$legacy_secrets_file" ] && [ ! -L "$legacy_secrets_file" ]; then
    secrets_file=$legacy_secrets_file
fi

saved_jwt=
saved_setup=
if [ -e "$secrets_file" ]; then
    if [ ! -f "$secrets_file" ] || [ -L "$secrets_file" ]; then
        fail "$secrets_file must be a regular file"
    fi
    while IFS='=' read -r name value || [ -n "$name$value" ]; do
        case "$name" in
            ''|'#'*) ;;
            JWT_SECRET) saved_jwt=$value ;;
            SETUP_TOKEN) saved_setup=$value ;;
            *) fail "unexpected entry in $secrets_file: $name" ;;
        esac
    done < "$secrets_file"
fi

setup_generated=0
if [ -n "${JWT_SECRET:-}" ]; then jwt=$JWT_SECRET
elif [ -n "$saved_jwt" ]; then jwt=$saved_jwt
else jwt=$(generate_secret)
fi

if [ -n "${SETUP_TOKEN:-}" ]; then setup=$SETUP_TOKEN
elif [ -n "$saved_setup" ]; then setup=$saved_setup
else
    setup=$(generate_secret)
    setup_generated=1
fi

validate_secret JWT_SECRET "$jwt"
validate_secret SETUP_TOKEN "$setup"
[ "$jwt" != "$setup" ] || fail "SETUP_TOKEN must differ from JWT_SECRET"

umask 077
mkdir -p "$data_dir"
tmp_file=$(mktemp "${secrets_file}.tmp.XXXXXX") \
    || fail "failed to create a temporary secrets file"
trap 'rm -f -- "$tmp_file"' EXIT HUP INT TERM
{
    printf '# Generated and managed by the StreamDock Docker entrypoint.\n'
    printf 'JWT_SECRET=%s\n' "$jwt"
    printf 'SETUP_TOKEN=%s\n' "$setup"
} > "$tmp_file"
chmod 0600 "$tmp_file"
mv -f "$tmp_file" "$secrets_file"
trap - EXIT HUP INT TERM

export "JWT_SECRET=$jwt"
export "SETUP_TOKEN=$setup"

if [ "$(id -u)" -eq 0 ]; then
    chown -R "$streamdock_user:$streamdock_group" "$data_dir"
    chmod 0700 "$data_dir"
    chmod 0600 "$secrets_file"
fi

if [ "$setup_generated" -eq 1 ]; then
    printf '\nStreamDock 首次初始化令牌: %s\n' "$setup"
    printf '请打开面板创建管理员；以后可用 docker logs 查找此令牌。\n\n'
fi

if [ "$(id -u)" -eq 0 ]; then
    exec su-exec "$streamdock_user:$streamdock_group" "$@"
fi
exec "$@"
