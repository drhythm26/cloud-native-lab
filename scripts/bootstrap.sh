#!/usr/bin/env bash
set -euo pipefail

echo "==> 添加 Argo CD Helm 仓库"
helm repo add argo https://argoproj.github.io/argo-helm
helm repo update

echo "==> 安装 Argo CD "
kubectl create namespace argocd --dry-run=client -o yaml | kubectl apply -f -
helm upgrade --install argocd argo/argo-cd \
    -n argocd \
    -f argo-cd/argocd-values.yaml \
    --wait

echo "==> 注册 Application"
kubectl apply -f argo-cd/application.yaml

echo "==> 获取初始密码"
kubectl -n argocd get secret argocd-initial-admin-secret \
    -o jsonpath="{.data.password}" | base64 -d
echo ""

echo "==> 获取 Argo CD UI 地址"
kubectl get svc argocd-server -n argocd
