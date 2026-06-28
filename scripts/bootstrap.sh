#!/usr/bin/env bash
set -euo pipefail

GREEN='\033[0;32m'
YELLOW='\033[0;33m'
RED='\033[0;31m'
NC='\033[0m'

REQUIRED_CMDS=("kubectl" "terraform" "gcloud")

log() {
    echo -e "${GREEN}[bootstrap] LOG: $*${NC}"
}

error() {
    echo -e "${RED}[bootstrap] ERROR: $*${NC}" >&2
    exit 1
}   

require_cmd() {
    if ! command -v "$1" &>/dev/null; then
        error "命令 $1 不存在"
    fi
}


precheck() {
    for cmd in "${REQUIRED_CMDS[@]}"; do
        require_cmd "$cmd"
    done
    log "所有必需命令都已安装"
    gcloud auth application-default print-access-token &>/dev/null || error "未配置 Google Cloud 凭证"
    log "已配置 Google Cloud 凭证"
    timeout 20s kubectl cluster-info &>/dev/null || error "无法连接到 Kubernetes 集群"
    log "已连接到 Kubernetes 集群"
}   



main() {
    log "开始 bootstrap"
    precheck
}

main "$@"