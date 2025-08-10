# GitHub Issues for Plan.md Implementation

Based on the analysis of plan.md, here are the detailed GitHub issues that should be created for each step. These issues are broken down into Claude Code implementable tasks.

## Phase 1: プロジェクト初期設定

### Issue 1.1: リポジトリ構造の整備
**Title**: Phase 1.1: リポジトリ構造の整備
**Labels**: Phase 1, setup
**Body**:
```
## 概要
ルートディレクトリに基本的なディレクトリ構造を作成し、プロジェクトの基盤を整備します。

## タスク詳細
- [ ] `backend/` ディレクトリの作成
- [ ] `frontend/` ディレクトリの作成  
- [ ] `infrastructure/` ディレクトリの作成
- [ ] `.gitignore` の設定（Go/Node.js/IDE用の適切な設定）
- [ ] `docs/` ディレクトリの作成
- [ ] `README.md` の基本構成作成

## 受け入れ基準
- [ ] 各ディレクトリが適切に作成されている
- [ ] .gitignoreが Go と Node.js の開発に適した設定になっている
- [ ] docs ディレクトリが作成されている

## 参考
- CLAUDE.md のディレクトリ構成に従うこと
- DDD + Clean Architecture を前提とした構造にすること

## 実装者向け注意
このタスクは比較的単純なファイル/ディレクトリ作成作業のため、Claude Codeで実装可能です。
```

### Issue 1.2: バックエンド基盤構築
**Title**: Phase 1.2: バックエンドクリーンアーキテクチャディレクトリ構成作成
**Labels**: Phase 1, backend, architecture
**Body**:
```
## 概要
backend/ 配下にクリーンアーキテクチャ準拠のディレクトリ構成を作成します。

## タスク詳細
- [ ] `backend/cmd/server/` ディレクトリ作成
- [ ] `backend/internal/domain/` ディレクトリ作成
- [ ] `backend/internal/usecase/` ディレクトリ作成
- [ ] `backend/internal/interface/http/` ディレクトリ作成
- [ ] `backend/internal/interface/persistence/` ディレクトリ作成
- [ ] `backend/internal/platform/` ディレクトリ作成
- [ ] `backend/migrations/` ディレクトリ作成
- [ ] 各ディレクトリに適切な README.md を配置

## 受け入れ基準
- [ ] CLAUDE.md に記載されたディレクトリ構成が作成されている
- [ ] 各層の責務が README.md で説明されている

## 参考
- CLAUDE.md のバックエンド構成サンプルを参照
- DDD + Clean Architecture の層分離原則に従うこと
```

### Issue 1.3: Goモジュール初期化と基本設定
**Title**: Phase 1.3: Goモジュール初期化と依存関係設定
**Labels**: Phase 1, backend, go
**Body**:
```
## 概要
Goモジュールの初期化と基本的な依存関係の設定を行います。

## タスク詳細
- [ ] `go mod init` によるモジュール初期化
- [ ] go.mod での基本依存関係追加（echo, pgx, uuid, zap等）
- [ ] `Makefile` 作成（ビルド/テスト/マイグレーション用コマンド）
- [ ] 基本的な main.go テンプレート作成

## 受け入れ基準
- [ ] go.mod が適切に作成されている
- [ ] 必要なライブラリが定義されている
- [ ] Makefileで基本的なコマンドが実行可能

## 参考
- CLAUDE.md のサンプルコードを参照
- Echo v4, pgx v5, uuid, zap などの推奨ライブラリを使用

## 実装者向け注意
依存関係のバージョンは最新安定版を使用してください。
```

