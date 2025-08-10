たしかに、先ほどの翻訳は「Claudeという製品の一般ポリシー」まで含まれており、**Claude.md に入れるべき“実装ガードレール”の趣旨から外れていました**。
ここでは目的を絞り直し、**React(Next.js) + Tailwind CSS + Go(Echo)**、**AWSデプロイ前提**、**DDD + クリーンアーキテクチャ**での実装をClaude Codeに依頼するための、**そのまま使える Claude.md サンプル**をお渡しします。

---

# Claude.md（実装方針 / ガードレール）

**対象**: Next.js (React) + Tailwind CSS / Go (echo) / AWS
**設計原則**: Domain-Driven Design（DDD） + Clean Architecture
**運用**: GitHub Issues を @claude にアサイン → GitHub Actions で自動実装・PR化

## 1) このドキュメントの目的

* Claude Code（@claude）が行う**実装の一貫性**と**品質**を担保するためのガードレール。
* **層の責務・依存方向・命名規約・フォルダ構成・テスト/セキュリティ/CI**を明文化。
* 本ドキュメントは、Issue に添付（リンク）し、**Claude は必ずこれを遵守**して実装すること。

---

## 2) 技術スタックと共通方針

* **フロント**: Next.js（原則 App Router / 必要に応じ Pages Router 可）、TypeScript、Tailwind CSS、Recoil（クライアント状態）、（任意で）TanStack Query（サーバ状態）
* **バックエンド**: Go 1.22+、echo v4、DBは PostgreSQL（RDS）を想定
* **インフラ**:

  * フロント: S3 + CloudFront（原則 SSG/CSR。SSRが必要な場合は Amplify Hosting or ECS/EC2 Node を別途用意）
  * バック: ECS/Fargate + ALB + Private Subnet、RDS(PostgreSQL)、Secrets Manager、CloudWatch、WAF
* **認証**: 原則 **HTTP-only Cookie**（アクセストークンを localStorage に保存しない）
* **I18N/時刻**: 時刻は原則 UTC 保持・表示はロケール変換
* **エラーハンドリング**: ドメインエラー（ビジネス）とインフラエラー（技術）を区別して扱う

---

## 3) DDD + クリーンアーキテクチャ原則（必読）

**依存の向きは常に内向き**。フレームワークは外側、ドメインは内側。

* **Domain（エンティティ/値オブジェクト/ドメインサービス）**

  * ビジネスルールの唯一の所在。**外部技術に依存しない**（echo/SQL/HTTP等の import 禁止）。
  * **集約ルート**に対してのみ外部から書き込みを許容。\*\*不変条件（Invariant）\*\*は集約内で守る。
* **UseCase/Application**

  * **ユースケースオーケストレーション**。ドメイン操作の手順・トランザクション境界を定義。
  * 外部 I/O は **ポート（インターフェース）** 経由のみ。
* **Interface Adapters（Controller/Presenter/Gateway）**

  * 入出力の**変換層**。HTTPハンドラ/Echo、DTO、DB リポジトリ実装などはここ。
  * **ビジネスルール禁止**。変換・検証・マッピングに徹する。
* **Infrastructure**

  * DB / 外部API / メール / キャッシュ等の実装詳細。
  * **ドメインに依存は可**、逆は不可。

**禁止事項**（Claude は生成しないこと）

* ドメイン層で echo.Context / SQL / HTTP クライアントを import
* 直接的な SQL/HTTP 呼び出しをユースケースやエンティティに記述
* UI コンポーネント内でビジネスロジック・API呼び出し（サービス層経由にする）
* アナミックドメイン（振る舞いゼロの巨大 DTO）の氾濫
* トランザクションを複数ユースケースでまたぐ設計

---

## 4) ディレクトリ構成（バックエンド / Go + echo）

