# 项目状态

> 记录当前事实、学习进度和下一步。协作规则见 [`../AGENTS.md`](../AGENTS.md)。

**最后更新**：2026-07-15

**当前焦点**：冲刺简历作品集。第一梯队（Grafana Dashboard、PrometheusRule 告警、go-api 中间件改造、监控卫生）完成即开始投递，SRE/云原生 与 Go 后端两个方向并投。

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
- [x] Service / EndpointSlice / selector 排障（演练：selector 改错 → Endpoint 空）
- [ ] CrashLoopBackOff / ImagePullBackOff / Pending / Evicted 归类
- [ ] 完成至少 3 次“故障注入 → 定位 → 修复 → 复盘”（已完成 1/3）

### 2. Prometheus / Grafana / PromQL `[~]`

- [x] 确认 go-api 和 podinfo Targets
- [x] 掌握 `rate()`、`sum by()` / `without`、label 匹配、`histogram_quantile()`、`increase()`、`topk()`
- [x] 用 go-api / podinfo 指标计算 QPS、错误率和 P95 延迟
- [x] 掌握 RED 方法（QPS / 错误率 / 延迟全程手写）
- [ ] 创建 go-api Grafana Dashboard
- [ ] 编写 PrometheusRule，理解 `for` / labels / annotations
- [ ] 掌握四大黄金信号和 USE

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
- 掌握 PromQL 数据模型与核心函数：序列 = 指标名 + label 集合；`rate` / `increase` / `sum by` / `without` / `histogram_quantile`；向量除法按 label 集合配对，不匹配的序列被静默丢弃。
- 会读监控图：no data 和 0 是两种状态；爬坡宽度 ≈ rate 窗口宽度；histogram 分位数是桶内线性插值的估计值，精度受 bucket 布局限制。
- 排障方法论：顺数据流逐环检查（Prometheus → ServiceMonitor → Service → Endpoint → Pod）；比值异动先拆分子分母；`kubectl get endpoints` 是 Service 链路的分流判断点。
- 理解 GitOps 修复原则：集群状态与 Git 对齐，修复用 manifests apply 而非手动 patch；`kubectl port-forward` 钉死单个 pod，不做负载均衡。

## 下一步

冲刺至首投（第一梯队）：

1. Grafana RED Dashboard（go-api / podinfo），完成后导出 JSON 入库（dashboard as code）。
2. PrometheusRule 三条告警（错误率、P95、target down），理解 `for` / labels / annotations。
3. go-api metrics 中间件改造（统计 404、路由模板防基数爆炸），走完整交付链（CI 构建 → pin SHA → apply）。
4. 监控卫生：Grafana 收回 ClusterIP、关闭 GKE 托管控制面无效抓取（values 已改，待 push 生效）。

首投后（第二梯队）：go-api 纳入 Argo、MySQL（StatefulSet + PVC）+ Redis 接入、HPA / PDB / 优雅终止、`make check` + PR CI、根 README 作品集化。
排障演练（剩 2 次）与 note 07 复盘在面试前补足。

## 已知问题

- Terraform state 仍在本地，后续考虑迁移到 GCS backend。
- 仓库暂无统一 `make check` 和完整的 Pull Request 检查。
- 监控卫生问题（Grafana 公网默认口令、GKE 托管控制面抓取常驻 DOWN）已在 values 修复，待 commit + push 后由 Argo 同步生效。

## 最近记录

### 2026-07-15

- 打通监控链路：Targets 确认 go-api / podinfo UP，识别 GKE 托管控制面的 0/0 目标为预期噪音。
- 手写 RED 三件套 PromQL；两轮闭卷练习共 11 题，纠正“指标名 typo 静默返回空”“通用指标名要圈 namespace”两个习惯。
- 一次真实排障：错误率无故下台阶 → 拆分子分母 → 定位为 delay 序列迟到 9 分钟入场稀释分母。
- 第一次故障注入演练：Service selector 改错 → Endpoint `<none>` → 流量与指标双断；按 GitOps 原则用 manifests apply 修复，已确认还原干净。

### 2026-07-14

- 统一 GitOps 和 CI 的发布分支为 `main`。
- 补充根目录及各模块 README。
- Argo CD Helm Chart 版本更新为 `10.1.3`。
- 建立 AI 协作规则和项目状态记录。
