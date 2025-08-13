# Terraform AWS インフラ構築ガイド

このドキュメントでは、AWS上での本番環境構築についてTerraformを使用した手順を説明します。

## 概要

Terraformにより、以下のAWSリソースを自動構築します：
- **ECS Fargate** - コンテナ化されたバックエンドの実行環境
- **RDS PostgreSQL** - マネージドデータベース
- **S3 + CloudFront** - 静的フロントエンドの配信
- **ALB** - ロードバランサー
- **ECR** - コンテナイメージレジストリ
- **Secrets Manager** - 機密情報の管理

## ファイル構成

```
infra/terraform/
├── main.tf              # メインリソース定義
├── variables.tf         # 変数定義
├── outputs.tf           # 出力値定義
├── providers.tf         # プロバイダー設定
└── env/
    └── dev.tfvars       # 開発環境用変数
```

## 前提条件

### 1. AWSアカウントの準備
- AWSアカウントの作成
- IAMユーザーの作成（管理権限付与）
- アクセスキーの取得

### 2. ツールのインストール
```bash
# Terraformのインストール
brew install terraform  # macOS
# or
curl -fsSL https://apt.releases.hashicorp.com/gpg | sudo apt-key add -
sudo apt-add-repository "deb [arch=amd64] https://apt.releases.hashicorp.com $(lsb_release -cs) main"
sudo apt-get update && sudo apt-get install terraform

# AWS CLIのインストール
pip install awscli
# or
brew install awscli
```

### 3. AWS認証設定
```bash
# AWS CLIの設定
aws configure
# AWS Access Key ID: [YOUR_ACCESS_KEY]
# AWS Secret Access Key: [YOUR_SECRET_KEY]
# Default region name: ap-northeast-1
# Default output format: json

# 認証確認
aws sts get-caller-identity
```

## デプロイ手順

### 1. Terraform初期化

```bash
# infra/terraformディレクトリに移動
cd infra/terraform

# Terraformの初期化
terraform init
```

### 2. 環境変数の設定

```bash
# env/dev.tfvarsファイルを確認・編集
cat env/dev.tfvars
```

現在の設定:
```hcl
project_name = "lstep-automation"
region = "ap-northeast-1"
db_instance_class = "db.t3.micro"
db_allocated_storage = 20
enable_frontend_static = true
```

### 3. プランの確認

```bash
# デプロイプランの確認
terraform plan -var-file=env/dev.tfvars

# プラン結果の保存（推奨）
terraform plan -var-file=env/dev.tfvars -out=tfplan
```

### 4. インフラストラクチャのデプロイ

```bash
# 実際のデプロイ実行
terraform apply -var-file=env/dev.tfvars

# または事前に保存したプランを使用
terraform apply tfplan
```

**注意**: 初回デプロイには15-20分程度かかります。

## デプロイ後の設定

### 1. 出力情報の確認

```bash
# デプロイ結果の確認
terraform output

# 個別の出力値確認
terraform output ecr_backend_url
terraform output alb_dns_name
terraform output cloudfront_domain_name
```

### 2. ECRへのイメージプッシュ

```bash
# ECRログイン
aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin $(terraform output -raw ecr_backend_url)

# バックエンドイメージのビルドとプッシュ
cd ../../  # プロジェクトルートに戻る
docker build -f backend/Dockerfile -t lstep-backend .
docker tag lstep-backend:latest $(terraform output -raw ecr_backend_url):latest
docker push $(terraform output -raw ecr_backend_url):latest

# フロントエンド（SSR）イメージのプッシュ（必要に応じて）
docker build -f frontend/Dockerfile.ssr -t lstep-frontend-ssr .
docker tag lstep-frontend-ssr:latest $(terraform output -raw ecr_frontend_ssr_url):latest
docker push $(terraform output -raw ecr_frontend_ssr_url):latest
```

### 3. ECSサービスの更新

```bash
# ECSサービスの強制更新（新しいイメージをデプロイ）
aws ecs update-service \
  --cluster lstep-automation-cluster \
  --service lstep-automation-backend \
  --force-new-deployment \
  --region ap-northeast-1
```

### 4. フロントエンド（静的サイト）のデプロイ

```bash
# Next.jsの静的書き出し
cd frontend
npm run build
npm run export

# S3への同期
aws s3 sync out/ s3://$(terraform output -raw s3_frontend_bucket)/

# CloudFrontキャッシュの無効化
aws cloudfront create-invalidation \
  --distribution-id $(terraform output -raw cloudfront_distribution_id) \
  --paths "/*"
```

## アクセス確認

### バックエンドAPI
```bash
# ALB経由でのアクセステスト
curl http://$(terraform output -raw alb_dns_name)/healthz
```

### フロントエンド
```bash
# CloudFront経由でのアクセス
open https://$(terraform output -raw cloudfront_domain_name)
```

## リソース詳細

