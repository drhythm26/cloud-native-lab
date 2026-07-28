# Cloud Native Lab 安全规则

学习求职向的 Cloud Native 实验仓：Terraform 建 GKE → bootstrap 装 Argo CD → App of Apps 管平台组件与应用。进度和下一步见 [`docs/STATUS.md`](docs/STATUS.md)，各目录用法见对应 README，排障记录在 `apps/*/note/`。

## 安全边界

push 到 `main` 即部署：`gitops/`、`apps/applications/`、`apps/podinfo/` 由 Argo CD 自动同步（`prune` + `selfHeal` 已开），不要 `kubectl edit` 这些资源。`apps/go-api/manifests/` 和 `sing-box/` 手动 apply。

以下操作必须先经用户确认：

- `terraform apply` / `destroy`
- `kubectl apply` / `delete` / `replace` / `edit` / `patch` / `scale`
- `helm install` / `upgrade` / `uninstall`，Argo CD sync / delete
- 删除 namespace、PVC、节点池等有状态资源
- `git commit` / `push`

敏感信息：不提交真实密钥、Token、证书、kubeconfig、`tfstate`、`tfvars`；不放宽 `.gitignore` 现有规则；本地 Terraform state 不删不改。

## 故障注入

仅在用户明确要求时注入，限 `go-api` 或 `podinfo`；单一、可控、可逆；注入前预写复原命令，结束后只读核验恢复干净。
