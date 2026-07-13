#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SCRIPT_NAME="$(basename "${BASH_SOURCE[0]}")"
ENV_PATH="${SCRIPT_DIR}/secrets"
ENV_FILE="${ENV_PATH}/sing-box.env"
CLIENT_PATH="/home/lan/.config/sing-box"
CLIENT_FILE="${CLIENT_PATH}/05-outbound.json"
VERSION="v1.13.14"
TEMPLATE="${SCRIPT_DIR}/config/config.json.template"
OUTPUT="${SCRIPT_DIR}/secrets/config.json"


gen_env() {
    mkdir -p "$ENV_PATH"
    if [[ -f "${ENV_FILE}" ]];then
        echo "已存在 $ENV_FILE 复用"
        return
    fi
    uuid="$(docker run --rm ghcr.io/sagernet/sing-box:"$VERSION" generate uuid)"
    keypair="$(docker run --rm ghcr.io/sagernet/sing-box:"$VERSION" generate reality-keypair)"
    private_key="$(awk '/PrivateKey/ {print $2}' <<< "$keypair")"
    public_key="$(awk '/PublicKey/ {print $2}' <<< "$keypair")"
    short_id="$(openssl rand -hex 4)"
    cat > "$ENV_FILE" << EOF
SB_UUID="$uuid"
SB_PRIVATE_KEY="$private_key"
SB_SHORT_ID="$short_id"
SB_PUBLIC_KEY="$public_key"
EOF
}

render_config() {
    if [[ ! -f "$ENV_FILE" ]]; then
        echo "先运行: $SCRIPT_NAME env"
        exit 1
    fi
    if [[ ! -f "$TEMPLATE" ]]; then
        echo "缺少模板: $TEMPLATE"
        exit 1
    fi
    if ! command -v envsubst >/dev/null; then
        echo "需要 envsubst，Arch 上: sudo pacman -S gettext"
        exit 1
    fi
    set -a
    source "$ENV_FILE"
    set +a
    envsubst '${SB_UUID} ${SB_PRIVATE_KEY} ${SB_SHORT_ID}' \
        < "$TEMPLATE" > "$OUTPUT"
}

gen_client_config() {
    if [[ ! -f "$ENV_FILE" ]]; then
        echo "先运行: $SCRIPT_NAME env 生成密钥"
        exit 1
    fi
    source "$ENV_FILE"
    local ip="$(kubectl get svc sing-box -n sing-box \
        -o jsonpath='{.status.loadBalancer.ingress[0].ip}' 2>/dev/null || true)"
    if [[ -z "$ip" ]]; then
        echo "还没拿到 Service 外网 IP。请确认："
        echo "  1. 已经 kubectl apply -k sing-box/"
        echo "  2. kubectl get svc -n sing-box sing-box 显示 EXTERNAL-IP 不是 <pending>"
        exit 1
    fi
    mkdir -p "$CLIENT_PATH"
    cat > "$CLIENT_FILE" <<EOF
{
  "outbounds": [
    {
      "type": "vless",
      "tag": "hk",
      "server": "$ip",
      "server_port": 443,
      "uuid": "$SB_UUID",
      "flow": "xtls-rprx-vision",
      "tls": {
        "enabled": true,
        "server_name": "www.cloudflare.com",
        "utls": { "enabled": true, "fingerprint": "chrome" },
        "reality": {
          "enabled": true,
          "public_key": "$SB_PUBLIC_KEY",
          "short_id": "$SB_SHORT_ID"
        }
      }
    }
  ]
}
EOF
}


main() {
    case "${1:-}" in
        env)
            gen_env
        ;;
        render)
            render_config
        ;;
        config)
            gen_client_config
        ;;
        *)
            echo "请使用："
            echo "$SCRIPT_NAME env"
            echo "$SCRIPT_NAME config"
            exit 1
    esac
}

main "$@"