### ネットワーク構成
- **VPC**: デフォルトVPCを使用（本番では専用VPC推奨）
- **サブネット**: マルチAZ配置
- **セキュリティグループ**: 最小権限アクセス

### コンピューティング
- **ECS Cluster**: Fargateクラスター
- **Task Definition**: CPU 512, Memory 1024
- **Auto Scaling**: 設定可能（現在は固定1タスク）

### データベース
- **Engine**: PostgreSQL 15
- **Instance Class**: db.t3.micro（開発用）
- **Storage**: 20GB
- **Backup**: 自動バックアップ有効

### セキュリティ
- **Secrets Manager**: DB認証情報の暗号化保存
- **IAM Roles**: 最小権限の原則
- **Security Groups**: ポート制限

## 運用管理

### 監視・ログ
```bash
# CloudWatchログの確認
aws logs describe-log-groups --log-group-name-prefix "/ecs/lstep-automation"

# ログストリームの確認
aws logs describe-log-streams --log-group-name "/ecs/lstep-automation-backend"
```

### スケーリング
```bash
# ECSサービスのスケーリング
aws ecs update-service \
  --cluster lstep-automation-cluster \
  --service lstep-automation-backend \
  --desired-count 2
```

### メンテナンス
```bash
# データベースのメンテナンス時間設定
aws rds modify-db-instance \
  --db-instance-identifier lstep-automation-pg \
  --preferred-maintenance-window "sun:04:00-sun:05:00"
```

## CI/CD統合

### GitHub Actions例
```yaml
name: Deploy to AWS
on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v2
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: ap-northeast-1
      
      - name: Build and push to ECR
        run: |
          aws ecr get-login-password | docker login --username AWS --password-stdin $ECR_REGISTRY
          docker build -f backend/Dockerfile -t $ECR_REGISTRY/lstep-automation-backend:$GITHUB_SHA .
          docker push $ECR_REGISTRY/lstep-automation-backend:$GITHUB_SHA
          
      - name: Update ECS service
        run: |
          aws ecs update-service --cluster lstep-automation-cluster --service lstep-automation-backend --force-new-deployment
```

## コスト最適化

### 推定月額費用（東京リージョン）
- **ECS Fargate**: ~$15/月（1タスク常時稼働）
- **RDS t3.micro**: ~$13/月（Single-AZ）
- **ALB**: ~$22/月
- **S3**: ~$1/月（1GB想定）
- **CloudFront**: ~$1/月（1GB転送）
- **その他**: ~$5/月

**合計**: 約$60/月

### コスト削減Tips
```bash
# 開発環境の自動停止設定
aws ecs update-service --cluster lstep-automation-cluster --service lstep-automation-backend --desired-count 0

# 不要なリソースのクリーンアップ
terraform destroy -var-file=env/dev.tfvars
```

## トラブルシューティング

### よくある問題

#### 1. ECRプッシュエラー
```bash
# ECRリポジトリの存在確認
aws ecr describe-repositories --repository-names lstep-automation-backend

# 認証の再実行
aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin $(terraform output -raw ecr_backend_url)
```

#### 2. ECSタスクの起動失敗
```bash
# ECSタスクの状態確認
aws ecs list-tasks --cluster lstep-automation-cluster --service-name lstep-automation-backend

# タスクの詳細確認
aws ecs describe-tasks --cluster lstep-automation-cluster --tasks [TASK_ARN]
```

#### 3. RDS接続エラー
```bash
# セキュリティグループの確認
aws ec2 describe-security-groups --group-ids [SECURITY_GROUP_ID]

# RDSエンドポイントの確認
aws rds describe-db-instances --db-instance-identifier lstep-automation-pg
```

### ログ確認

```bash
# Terraformの詳細ログ
export TF_LOG=DEBUG
terraform apply -var-file=env/dev.tfvars

# AWSリソースの状態確認
aws ecs describe-services --cluster lstep-automation-cluster --services lstep-automation-backend
```

## 環境の削除

```bash
# インフラストラクチャの完全削除
terraform destroy -var-file=env/dev.tfvars

# ECRイメージの削除（オプション）
aws ecr delete-repository --repository-name lstep-automation-backend --force
aws ecr delete-repository --repository-name lstep-automation-frontend-ssr --force
```

**注意**: RDSデータは完全に削除されます。重要なデータは事前にバックアップしてください。

## セキュリティベストプラクティス

### 実装済み
- Secrets Managerでの機密情報管理
- 最小権限IAMロール
- VPCセキュリティグループによるネットワーク制限
- 非rootユーザーでのコンテナ実行

### 推奨追加設定
- **WAF**の設定（DDoS攻撃対策）
- **AWS Config**でのコンプライアンス監視
- **GuardDuty**での脅威検知
- **CloudTrail**での操作ログ記録

これでAWS上での本番環境が構築できます。開発環境での検証が完了したら、このガイドに従ってステップバイステップでデプロイを進めてください。