### Issue 1.4: フロントエンド基盤構築
**Title**: Phase 1.4: Next.jsプロジェクト初期化とTypeScript設定
**Labels**: Phase 1, frontend, nextjs
**Body**:
```
## 概要
Next.jsプロジェクトの初期化とTypeScript、Tailwind CSSの導入を行います。

## タスク詳細
- [ ] frontend/ ディレクトリでNext.jsプロジェクト初期化
- [ ] TypeScript設定（tsconfig.json）
- [ ] Tailwind CSS導入と設定
- [ ] ディレクトリ構成作成（app/, components/, lib/, services/）
- [ ] 基本的なpackage.json スクリプト設定

## 受け入れ基準
- [ ] Next.js App Routerが使用可能な状態になっている
- [ ] TypeScriptが適切に設定されている
- [ ] Tailwind CSSが動作している
- [ ] CLAUDE.md で推奨されるディレクトリ構成が作成されている

## 参考
- CLAUDE.md のフロントエンド構成を参照
- App Router + TypeScript + Tailwind CSS の組み合わせ
```

### Issue 1.5: 開発環境のDocker化
**Title**: Phase 1.5: 開発環境Docker Compose設定
**Labels**: Phase 1, docker, devops
**Body**:
```
## 概要
ローカル開発環境用のDocker Compose設定を作成します。

## タスク詳細
- [ ] docker-compose.yml の作成（PostgreSQL, Redis含む）
- [ ] 開発用環境変数設定（.env.example）
- [ ] データベース初期化用スクリプト
- [ ] ローカル開発環境構築手順のドキュメント作成

## 受け入れ基準
- [ ] docker-compose up でローカル開発環境が起動する
- [ ] PostgreSQLにアクセス可能
- [ ] 環境変数が適切に設定されている

## 参考
- PostgreSQL 15+ を使用
- 開発用のポート設定（5432, 6379等）
- ヘルスチェック機能の実装
```

## Phase 2: ドメイン設計・実装

### Issue 2.1: ドメイン分析とエンティティ設計
**Title**: Phase 2.1: LSTEP自動化システムドメイン分析
**Labels**: Phase 2, domain, analysis
**Body**:
```
## 概要
LSTEP自動化システムの業務要件を分析し、ドメインモデルを設計します。

## タスク詳細
- [ ] 業務要件の整理とドメインの特定
- [ ] 集約（Aggregate）の境界設定
- [ ] エンティティとバリューオブジェクトの設計
- [ ] ユビキタス言語の定義
- [ ] ドメインモデル図の作成

## 受け入れ基準
- [ ] 主要な集約が特定されている
- [ ] エンティティとVOの責務が明確に定義されている
- [ ] ユビキタス言語がドキュメント化されている

## 参考
- DDD戦術設計パターンに従う
- CLAUDE.mdのドメイン実装サンプルを参考

## 実装者向け注意
実装前の設計段階のため、実装よりもドキュメント作成が中心となります。
```

### Issue 2.2: ユーザー管理ドメイン実装
**Title**: Phase 2.2: User エンティティとリポジトリインターフェース実装
**Labels**: Phase 2, domain, user
**Body**:
```
## 概要
ユーザー管理ドメインの中核となるUserエンティティとリポジトリインターフェースを実装します。

## タスク詳細
- [ ] User エンティティの実装（domain/user.go）
- [ ] ユーザー関連ドメインエラーの定義
- [ ] UserRepository インターフェースの定義
- [ ] ユーザー作成時のバリデーションルール実装
- [ ] ドメイン層のユニットテスト作成

## 受け入れ基準
- [ ] User エンティティが不変条件を維持している
- [ ] ビジネスルールが適切に実装されている
- [ ] 外部依存のないピュアなドメインオブジェクト
- [ ] 全てのテストが通過している

## 参考
- CLAUDE.md の User エンティティサンプル
- ドメイン駆動設計の原則に従った実装

## 実装ファイル
- `backend/internal/domain/user.go`
- `backend/internal/domain/user_repository.go`
- `backend/internal/domain/user_test.go`
```

