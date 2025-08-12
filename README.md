# Lstep-automation

## 概要

このプロジェクトは**Lステップ構築プロジェクト**です。Lステップとは、顧客との関係構築を自動化し、段階的にアプローチすることでコンバージョン率を向上させるマーケティング手法の実装システムです。

## 機能

- **ステップメール配信**: 顧客の行動に応じた段階的なメール配信
- **顧客セグメンテーション**: 顧客属性や行動データによる自動セグメント分類
- **自動フォローアップ**: 設定したルールに基づく自動的な顧客フォローアップ
- **分析・レポート**: Lステップの効果測定とパフォーマンス分析

## 必要な環境

- Go 1.22+
- Node.js 18+
- PostgreSQL 14+
- Docker & Docker Compose（開発環境用）

## インストール

```bash
# リポジトリをクローン
git clone https://github.com/jumpei-yamamoto/Lstep-automation.git
cd Lstep-automation

# 開発環境をセットアップ（Docker使用）
docker-compose up -d

# フロントエンドの依存関係をインストール
cd frontend
npm install

# バックエンドの依存関係をインストール
cd ../backend
go mod download
```

## 使い方

### 基本的な使用方法

```bash
# バックエンドサーバーを起動
cd backend
go run cmd/server/main.go

# フロントエンドアプリケーションを起動（別ターミナル）
cd frontend
npm run dev
```

### 設定

1. `.env`ファイルを作成し、必要な環境変数を設定
2. データベース接続情報を設定
3. メール送信サービス（SMTP）の設定

## 設定項目

| 項目名 | 説明 | デフォルト値 |
|--------|------|--------------|
| DB_DSN | PostgreSQL接続文字列 | postgres://user:pass@localhost/lstep |
| SMTP_HOST | SMTPサーバーホスト | smtp.gmail.com |
| SMTP_PORT | SMTPサーバーポート | 587 |
| JWT_SECRET | JWT認証秘密鍵 | - |

## ディレクトリ構成

```
Lstep-automation/
├── README.md
├── CLAUDE.md                    # Claude.AI実装ガイドライン
├── .gitignore                   # Go & Node.js gitignore設定
├── backend/                     # Go + Echo バックエンド
│   ├── cmd/
│   │   └── server/main.go       # エントリポイント（DI/起動/ルーティング）
│   ├── internal/
│   │   ├── domain/              # 100% 純粋なドメイン（技術依存なし）
│   │   ├── usecase/             # アプリケーションサービス（ユースケース）
│   │   ├── interface/           # 入出力変換
│   │   │   ├── http/            # echo ハンドラ/ルーティング/DTO
│   │   │   └── persistence/     # DBアダプタ（リポジトリ実装）
│   │   └── platform/            # 共有基盤（DB接続、設定、ログ）
│   └── migrations/              # DBマイグレーション（goose等）
├── frontend/                    # Next.js + Tailwind フロントエンド
│   ├── app/                     # App Router 推奨（Pages Routerでも可）
│   ├── components/              # UIコンポーネント
│   ├── lib/                     # fetch ラッパ（/api 経由）、config等
│   ├── services/                # API 呼び出し集約（UIから直 fetch 禁止）
│   └── styles/                  # グローバルCSS
├── infrastructure/              # インフラ構成（AWS、Terraform等）
└── docs/                        # プロジェクト文書
```

## 開発

### 開発環境のセットアップ

```bash
# 開発用コンテナを起動
docker-compose -f docker-compose.dev.yml up -d

# データベースマイグレーション
cd backend
make migrate-up
```

### テスト

```bash
# バックエンドテスト
cd backend
go test ./...

# フロントエンドテスト
cd frontend
npm test
```

### Lステップの基本概念

- **リード（見込み客）**: メールアドレス等を登録した潜在顧客
- **ステップメール**: 段階的に配信される一連のメール
- **シナリオ**: 顧客の行動に応じた自動化フロー
- **セグメント**: 顧客属性や行動による分類

## 貢献

1. このリポジトリをフォーク
2. 機能ブランチを作成 (`git checkout -b feature/AmazingFeature`)
3. 変更をコミット (`git commit -m 'Add some AmazingFeature'`)
4. ブランチにプッシュ (`git push origin feature/AmazingFeature`)
5. プルリクエストを作成

## ライセンス

[ライセンス情報]

## 作者

- **jumpei-yamamoto** - *Initial work* - [jumpei-yamamoto](https://github.com/jumpei-yamamoto)

## 謝辞

- Lステップマーケティング手法の研究と実装に関する知見
- DDD + Clean Architecture 設計パターンの採用

## 変更履歴

### v1.0.0 (2025-08-10)
- Lステップ構築プロジェクト初期リリース
- ユーザー管理機能
- ステップメール配信基盤