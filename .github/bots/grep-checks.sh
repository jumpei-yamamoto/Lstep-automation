#!/usr/bin/env bash
set -euo pipefail

OUT=".github/bots/out/grep-violations.txt"
mkdir -p "$(dirname "$OUT")"
: > "$OUT"

# Claude.md の存在（任意: docs/Claude.md または 直下）
if [ ! -f docs/Claude.md ] && [ ! -f ./Claude.md ]; then
  echo "MISSING: docs/Claude.md が見つかりません（ガードレール文書の配置を推奨）" | tee -a "$OUT"
fi

# backend の CORS "*"（保険）
if grep -R --include='*.go' -n 'AllowOrigins: \[\]string{"\*"}' backend >/dev/null 2>&1; then
  echo 'CORS: AllowOrigins に "*" が設定されています（明示的なオリジンを指定してください）' | tee -a "$OUT"
fi

# localStorage の全体検知（警告）
if grep -R --include='*.{ts,tsx,js,jsx}' -n '\blocalStorage\b' frontend 2>/dev/null; then
  echo 'localStorage: 使用を検出（トークン保存には使用しない。HTTP-only Cookie を原則）。' | tee -a "$OUT"
fi

echo "grep checks complete. See $OUT"
exit 0
