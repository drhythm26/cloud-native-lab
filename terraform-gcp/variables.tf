variable "project" {
  type = string
}

variable "region" {
  type    = string
  default = "asia-east2"
}

variable "zone" {
  type    = string
  default = "asia-east2-a"
}

variable "cluster_name" {
  type    = string
  default = "release-tracker"
}

variable "repo_id" {
  type    = string
  default = "release-tracker"
}

variable "node_count" {
  type    = number
  default = 2
}

variable "machine_type" {
  type    = string
  default = "e2-standard-2"
}

variable "environment" {
  type    = string
  default = "dev"
}