```
backend/
├── cmd/
│   └── server/main.go          # エントリポイント（DI/起動/ルーティング）
├── internal/
│   ├── domain/                 # 100% 純粋なドメイン（技術依存なし）
│   │   ├── user.go
│   │   ├── errors.go
│   │   └── user_repository.go  # ポート（インターフェース）
│   ├── usecase/                # アプリケーションサービス（ユースケース）
│   │   └── register_user.go
│   ├── interface/              # 入出力変換
│   │   ├── http/               # echo ハンドラ/ルーティング/DTO
│   │   │   ├── router.go
│   │   │   └── user_handler.go
│   │   └── persistence/        # DBアダプタ（リポジトリ実装）
│   │       └── user_repo_pg.go
│   └── platform/               # 共有基盤（DB接続、設定、ログ）
│       ├── db.go
│       └── config.go
├── migrations/                 # DBマイグレーション（goose等）
├── go.mod / go.sum
└── Makefile / Taskfile.yml
```

### バックエンド・サンプルコード（最小骨格）

**domain/user.go**

```go
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrEmailAlreadyUsed = errors.New("email already used")
	ErrInvalidEmail     = errors.New("invalid email")
	ErrEmptyName        = errors.New("name is empty")
)

type User struct {
	ID        uuid.UUID
	Name      string
	Email     string
	CreatedAt time.Time
}

func NewUser(name, email string, now func() time.Time) (*User, error) {
	if name == "" {
		return nil, ErrEmptyName
	}
	if !isValidEmail(email) {
		return nil, ErrInvalidEmail
	}
	return &User{
		ID:        uuid.New(),
		Name:      name,
		Email:     email,
		CreatedAt: now(),
	}, nil
}

func isValidEmail(s string) bool {
	// ここでは簡易チェック。詳細は別途 ValueObject に切り出してもよい。
	return len(s) >= 3 && len(s) <= 255 && containsAt(s)
}

func containsAt(s string) bool {
	for _, ch := range s {
		if ch == '@' {
			return true
		}
	}
	return false
}
```

**domain/user\_repository.go**

```go
package domain

import "context"

type UserRepository interface {
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	Save(ctx context.Context, u *User) error
}
```

**usecase/register\_user.go**

```go
package usecase

import (
	"context"
	"time"

	"example.com/app/internal/domain"
)

type RegisterUserInput struct {
	Name  string
	Email string
}

type RegisterUser struct {
	Repo domain.UserRepository
	Now  func() time.Time
}

func (uc *RegisterUser) Execute(ctx context.Context, in RegisterUserInput) (*domain.User, error) {
	if exists, err := uc.Repo.ExistsByEmail(ctx, in.Email); err != nil {
		return nil, err
	} else if exists {
		return nil, domain.ErrEmailAlreadyUsed
	}

	user, err := domain.NewUser(in.Name, in.Email, uc.Now)
	if err != nil {
		return nil, err
	}
	if err := uc.Repo.Save(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}
```

**interface/persistence/user\_repo\_pg.go**

```go
package persistence

import (
	"context"
	"database/sql"

	"example.com/app/internal/domain"
)

type UserRepoPG struct {
	DB *sql.DB
}

func (r *UserRepoPG) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	const q = `SELECT 1 FROM users WHERE email = $1 LIMIT 1`
	var x int
	err := r.DB.QueryRowContext(ctx, q, email).Scan(&x)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return err == nil, err
}

func (r *UserRepoPG) Save(ctx context.Context, u *domain.User) error {
	const q = `
		INSERT INTO users(id, name, email, created_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.DB.ExecContext(ctx, q, u.ID, u.Name, u.Email, u.CreatedAt)
	return err
}
```

**interface/http/user\_handler.go**

```go
package http

import (
	"net/http"

	"example.com/app/internal/domain"
	"example.com/app/internal/usecase"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	Register *usecase.RegisterUser
}