### Issue 2.3: 自動化ワークフロードメイン実装
**Title**: Phase 2.3: Workflow エンティティとStep バリューオブジェクト実装
**Labels**: Phase 2, domain, workflow
**Body**:
```
## 概要
自動化ワークフローの中核となるWorkflowエンティティとStepバリューオブジェクトを実装します。

## タスク詳細
- [ ] Workflow エンティティの実装
- [ ] Step バリューオブジェクトの実装
- [ ] ワークフロー実行順序に関するビジネスルール
- [ ] WorkflowRepository インターフェース定義
- [ ] ワークフロー関連ドメインエラー定義
- [ ] ドメイン層ユニットテスト作成

## 受け入れ基準
- [ ] ワークフローの状態管理が適切に実装されている
- [ ] ステップ間の依存関係が管理されている
- [ ] 不変条件が守られている
- [ ] テストカバレッジが十分である

## 実装ファイル
- `backend/internal/domain/workflow.go`
- `backend/internal/domain/step.go`
- `backend/internal/domain/workflow_repository.go`
- `backend/internal/domain/workflow_test.go`
```

### Issue 2.4: プロジェクト管理ドメイン実装
**Title**: Phase 2.4: Project エンティティとProjectStatus 実装
**Labels**: Phase 2, domain, project
**Body**:
```
## 概要
プロジェクト管理ドメインのProjectエンティティとProjectStatusバリューオブジェクトを実装します。

## タスク詳細
- [ ] Project エンティティの実装
- [ ] ProjectStatus バリューオブジェクトの実装
- [ ] プロジェクト完了条件のビジネスルール
- [ ] ProjectRepository インターフェース定義
- [ ] プロジェクト関連ドメインエラー定義
- [ ] ドメイン層ユニットテスト作成

## 受け入れ基準
- [ ] プロジェクトのライフサイクル管理が実装されている
- [ ] ステータス遷移のビジネスルールが守られている
- [ ] 適切なドメインエラーが定義されている

## 実装ファイル
- `backend/internal/domain/project.go`
- `backend/internal/domain/project_status.go`
- `backend/internal/domain/project_repository.go`
- `backend/internal/domain/project_test.go`
```

### Issue 2.5: 実行履歴ドメイン実装
**Title**: Phase 2.5: ExecutionHistory エンティティ実装
**Labels**: Phase 2, domain, execution
**Body**:
```
## 概要
実行履歴管理のためのExecutionHistoryエンティティとExecutionStatusを実装します。

## タスク詳細
- [ ] ExecutionHistory エンティティの実装
- [ ] ExecutionStatus バリューオブジェクトの実装
- [ ] 実行結果の記録に関するビジネスルール
- [ ] ExecutionHistoryRepository インターフェース定義
- [ ] 実行履歴関連ドメインエラー定義
- [ ] ドメイン層ユニットテスト作成

## 実装ファイル
- `backend/internal/domain/execution_history.go`
- `backend/internal/domain/execution_status.go`
- `backend/internal/domain/execution_history_repository.go`
- `backend/internal/domain/execution_history_test.go`
```

## Phase 3: アプリケーション層実装

### Issue 3.1: ユーザー管理ユースケース実装
**Title**: Phase 3.1: ユーザー管理ユースケース実装
**Labels**: Phase 3, usecase, user
**Body**:
```
## 概要
ユーザー管理に関するユースケースを実装します。

## タスク詳細
- [ ] RegisterUser ユースケース実装
- [ ] AuthenticateUser ユースケース実装
- [ ] UpdateUserProfile ユースケース実装
- [ ] 各ユースケースのエラーハンドリング
- [ ] ユースケース層のユニットテスト作成

## 受け入れ基準
- [ ] トランザクション境界が適切に設定されている
- [ ] ドメインオブジェクトを通じてビジネスロジックが実行されている
- [ ] 適切なエラーハンドリングが実装されている

## 実装ファイル
- `backend/internal/usecase/register_user.go`
- `backend/internal/usecase/authenticate_user.go`
- `backend/internal/usecase/update_user_profile.go`
- テストファイル群
```

### Issue 3.2: ワークフロー管理ユースケース実装
**Title**: Phase 3.2: ワークフロー管理ユースケース実装
**Labels**: Phase 3, usecase, workflow
**Body**:
```
## 概要
ワークフロー管理に関するユースケースを実装します。

## タスク詳細
- [ ] CreateWorkflow ユースケース実装
- [ ] ExecuteWorkflow ユースケース実装
- [ ] UpdateWorkflowStatus ユースケース実装
- [ ] ワークフロー実行時のエラーハンドリング
- [ ] ユースケース層のユニットテスト作成

## 実装ファイル
- `backend/internal/usecase/create_workflow.go`
- `backend/internal/usecase/execute_workflow.go`
- `backend/internal/usecase/update_workflow_status.go`
- テストファイル群
```

