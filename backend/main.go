package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/webauthn"
)

// ユーザー構造体
type User struct {
	ID                      string
	Name                    string
	DisplayName             string
	Credentials             []webauthn.Credential
	RegistrationSessionData *webauthn.SessionData // セッションデータを保持
}

// WebAuthnのインターフェースを実装
func (u *User) WebAuthnID() []byte {
	return []byte(u.ID)
}

func (u *User) WebAuthnName() string {
	return u.Name
}

func (u *User) WebAuthnDisplayName() string {
	return u.DisplayName
}

func (u *User) WebAuthnCredentials() []webauthn.Credential {
	return u.Credentials
}

func (u *User) WebAuthnIcon() string {
	return ""
}

// グローバル変数
var (
	webAuthn      *webauthn.WebAuthn
	userDB        = make(map[string]*User)             // メモリ内のユーザーデータベース
	allowedEmails []string                             // 許可されたメールアドレスのリスト
	sessionStore  = map[string]*webauthn.SessionData{} // セッションデータを保存
)

func main() {
	// .envファイルからメールアドレスを読み込む
	loadAllowedEmails()

	// WebAuthn の初期化
	var err error
	webAuthn, err = webauthn.New(&webauthn.Config{
		RPDisplayName: "Passkey Sample App",    // Relying Party（サービス提供者）の表示名
		RPID:          "localhost",             // RPのドメイン
		RPOrigin:      "http://localhost:3000", // RPの起点（オリジン）URL
	})

	if err != nil {
		log.Fatal("WebAuthn initialization error:", err)
	}

	// Ginルーターの初期化
	r := gin.Default()

	// CORS設定
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"}, // Acceptヘッダーを追加
		AllowCredentials: true,
	}))

	// APIルートグループ
	api := r.Group("/api")
	{
		// 登録関連エンドポイント
		api.POST("/register/options", getRegistrationOptions)
		api.POST("/register/verify", verifyRegistration)

		// ログイン関連エンドポイント
		api.POST("/login/options", getAuthenticationOptions)
		api.POST("/login/verify", verifyAuthentication)
	}

	// サーバー起動
	fmt.Println("サーバーを起動しました。http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Server start error:", err)
	}
}

