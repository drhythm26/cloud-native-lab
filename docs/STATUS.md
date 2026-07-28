# 项目状态

> 只记当前事实、路线进度、下一步和已知问题。安全规则见 [`../AGENTS.md`](../AGENTS.md)，历史变更看 git log。

**最后更新**：2026-07-28

**当前焦点**：podinfo K8s 生产化（HPA、PDB、故障演练）。

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

- `main` 是 GitOps 同步源；Argo CD App of Apps，`prune` + `selfHeal` 开启。
- podinfo 改 `apps/podinfo/values.yaml` → push → Argo 自动同步。
- Terraform state 保存在本地。

## 路线进度

1. **K8s 应用生产化** `[ ]`：podinfo 多副本 + 滚动观察、HPA、PDB。
2. **排障与可观测** `[~]`：排障/PromQL/告警已练；剩故障演练、PromQL 复测。
3. **平台开发** `[ ]`：kubebuilder 最小 CRD + reconcile；controller 自带 metrics。
4. **可选后续**：Ingress/TLS、Redis、NetworkPolicy/RBAC、EFK/Loki、`make check` + PR CI、根 README 整理。

## 下一步

1. podinfo 开 HPA + PDB，压测观察扩缩容。
2. podinfo 故障演练（faults / delay / status），笔记写入 `apps/podinfo/note/lab.md`。
3. PromQL 复测。

## 已知问题

- Terraform state 在本地，未迁 GCS backend。
- 无统一 `make check` 和 PR CI。
- Grafana Dashboard 未自动供应。