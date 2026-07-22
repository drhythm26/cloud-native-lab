# 项目状态

> 只记当前事实、路线进度、下一步和已知问题。协作规则见 [`../AGENTS.md`](../AGENTS.md)，历史变更看 git log。

**最后更新**：2026-07-20

**当前焦点**：平台开发 / SRE 岗。主线 Go → K8s 生产化 → operator；旁线排障与可观测复测。

## 当前系统

| 层级 | 组件 | 部署方式 | 状态 |
|---|---|---|---|
| 基础设施 | GCP VPC + zonal GKE | Terraform | 运行中，计费资源 |
| GitOps | Argo CD | bootstrap 后自管理 | 已部署 |
| 监控 | kube-prometheus-stack | Argo CD | 已部署 |
| 应用 | podinfo | Argo CD 自动同步 | 已部署 |
| 应用 | go-api | CI 构建 + 手动 apply | 已部署 |
| 实验 | sing-box | Kustomize 手动部署 | 已部署 |

## 已确定的约定

- `main` 是 GitOps 同步源和 go-api 镜像构建分支；Argo CD App of Apps，`prune` + `selfHeal` 开启。
- go-api 镜像由 GitHub Actions 推 GHCR，Deployment 手动 pin commit SHA。
- Terraform state 保存在本地。
- go-api RED Dashboard：`gitops/prometheus/dashboard/go-api-red.json`（手动 Import）。

## 学习路线

1. **Go 服务深化** `[~]`：metrics 中间件已上线；待 `context`/超时、并发安全、单测、优雅终止。
2. **K8s 应用生产化** `[ ]`：多副本 + 滚动观察、HPA、PDB；go-api 纳入 GitOps。
3. **平台开发** `[ ]`：kubebuilder 最小 CRD + reconcile（管 go-api）；controller 自带 metrics。
4. **排障与可观测** `[~]`：排障/PromQL/告警已练；剩故障演练 2 次、错题复测；黄金信号/USE 可选。

首投后可选（不占主线）：MySQL/Redis、Ingress/TLS、External Secrets、NetworkPolicy/RBAC/WI、`make check` + PR CI、根 README 作品集化。

## 已掌握（面试复习索引）

- K8s 排障：探针、requests/limits、OOMKilled、ConfigMap 投影、Service → Endpoint。
- GitOps：App of Apps、多源 Application、prune/selfHeal；修复对齐 Git。
- PromQL / RED / 告警：核心函数、手写 RED、Grafana、PrometheusRule（pending → firing）。
- 排障方法论：顺数据流逐环查；比值先拆分子分母。
- Go metrics：包整个 ServeMux、`r.Pattern` 防基数爆炸、跳过 `/metrics`；CI → pin SHA → apply。

## 下一步

1. go-api：`context` + 并发安全 + 测试。
2. go-api：replicas ≥ 2 + HPA 压测观察。
3. 最小 operator（kubebuilder）：CR 管 go-api Deployment。
4. 穿插：故障演练 2 次；错题本复测。

## 已知问题

- Terraform state 在本地，未迁 GCS backend。
- 无统一 `make check` 和 PR CI。
- Grafana Dashboard 未自动供应。
