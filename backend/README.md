# Backend

Go + Echo サーバーアプリケーション

## アーキテクチャ

Domain-Driven Design（DDD）+ Clean Architecture に基づく構成

## ディレクトリ構成

```
backend/
├── cmd/server/          # エントリポイント（DI/起動/ルーティング）
├── internal/
│   ├── domain/          # 100% 純粋なドメイン（技術依存なし）
│   ├── usecase/         # アプリケーションサービス（ユースケース）
│   ├── interface/       # 入出力変換
│   │   ├── http/        # echo ハンドラ/ルーティング/DTO
│   │   └── persistence/ # DBアダプタ（リポジトリ実装）
│   └── platform/        # 共有基盤（DB接続、設定、ログ）
└── migrations/          # DBマイグレーション
```

## 技術スタック

- **言語**: Go 1.22+
- **Webフレームワーク**: Echo v4
- **データベース**: PostgreSQL（RDS想定）
- **マイグレーション**: goose または golang-migrate

## 依存関係の原則

**依存の向きは常に内向き**

- `Domain` ← `UseCase` ← `Interface` ← `Infrastructure`
- フレームワークは外側、ドメインは内側
- ドメイン層は外部技術に依存しない（echo/SQL/HTTP等の import 禁止）

## 開発ガイド

詳細な実装方針については、プロジェクトルートの `CLAUDE.md` を参照してください。