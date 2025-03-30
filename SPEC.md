# 要件

- React SPA と Go backendで構成する。
- macosx ローカルで動作すれば良い。
- reactは、passkeyによるログインページと ログイン後の welcomeページからなる。
- welcomeページにはログアウトボタンがある。
- ログインページには以下の要素を含める：
  - パスキーによるログインボタン
  - パスキー登録ボタン（初回ユーザー向け）
- ユーザー制限：
  - パスキー登録は go の読み取れる .envファイルに記載されたemailによるものだけが行うことができる。
  - 登録時にユーザーが入力したemailと.envファイルのemailを照合して、一致する場合のみPasskey登録を許可する。
- エラーハンドリング：
  - Passkey登録・認証失敗時には、適切なエラーメッセージをユーザーに表示する。
  - 未登録ユーザーのログイン試行や、.envファイルに存在しないemailでの登録試行に対するエラー処理を実装する。
- 構築用のスクリプトは必要ないが、 rootに Makefileをおき、react の実行 (yarn dev / yarn build / yarn start) と　go の実行 (go run main.go) をそこから行えるようにする。
- 使い方は README.md に記載する。
- node: v22 を nodeenv から使用する
- react は typescript で記載する
- reactのUIコンポーネントは radix-ui / css-moduleを使用する
