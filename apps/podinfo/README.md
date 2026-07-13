# podinfo

上游 [podinfo](https://github.com/stefanprodan/podinfo) Helm chart(v6.14.0)的 vendor 副本,作为 GitOps 同步与 Prometheus 监控的练习对象。HTTP 端口 `9898`。

## 部署方式(GitOps)

由 Argo CD 自动同步,**不要手动 apply**:

- Application 定义:`apps/applications/podinfo.yaml`
- source:`main` 分支的 `apps/podinfo`,valueFiles 为 `default-values.yaml` + `values.yaml`
- 目标 namespace:`podinfo`(CreateNamespace),`prune` + `selfHeal` 开启

改配置的方式:改 `values.yaml` → commit → push → Argo 自动同步。

## values 覆盖(`values.yaml`)

上游默认值在 `default-values.yaml`,本仓覆盖项:

| 配置 | 值 |
|------|-----|
| `replicaCount` | 2 |
| `image` | `ghcr.io/stefanprodan/podinfo:6.14.0` |
| `resources` | requests 100m/64Mi,limits 200m/128Mi |
| `serviceMonitor.enabled` | true,interval 15s |
| `serviceMonitor.additionalLabels` | `release: prometheus`(kube-prometheus-stack 选择器要求) |

## 验证

```bash
kubectl get application podinfo -n argocd    # Synced / Healthy
kubectl -n podinfo get pods                  # 2 副本 Running
kubectl -n podinfo port-forward svc/podinfo 9898:9898
curl -s localhost:9898/
curl -s localhost:9898/metrics | head
```

常用练习端点:`/delay/{seconds}`(模拟慢请求)、`/status/{code}`(模拟指定状态码)。

## 实验笔记

见 [`note/lab.md`](note/lab.md):replicaCount 扩缩容实验与 Argo / ReplicaSet 事件观察。
