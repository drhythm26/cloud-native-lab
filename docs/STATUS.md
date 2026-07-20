# 项目状态

> 只记当前事实、路线进度、下一步和已知问题。协作规则见 [`../AGENTS.md`](../AGENTS.md)，历史变更看 git log。

**最后更新**：2026-07-20

**当前焦点**：冲刺简历首投（SRE/云原生 与 Go 后端并投）。第一梯队剩余：go-api 中间件改造。PrometheusRule 已结课（2026-07-19，错题 4 条待复测，见 `docs/mistakes.md`）。

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
- go-api RED Dashboard JSON：`gitops/prometheus/dashboard/go-api-red.json`（手动 Import，未自动供应）。

## 学习路线

1. **K8s 排障** `[~]`：probes / OOMKilled / FailedScheduling / 挂载 / Service-selector 已练；剩故障演练 2 次、异常状态归类、固化排障顺序。
2. **Prometheus / PromQL** `[~]`：PromQL 核心函数、RED 手写、Grafana Dashboard、PrometheusRule 已完成；黄金信号 / USE 未学。
3. **应用与 K8s 对象** `[ ]`：go-api 路由与测试、Deployment/Service/ConfigMap/Secret/HPA/PDB、优雅终止、MySQL + Redis（StatefulSet / PVC）。
4. **交付与安全** `[ ]`：Ingress / TLS、go-api 纳入 GitOps、External Secrets、NetworkPolicy / RBAC / Workload Identity。

## 已掌握（面试复习索引）

- K8s：liveness / readiness 排查、requests 调度与 limits 运行时约束、OOMKilled、ConfigMap symlink 投影、Service → Endpoint 链路排障。
- GitOps：App of Apps、多源 Application、prune / selfHeal 行为；修复对齐 Git 而非手动 patch。
- PromQL：rate / increase / sum by / without / histogram_quantile；向量除法按 label 配对；「空 vs 0」与指标名 typo 静默返回空。
- RED 实战：手写 QPS / 5xx 错误率 / P95，Grafana 三面板；基线 QPS 与探针周期交叉验证。
- 排障方法论：顺数据流逐环检查（Prometheus → ServiceMonitor → Service → Endpoint → Pod）；比值异动先拆分子分母。
- 告警（PrometheusRule 已结课）：三条规则亲手写、注入验证 pending → firing、for 与求值节奏（测验全对）；薄弱点在错题本待复测：up 自报机制、absent 兜底、release label 加载暗号、空结果不参与比较。

## 下一步

1. go-api metrics 中间件（统计 404、路由模板防基数爆炸），走完整交付链：CI 构建 → pin SHA → apply。
2. 首投前补齐：剩余 2 次故障演练；错题本 4 条复测。
3. 首投后：go-api 纳入 Argo、MySQL / Redis、HPA / PDB / 优雅终止、`make check` + PR CI、根 README 作品集化。

## 已知问题

- Terraform state 在本地，未迁 GCS backend。
- 无统一 `make check` 和 PR CI。
- Grafana Dashboard 未自动供应，换环境需手动 Import。
- 本地领先 origin/main 3 个 commit，待 push。
