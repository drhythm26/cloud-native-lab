# gitops — App of Apps

Argo CD 的 App of Apps 结构:两个根应用分别管平台组件和业务应用,全部从 `main` 分支同步。

```
platform-root ──→ gitops/applications/ (recurse)
                    ├── argocd.yaml      → Argo CD 自管(argo-helm chart 10.1.3)
                    └── prometheus.yaml  → kube-prometheus-stack 72.6.2
applications-root ──→ apps/applications/ (recurse)
                    └── podinfo.yaml     → apps/podinfo chart
```

## 根应用

| 文件 | 名称 | path | 作用 |
|------|------|------|------|
| `platform-root.yaml` | platform-root | `gitops/applications` | 平台组件(argocd、prometheus) |
| `application-root.yaml` | applications-root | `apps/applications` | 业务应用(podinfo) |

两者都开 `prune` + `selfHeal`:清单删了资源会被删,集群里手改会被改回。

## 平台组件

### argocd(自管)

多源 Application:chart 来自 `argoproj.github.io/argo-helm`(版本 `10.1.3`),values 来自本仓 `gitops/argocd/values.yaml`(`$values` 引用)。bootstrap 时先用 helm 手装同一版本,注册根应用后由这个 Application 接管——Argo CD 从此自己管理自己。

### prometheus(kube-prometheus-stack)

同样是多源:chart `72.6.2` + 本仓 `gitops/prometheus/values.yaml`。要点:

- `ServerSideApply=true`:CRD 注解超过 262KB 限制,必须 server-side apply(踩坑见 commit `c4d449b`)
- values 里 `serviceMonitorSelectorNilUsesHelmValues: false`:让 Prometheus 选中所有带 `release: prometheus` 标签的 ServiceMonitor,而不仅是 chart 自带的
- Grafana Service 为 ClusterIP,访问用 `kubectl -n prometheus port-forward svc/prometheus-grafana 3000:80`
- GKE 托管控制面组件(scheduler / controller-manager / etcd / kube-proxy / coredns)的抓取已关闭,避免常驻 DOWN 目标

## Bootstrap 顺序

`scripts/bootstrap.sh` 负责从零拉起(顺序重要):

1. precheck:kubectl / terraform / gcloud / helm 存在、gcloud 凭证可用、集群可达
2. helm 安装 Argo CD(带 `gitops/argocd/values.yaml`)
3. 等待 Application CRD 就绪
4. `kubectl apply` 两个根应用 → Argo 接管一切

## 验证

```bash
kubectl get application -n argocd
# platform-root / applications-root / argocd / prometheus / podinfo 均应 Synced + Healthy
kubectl -n argocd port-forward svc/argo-cd-argocd-server 8080:443   # UI: http://localhost:8080
```