## Phase 4: インフラストラクチャ層実装

### Issue 4.1: データベースマイグレーション作成
**Title**: Phase 4.1: データベース設計とマイグレーション作成
**Labels**: Phase 4, database, migration
**Body**:
```
## 概要
PostgreSQL用のテーブル設計とマイグレーションスクリプトを作成します。

## タスク詳細
- [ ] users テーブルのマイグレーション作成
- [ ] workflows テーブルのマイグレーション作成
- [ ] projects テーブルのマイグレーション作成
- [ ] execution_histories テーブルのマイグレーション作成
- [ ] 適切なインデックス設計
- [ ] 初期データ投入スクリプト作成

## 実装ファイル
- `backend/migrations/` 配下のマイグレーションファイル群
```

### Issue 4.2: PostgreSQLリポジトリ実装
**Title**: Phase 4.2: PostgreSQL リポジトリ実装
**Labels**: Phase 4, repository, postgresql
**Body**:
```
## 概要
ドメイン層で定義されたリポジトリインターフェースのPostgreSQL実装を行います。

## タスク詳細
- [ ] UserRepository PostgreSQL実装
- [ ] WorkflowRepository PostgreSQL実装
- [ ] ProjectRepository PostgreSQL実装
- [ ] ExecutionHistoryRepository PostgreSQL実装
- [ ] 統合テストの作成

## 実装ファイル
- `backend/internal/interface/persistence/` 配下の実装ファイル群
```

## Phase 5: Webインターフェース実装

### Issue 5.1: HTTPハンドラー・ルーター実装
**Title**: Phase 5.1: Echo HTTPハンドラー・ルーター実装
**Labels**: Phase 5, http, api
**Body**:
```
## 概要
Echo フレームワークを使用したHTTPハンドラーとルーティング実装します。

## タスク詳細
- [ ] Echo サーバー設定
- [ ] ルーティング実装
- [ ] ミドルウェア実装（CORS, Logger, Recovery）
- [ ] DTOの定義
- [ ] エラーレスポンスの標準化

## 実装ファイル
- `backend/internal/interface/http/` 配下の実装ファイル群
- `backend/cmd/server/main.go` の更新
```

### Issue 5.2: RESTful API実装
**Title**: Phase 5.2: RESTful API エンドポイント実装
**Labels**: Phase 5, api, rest
**Body**:
```
## 概要
各ドメインのRESTful APIエンドポイントを実装します。

## タスク詳細
- [ ] ユーザー管理API実装
- [ ] ワークフロー管理API実装
- [ ] プロジェクト管理API実装
- [ ] 実行管理API実装
- [ ] APIのドキュメント作成

## 実装ファイル
- 各種ハンドラー実装ファイル
- OpenAPI仕様書
```

## Phase 6: フロントエンド実装

### Issue 6.1: Next.js基盤・レイアウト実装
**Title**: Phase 6.1: Next.js App Router基盤とレイアウト実装
**Labels**: Phase 6, frontend, nextjs
**Body**:
```
## 概要
Next.js App Routerの基盤設定とアプリケーション全体のレイアウトを実装します。

## タスク詳細
- [ ] App Router設定
- [ ] レイアウトコンポーネント実装
- [ ] プロバイダー実装（Recoil, React Query）
- [ ] グローバルCSS設定
- [ ] 基本的なナビゲーション実装

## 実装ファイル
- `frontend/app/layout.tsx`
- `frontend/app/providers.tsx`
- `frontend/styles/globals.css`
```

