# インフラストラクチャ概要

このドキュメントでは、Lstep-automationプロジェクトのインフラストラクチャ構成について説明します。

## アーキテクチャ概要

### 技術スタック
- **フロントエンド**: Next.js (React) + TypeScript + Tailwind CSS
- **バックエンド**: Go (Echo framework) + PostgreSQL
- **コンテナ化**: Docker + Docker Compose
- **本番環境**: AWS (ECS Fargate + RDS + S3 + CloudFront)

### 設計原則
- **DDD (Domain-Driven Design)**
- **Clean Architecture**
- **マイクロサービス指向**
- **コンテナファースト**

## デプロイメント戦略

### 開発環境
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │    Backend      │    │   PostgreSQL    │
│   (Next.js)     │◄──►│   (Go/Echo)     │◄──►│   (Container)   │
│   Port: 3000    │    │   Port: 8080    │    │   Port: 5432    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

- Docker Composeによる統合開発環境
- ホットリロード対応
- ローカルDB (PostgreSQL) 

### 本番環境 (AWS)
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   CloudFront    │    │       ALB       │    │   ECS Fargate   │
│   (CDN)         │    │  (Load Balancer)│    │   (Backend)     │
│                 │    │                 │    │                 │
│   S3 Bucket     │    │                 │    │                 │
│   (Frontend)    │    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────┬───────┘
                                                        │
                                              ┌─────────▼───────┐
                                              │       RDS       │
                                              │   (PostgreSQL)  │
                                              │                 │
                                              └─────────────────┘
```

## ファイル構成

```
project/
├── backend/
│   └── Dockerfile                    # Goアプリケーション用
├── frontend/
│   ├── Dockerfile.ssr               # SSR用Next.js
│   └── Dockerfile.static            # 静的サイト用
├── docker-compose.dev.yml           # 開発環境
├── docker-compose.prod.yml          # 本番プレビュー環境
└── infra/terraform/                 # AWS本番インフラ
    ├── main.tf                      # メインリソース定義
    ├── variables.tf                 # 変数定義
    ├── outputs.tf                   # 出力定義
    ├── providers.tf                 # プロバイダー設定
    └── env/dev.tfvars              # 環境別設定
```

## セキュリティ考慮事項

### 開発環境
- データベース認証情報はコンテナ内にのみ存在
- CORS設定により許可されたオリジンからのみアクセス可能

### 本番環境
- **AWS Secrets Manager**でDB認証情報を管理
- **Security Groups**によるネットワーク制限
- **IAM**による最小権限アクセス
- **HTTP-only Cookie**による認証トークン管理
- **WAF**による攻撃防御（Terraform設定に含む予定）

## 環境別構成

| 項目 | 開発環境 | 本番環境 |
|------|----------|----------|
| フロントエンド | Docker (Next.js SSR) | S3 + CloudFront |
| バックエンド | Docker (Go) | ECS Fargate |
| データベース | PostgreSQL Container | RDS PostgreSQL |
| 認証 | 簡易設定 | Secrets Manager |
| HTTPS | なし | CloudFront SSL |
| 監視 | ローカルログ | CloudWatch |

## 次のステップ

1. **開発環境セットアップ**: [docker-compose.md](docker-compose.md)
2. **本番デプロイ準備**: [terraform.md](terraform.md)
3. **個別コンテナ設定**: [docker.md](docker.md)

各環境の詳細な設定方法については、該当するドキュメントを参照してください。