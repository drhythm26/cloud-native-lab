#!/usr/bin/env bash
set -euo pipefail

GREEN='\033[0;32m'
YELLOW='\033[0;33m'
RED='\033[0;31m'
NC='\033[0m'

REQUIRED_CMDS=("kubectl" "terraform" "gcloud" "helm")
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

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

install_argocd() {
    local namespace="argocd"
    local release="argo-cd"
    local chart="argo/argo-cd"
    local values="${ROOT_DIR}/gitops/platform/argocd/values.yaml"
    log "开始安装 Argo CD"
    [[ -f "${values}" ]] || error "找不到values文件: ${values}"
    helm repo add argo https://argoproj.github.io/argo-helm 2>/dev/null || true
    helm repo update
    helm upgrade --install "${release}" "${chart}" \
        --namespace "${namespace}" --create-namespace \
        -f "${values}" \
        --wait --timeout 10m
    log "Argo CD 安装成功"
    log "UI 访问:"
    log "kubectl port-forward svc/argo-cd-argocd-server -n ${namespace} 8080:443"
    log "浏览器访问: http://localhost:8080"
    log "管理员用户名: admin"
    log "管理员密码: $(kubectl get secret -n ${namespace} argocd-initial-admin-secret -o jsonpath='{.data.password}' | base64 -d)"
}



main() {
    log "开始 bootstrap"
    precheck
    install_argocd
}

main "$@"