type registerUserReq struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type userRes struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (h *UserHandler) RegisterUser(c echo.Context) error {
	var req registerUserReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}

	u, err := h.Register.Execute(c.Request().Context(), usecase.RegisterUserInput{
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		switch err {
		case domain.ErrEmptyName, domain.ErrInvalidEmail:
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		case domain.ErrEmailAlreadyUsed:
			return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
		}
	}

	return c.JSON(http.StatusCreated, userRes{
		ID:    u.ID.String(),
		Name:  u.Name,
		Email: u.Email,
	})
}
```

**interface/http/router.go**

```go
package http

import "github.com/labstack/echo/v4"

func RegisterRoutes(e *echo.Echo, uh *UserHandler) {
	v1 := e.Group("/api/v1")
	v1.POST("/users", uh.RegisterUser)
}
```

**cmd/server/main.go**

```go
package main

import (
	"database/sql"
	"log"
	"time"

	"example.com/app/internal/interface/http"
	"example.com/app/internal/interface/persistence"
	"example.com/app/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	// DB接続
	dsn := mustEnv("DB_DSN")
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// DI
	userRepo := &persistence.UserRepoPG{DB: db}
	registerUC := &usecase.RegisterUser{Repo: userRepo, Now: time.Now}
	uh := &http.UserHandler{Register: registerUC}

	// Echo
	e := echo.New()
	e.Use(middleware.Recover(), middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://your-frontend-domain.example"},
		AllowMethods: []string{echo.POST, echo.GET, echo.PUT, echo.DELETE},
		AllowHeaders: []string{"Content-Type", "Authorization"},
	}))
	http.RegisterRoutes(e, uh)

	e.Logger.Fatal(e.Start(":8080"))
}

func mustEnv(key string) string {
	v := getenv(key, "")
	if v == "" {
		log.Fatalf("missing env: %s", key)
	}
	return v
}

func getenv(key, def string) string {
	if v, ok := syscall.Getenv(key); ok {
		return v
	}
	return def
}
```

---

## 5) フロントエンド構成（Next.js + Tailwind + Recoil）

```
frontend/
├── app/                         # App Router 推奨（Pages Routerでも可）
│   ├── layout.tsx
│   ├── providers.tsx            # Recoil/QueryClient のプロバイダ
│   └── users/
│       └── new/page.tsx
├── components/
│   └── users/UserForm.tsx
├── lib/
│   ├── apiClient.ts             # fetch ラッパ（/api 経由）
│   └── config.ts                # APIベースURL等
├── services/
│   └── users.ts                 # API 呼び出し集約（UIから直 fetch 禁止）
├── styles/globals.css
├── tailwind.config.ts
└── next.config.ts
```

**Tailwind セットアップ（抜粋）**

```ts
// tailwind.config.ts
import type { Config } from 'tailwindcss'
export default {
  content: ['./app/**/*.{ts,tsx}', './components/**/*.{ts,tsx}'],
  theme: { extend: {} },
  plugins: [],
} satisfies Config
```

**グローバルレイアウト / プロバイダ**

```tsx
// app/layout.tsx
import './globals.css'
import Providers from './providers'

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="ja">
      <body className="min-h-screen bg-gray-50 text-gray-900">
        <Providers>{children}</Providers>
      </body>
    </html>
  )
}
```

```tsx
// app/providers.tsx
'use client'
import { RecoilRoot } from 'recoil'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { useState } from 'react'

export default function Providers({ children }: { children: React.ReactNode }) {
  const [qc] = useState(() => new QueryClient())
  return (
    <RecoilRoot>
      <QueryClientProvider client={qc}>{children}</QueryClientProvider>
    </RecoilRoot>
  )
}
```

**API クライアント**

```ts
// lib/apiClient.ts
export async function api<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(path, {
    credentials: 'include',
    headers: { 'Content-Type': 'application/json', ...(init?.headers || {}) },
    ...init,
  })
  if (!res.ok) {
    const body = await res.json().catch(() => ({}))
    throw new Error(body?.error || `HTTP ${res.status}`)
  }
  return res.json() as Promise<T>
}
```

**サービス層**

```ts
// services/users.ts
type CreateUserReq = { name: string; email: string }
type UserRes = { id: string; name: string; email: string }

