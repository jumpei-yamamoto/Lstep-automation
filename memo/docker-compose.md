# Docker Compose開発環境ガイド

このドキュメントでは、Docker Composeを使用した開発環境のセットアップと運用方法について説明します。

## 概要

Docker Composeにより、フロントエンド・バックエンド・データベースを統合した開発環境を簡単に構築できます。AWSなどのクラウド環境が未構築でも、ローカルで完全な開発が可能です。

## ファイル構成

- **`docker-compose.dev.yml`**: 開発環境用設定
- **`docker-compose.prod.yml`**: 本番プレビュー用設定

## 開発環境セットアップ

### 前提条件
- Docker Desktop のインストール
- Git でリポジトリをクローン済み

### 起動手順

```bash
# プロジェクトルートに移動
cd Lstep-automation

# 開発環境の起動 (初回はイメージビルドのため時間がかかります)
docker-compose -f docker-compose.dev.yml up -d

# ログの確認
docker-compose -f docker-compose.dev.yml logs -f

# 状態確認
docker-compose -f docker-compose.dev.yml ps
```

### アクセス情報

| サービス | URL | 説明 |
|----------|-----|------|
| フロントエンド | http://localhost:3000 | Next.js SSRアプリケーション |
| バックエンドAPI | http://localhost:8080 | Go (Echo) REST API |
| PostgreSQL | localhost:5432 | データベース (外部クライアント接続可) |

### データベース接続情報

```
Host: localhost
Port: 5432
Database: app
Username: app
Password: app
```

## サービス詳細

### フロントエンド (`frontend-ssr`)
- **基盤**: Next.js SSR
- **ポート**: 3000
- **環境変数**:
  - `NEXT_PUBLIC_API_BASE_URL`: バックエンドAPIのURL
- **依存関係**: backend サービス

### バックエンド (`backend`)
- **基盤**: Go (Echo framework)
- **ポート**: 8080
- **環境変数**:
  - `DB_DSN`: PostgreSQL接続文字列
  - `CORS_ALLOW_ORIGINS`: CORS許可オリジン
- **依存関係**: db サービス

### データベース (`db`)
- **基盤**: PostgreSQL 15 Alpine
- **ポート**: 5432
- **データ永続化**: `db-data` Dockerボリューム

## 開発ワークフロー

### 1. 日常的な開発

```bash
# 環境起動
docker-compose -f docker-compose.dev.yml up -d

# コード変更
# フロントエンド: frontend/ ディレクトリ
# バックエンド: backend/ ディレクトリ

# 変更反映 (ホットリロード対応)
# フロントエンド: 自動反映
# バックエンド: コンテナ再ビルドが必要

# バックエンド変更時の再ビルド
docker-compose -f docker-compose.dev.yml up --build backend

# 環境停止
docker-compose -f docker-compose.dev.yml down
```

### 2. データベース操作

```bash
# PostgreSQLクライアントでの接続
psql -h localhost -p 5432 -U app -d app

# pgAdminなどのGUIツールも使用可能

# マイグレーション実行 (例)
docker-compose -f docker-compose.dev.yml exec backend /app/server migrate

# データリセット
docker-compose -f docker-compose.dev.yml down -v
docker-compose -f docker-compose.dev.yml up -d
```

### 3. ログとデバッグ

```bash
# 全サービスのログ
docker-compose -f docker-compose.dev.yml logs -f

# 特定サービスのログ
docker-compose -f docker-compose.dev.yml logs -f backend

# コンテナ内でのコマンド実行
docker-compose -f docker-compose.dev.yml exec backend sh
docker-compose -f docker-compose.dev.yml exec frontend-ssr sh
```

## 本番プレビュー環境

本番環境に近い状態での検証用:

```bash
# 本番プレビュー環境の起動
docker-compose -f docker-compose.prod.yml up -d

# アクセス: http://localhost:3000 (nginx経由)
```

### プレビュー環境の特徴
- 本番用Dockerfileを使用
- 最適化されたビルド
- 本番に近いパフォーマンス

## 環境変数のカスタマイズ

### `.env` ファイルの作成

```bash
# プロジェクトルートに .env ファイルを作成
cat > .env << EOF
# Database
POSTGRES_USER=app
POSTGRES_PASSWORD=app
POSTGRES_DB=app

# Backend
CORS_ALLOW_ORIGINS=http://localhost:3000,http://localhost:3001

# Frontend
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
EOF
```

### docker-compose.dev.yml での参照

```yaml
services:
  backend:
    environment:
      CORS_ALLOW_ORIGINS: ${CORS_ALLOW_ORIGINS:-http://localhost:3000}
```

## トラブルシューティング

### よくある問題と解決方法

#### 1. ポートが既に使用されている
```bash
# 使用中のポートを確認
lsof -i :3000
lsof -i :8080
lsof -i :5432

# プロセスを停止してから再起動
docker-compose -f docker-compose.dev.yml down
docker-compose -f docker-compose.dev.yml up -d
```

#### 2. データベース接続エラー
```bash
# データベースコンテナの状態確認
docker-compose -f docker-compose.dev.yml logs db

# データベース再起動
docker-compose -f docker-compose.dev.yml restart db
```

#### 3. ビルドエラー
```bash
# キャッシュをクリアして再ビルド
docker-compose -f docker-compose.dev.yml build --no-cache

# 不要なイメージ・コンテナの削除
docker system prune -f
```

#### 4. CORS エラー
```bash
# CORS設定の確認
docker-compose -f docker-compose.dev.yml logs backend | grep CORS

# フロントエンドのAPIベースURL確認
docker-compose -f docker-compose.dev.yml exec frontend-ssr env | grep API
```

## パフォーマンス最適化

### 開発環境での推奨設定

```bash
# Docker Desktopのリソース設定
# Memory: 4GB以上
# CPU: 2コア以上
# Disk: 20GB以上の空き容量
```

### ビルド時間短縮

```bash
# 並列ビルド
docker-compose -f docker-compose.dev.yml build --parallel

# 特定サービスのみビルド
docker-compose -f docker-compose.dev.yml build backend
```

## 他の開発者との連携

### 設定の共有
1. `.env.example` ファイルで環境変数のテンプレートを共有
2. `docker-compose.override.yml` で個人用設定を管理
3. データベースの初期データは `migrations/` で管理

### チーム開発のベストプラクティス
- コンテナイメージのバージョンを固定
- 環境固有の設定は `.env` で管理
- データベーススキーマはマイグレーションで管理

## 次のステップ

開発環境が正常に動作したら:
1. [terraform.md](terraform.md) で本番環境の構築準備
2. CLAUDE.mdの開発ガイドラインに従った実装
3. テスト環境でのCI/CD設定

これで完全なローカル開発環境が構築できました。AWS環境の準備ができたら、Terraformでの本番デプロイに進んでください。