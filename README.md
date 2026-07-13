# cloud-native-lab

个人 Cloud Native 实验仓:在 GKE 上用 Argo CD(App of Apps)做 GitOps,跑自研 go-api、podinfo 和 kube-prometheus-stack,配套 Terraform 管基础设施。所有组件的同步源统一为 `main` 分支。

## 仓库结构

```
├── gitops/                  # 平台层:App of Apps 根 + 平台组件
│   ├── platform-root.yaml   #   根应用:装 gitops/applications/ 下所有平台组件
│   ├── application-root.yaml#   根应用:装 apps/applications/ 下所有业务应用
│   ├── applications/        #   平台组件的 Argo Application(argocd 自管、prometheus)
│   ├── argocd/values.yaml   #   Argo CD Helm values
│   └── prometheus/values.yaml
├── apps/                    # 应用层
│   ├── applications/        #   业务应用的 Argo Application(podinfo)
│   ├── go-api/              #   自研 Go HTTP 服务(手动部署,不走 Argo)
│   └── podinfo/             #   podinfo Helm chart(Argo 自动同步)
├── sing-box/                # sing-box 代理实验(kustomize,手动部署)
├── infrastructure/gcp/      # Terraform:VPC + GKE
└── scripts/bootstrap.sh     # 从零拉起 Argo CD 并注册两个根应用
```

各部分详细说明见对应目录的 README。

## 从零搭建

```bash
# 1. 创建基础设施(需先配置 gcloud 凭证,见 infrastructure/gcp/README.md)
cd infrastructure/gcp && terraform init && terraform apply

# 2. 获取集群凭证
gcloud container clusters get-credentials cloud-native-lab-gke --zone asia-east2-a

# 3. 拉起 Argo CD 并注册 App of Apps
./scripts/bootstrap.sh
```

bootstrap 完成后 Argo CD 接管一切:它会同步 `platform-root`(argocd 自管 + prometheus)和 `applications-root`(podinfo)。脚本最后会打印 UI 端口转发命令和 admin 初始密码。

## GitOps 工作流

改 `apps/` 或 `gitops/` 下的清单 → commit → push 到 `main` → Argo 自动同步(`prune` + `selfHeal` 均开启)。**不要**直接 `kubectl edit` Argo 管理的资源,会被 selfHeal 冲掉。

例外:

- **go-api** 手动部署:`kubectl apply -f apps/go-api/manifests/`(镜像由 GitHub Actions 构建,见 `.github/workflows/go-api.yaml`)
- **sing-box** 手动部署:`kubectl apply -k sing-box/`

## 验证

```bash
kubectl get application -n argocd          # 全部应为 Synced / Healthy
kubectl get pods -n go-api -n podinfo
kubectl get servicemonitor -A              # go-api、podinfo 均带 release=prometheus
```

## 实验笔记

排障与实验记录跟随组件存放:

- `apps/go-api/note/go-api.md` — 探针命名坑、readiness 假绿、OOM、ConfigMap symlink 投影
- `apps/podinfo/note/lab.md` — 扩缩容与 Argo 同步行为观察
