# Cloud Native Release Tracker（云原生发布追踪系统）

![CI](https://github.com/drhythm26/cloud-native-lab/actions/workflows/ci.yaml/badge.svg)

> 一个内部**发布管理（release management）**服务，演示从本地 Go API + PostgreSQL，到容器化、GCP 基础设施即代码、Kubernetes/Helm、CI/CD、GitOps、可观测性的**完整云原生交付链路**。

---

## 架构

```text
开发者
  │ git push 代码
  ▼
GitHub Actions（CI）
  │ ① go test
  │ ② 构建镜像（tag = commit SHA），推送到 Artifact Registry
  │ ③ 把新 tag 写回 Helm values 并提交回 Git（带 [skip ci] 防 CI 自触发）
  ▼
Git（唯一事实来源 / 期望状态）
  │ Argo CD 监听并自动同步（prune + selfHeal）
  ▼
GKE / Kubernetes ── Deployment 滚动更新（从 Artifact Registry 拉取新镜像）
  │ Service（LoadBalancer）
  ▼
Release Tracker API ──SQL──► PostgreSQL
  ▲
  └─ Prometheus 抓取 /metrics
```

- **Git 是唯一事实来源（source of truth）。** 从 `git push` 到 Pod 上线全自动：CI 测试、构建并推送镜像，再把新镜像 tag 写回 Git；Argo CD 监听 Git 的期望状态并自动调谐（reconcile）到集群，属于拉取式部署（pull-based deployment）。每次上线都是一次可追溯、可 `git revert` 回滚的提交。
- **Prometheus** 通过 `ServiceMonitor` 抓取 API 的 `/metrics` 指标。

---

## 技术栈

| 层 | 技术 |
|---|---|
| API | Go（标准库 `net/http`，`pgx` 驱动） |
| 数据库 | PostgreSQL 16 |
| 容器 | Docker（多阶段构建、distroless、非 root） |
| 本地编排 | Docker Compose |
| 基础设施即代码 | Terraform（GCP provider） |
| 云平台 | Google Cloud Platform（GCP） |
| Kubernetes | GKE（VPC-native，独立节点池） |
| 镜像仓库 | Artifact Registry |
| 打包 | Helm |
| CI/CD | GitHub Actions（通过 Workload Identity Federation 实现无密钥认证） |
| GitOps | Argo CD |
| 可观测性 | Prometheus + Grafana（kube-prometheus-stack） |

---

## API

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/healthz` | 存活探针（liveness）——进程是否存活 |
| GET | `/readyz` | 就绪探针（readiness）——检查数据库连接 |
| GET | `/metrics` | Prometheus 指标 |
| POST | `/api/v1/releases` | 创建一条发布记录 |
| GET | `/api/v1/releases` | 查询发布记录列表 |
| GET | `/api/v1/releases/{id}` | 查询单条发布记录 |

发布记录示例：

```json
{
  "id": "rel_1718000000000000000",
  "serviceName": "payments",
  "version": "v1.4.2",
  "environment": "staging",
  "status": "pending",
  "owner": "alice",
  "createdAt": "2026-06-24T10:00:00Z",
  "updatedAt": "2026-06-24T10:00:00Z"
}
```

状态取值：`pending`、`deploying`、`succeeded`、`failed`、`rolled_back`。

---

## 目录结构

```text
cmd/api/             Go API 源码 + 测试
docker/postgres/     PostgreSQL 初始化 SQL
Dockerfile           多阶段构建（distroless 运行时）
docker-compose.yaml  本地 API + PostgreSQL 服务栈
terraform-gcp/       IaC：VPC、GKE、Artifact Registry、IAM、WIF
kubernetes/          原始 Kubernetes 清单
helm-charts/         Helm chart（参数化清单）
argo-cd/             Argo CD Application + 安装 values
prometheus/          ServiceMonitor + Prometheus values
scripts/bootstrap.sh 安装 Argo CD + kube-prometheus-stack
.github/workflows/   CI 流水线
```

---

## 本地运行

需要 Docker。

```bash
docker compose up -d --build

# 健康检查与就绪检查
curl localhost:8080/healthz
curl localhost:8080/readyz

# 创建一条发布记录
curl -X POST localhost:8080/api/v1/releases \
  -H 'Content-Type: application/json' \
  -d '{"serviceName":"payments","version":"v1.4.2","environment":"staging","owner":"alice"}'

# 查询发布记录（从 PostgreSQL 读回）
curl localhost:8080/api/v1/releases
```

数据存放在命名卷（named volume）中，因此 API 容器重启后数据仍然保留。

### 配置

API 完全通过环境变量配置（可以干净地映射到 Kubernetes 的 ConfigMap + Secret）：

| 变量 | 默认值 | 说明 |
|---|---|---|
| `PORT` | `8080` | HTTP 监听端口 |
| `DB_HOST` | `localhost` | PostgreSQL 主机 |
| `DB_PORT` | `5432` | |
| `DB_NAME` | `release_tracker` | |
| `DB_USER` | `release_tracker` | |
| `DB_PASSWORD` | `release_tracker` | 在 Kubernetes 上存放于 Secret |
| `DB_SSLMODE` | `disable` | |

---

## 部署到 GCP / GKE

### 1. 用 Terraform 创建基础设施

创建一个带 VPC-native 二级网段的自定义 VPC、一个带独立节点池的 GKE 集群、一个 Artifact Registry 仓库、一个最小权限的节点服务账号，以及供 GitHub Actions 使用的 Workload Identity Federation。

```bash
cd terraform-gcp
terraform init
terraform plan      # 审查将要创建的资源
terraform apply
```

### 2. 连接 kubectl

```bash
gcloud container clusters get-credentials release-tracker --zone asia-east2-a
kubectl get nodes
```

### 3. 构建并推送镜像

每次推送到 `main`，CI 会自动完成这一步（见下文）。手动操作：

```bash
REGION=asia-east2; PROJECT=dev-2026522; REPO=release-tracker
gcloud auth configure-docker $REGION-docker.pkg.dev
docker build -t $REGION-docker.pkg.dev/$PROJECT/$REPO/release-tracker-api:v0.1.0 .
docker push $REGION-docker.pkg.dev/$PROJECT/$REPO/release-tracker-api:v0.1.0
```

### 4. 用 Helm 部署

```bash
helm upgrade --install release-tracker helm-charts/release-tracker
```

API 通过 `LoadBalancer` 类型的 Service 暴露，用 `kubectl get svc -n release-tracker` 查看外部 IP。

---

## CI/CD

GitHub Actions（`.github/workflows/ci.yaml`）：

1. **test** —— 运行 `go test ./...`。
2. **build-push** —— 通过 **Workload Identity Federation** 认证到 GCP（**无需长期有效的服务账号密钥**），构建以短 commit SHA 为 tag 的镜像并推送到 Artifact Registry。
3. **写回** —— 把新镜像 tag 写回 `helm-charts/release-tracker/values.yaml` 并提交回 Git。提交信息带 `[skip ci]`，避免这次自动提交再次触发 CI 形成死循环。

用 WIF 认证意味着 CI 用一个**短期有效的 GitHub OIDC 令牌**换取 GCP 凭证，没有静态 secret 会泄露或需要轮换。镜像 tag 一旦写回 Git，剩下的上线交给 Argo CD（见下）——CI 只负责"构建"，部署由 GitOps 负责。

---

## GitOps（Argo CD）

`argo-cd/application.yaml` 定义了一个 Argo CD `Application`，指向本仓库里的 Helm chart，并启用自动同步策略（`prune` + `selfHeal`）：

- **全自动交付闭环** —— `git push` 代码 → CI 构建镜像并写回新 tag → Argo CD 检测到 Git 变化自动同步 → Deployment 滚动更新拉取新镜像，整条链路无需任何手动 `kubectl` 操作。
- **拉取式部署** —— 集群主动从 Git 拉取期望状态，而不是 CI 把变更推进集群。
- **漂移纠正** —— 手动改动被 Argo 管理的资源后，会被自动还原回 Git 声明的状态。

安装 Argo CD（和 Prometheus）：

```bash
./scripts/bootstrap.sh
```

---

## 可观测性

`prometheus/` 安装 `kube-prometheus-stack`，并注册一个 `ServiceMonitor`，每 30 秒抓取一次 API 的 `/metrics` 端点，指标在内置的 Grafana 中可视化。

---

## 这个项目展示了什么

- 用小巧、非 root 的 distroless 镜像容器化 Go 服务。
- 在部署任何东西之前，先用可复现的代码（Terraform）描述云基础设施。
- 在托管 Kubernetes 上运行服务，配有健康探针、资源请求/限制、ConfigMap/Secret 分离和持久化存储。
- 用 Helm 把清单参数化，支持按环境区分。
- 无密钥 CI/CD 与基于 GitOps 的交付，具备漂移纠正能力。
- 以指标抓取作为生产排障的基础。
