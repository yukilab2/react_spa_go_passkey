.PHONY: frontend-dev frontend-build frontend-start backend setup-nodeenv

# Node.js環境をセットアップ
setup-nodeenv:
	@echo "nodeenv を使用してNode.js v22環境をセットアップ中..."
	@if ! command -v nodeenv > /dev/null; then \
		echo "nodeenvをインストール中..."; \
		pip install nodeenv; \
	fi
	@if [ ! -d "node_env" ]; then \
		echo "Node.js環境を作成中..."; \
		nodeenv --node=22.0.0 node_env; \
	fi
	@echo "Node.js環境のセットアップが完了しました。使用するには 'source node_env/bin/activate' を実行してください。"

# React関連コマンド
frontend-dev:
		cd frontend && yarn dev

frontend-build:
	cd frontend && yarn build

frontend-start:
	cd frontend && yarn start

# Goバックエンド関連コマンド
backend:
	cd backend && go run main.go

# 開発環境のセットアップ
setup: setup-nodeenv
	@echo "フロントエンドの依存関係をインストール中..."
	cd frontend && yarn install
	@echo "バックエンドの依存関係をインストール中..."
	cd backend && go mod tidy

# ヘルプメッセージ
help:
	@echo "使用可能なコマンド:"
	@echo "  make setup-nodeenv    - nodeenvを使用してNode.js v22環境をセットアップ"
	@echo "  make frontend-dev     - React開発サーバーを起動 (開発モード)"
	@echo "  make frontend-build   - Reactアプリをビルド"
	@echo "  make frontend-start   - ビルドしたReactアプリを起動"
	@echo "  make backend          - Goバックエンドを起動"
	@echo "  make setup            - 開発環境をセットアップ" 