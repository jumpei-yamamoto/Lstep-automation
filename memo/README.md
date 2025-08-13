# Infrastructure Documentation

このフォルダには、プロジェクトのインフラストラクチャ関連の設定ファイルとその使用方法について説明したドキュメントが含まれています。

## ドキュメント構成

- **[overview.md](overview.md)** - インフラストラクチャ全体の概要
- **[docker.md](docker.md)** - Dockerfileの説明と使用方法
- **[docker-compose.md](docker-compose.md)** - 開発環境のセットアップと運用
- **[terraform.md](terraform.md)** - AWS本番環境のデプロイ設定

## 開発環境クイックスタート

AWS環境が未構築でも、以下の手順でローカル開発環境を構築できます：

```bash
# リポジトリをクローン
git clone <repository-url>
cd Lstep-automation

# 開発環境の起動
docker-compose -f docker-compose.dev.yml up -d

# アクセス確認
# Frontend: http://localhost:3000
# Backend API: http://localhost:8080
# PostgreSQL: localhost:5432
```

詳細な手順は [docker-compose.md](docker-compose.md) を参照してください。

## 目的

- **開発環境の一貫性**: 全ての開発者が同じ環境で作業できる
- **本番環境への円滑な移行**: Docker化により開発と本番の差異を最小化
- **スケーラブルなインフラ**: AWS上でのコンテナベースな本番環境
- **DDD/Clean Architecture準拠**: CLAUDE.mdの方針に従った設計

各ドキュメントを参照して、目的に応じた環境を構築してください。