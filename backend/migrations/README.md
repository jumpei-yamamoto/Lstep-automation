# migrations

データベースマイグレーション

## 責務

- **データベーススキーマの管理**: テーブル、インデックス、制約の作成・変更・削除
- **データマイグレーション**: データの移行・変換処理
- **バージョン管理**: スキーマ変更の履歴と適用状況の管理
- **ロールバック機能**: 問題が発生した場合の巻き戻し処理

## 推奨ツール

- **golang-migrate**: Go言語標準的なマイグレーションツール
- **goose**: シンプルで使いやすいマイグレーションツール

## ファイル構成例

```
migrations/
├── 20250101000001_create_users_table.up.sql      # ユーザーテーブル作成
├── 20250101000001_create_users_table.down.sql    # ユーザーテーブル削除（ロールバック用）
├── 20250101000002_add_user_indexes.up.sql        # インデックス追加
├── 20250101000002_add_user_indexes.down.sql      # インデックス削除（ロールバック用）
└── 20250101000003_add_auth_tables.up.sql         # 認証テーブル追加
```

## 命名規約

### ファイル名
```
{timestamp}_{description}.{direction}.sql
```

- `timestamp`: YYYYMMDDHHMMSS形式
- `description`: 変更内容を表す英語の説明（スネークケース）
- `direction`: `up`（適用）または `down`（ロールバック）

### テーブル設計指針
- **主キー**: UUIDを推奨（`id UUID PRIMARY KEY DEFAULT gen_random_uuid()`）
- **タイムスタンプ**: `created_at`、`updated_at`を基本として追加
- **外部キー制約**: データ整合性を保証する制約を適切に設定
- **インデックス**: クエリパフォーマンスを考慮した設計

## サンプルマイグレーション

```sql
-- 20250101000001_create_users_table.up.sql
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_created_at ON users(created_at);
```

```sql
-- 20250101000001_create_users_table.down.sql
DROP TABLE IF EXISTS users CASCADE;
```

## 実行方法

```bash
# マイグレーション適用
migrate -path ./migrations -database "postgres://user:pass@localhost/db?sslmode=disable" up

# ロールバック
migrate -path ./migrations -database "postgres://user:pass@localhost/db?sslmode=disable" down 1
```