export const createUser = (payload: CreateUserReq) =>
  api<UserRes>('/api/v1/users', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
```

**UI コンポーネント**

```tsx
// components/users/UserForm.tsx
'use client'
import { useState } from 'react'

type Props = { onSubmit: (v: { name: string; email: string }) => Promise<void>; submitting?: boolean }

export default function UserForm({ onSubmit, submitting }: Props) {
  const [name, setName] = useState('')
  const [email, setEmail] = useState('')
  const [error, setError] = useState<string | null>(null)

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError(null)
    try {
      await onSubmit({ name, email })
    } catch (err: any) {
      setError(err.message ?? 'エラーが発生しました')
    }
  }

  return (
    <form onSubmit={handleSubmit} className="max-w-md space-y-4 rounded-lg bg-white p-6 shadow">
      <div>
        <label className="block text-sm font-medium">名前</label>
        <input className="mt-1 w-full rounded border px-3 py-2" value={name} onChange={e => setName(e.target.value)} />
      </div>
      <div>
        <label className="block text-sm font-medium">メール</label>
        <input className="mt-1 w-full rounded border px-3 py-2" value={email} onChange={e => setEmail(e.target.value)} />
      </div>
      {error && <p className="text-sm text-red-600">{error}</p>}
      <button
        disabled={submitting}
        className="inline-flex items-center rounded bg-black px-4 py-2 text-white disabled:opacity-50"
      >
        {submitting ? '送信中...' : '登録'}
      </button>
    </form>
  )
}
```

**ページ**

```tsx
// app/users/new/page.tsx
'use client'
import { useState } from 'react'
import UserForm from '@/components/users/UserForm'
import { createUser } from '@/services/users'
import { useRouter } from 'next/navigation'

