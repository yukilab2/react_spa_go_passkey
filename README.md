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

## 改善点と今後の検討事項

### セキュリティ関連
- **データの永続化**: 現在ユーザー情報はメモリ上に保存されていますが、本番環境では以下を検討する必要があります：
  - データベースなどの永続的なストレージの使用
  - パスワードなどの機密情報の適切なハッシュ化・暗号化
  - セッションストアをRedisなどの外部ストアに移行

- **セッション管理の改善**:
  - フロントエンド: localStorageからHTTP Only Cookieを使用したセッショントークン管理への移行
  - バックエンド: メモリベースのセッションストアから永続的なストレージへの移行

### コード品質
- **エラーハンドリングの改善**:
  - `log.Fatal`/`log.Fatalf`の使用を見直し、適切なエラーレスポンスの実装
  - より詳細なエラーメッセージとログの提供

- **コード構成の最適化**:
  - バックエンドコードの機能別モジュール分割（ユーザー管理、認証、WebAuthn設定など）
  - フロントエンドの型定義の整理（`types`ディレクトリの活用）
  - APIレスポンス/リクエストの型定義の共通化

### 運用面
- **設定管理**:
  - `.env`ファイルの直接読み込みから環境変数経由の設定管理への移行
  - 本番/開発/テスト環境ごとの設定分離

- **依存関係管理**:
  - `frontend/package.json`と`backend/go.mod`の定期的な更新確認
  - セキュリティアップデートの適用

- **テスト整備**:
  - 単体テストの追加
  - 結合テストの実装
  - E2Eテストの導入

### デプロイメント
- CI/CDパイプラインの構築
- コンテナ化対応
- スケーリング戦略の検討

これらの改善は、アプリケーションを本番環境で運用する際に考慮すべき重要な要素です。現在の実装はあくまでもサンプルとしての位置づけであり、本番環境での使用には上記の対応が推奨されます。 