### Issue 6.2: 共通UIコンポーネント実装
**Title**: Phase 6.2: 共通UIコンポーネントライブラリ実装
**Labels**: Phase 6, frontend, components
**Body**:
```
## 概要
アプリケーション全体で使用する共通UIコンポーネントを実装します。

## タスク詳細
- [ ] Button, Input, Modal などの基本コンポーネント
- [ ] フォームバリデーション機能
- [ ] エラーハンドリングコンポーネント
- [ ] ローディング表示コンポーネント
- [ ] Tailwind CSSを使用したスタイリング

## 実装ファイル
- `frontend/components/ui/` 配下のコンポーネント群
```

### Issue 6.3: 画面実装（認証系）
**Title**: Phase 6.3: ログイン・認証画面実装
**Labels**: Phase 6, frontend, auth
**Body**:
```
## 概要
ユーザー認証に関連する画面を実装します。

## タスク詳細
- [ ] ログイン画面の実装
- [ ] ユーザー登録画面の実装
- [ ] 認証状態管理（Recoil）
- [ ] 認証ガードの実装
- [ ] HTTP-only Cookie との連携

## 実装ファイル
- `frontend/app/login/page.tsx`
- `frontend/app/register/page.tsx`
- 認証関連のコンポーネント・サービス
```

### Issue 6.4: ダッシュボード画面実装
**Title**: Phase 6.4: ダッシュボード画面実装
**Labels**: Phase 6, frontend, dashboard
**Body**:
```
## 概要
アプリケーションのメインダッシュボード画面を実装します。

## タスク詳細
- [ ] ダッシュボードレイアウト実装
- [ ] プロジェクト一覧表示
- [ ] ワークフロー実行状況表示
- [ ] 統計情報の表示
- [ ] リアルタイム更新機能

## 実装ファイル
- `frontend/app/dashboard/page.tsx`
- 関連コンポーネント群
```

### Issue 6.5: サービス層・状態管理実装
**Title**: Phase 6.5: API呼び出しサービス・状態管理実装
**Labels**: Phase 6, frontend, services
**Body**:
```
## 概要
バックエンドAPIとの通信およびクライアント状態管理を実装します。

## タスク詳細
- [ ] API呼び出しサービス実装
- [ ] Recoil による状態管理実装
- [ ] React Query によるサーバー状態管理
- [ ] エラーハンドリング機能
- [ ] キャッシュ戦略の実装

## 実装ファイル
- `frontend/services/` 配下のサービス群
- `frontend/lib/apiClient.ts`
- Recoil atoms/selectors
```

## 品質保証・デプロイメント関連

### Issue 7.1: CI/CD パイプライン設定
**Title**: CI/CD: GitHub Actions パイプライン設定
**Labels**: ci-cd, devops
**Body**:
```
## 概要
GitHub Actions による CI/CD パイプラインを設定します。

## タスク詳細
- [ ] 自動テスト実行のワークフロー
- [ ] コードカバレッジ測定
- [ ] 静的解析（golangci-lint, ESLint）
- [ ] セキュリティスキャン
- [ ] ビルド・デプロイワークフロー

## 実装ファイル
- `.github/workflows/` 配下のワークフローファイル群
```

### Issue 7.2: AWS インフラストラクチャ構築
**Title**: AWS: Terraform による Infrastructure as Code 実装
**Labels**: aws, terraform, infrastructure
**Body**:
```
## 概要
AWS インフラストラクチャを Terraform で構築します。

## タスク詳細
- [ ] ECS/Fargate + ALB 設定
- [ ] RDS PostgreSQL 設定
- [ ] S3 + CloudFront 設定
- [ ] WAF + Security Group 設定
- [ ] 環境別設定（dev/staging/prod）

## 実装ファイル
- `infrastructure/terraform/` 配下のTerraformファイル群

## 実装者向け注意
この作業は AWS リソースの作成を伴うため、実際のデプロイは慎重に行ってください。
```

---

## 総括

計35ステップを以下の粒度で分割しました：
- **実装可能なタスク**: 25個のIssue（コード実装中心）
- **設計・分析タスク**: 7個のIssue（ドキュメント作成中心）  
- **インフラ・デプロイタスク**: 5個のIssue（環境構築中心）

各IssueはClaude Codeで実装可能な粒度まで詳細化し、具体的なファイルパスと受け入れ基準を明記しました。