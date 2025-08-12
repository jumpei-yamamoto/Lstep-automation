# internal/interface/http

HTTP インターフェース層

## 責務

- **HTTPリクエストの受信と解析**: JSON、クエリパラメータ等の解析
- **リクエスト検証**: 基本的な形式チェック（詳細な業務検証はドメイン層）
- **ユースケースの呼び出し**: HTTPリクエストをユースケース実行に変換
- **HTTPレスポンスの生成**: ユースケース結果をJSON、HTTPステータスコードに変換
- **エラーハンドリング**: ドメインエラーを適切なHTTPエラーレスポンスに変換

## ファイル構成例

```
http/
├── router.go           # ルーティング設定
├── user_handler.go     # ユーザー関連のHTTPハンドラー
├── auth_handler.go     # 認証関連のHTTPハンドラー
├── middleware.go       # カスタムミドルウェア
└── dto/
    ├── user_dto.go     # ユーザー関連のリクエスト/レスポンスDTO
    └── error_dto.go    # エラーレスポンス用DTO
```

## 実装方針

### リクエストハンドリング
```go
func (h *UserHandler) CreateUser(c echo.Context) error {
    // 1. リクエスト解析
    var req CreateUserRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(400, ErrorResponse{Error: "invalid request"})
    }
    
    // 2. ユースケース実行
    user, err := h.useCase.Execute(c.Request().Context(), req.ToUseCaseInput())
    
    // 3. エラーハンドリング
    if err != nil {
        return h.handleError(c, err)
    }
    
    // 4. レスポンス生成
    return c.JSON(201, UserResponse{}.FromDomain(user))
}
```

### エラーマッピング
- ドメインエラー → 4xx系レスポンス
- インフラエラー → 5xx系レスポンス
- 詳細なエラー情報の適切なマスキング

## 依存関係

- **依存する**: UseCase層、Domain層エラー、Echo フレームワーク
- **依存される**: cmd/server (エントリポイント)