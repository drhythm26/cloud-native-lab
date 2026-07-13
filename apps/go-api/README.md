# go-api

自研 Go HTTP 服务,用来练习 K8s 部署、配置注入与 Prometheus 指标接入。监听 `:8080`。

## 路由

| 路径 | 说明 |
|------|------|
| `/healthz` | liveness 探针,返回 `{"status":"ok"}` |
| `/readyz` | readiness 探针 |
| `/env` | 回显 `APP_ENV` / `LOG_LEVEL` / `API_TOKEN` 环境变量(来自 ConfigMap / Secret) |
| `/file` | 读取 `/etc/go-api` 下挂载的 ConfigMap 文件内容 |
| `/metrics` | Prometheus 指标 |

注意:**没有 `/` 根路由**——未知路径返回 404,避免探针 path 写错时被根路由兜底造成"假绿"(踩坑记录见 `note/go-api.md`)。

## 指标

- `go_api_http_requests_total{path, method, code}` — Counter,请求计数
- `go_api_http_request_duration_seconds{path}` — Histogram,请求延迟(默认桶)

ServiceMonitor 抓取 `/metrics`,间隔 15s,带 `release: prometheus` 标签(kube-prometheus-stack 的选择器要求)。

## 构建与发布

- 多阶段 Dockerfile:`golang:1.26-alpine` 编译 → `distroless/static` 运行,nonroot
- CI(`.github/workflows/go-api.yaml`):push 到 `main` 且 `apps/go-api/**` 有变更时,构建并推送 `ghcr.io/drhythm26/go-api:<commit-sha>`
- Deployment 中镜像 tag **手动 pin 到具体 sha**,升级时改 `manifests/deployment.yaml` 再 apply

本地运行:`docker compose up`(见 `docker-compose.yaml`,挂载 `config/` 目录)。

## 部署(手动,不走 Argo)

```bash
kubectl apply -f manifests/
```

manifests 包含:Namespace、Deployment、Service、ConfigMap、Secret、ServiceMonitor,全部在 `go-api` namespace。

Deployment 要点:

- 端口命名 `http`(8080),探针 `port: http` 依赖这个名字——**端口名和探针引用必须一致**
- 环境变量:`APP_ENV` / `LOG_LEVEL` 来自 ConfigMap `go-api-config`,`API_TOKEN` 来自 Secret `go-api-secret`(实验用假值)
- ConfigMap 的 `greeting.txt` 挂载到 `/etc/go-api/`(K8s 用 symlink 投影,`/file` 接口处理过这个细节)
- resources:requests 50m/64Mi,limits 200m/128Mi

## 验证

```bash
kubectl -n go-api port-forward svc/go-api 8080:8080
curl -s localhost:8080/healthz
curl -s localhost:8080/metrics | grep go_api
```

## 踩坑记录

见 [`note/go-api.md`](note/go-api.md):探针端口未命名导致 Ready 失败、readiness path 拼错仍 Ready(假绿)、memory limit 过低 OOMKilled、requests 过大 Pending、ConfigMap 不存在 FailedMount、ConfigMap symlink 投影读取。
