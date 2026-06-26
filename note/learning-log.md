# 01 docker-compose-lab

## 完整配置

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

## 实验
```sh
更改./docker-compose.yaml中
DB_HOST: postgres
DB_HOST: wrong-host
会导致api 连不上数据库直接退出
原因：main() 在启动HTTP之前就是尝试连接数据库，失败就Fatalf，所以8080从未监听
现象：
curl 127.0.0.1:8080/healthz # 会连不上
curl: (7) Failed to connect to 127.0.0.1 port 8080 after 4 ms: Could not connect to server
ss -tnlp | grep :8080 # 空的
排障：
docker compose ps # 发现api没有起来

docker compose logs api # 查看具体原因，数据库连接失败，hostname resolving error 明显是主机名解析问题
lan@DESKTOP-7OAIVON:~/cloud-native-lab$ docker compose logs api
api-1  | 2026/06/26 12:38:54 connect database failed: failed to connect to `user=release_tracker database=release_tracker`: hostname resolving error: lookup wrong-host on 127.0.0.11:53: server misbehaving
```