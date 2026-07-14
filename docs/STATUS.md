# 项目状态

> 记录当前事实、学习进度和下一步。协作规则见 [`../AGENTS.md`](../AGENTS.md)。

**最后更新**：2026-07-14

**当前焦点**：完成 Kubernetes 基础排障练习，下一步进入 Prometheus、Grafana 和 PromQL。

## 当前系统

| 层级 | 组件 | 部署方式 | 状态 |
|---|---|---|---|
| 基础设施 | GCP VPC + zonal GKE | Terraform | 运行中，计费资源 |
| GitOps | Argo CD | bootstrap 后自管理 | 已部署 |
| 监控 | kube-prometheus-stack | Argo CD | 已部署 |
| 应用 | podinfo | Argo CD 自动同步 | 已部署 |
| 应用 | go-api | CI 构建 + 手动 apply | 已部署 |
| 实验 | sing-box | Kustomize 手动部署 | 已部署 |

架构详细见根目录和各组件 `README.md`。

## 已确定的约定

- `main` 是 GitOps 同步源和 go-api 镜像构建分支。
- Argo CD 使用 App of Apps，`prune` 和 `selfHeal` 均已开启。
- Argo CD 和 Prometheus 使用多源 Application 引用本仓库 values。
- podinfo 由 Argo CD 管理；go-api 和 sing-box 目前手动部署。
- go-api 镜像由 GitHub Actions 推送到 GHCR，Deployment 手动固定 commit SHA。
- Terraform state 目前保存在本地。

## 学习路线

### 1. Kubernetes 排障 `[~]`

- [x] probes、OOMKilled、FailedScheduling、FailedMount、ConfigMap 挂载
- [ ] 形成固定排障顺序：Pod 状态 → Events → describe → logs → YAML
- [ ] Service / EndpointSlice / selector 排障
- [ ] CrashLoopBackOff / ImagePullBackOff / Pending / Evicted 归类
- [ ] 完成至少 3 次“故障注入 → 定位 → 修复 → 复盘”

### 2. Prometheus / Grafana / PromQL `[ ]`

- [ ] 确认 go-api 和 podinfo Targets
- [ ] 掌握 `rate()`、`sum by()`、label 匹配和 `histogram_quantile()`
- [ ] 用 go-api 指标计算 QPS、错误率和 P95/P99 延迟
- [ ] 创建 go-api Grafana Dashboard
- [ ] 编写 PrometheusRule，理解 `for` / labels / annotations
- [ ] 掌握四大黄金信号、RED 和 USE

### 3. 应用与 Kubernetes 对象 `[ ]`

- [ ] 完善 go-api 业务路由和测试
- [ ] 掌握 Deployment、Service、ConfigMap、Secret、HPA、PDB
- [ ] 练习优雅终止和自动扩缩容
- [ ] 接入 MySQL 和 Redis，练习 StatefulSet、PVC 和连接排障

### 4. 交付与安全 `[ ]`

- [ ] Gateway API / Ingress 与 TLS
- [ ] go-api 纳入 GitOps 或自动更新镜像
- [ ] External Secrets / GCP Secret Manager
- [ ] NetworkPolicy、RBAC 和 Workload Identity

## 已掌握

- 理解 liveness 和 readiness 的用途，以及探针配置错误的排查方法。
- 理解 requests 影响调度、limits 影响运行时约束，排查过 OOMKilled 和 FailedScheduling。
- 理解 ConfigMap volume 的 symlink 投影机制。
- 理解 Argo CD App of Apps、多源 Application、`prune` 和 `selfHeal`。

## 下一步

1. 连接 Prometheus，查看 Targets 并确认 go-api / podinfo 抓取正常。
2. 用 go-api 现有指标练习 QPS、错误率和延迟 PromQL。
3. 建立统一的本地检查命令，再接入 Pull Request CI。

## 已知问题

- `gitops/README.md` 中“bootstrap 先安装旧版再升级”的说明已过时。
- Terraform state 仍在本地，后续考虑迁移到 GCS backend。
- 仓库暂无统一 `make check` 和完整的 Pull Request 检查。

## 最近记录

### 2026-07-14

- 统一 GitOps 和 CI 的发布分支为 `main`。
- 补充根目录及各模块 README。
- Argo CD Helm Chart 版本更新为 `10.1.3`。
- 建立 AI 协作规则和项目状态记录。
