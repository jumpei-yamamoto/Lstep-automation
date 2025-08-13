variable "project_name" { type = string }
variable "region"       { type = string }

# GitHub OIDC 用
variable "github_org"   { type = string }
variable "github_repo"  { type = string } # 例: my-org/my-repo の "my-repo" ではなくリポジトリ名のみ

# フロント（S3+CF）
variable "enable_frontend_static" {
  type    = bool
  default = true
}

# DB スペック（簡易）
variable "db_instance_class" {
  type    = string
  default = "db.t4g.micro"
}
variable "db_allocated_storage" {
  type    = number
  default = 20
}