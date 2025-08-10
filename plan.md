# LSTEP自動化システム構築実行計画

## 概要

本計画は、CLAUDE.mdで定義されたDDD + Clean Architecture原則に基づき、Next.js + Go + AWSでLSTEP（エルステップ）自動化システムを構築するための実行計画です。

## 技術スタック

- **フロントエンド**: Next.js (App Router) + TypeScript + Tailwind CSS + Recoil
- **バックエンド**: Go + Echo v4 + PostgreSQL
- **インフラ**: AWS (ECS/Fargate + ALB + RDS + S3 + CloudFront)
- **認証**: HTTP-only Cookie
- **アーキテクチャ**: DDD + Clean Architecture

---

## 実行計画（35ステップ）

### Phase 1: プロジェクト初期設定 (1-5)

#### 1. リポジトリ構造の整備
- ルートディレクトリに `backend/`, `frontend/`, `infrastructure/` を作成
- `.gitignore` の設定（Go/Node.js/IDE用）
- `docs/` ディレクトリの作成（設計書/API仕様書用）

#### 2. バックエンド基盤構築
- `backend/` 配下にクリーンアーキテクチャ準拠のディレクトリ構成を作成
  - `cmd/server/`
  - `internal/domain/`
  - `internal/usecase/`
  - `internal/interface/http/`
  - `internal/interface/persistence/`
  - `internal/platform/`
  - `migrations/`

#### 3. Goモジュール初期化と依存関係設定
- `go mod init` によるモジュール初期化
- 必要なライブラリのインストール（echo, pgx, uuid, zap等）
- `Makefile` 作成（ビルド/テスト/マイグレーション用コマンド）

#### 4. フロントエンド基盤構築
- Next.js プロジェクトの初期化
- TypeScript設定
- Tailwind CSS導入
- ディレクトリ構成作成（app/, components/, lib/, services/）

#### 5. 開発環境のDocker化
- Docker Compose設定（PostgreSQL, Redis, Go dev server）
- 開発用環境変数設定
- ローカル開発環境の構築・テスト

### Phase 2: ドメイン設計・実装 (6-15)

#### 6. ドメイン分析とエンティティ設計
- LSTEP自動化システムの業務要件分析
- 集約（Aggregate）の特定と境界設定
- エンティティとバリューオブジェクトの設計

#### 7. ユーザー管理ドメインの実装
- User エンティティ実装
- ユーザー関連のドメインエラー定義
- UserRepository インターフェース定義

#### 8. 自動化ワークフロー ドメインの実装
- Workflow エンティティ実装
- Step バリューオブジェクト実装
- WorkflowRepository インターフェース定義

#### 9. プロジェクト管理ドメインの実装
- Project エンティティ実装
- ProjectStatus バリューオブジェクト実装
- ProjectRepository インターフェース定義

#### 10. 実行履歴ドメインの実装
- ExecutionHistory エンティティ実装
- ExecutionStatus バリューオブジェクト実装
- ExecutionHistoryRepository インターフェース定義

#### 11. 通知・アラートドメインの実装
- Notification エンティティ実装
- NotificationChannel バリューオブジェクト実装
- NotificationRepository インターフェース定義

#### 12. ドメインサービス実装
- ワークフロー実行順序決定サービス
- プロジェクト重複チェックサービス
- 通知配信ルール決定サービス

#### 13. ドメインイベント設計・実装
- WorkflowExecuted イベント実装
- ProjectCompleted イベント実装
- NotificationSent イベント実装

#### 14. ドメイン層の単体テスト作成
- 各エンティティのテストケース実装
- ドメインサービスのテストケース実装
- 不変条件（Invariant）のテストケース実装

#### 15. ドメイン層のドキュメント作成
- ユビキタス言語の定義書作成
- ドメインモデル図の作成
- 集約境界図の作成

### Phase 3: アプリケーション層実装 (16-20)

#### 16. ユーザー管理ユースケース実装
- RegisterUser ユースケース実装
- AuthenticateUser ユースケース実装
- UpdateUserProfile ユースケース実装

#### 17. ワークフロー管理ユースケース実装
- CreateWorkflow ユースケース実装
- ExecuteWorkflow ユースケース実装
- UpdateWorkflowStatus ユースケース実装

#### 18. プロジェクト管理ユースケース実装
- CreateProject ユースケース実装
- AssignWorkflowToProject ユースケース実装
- CompleteProject ユースケース実装

