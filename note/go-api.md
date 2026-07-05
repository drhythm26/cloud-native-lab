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

