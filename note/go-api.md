# go-api 实验记录

## 01 探针 port 未命名导致 Ready 失败

操作：
在 `go-api/manifests/deployment.yaml` 里给 Pod 加了 readiness/liveness 探针，探针写 `port: http`，但 `containers.ports` 里只有 `containerPort: 8080`，没有 `name: http`。

现象：

```sh
> kubectl get pods -n go-api
NAME                      READY   STATUS    RESTARTS   AGE
go-api-6cd9c47c75-pq269   0/1     Running   0          25m

> kubectl get events -n go-api
39m         Warning   Unhealthy           pod/go-api-6cd9c47c75-pq269    Readiness probe errored and resulted in unknown state: strconv.Atoi: parsing "http": invalid syntax
4m52s       Warning   Unhealthy           pod/go-api-6cd9c47c75-pq269    Liveness probe errored and resulted in unknown state: strconv.Atoi: parsing "http": invalid syntax

> kubectl get endpointslice -n go-api -l kubernetes.io/service-name=go-api -o yaml
conditions:
  ready: false      # Pod 未 Ready，Endpoint 不会对外 serving
  serving: false

> kubectl rollout status deploy -n go-api
error: deployment "go-api" exceeded its progress deadline
```

修复：
给 container port 加上名字，和探针里的 `port: http` 对应：

```yaml
ports:
  - name: http
    containerPort: 8080
readinessProbe:
  httpGet:
    path: /healthz
    port: http
livenessProbe:
  httpGet:
    path: /healthz
    port: http
```

结论：

- 探针 `port` 可以是**端口号**（如 `8080`）或**端口名**（如 `http`）；用名字时，必须在 `containers.ports` 里声明同名 `name`
- 没命名时 kubelet 会把 `"http"` 当数字解析 → `strconv.Atoi: parsing "http": invalid syntax`
- Pod 不 Ready → Service 的 EndpointSlice 里 `ready: false` / `serving: false`，流量不会打到该 Pod

## 02 readiness path 拼错仍 Ready（`/` 兜底导致探针假绿）

操作：
修复 port 命名后，readiness 探针 path 误写为 `/heathz`（少了一个 `l`），liveness 写的是正确的 `/healthz`。
`main.go` 里同时注册了 `/healthz` 和 `/`（rootHandler 兜底）。

现象：

```sh
> kubectl get deploy -n go-api -o yaml
readinessProbe:
  httpGet:
    path: /heathz      # 拼写错误
livenessProbe:
  httpGet:
    path: /healthz    # 正确

> kubectl get pods -n go-api
NAME                      READY   STATUS    RESTARTS   AGE
go-api-9fff95c8c-xxxxx    1/1     Running   0          ...

> kubectl rollout status deploy -n go-api
deployment "go-api" successfully rolled out
```

Pod 显示 `1/1 Ready`，rollout 也成功，但 readiness 实际没有测到 `/healthz`。

验证（port-forward 或 exec 进 Pod 后）：

```sh
> curl -s http://127.0.0.1:8081/heathz
{"server":"go-api","status":"0.1.0"}     # 走 rootHandler，不是 healthz

> curl -s http://127.0.0.1:8081/healthz
{"server":"healthz","status":"ok"}       # 才是探针 intended 的响应
```

原因：

```text
kubelet GET /heathz
  → 不匹配 /healthz handler
  → 无更具体匹配时，回退到注册的 "/"
  → rootHandler 返回 200
  → readiness 探针认为成功（HTTP 探针只看 2xx，不检查 body）
```

结论：

- 有 `/` 默认路由兜底时，**探针 path 配错也可能 Ready**，属于「假绿」
- liveness 和 readiness path 不一致时，两个探针测到的可能是不同 handler
- 规避：探针只用专用 path（如 `/healthz`）；生产里未知 path 返回 404 而非走 `/` 兜底；CI 里 curl 验证探针 path 的响应 body

修复：
readiness 与 liveness 统一为 `/healthz`：

```yaml
readinessProbe:
  httpGet:
    path: /healthz
    port: http
```

