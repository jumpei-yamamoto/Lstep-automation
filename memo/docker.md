# Docker設定ガイド

このドキュメントでは、各Dockerfileの詳細と使用方法について説明します。

## バックエンド (Go) - `backend/Dockerfile`

### 特徴
- **マルチステージビルド**による最適化
- **静的リンクバイナリ**の生成 (CGO無効)
- **Alpine Linux**による軽量化
- **非rootユーザー**での実行

### 構成

#### ビルドステージ
```dockerfile
FROM golang:1.22-alpine AS build
```
- Go 1.22の公式イメージを使用
- 必要なパッケージ (git, ca-certificates, tzdata) をインストール
- Go modulesのダウンロードと依存関係の解決
- 静的リンクバイナリの生成

#### ランタイムステージ
```dockerfile
FROM alpine:3.20
```
- 軽量なAlpine Linuxベース
- セキュリティ強化のため非rootユーザーで実行
- ポート8080で待機

### ビルドと実行

```bash
# 単体でのビルド
docker build -f backend/Dockerfile -t lstep-backend .

# 実行
docker run -p 8080:8080 -e DB_DSN="postgres://..." lstep-backend
```

## フロントエンド SSR - `frontend/Dockerfile.ssr`

### 特徴
- **Next.js standalone出力**によるサーバーサイドレンダリング
- **3段階ビルドプロセス**による最適化
- **Node.js 20 Alpine**ベース

### 構成

#### 依存関係ステージ
```dockerfile
FROM node:20-alpine AS deps
```
- package.jsonベースの依存関係インストール
- npm ciによる決定論的インストール

#### ビルドステージ
```dockerfile
FROM node:20-alpine AS builder
```
- Next.jsアプリケーションのビルド
- standalone出力の生成

#### ランタイムステージ
```dockerfile
FROM node:20-alpine AS runner
```
- 本番環境用の最小構成
- 非rootユーザーでの実行
- ポート3000で待機

### ビルドと実行

```bash
# SSR用ビルド
docker build -f frontend/Dockerfile.ssr -t lstep-frontend-ssr .

# 実行
docker run -p 3000:3000 -e NEXT_PUBLIC_API_BASE_URL="http://localhost:8080" lstep-frontend-ssr
```

## フロントエンド 静的サイト - `frontend/Dockerfile.static`

### 特徴
- **静的HTML出力**による高速配信
- **Nginx**による軽量サーブ
- **CDN配信**に最適化

### 構成

#### ビルドステージ
```dockerfile
FROM node:20-alpine AS builder
```
- Next.jsの静的出力生成 (`npm run export`)
- out/ディレクトリに静的ファイルを生成

#### サーブステージ
```dockerfile
FROM nginx:alpine
```
- Nginxによる静的ファイル配信
- ポート80で待機

### ビルドと実行

```bash
# 静的サイト用ビルド
docker build -f frontend/Dockerfile.static -t lstep-frontend-static .

# 実行
docker run -p 3000:80 lstep-frontend-static
```

## 使い分けガイド

### 開発環境
- **SSR Dockerfile**: `frontend/Dockerfile.ssr`
- リアルタイムなAPI通信が必要な場合
- サーバーサイドレンダリングが必要な場合

### 本番環境

#### 高度なSEO・SSRが必要な場合
- **SSR Dockerfile**: `frontend/Dockerfile.ssr`
- ECS Fargate等のコンテナ環境で運用

#### 静的サイトで十分な場合
- **Static Dockerfile**: `frontend/Dockerfile.static`
- S3 + CloudFrontでの配信推奨

## 環境変数

### バックエンド
| 変数名 | 説明 | 例 |
|--------|------|-----|
| `PORT` | サーバーポート | `8080` |
| `DB_DSN` | PostgreSQL接続文字列 | `postgres://user:pass@host:5432/db` |
| `CORS_ALLOW_ORIGINS` | CORS許可オリジン | `http://localhost:3000` |

### フロントエンド (SSR)
| 変数名 | 説明 | 例 |
|--------|------|-----|
| `NODE_ENV` | Node.js環境 | `production` |
| `NEXT_PUBLIC_API_BASE_URL` | APIベースURL | `http://localhost:8080` |
| `NEXT_TELEMETRY_DISABLED` | Next.jsテレメトリ無効化 | `1` |

## セキュリティベストプラクティス

### 実装済み
- **非rootユーザー**での実行
- **最小権限**でのプロセス実行
- **静的リンクバイナリ**によるライブラリ脆弱性の回避
- **軽量ベースイメージ**による攻撃面の縮小

### 推奨事項
- 定期的なベースイメージの更新
- セキュリティスキャンの実行 (ECRで自動実行)
- 機密情報の環境変数での管理

## トラブルシューティング

### よくある問題

#### ビルドエラー
```bash
# キャッシュクリア
docker builder prune

# ネットワーク確認
docker network ls
```

#### 接続エラー
```bash
# コンテナ間通信確認
docker network inspect <network_name>

# ログ確認
docker logs <container_name>
```

#### パフォーマンス問題
```bash
# リソース使用量確認
docker stats

# イメージサイズ確認
docker images
```

次は [docker-compose.md](docker-compose.md) で統合開発環境の設定方法を確認してください。