#### 19. 実行管理ユースケース実装
- StartExecution ユースケース実装
- MonitorExecution ユースケース実装
- HandleExecutionError ユースケース実装

#### 20. アプリケーション層のテスト作成
- 各ユースケースの単体テスト実装
- モック/フェイクリポジトリを使用したテスト
- エラーハンドリングのテストケース実装

### Phase 4: インフラストラクチャ層実装 (21-25)

#### 21. データベース設計・マイグレーション作成
- テーブル設計（正規化、インデックス設計）
- マイグレーションスクリプト作成
- 初期データ投入スクリプト作成

#### 22. リポジトリ実装（PostgreSQL）
- UserRepository 実装
- WorkflowRepository 実装
- ProjectRepository 実装
- ExecutionHistoryRepository 実装

#### 23. 外部API統合実装
- 通知サービス（メール/Slack）のアダプタ実装
- ファイルストレージ（S3）のアダプタ実装
- ログ出力（CloudWatch）のアダプタ実装

#### 24. 設定管理・環境変数実装
- 環境別設定ファイル実装
- AWS Secrets Manager連携実装
- 設定値バリデーション実装

#### 25. インフラ層の統合テスト作成
- データベース接続テスト
- リポジトリの統合テスト
- 外部API連携テスト

### Phase 6: Webインターフェース実装 (26-30)

#### 26. HTTPハンドラー・ルーター実装
- Echo サーバー設定
- ルーティング実装
- ミドルウェア実装（CORS, Logger, Recovery）

#### 27. RESTful API実装
- ユーザー管理API実装
- ワークフロー管理API実装
- プロジェクト管理API実装
- 実行管理API実装

#### 28. 認証・認可システム実装
- JWT トークン発行/検証実装
- HTTP-only Cookie実装
- 権限チェックミドルウェア実装

#### 29. APIドキュメント作成
- OpenAPI 3.0 仕様書作成
- Swagger UI 設定
- APIテスト用Postmanコレクション作成

#### 30. API層の統合テスト作成
- HTTPハンドラーのテスト実装
- 認証・認可のテスト実装
- エラーレスポンスのテスト実装

### Phase 6: フロントエンド実装 (31-35)

#### 31. Next.js基盤実装
- App Router設定
- レイアウトコンポーネント実装
- プロバイダー実装（Recoil, React Query）

#### 32. 共通コンポーネント実装
- UIコンポーネントライブラリ構築
- フォームバリデーション実装
- エラーハンドリングコンポーネント実装

#### 33. 画面実装
- ログイン/ログアウト画面
- ダッシュボード画面
- ワークフロー管理画面
- プロジェクト管理画面
- 実行履歴画面

#### 34. サービス層実装
- API呼び出しサービス実装
- 状態管理（Recoil）実装
- キャッシュ戦略（React Query）実装

#### 35. フロントエンドテスト実装
- コンポーネントのユニットテスト
- 統合テスト（MSW使用）
- E2Eテスト（Playwright）

---

## 品質保証・デプロイメント

### CI/CD パイプライン設定
- GitHub Actions による自動テスト実行
- コードカバレッジ測定
- 静的解析（golangci-lint, ESLint）
- セキュリティスキャン

### AWS インフラストラクチャ構築
- Terraform による IaC 実装
- ECS/Fargate + ALB 設定
- RDS PostgreSQL 設定
- S3 + CloudFront 設定
- WAF + Security Group 設定

### 監視・運用設定
- CloudWatch ログ/メトリクス設定
- アラート設定
- ヘルスチェック実装
- バックアップ戦略実装

---

## 成果物

- 動作するLSTEP自動化システム
- 完全なソースコード（テスト含む）
- API仕様書
- アーキテクチャ設計書
- デプロイメント手順書
- 運用マニュアル

---

## 注意事項

1. **アーキテクチャ境界の遵守**: 各層の責務を明確に分離し、依存方向を内向きに保つ
2. **セキュリティ**: HTTP-only Cookie使用、入力検証、CORS設定の徹底
3. **テスト**: 各層のテストを必須とし、CI/CDで品質を担保
4. **ドキュメント**: コードと合わせて設計書・仕様書の更新を行う
5. **レビュー**: 各フェーズ完了時にコードレビューと設計レビューを実施

この実行計画に従って、品質の高いLSTEP自動化システムを段階的に構築していきます。