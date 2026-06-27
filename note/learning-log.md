# docker-compose-lab

完整配置：

```yaml
services:
  postgres:
    image: postgres:17
    environment:
      POSTGRES_DB: release_tracker
      POSTGRES_USER: release_tracker
      POSTGRES_PASSWORD: release_tracker
    ports:
      - "5432:5432"
    volumes:
      - postgres-data:/var/lib/postgresql/data
      - ./docker/postgres/init.sql:/docker-entrypoint-initdb.d/init.sql:ro
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U release_tracker -d release_tracker"]
      interval: 5s
      timeout: 3s
      retries: 10

  api:
    image: release-tracker-api:dev
    build:
      context: .
    environment:
      PORT: "8080"
      DB_HOST: postgres
      DB_PORT: "5432"
      DB_USER: release_tracker
      DB_NAME: release_tracker
      DB_PASSWORD: release_tracker
      DB_SSLMODE: disable
    ports:
      - "8080:8080"
    depends_on:
      postgres:
        condition: service_healthy

volumes:
  postgres-data:
```

## 1 hostname resolving error

更改`./docker-compose.yaml`中的`DB_HOST`
```yaml
DB_HOST: postgres
DB_HOST: wrong-host # 更改后 
```
现象：会导致api 连不上数据库直接退出
```bash
curl 127.0.0.1:8080/healthz # 会连不上
curl: (7) Failed to connect to 127.0.0.1 port 8080 after 4 ms: Could not connect to server
ss -tnlp | grep :8080 # 空的
```
原因：改变`DB_HOST`导致无法正确解析到ip，`main()`在启动HTTP之前就尝试连接数据库，失败就`Fatalf`，所以`8080`从未监听
排障：
```bash
docker compose ps # 查看api容器是否在运行
```
```bash
docker compose logs api # 看日志
```
```bash
# 数据库连接失败，hostname resolving error 是数据库主机名解析问题，排查DB_HOST
lan@DESKTOP-7OAIVON:~/cloud-native-lab$ docker compose logs api
api-1  | 2026/06/26 12:38:54 connect database failed: failed to connect to `user=release_tracker database=release_tracker`: hostname resolving error: lookup wrong-host on 127.0.0.11:53: server misbehaving
```
docker-compose 中service名就是dns记录，容器间通过名字访问


## 02 undefined volume

在`docker-compose.yaml`中使用卷挂载，如果不使用顶级`volumes`声明，在`docker compose up`就会报错
```bash
lan@DESKTOP-7OAIVON:~/cloud-native-lab$ docker compose up -d
service "postgres" refers to undefined volume postgres-data: invalid compose project
```