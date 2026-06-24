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

echo "==> 添加 Prometheus Helm 仓库"
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

echo "==> 安装 kube-prometheus-stack"
kubectl create namespace monitoring --dry-run=client -o yaml | kubectl apply -f -

helm upgrade --install prometheus prometheus-community/kube-prometheus-stack \
-n monitoring \
-f prometheus/values.yaml \
--wait

echo "==> 注册 ServiceMonitor"
kubectl apply -f prometheus/servicemonitor.yaml

echo "==> 获取 Grafana UI 地址和密码"
kubectl get svc -n monitoring prometheus-grafana
echo "Grafana 默认密码："
kubectl get secret -n monitoring prometheus-grafana \
-o jsonpath="{.data.admin-password}" | base64 -d
echo ""
