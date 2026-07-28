# podinfo

使用上游 [podinfo](https://github.com/stefanprodan/podinfo) 官方镜像，手写 Kubernetes manifests，用 Kustomize 组织，手动部署（不走 Argo CD）。HTTP 端口 `9898`。

练习 K8s 基础资源、探针、扩缩容与故障演练（`/delay`、`/status` 等端点）。

## 结构

```
├── kustomization.yaml   # namespace + resources + labels
├── manifests/
│   ├── namespace.yaml
│   ├── deployment.yaml
│   └── service.yaml
└── note/lab.md          # 实验笔记
```

## 部署（手动）

```bash
# 预览渲染结果
kubectl kustomize .

# 客户端校验（勿对「新建 Namespace + 同批 namespaced 资源」用 --dry-run=server，会假失败）
kubectl apply -k . --dry-run=client

# 部署
kubectl apply -k .
```

卸载：

```bash
kubectl delete -k .
# 或
kubectl delete namespace podinfo
```

## 当前配置

| 项 | 值 |
|---|---|
| 镜像 | `ghcr.io/stefanprodan/podinfo:6.14.1` |
| 副本 | 3 |
| 端口 | `http` → 9898 |
| 探针 | readiness `/readyz`，liveness `/healthz` |
| resources | requests 100m/64Mi，limits 200m/128Mi |

改配置：编辑 `manifests/*.yaml` 或 `kustomization.yaml`，再 `kubectl apply -k .`。

## 验证

```bash
kubectl -n podinfo get pods,svc
kubectl -n podinfo port-forward svc/podinfo 9898:9898
curl -s localhost:9898/
curl -s localhost:9898/healthz
curl -s localhost:9898/metrics | head
```

常用练习端点：`/delay/{seconds}`（模拟慢请求）、`/status/{code}`（模拟指定状态码）。

## 实验笔记

见 [`note/lab.md`](note/lab.md)。