export default function Page() {
  const [loading, setLoading] = useState(false)
  const router = useRouter()

  return (
    <div className="mx-auto max-w-2xl p-6">
      <h1 className="mb-6 text-2xl font-semibold">ユーザー登録</h1>
      <UserForm
        submitting={loading}
        onSubmit={async (v) => {
          setLoading(true)
          try {
            await createUser(v)
            router.push('/users') // 一覧へ遷移（実装済とする）
          } finally {
            setLoading(false)
          }
        }}
      />
    </div>
  )
}
```

**Next のリライト（ローカル開発でバックエンドへプロキシ）**

```ts
// next.config.ts
import type { NextConfig } from 'next'
const config: NextConfig = {
  async rewrites() {
    return [{ source: '/api/:path*', destination: 'http://localhost:8080/api/:path*' }]
  },
}
export default config
```

**ガードレール（フロント）**

* **UI コンポーネントは I/O 禁止**（API呼び出しは `services/*` へ）
* **ビジネスルールはバックエンド優先**（フロントでの重複実装禁止。UI側は軽量バリデーションのみ）
* **状態管理**:

  * Recoil＝UI/セッション用の**クライアント状態**
  * サーバ状態は（任意で）**TanStack Query**を利用（キャッシュ・リトライ・並列制御）
* **CSS**: Tailwind ユーティリティを基本に、複雑化時は**小さなコンポーネントへ分割**

---

## 6) 例: エンドツーエンドのデータフロー（ユーザー登録）

1. `POST /api/v1/users`（Echo Handler → UseCase → Repository）
2. DB へ INSERT（RepoPG）
3. 成功時 201 + JSON を返却
4. フロント `services/users.createUser()` が受け取り、UI に反映/遷移

**境界確認**

* バリデーションの中核（空文字/メール形式/重複）は **ドメイン or UseCase**
* Handler は **変換とHTTP応答** のみ
* UI は **エラー表示とUX** のみ（ビジネス判断は禁止）

---

## 7) テスト方針

* **ドメイン**: `*_test.go` でユニットテスト（ビジネスルールを網羅、外部依存なしで高速）
* **ユースケース**: リポジトリをモック/フェイク化し、分岐をテスト
* **インフラ**: リポジトリは**統合テスト**（テスト用DB、Txロールバック）
* **フロント**: 重要 UI のレンダ/送信/エラー表示を React Testing Library でテスト
* **DoD（Definition of Done）**

  * 新規/変更ロジックに**テストがある**
  * CI（lint/format/test/build）が**全て緑**
  * **アーキテクチャ境界の違反がない**

---

## 8) セキュリティ / エラー / ロギング

* **トークンは HTTP-only Cookie**。localStorage 保存は禁止
* **入力検証**: バックエンドで必ず検証（UIのバリデは補助）
* **CORS**: 許可オリジンを明示。`*` は使わない
* **秘密情報**: AWS Secrets Manager / SSM Parameter Store 管理。コード直書き禁止
* **ログ**: CloudWatch へ構造化出力（PIIはマスク）
* **エラー**: ドメインエラー（4xx相当）、技術エラー（5xx）を明確にマッピング

---

## 9) AWS デプロイ指針（最小構成）

* **フロント**: `next build && next export`（CSR/SSG主軸）→ S3 → CloudFront（Cache/Compression/HTTPS）
  ※ SSR 必須の場合は Amplify Hosting か ECS/EC2 Node を別途構築
* **バック**:

  * **ECS Fargate**（private subnet） + **ALB**（public）
  * **RDS(PostgreSQL)**（private）
  * **Secrets Manager**（DB資格情報）
  * **Security Group**: ALB→ECS のみ / ECS→RDS のみ
  * **WAF**（OWASP基準）
  * **CloudWatch Logs / Metrics / Alarms**
* **CI/CD**: GitHub Actions → ECR push → ECS service update / S3 sync + CloudFront Invalidation

---

## 10) 依頼フロー（@claude / GitHub Actions）

* **Issue を @claude にアサイン**（または `ai` ラベル付与）→ Actions が起動 → 本 `Claude.md` を文脈投入
* Claude は**本ガイドに準拠**して実装 → ブランチ作成 → PR 作成 → CI 実行
* レビューで指摘があれば、**同Issueに追記**し再実装を依頼

**Issue テンプレ（サンプル）**

```md
### 背景/目的
（ドメイン言語で簡潔に）

### 要件
- Domain: （エンティティ/VO/ルール）
- UseCase: （入出力/トランザクション境界）
- Interface: （HTTP ルート・DTO）
- Infra: （リポジトリ/外部API/マイグレーション）

### 受け入れ基準
- [ ] ドメイン不変条件が満たされる
- [ ] 4xx/5xx のハンドリングが要件通り
- [ ] テストが追加され CI 緑

### 参照
- このファイル: /docs/Claude.md
- 該当モジュール: （パス）
```

**PR チェックリスト（レビュア向け）**

* [ ] 依存方向（内向き）の遵守
* [ ] ドメイン層が技術から独立している
* [ ] ユースケースが I/O を直接触っていない
* [ ] Handler は変換/応答に限定されている
* [ ] フロントがサービス層経由で API を呼び出している
* [ ] テスト/リンタ/型チェックが通過

---

## 11) 命名規約 / コーディング規約（抜粋）

* **Go**: パッケージ小文字、公開型/関数は PascalCase、エラーは `ErrXxx`
* **TS/React**: コンポーネントは PascalCase、hook は `useXxx`、型は `PascalCase`
* **HTTP**: RESTful に `/api/v1/{resources}`。バージョンはパスで管理
* **コミット**: Conventional Commits（例: `feat(user): add register use case`）

---

## 12) 許可ライブラリ（推奨）

* **Go**: `echo`, `pgx` or `database/sql`, `golang-migrate`/`goose`, `zap`（ログ）
* **TS**: `recoil`, `@tanstack/react-query`（任意）, `zod`（バリデーション任意）

---

## 13) 付録：DB マイグレーション（例）

```sql
-- migrations/20250101000001_create_users.sql
CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY,
  name TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL
);
```

