# infrastructure/gcp — Terraform

GKE 实验集群及其网络的 Terraform 定义。

## 资源清单

| 文件 | 资源 |
|------|------|
| `apis.tf` | 启用项目 API(compute、container) |
| `network.tf` | 自定义 VPC `cloud-native-lab-vpc` + 子网 `10.0.0.0/20`,secondary ranges:pods `10.16.0.0/14`、services `10.20.0.0/20` |
| `gke.tf` | VPC-native GKE 集群 `cloud-native-lab-gke`(zonal,`asia-east2-a`,移除默认节点池)+ 节点池:1 × `e2-standard-2`,30GB pd-standard |

`deletion_protection = false`——实验集群,允许 `terraform destroy`。

## 变量

| 变量 | 默认 | 说明 |
|------|------|------|
| `project` | (必填) | GCP 项目 ID |
| `region` | `asia-east2` | |
| `zone` | `asia-east2-a` | |

`project` 通过 `terraform.tfvars` 提供(**不入库**,见下)。

## 使用

```bash
gcloud auth application-default login
terraform init
terraform plan
terraform apply
# 完成后获取 kubeconfig:
gcloud container clusters get-credentials cloud-native-lab-gke --zone asia-east2-a
```

## 敏感文件约定

`*.tfvars`、`*.tfstate`、`.terraform/` 均被 `.gitignore` 排除,**不得提交**;`*.tfvars.example` 允许入库作模板。state 目前在本地——换机器前注意备份或迁移到远端 backend。
