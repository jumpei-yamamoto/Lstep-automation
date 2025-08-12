# internal/platform

プラットフォーム層 - 共有基盤

## 責務

- **データベース接続管理**: 接続プール、設定、ヘルスチェック
- **設定管理**: 環境変数、設定ファイルの読み込みと検証
- **ログ管理**: 構造化ログ、ログレベル制御
- **監視・メトリクス**: アプリケーションメトリクスの収集
- **共通ユーティリティ**: アプリケーション横断的な便利機能

## ファイル構成例

```
platform/
├── db.go              # データベース接続とプール管理
├── config.go          # 設定値の読み込みと検証
├── logger.go          # 構造化ログの設定
├── metrics.go         # メトリクス収集
└── util/
    ├── validator.go   # バリデーションユーティリティ
    └── crypto.go      # 暗号化ユーティリティ
```

## 実装方針

### データベース接続
```go
type DB struct {
    *sql.DB
    config DBConfig
}

func NewDB(config DBConfig) (*DB, error) {
    db, err := sql.Open("pgx", config.DSN)
    if err != nil {
        return nil, err
    }
    
    // コネクションプール設定
    db.SetMaxOpenConns(config.MaxOpenConns)
    db.SetMaxIdleConns(config.MaxIdleConns)
    db.SetConnMaxLifetime(config.ConnMaxLifetime)
    
    return &DB{DB: db, config: config}, nil
}
```

### 設定管理
- 環境変数からの設定読み込み
- デフォルト値の定義
- 設定値の検証とエラーハンドリング
- 開発/ステージング/本番環境の切り替え

### ログ管理
- 構造化ログ（JSON形式）
- コンテキスト情報の自動付与
- CloudWatch等の外部ログサービス連携
- 個人情報のマスキング

### セキュリティ
- パスワードハッシュ化
- JWTトークンの生成・検証
- CORS設定
- レート制限

## 依存関係

- **依存する**: 外部ライブラリ（データベースドライバー、ログライブラリ等）
- **依存される**: 全ての層から利用される共通基盤