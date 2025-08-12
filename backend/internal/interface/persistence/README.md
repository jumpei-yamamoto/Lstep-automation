# internal/interface/persistence

永続化インターフェース層

## 責務

- **リポジトリインターフェースの実装**: ドメイン層で定義されたリポジトリの具体実装
- **データベースアクセス**: SQL実行、接続管理、トランザクション処理
- **データマッピング**: ドメインオブジェクト ↔ データベースレコードの相互変換
- **クエリ最適化**: 効率的なSQL文の作成とインデックス活用

## ファイル構成例

```
persistence/
├── user_repository_pg.go    # ユーザーリポジトリのPostgreSQL実装
├── auth_repository_pg.go    # 認証リポジトリのPostgreSQL実装
├── mapper/
│   ├── user_mapper.go       # ユーザーデータマッピング
│   └── auth_mapper.go       # 認証データマッピング
└── migration/
    └── schema.sql           # テーブル定義（参考用）
```

## 実装方針

### リポジトリパターン
```go
type UserRepositoryPG struct {
    DB *sql.DB
}

func (r *UserRepositoryPG) Save(ctx context.Context, user *domain.User) error {
    const query = `
        INSERT INTO users(id, name, email, created_at) 
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (id) DO UPDATE SET 
        name = $2, email = $3
    `
    _, err := r.DB.ExecContext(ctx, query, 
        user.ID, user.Name, user.Email, user.CreatedAt)
    return err
}
```

### データマッピング
- ドメインオブジェクト → データベーステーブル行
- NULL値、デフォルト値の適切な処理
- 型変換（UUID、時刻フォーマット等）の統一

### エラーハンドリング
- SQLエラーの適切なラッピング
- 制約違反エラーのドメインエラーへの変換
- 接続エラー、タイムアウトの処理

### パフォーマンス
- クエリの最適化とインデックス設計
- バッチ処理の活用
- コネクションプールの効率的な利用

## 依存関係

- **依存する**: Domain層インターフェース、データベースドライバー
- **依存される**: cmd/server (DI時)