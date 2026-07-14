# Cloud Native Lab 协作规则

> 所有 AI agent 进入仓库后，先读本文件，再读 [`docs/STATUS.md`](docs/STATUS.md)。

## 项目目标

这是一个面向学习和求职的 Cloud Native 实验仓，重点练习 Kubernetes 排障、可观测性和 GitOps。

```text
Terraform 创建 GKE
  ↓
bootstrap 安装 Argo CD
  ↓
App of Apps 管理 Argo CD、Prometheus 和 podinfo
```

- `apps/go-api/`：自研 Go 实验应用，CI 构建镜像，目前手动部署。
- `apps/podinfo/`：Helm Chart，由 Argo CD 自动同步。
- `gitops/`：Argo CD App of Apps 和平台组件配置。
- `infrastructure/gcp/`：GCP VPC 与 GKE Terraform。
- `sing-box/`：独立的 Kustomize 实验，不属于学习主线。

详细用法见各目录 `README.md`，排障记录见 `apps/*/note/`。

## 协作方式

- 沟通和文档使用中文，代码、命令和专有名词保留英文。
- 学习与排障任务优先讲清目标、原理、命令和输出看点，再让用户操作。
- 用户明确要求“直接做”时，可以实施修改并完成相应验证。
- 文档、注释、typo 和其他不影响学习过程的杂务可直接处理。
- 排障时优先教如何定位，提示按“方向 → 线索 → 关键命令 → 根因”递进。

## 部署与安全边界

### GitOps 自动同步

`main` 是 Argo CD 同步源和 CI 发布分支。以下内容推送到 `main` 后可能直接影响集群：

- `gitops/`
- `apps/applications/`
- `apps/podinfo/`

Argo CD 已开启 `prune` 和 `selfHeal`：

- 不要用 `kubectl edit` 长期修改 Argo CD 管理的资源。
- 删除、改名或更改 selector / namespace 前，说明可能的删除和中断影响。
- 未经用户要求，不主动 commit 或 push。

### 手动部署

- `apps/go-api/manifests/`：`kubectl apply -f`
- `sing-box/`：`kubectl apply -k`
- `scripts/bootstrap.sh`：安装 Argo CD 并注册根 Application

这些都会修改集群，不属于普通检查命令。

### 需要用户确认的操作

- `terraform apply` / `terraform destroy`
- `kubectl apply` / `delete` / `replace`
- `helm install` / `upgrade` / `uninstall`
- Argo CD sync / delete
- 删除 namespace、PVC、节点池或其他有状态资源
- commit 或 push 到 `main`

### 敏感信息

- 不提交真实密钥、Token、证书、kubeconfig、`tfstate` 或 `tfvars`。
- 不取消 `.gitignore` 中现有的敏感文件规则。
- 实验性假 Secret 可保留，但不得替换为真实凭证。
- Terraform state 目前保存在本地，不要删除或手动编辑。

## 修改与验证

- 只修改当前任务相关文件，不顺手大规模重构。
- 优先沿用仓库现有结构、命名和工具。
- 不擅自升级 Terraform Provider、Helm Chart、GitHub Action 或镜像大版本。
- Go 修改运行 `gofmt`、`go vet ./...` 和 `go test ./...`。
- Terraform 修改运行 `terraform fmt -check` 和 `terraform validate`。
- Helm 修改运行 `helm lint` 和 `helm template`。
- Shell 修改至少运行 `bash -n`。
- 仓库提供统一检查命令后，优先运行 `make check`。
- 不把 deploy、apply、sync 或 push 放入普通检查流程。

## 排障演练

只在用户明确要求时注入故障。故障必须：

- 限定在 `go-api` 或 `podinfo`，不影响 Argo CD、Prometheus、Terraform、网络或节点池。
- 单一、可控、可逆，并在结束后恢复干净状态。
- 用户主导排查，agent 递进提示。
- 复盘按“操作 / 现象 / 原因 / 结论”记入对应 `note/`。

## `docs/STATUS.md` 维护

- 只记录当前事实、已确认决策、路线进度、下一步和已知问题。
- 不复制 README 中的使用说明，不维护大段命令速查。
- 不确定的集群实时状态标记为“待确认”，不根据代码猜测。
- 完成阶段或改变路线时更新；普通 typo 不必写入日志。

## 开始工作前

1. 读 `AGENTS.md` 和 `docs/STATUS.md`。
2. 用 `git status` 和 `git log -10 --oneline` 了解本地状态。
3. 操作集群前确认 `kubectl config current-context`。
4. 操作 GCP 前确认当前账号、project 和 region/zone。
