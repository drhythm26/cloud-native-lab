# Cloud Native Lab 协作规则

学习求职向的 Cloud Native 实验仓：Terraform 建 GKE → bootstrap 装 Argo CD → App of Apps 管平台组件与应用。进度和下一步见 [`docs/STATUS.md`](docs/STATUS.md)，各目录用法见对应 README，排障记录在 `apps/*/note/`。

## 协作方式

- 中文沟通，代码、命令和专有名词保留英文。
- 学习优先：排障和教学产出让用户亲手做，按「方向 → 线索 → 关键命令 → 根因」递进提示，不直接给根因。用户明确说「直接做」时才代劳并完成验证。
- `apps/*/note/` 的笔记由用户写，不代笔；Claude 只做只读核验和指出遗漏。
- 文档、注释、typo 等杂务可直接处理。

## 安全边界

push 到 `main` 即部署：`gitops/`、`apps/applications/`、`apps/podinfo/` 由 Argo CD 自动同步（`prune` + `selfHeal` 已开），不要 `kubectl edit` 这些资源。`apps/go-api/manifests/` 和 `sing-box/` 手动 apply。

以下操作必须先经用户确认：

- `terraform apply` / `destroy`
- `kubectl apply` / `delete` / `replace` / `edit` / `patch` / `scale`
- `helm install` / `upgrade` / `uninstall`，Argo CD sync / delete
- 删除 namespace、PVC、节点池等有状态资源
- `git commit` / `push`

敏感信息：不提交真实密钥、Token、证书、kubeconfig、`tfstate`、`tfvars`；不放宽 `.gitignore` 现有规则；本地 Terraform state 不删不改。

## 修改与验证

- 只改当前任务相关文件，不顺手重构；不擅自升级 Provider、Chart、Action 或镜像大版本。
- Go 跑 `gofmt` / `go vet ./...` / `go test ./...`；Terraform 跑 `fmt -check` / `validate`；Helm 跑 `lint` / `template`；Shell 跑 `bash -n`。
- deploy / apply / sync / push 不进普通检查流程。

## 故障注入

仅在用户明确要求时注入，限 `go-api` 或 `podinfo`；单一、可控、可逆；注入前预写复原命令，结束后只读核验恢复干净。

## 开始工作前

1. 读本文件和 `docs/STATUS.md`；用 `git status` 和 `git log --oneline` 了解本地状态。
2. 操作集群前确认 `kubectl config current-context`；操作 GCP 前确认账号和 project。