// .envファイルからメールアドレスを読み込む関数
func loadAllowedEmails() {
	file, err := os.Open("../.env") // project root
	if err != nil {
		// .envファイルが存在しない場合は警告のみ表示して続行
		if os.IsNotExist(err) {
			log.Println("警告: .envファイルが見つかりません。")
			return
		}
		log.Fatalf(".envファイルを開けませんでした: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		email := strings.TrimSpace(scanner.Text())
		if email != "" && strings.Contains(email, "@") {
			allowedEmails = append(allowedEmails, email)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("ファイル読み込みエラー: %v", err)
	}

	if len(allowedEmails) > 0 {
		log.Printf("許可されたメールアドレス: %v", allowedEmails)
	} else {
		log.Println("警告: .envファイルに許可されたメールアドレスが設定されていません。")
	}
}

// メールアドレスが許可リストにあるか確認
func isEmailAllowed(email string) bool {
	if len(allowedEmails) == 0 {
		// .envファイルが空または存在しない場合は、便宜上すべてのメールを許可（開発用）
		log.Println("警告: .envにメールアドレスが設定されていないため、すべてのメールアドレスを許可します。")
		return true
	}
	for _, allowed := range allowedEmails {
		if allowed == email {
			return true
		}
	}
	return false
}

// 登録オプションを取得するエンドポイント
func getRegistrationOptions(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	// fmt.Println("DEBUG: 登録オプション処理を開始")

	if err := c.ShouldBindJSON(&req); err != nil {
		// fmt.Printf("DEBUG: JSONバインドエラー: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "無効なリクエスト: " + err.Error()})
		return
	}

	// fmt.Printf("DEBUG: リクエスト内容: %+v\n", req)

	// メールアドレスが許可リストにあるか確認
	if !isEmailAllowed(req.Email) {
		// fmt.Printf("DEBUG: 許可されていないメールアドレス: %s\n", req.Email)
		c.JSON(http.StatusForbidden, gin.H{"message": "このメールアドレスは登録が許可されていません。"})
		return
	}

	// ユーザーの存在確認
	user, exists := userDB[req.Email]
	if !exists {
		// 新規ユーザー作成
		user = &User{
			ID:          req.Email, // IDはEmailと同じにする
			Name:        req.Email, // NameもEmailと同じにする
			DisplayName: req.Email,
			Credentials: []webauthn.Credential{},
		}
		userDB[req.Email] = user
		// fmt.Printf("DEBUG: 新規ユーザー作成: %s\n", req.Email)
	} else {
		// fmt.Printf("DEBUG: 既存ユーザー: %s\n", req.Email)
	}

	// 登録オプションを生成
	options, sessionData, err := webAuthn.BeginRegistration(user)
	if err != nil {
		// fmt.Printf("DEBUG: 登録オプション生成エラー: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "登録オプションの生成に失敗しました: " + err.Error()})
		return
	}

	// セッションデータをユーザーに関連付けて一時的に保存
	user.RegistrationSessionData = sessionData
	// fmt.Printf("DEBUG: 登録セッション開始: User=%s\n", user.ID)
	// fmt.Printf("DEBUG: Challenge: %x\n", options.Response.Challenge)
	// fmt.Printf("DEBUG: User.ID型: %T, 値: %v\n", options.Response.User.ID, options.Response.User.ID)

	// レスポンスを作成
	publicKeyMap := gin.H{
		"challenge": base64.RawURLEncoding.EncodeToString(options.Response.Challenge),
		"rp": gin.H{
			"name": options.Response.RelyingParty.Name,
			"id":   options.Response.RelyingParty.ID,
		},
		"user": gin.H{
			"id":          base64.RawURLEncoding.EncodeToString([]byte(user.ID)), // ユーザーIDをバイト配列に変換してbase64エンコード
			"name":        user.Name,                                             // オリジナルのUser構造体のフィールドを使用
			"displayName": user.DisplayName,                                      // オリジナルのUser構造体のフィールドを使用（空値を避ける）
		},
		"pubKeyCredParams":       options.Response.Parameters,
		"timeout":                options.Response.Timeout,
		"authenticatorSelection": options.Response.AuthenticatorSelection,
	}

	// Attestation Conveyance Preference を処理
	attestationPref := string(options.Response.Attestation) // stringに変換
	switch attestationPref {
	case "direct", "indirect", "enterprise":
		publicKeyMap["attestation"] = attestationPref // 有効な場合はそのまま使用
	case "none", "": // "none" または空文字列の場合
		publicKeyMap["attestation"] = "none"
	default:
		// 不明な値の場合も "none" として扱う (より安全)
		// fmt.Printf("DEBUG: 不明なAttestation Preference: '%s'、'none'として処理\n", attestationPref)
		publicKeyMap["attestation"] = "none"
	}

	responseJson := publicKeyMap

	// fmt.Printf("DEBUG: レスポンスJSON: %+v\n", responseJson)
	c.JSON(http.StatusOK, responseJson)
}

// 登録検証のエンドポイント
func verifyRegistration(c *gin.Context) {
	// fmt.Println("DEBUG: 登録検証処理を開始")

	// リクエストボディを出力
	rawData, _ := c.GetRawData()
	// fmt.Printf("DEBUG: 登録検証リクエストボディ: %s\n", string(rawData))

	// リクエストをバインド
	c.Request.Body = io.NopCloser(bytes.NewBuffer(rawData)) // ボディを再設定

	// フロントエンドから送られてくるデータ構造をパース
	var req struct {
		Email               string                 `json:"email"`
		AttestationResponse map[string]interface{} `json:"attestationResponse"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		// fmt.Printf("DEBUG: リクエストJSONパースエラー: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "リクエスト形式が不正です: " + err.Error()})
		return
	}

	// fmt.Printf("DEBUG: パースしたリクエスト: %+v\n", req)

	// ユーザーを検索
	user, exists := userDB[req.Email]
	if !exists || user.RegistrationSessionData == nil {
		// fmt.Printf("DEBUG: ユーザーまたはセッションが見つかりません: email=%s\n", req.Email)
		c.JSON(http.StatusBadRequest, gin.H{"message": "ユーザーまたは登録セッションが見つかりません。"})
		return
	}

	foundSession := user.RegistrationSessionData
	// fmt.Printf("DEBUG: 検証用セッションデータ: %+v\n", foundSession)

	// フロントエンドから送られてきたattestationResponseを直接WebAuthnライブラリに渡せる形に変換
	// WebAuthnライブラリが期待する構造に変換
	attestationResp := map[string]interface{}{
		"id":       req.AttestationResponse["id"],
		"rawId":    req.AttestationResponse["rawId"],
		"type":     req.AttestationResponse["type"],
		"response": req.AttestationResponse["response"],
	}

	// JSONに変換してhttpリクエストのボディに設定
	jsonBytes, err := json.Marshal(attestationResp)
	if err != nil {
		// fmt.Printf("DEBUG: attestationデータのJSON化エラー: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "内部サーバーエラー"})
		return
	}

	// 新しいリクエストを作成
	newReq, err := http.NewRequest("POST", c.Request.URL.String(), bytes.NewBuffer(jsonBytes))
	if err != nil {
		// fmt.Printf("DEBUG: リクエスト作成エラー: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "内部サーバーエラー"})
		return
	}
	newReq.Header.Set("Content-Type", "application/json")

	// fmt.Printf("DEBUG: 変換後のリクエストボディ: %s\n", string(jsonBytes))

	// リクエストをそのままFinishRegistrationに渡す
	credential, err := webAuthn.FinishRegistration(user, *foundSession, newReq)
	if err != nil {
		// fmt.Printf("DEBUG: 登録検証失敗: User=%s, Error=%v\n", user.ID, err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "登録検証に失敗しました: " + err.Error()})
		return
	}

	// ユーザーのクレデンシャルを保存
	user.Credentials = append(user.Credentials, *credential)
	user.RegistrationSessionData = nil // セッションデータをクリア
	// fmt.Printf("DEBUG: Passkey登録成功: User=%s, CredentialID=%x\n", user.ID, credential.ID)
	// fmt.Printf("DEBUG: 登録完了後のユーザーDB状態: %+v\n", userDB)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Passkey登録が完了しました。",
	})
}

// 認証オプションを取得するエンドポイント
func getAuthenticationOptions(c *gin.Context) {
	// fmt.Println("DEBUG: 認証オプション処理を開始")
	// 登録されたユーザーがいない場合
	if len(userDB) == 0 {
		// fmt.Println("DEBUG: ユーザーDBが空です")
		c.JSON(http.StatusBadRequest, gin.H{"message": "登録されたユーザーがいません。"})
		return
	}

	// fmt.Printf("DEBUG: 現在のユーザーDB: %+v\n", userDB)

	// ユーザーDBから最初のユーザーを取得
	var firstUser *User
	for _, u := range userDB {
		firstUser = u
		break // 最初のユーザーを取得したらループを抜ける
	}

	// fmt.Printf("DEBUG: ユーザー %s を使用して認証オプションを生成\n", firstUser.ID)

	// BeginLoginに最初のユーザーを渡す
	options, sessionData, err := webAuthn.BeginLogin(firstUser)
	if err != nil {
		// fmt.Printf("DEBUG: 認証オプション生成エラー: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "認証オプションの生成に失敗しました: " + err.Error()})
		return
	}

	// グローバルなセッションストアに認証セッションデータを保存
	sessionStore["auth"] = sessionData
	// fmt.Println("DEBUG: 認証セッション開始")
	// fmt.Printf("DEBUG: Challenge: %x\n", options.Response.Challenge)

	// レスポンスを作成
	publicKeyMap := gin.H{
		"challenge":        base64.RawURLEncoding.EncodeToString(options.Response.Challenge),
		"timeout":          options.Response.Timeout,
		"rpId":             options.Response.RelyingPartyID,
		"allowCredentials": []gin.H{}, // 空の配列を設定
		"userVerification": options.Response.UserVerification,
	}

	responseJson := publicKeyMap

	// fmt.Printf("DEBUG: 認証レスポンスJSON: %+v\n", responseJson)
	c.JSON(http.StatusOK, responseJson)
}

// 認証検証のエンドポイント
func verifyAuthentication(c *gin.Context) {
	// fmt.Println("DEBUG: 認証検証処理を開始")

	// グローバルなセッションデータを取得
	session, exists := sessionStore["auth"]
	if !exists {
		// fmt.Println("DEBUG: 認証セッションが見つかりません")
		c.JSON(http.StatusBadRequest, gin.H{"message": "認証セッションが見つかりません。"})
		return
	}

	// fmt.Printf("DEBUG: 認証セッションデータ: %+v\n", session)

	// リクエストボディをログに出力
	rawData, _ := c.GetRawData()
	// fmt.Printf("DEBUG: 認証検証リクエストボディ: %s\n", string(rawData))
	c.Request.Body = io.NopCloser(bytes.NewBuffer(rawData)) // ボディを再設定

	// フロントエンドから送られてくるデータ構造をパース
	var req struct {
		AssertionResponse map[string]interface{} `json:"assertionResponse"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		// fmt.Printf("DEBUG: リクエストJSONパースエラー: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "リクエスト形式が不正です: " + err.Error()})
		return
	}

	// fmt.Printf("DEBUG: パースしたリクエスト: %+v\n", req)

	// 認証応答データを抽出
	authData := req.AssertionResponse

	// エラーチェック: リクエストが正しい形式であることを確認
	if authData == nil || authData["id"] == nil || authData["response"] == nil {
		// fmt.Println("DEBUG: 認証データの形式が無効です: 必須フィールドがありません")
		c.JSON(http.StatusBadRequest, gin.H{"message": "認証データの形式が無効です: 必須フィールドがありません"})
		return
	}

	// authDataを直接使用してJSONにシリアライズし、新しいリクエストにする
	jsonBytes, err := json.Marshal(authData)
	if err != nil {
		// fmt.Printf("DEBUG: assertionデータのJSON化エラー: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "内部サーバーエラー"})
		return
	}

	// fmt.Printf("DEBUG: 変換後のリクエストデータ: %s\n", string(jsonBytes))

	// 新しいリクエストを作成
	newReq, err := http.NewRequest("POST", c.Request.URL.String(), bytes.NewBuffer(jsonBytes))
	if err != nil {
		// fmt.Printf("DEBUG: リクエスト作成エラー: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "内部サーバーエラー"})
		return
	}
	newReq.Header.Set("Content-Type", "application/json")

	// IDを直接base64デコードしてユーザー検索に使用
	credentialID, ok := authData["id"].(string)
	if !ok {
		// fmt.Printf("DEBUG: credential ID の取得に失敗\n")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Credential IDの形式が無効です"})
		return
	}

	credentialIDBytes, err := base64.RawURLEncoding.DecodeString(credentialID)
	if err != nil {
		// fmt.Printf("DEBUG: 認証検証エラー: Credential ID のデコードに失敗: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Credential IDの形式が無効です"})
		return
	}

	// ユーザー検索関数を定義
	userHandler := func(rawID, userHandle []byte) (webauthn.User, error) {
		// fmt.Printf("DEBUG: ValidateLogin: ユーザー検索 rawID=%x\n", rawID)
		var foundUser *User
		// まず rawID (Credential ID) で検索
		for _, u := range userDB {
			for _, cred := range u.Credentials {
				// cred.ID と rawID はどちらも []byte なので直接比較
				if bytes.Equal(cred.ID, rawID) {
					foundUser = u
					// fmt.Printf("DEBUG: ValidateLogin: rawID でユーザー %s を発見\n", foundUser.ID)
					return foundUser, nil // 発見したら返す
				}
			}
		}

		// userHandleが存在し、空でない場合はそれで検索
		if len(userHandle) > 0 {
			// userStr := string(userHandle)
			// fmt.Printf("DEBUG: userHandle文字列: %s\n", userStr)
			u, exists := userDB[string(userHandle)]
			if exists {
				foundUser = u
				// fmt.Printf("DEBUG: ValidateLogin: userHandle でユーザー %s を発見\n", foundUser.ID)
				return foundUser, nil // 発見したら返す
			}
		}

		// ユーザーDBが空でなければ、最初のユーザーを返す（デバッグ用）
		if len(userDB) > 0 {
			for _, u := range userDB {
				foundUser = u
				// fmt.Printf("DEBUG: ValidateLogin: 最初のユーザー %s を返します\n", foundUser.ID)
				return foundUser, nil
			}
		}

		// fmt.Println("DEBUG: ValidateLogin: ユーザーが見つかりません")
		return nil, fmt.Errorf("ユーザーが見つかりません")
	}

	// userHandleの抽出を試みる
	var userHandle []byte
	if respMap, ok := authData["response"].(map[string]interface{}); ok {
		if uh, ok := respMap["userHandle"].(string); ok && uh != "" {
			userHandle = []byte(uh)
		}
	}

	// ユーザーを特定
	potentialUser, err := userHandler(credentialIDBytes, userHandle)
	if err != nil {
		// fmt.Printf("DEBUG: 認証検証前ユーザー検索エラー: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "登録されたユーザーがいません。"})
		return
	}

	// キャストして具体的なユーザー情報を取得
	userToValidate, ok := potentialUser.(*User)
	if !ok || userToValidate == nil {
		// fmt.Printf("DEBUG: 認証検証前ユーザーキャスト失敗\n")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "ユーザー情報の取得に失敗しました"})
		return
	}

	// go-webauthnライブラリの検証関数を使用
	credential, err := webAuthn.FinishLogin(userToValidate, *session, newReq)
	if err != nil {
		// fmt.Printf("DEBUG: 認証検証失敗(FinishLogin): Error=%v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "認証検証に失敗しました: " + err.Error()})
		return
	}
	_ = credential // credential変数は現時点では使わない

	// 検証成功。検証に使用したユーザー (userToValidate) が認証されたユーザーとなる

	// グローバルなセッションデータをクリア
	delete(sessionStore, "auth")
	// fmt.Printf("DEBUG: 認証成功: User=%s\n", userToValidate.ID) // userToValidate を使用

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"message":     "認証に成功しました。",
		"email":       userToValidate.ID,          // userToValidate を使用
		"displayName": userToValidate.DisplayName, // userToValidate を使用
	})
}
