# sing-box

在 GKE 上部署 [sing-box](https://github.com/SagerNet/sing-box)(VLESS + Reality)的实验清单。kustomize 组织,手动部署,不走 Argo。

## 结构

```
├── kustomization.yaml        # namespace sing-box + secretGenerator
├── config/
│   └── config.json.template  # 服务端配置模板(占位符注入密钥)
├── manifests/                # Namespace / Deployment / Service(LoadBalancer)
├── secrets/                  # 生成的密钥与渲染后配置(整目录 gitignore)
└── gen_config.sh             # 密钥生成 / 配置渲染 / 客户端配置生成
```

## 使用流程

```bash
# 1. 生成密钥(uuid + reality keypair + short_id,写入 secrets/sing-box.env,已存在则复用)
./gen_config.sh env

# 2. 渲染服务端配置(envsubst 模板 → secrets/config.json;需要 gettext 的 envsubst)
./gen_config.sh render

# 3. 部署(kustomize secretGenerator 把 secrets/config.json 打进 Secret)
kubectl apply -k .

# 4. 等 Service 拿到外网 IP 后,生成本机客户端 outbound 配置
#    (写入 ~/.config/sing-box/05-outbound.json)
./gen_config.sh config
```

## 要点

- 密钥用官方镜像生成:`docker run ghcr.io/sagernet/sing-box generate uuid / reality-keypair`
- `secrets/` 整目录被 `.gitignore` 排除(`**/secrets/`),**密钥和渲染后的配置不入库**;仓库里只有不含密钥的模板
- 配置有变更时:重新 `render` → `kubectl apply -k .`(secretGenerator 带内容 hash,Secret 变更会触发滚动更新)
- Service 是 LoadBalancer,客户端连 `EXTERNAL-IP:443`
