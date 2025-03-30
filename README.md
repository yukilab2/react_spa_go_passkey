# Passkey認証サンプルアプリケーション

このプロジェクトはReact SPAとGoバックエンドを使用したPasskey認証のサンプル実装です。React部分はTypeScriptで記述され、Node.js v22環境で動作します。macOS上でのローカル開発環境で動作します。

## 前提条件

- Python (nodeenvのインストールに必要)
- nodeenv (Node.js v22環境のセットアップに使用)
- Yarn
- Go (v1.16以上)
- macOS

## プロジェクト構成

- `frontend/` - React SPA (TypeScript)
- `backend/` - Goバックエンド
- `.env` - 許可されたEmailアドレスを含む設定ファイル
- `.nvmrc`, `.nodeenvrc` - Node.js v22の設定ファイル

## セットアップ手順

1. リポジトリをクローン
   ```
   git clone <repository-url>
   cd passkey_sample
   ```

2. nodeenvでNode.js v22環境のセットアップ
   ```
   make setup-nodeenv
   source node_env/bin/activate
   ```

3. 依存関係のインストール
   ```
   make setup
   ```

## 実行方法

### バックエンドの起動
```
make backend
```
バックエンドサーバーは `http://localhost:8080` で起動します。

### フロントエンドの起動（開発モード）
```
source node_env/bin/activate  # Node.js v22環境をアクティブ化
make frontend-dev
```
開発サーバーは `http://localhost:3000` で起動します。

### フロントエンドのビルドと実行
```
source node_env/bin/activate  # Node.js v22環境をアクティブ化
make frontend-build
make frontend-start
```

## 使用方法

1. フロントエンドとバックエンドの両方を起動します。
2. ブラウザで `http://localhost:3000` にアクセスします。
3. `.env` ファイルに登録されているメールアドレスを使用して、パスキーを登録できます。
4. 登録後、パスキーを使用してログインできます。

## 技術スタック

- フロントエンド: React 18 (TypeScript)、Material UI、React Router
- バックエンド: Go (Gin、WebAuthn)
- 認証方式: Passkey (WebAuthn)
- Node.js: v22 (nodeenv使用)

## 注意事項

- このサンプルはローカル開発環境での動作を目的としており、本番環境での使用は想定していません。
- `.env` ファイルには、パスキー登録が許可されているメールアドレスを1行に1つずつ記載してください。
- Node.js v22環境を使用するには、コマンド実行前に `source node_env/bin/activate` を実行してください。 