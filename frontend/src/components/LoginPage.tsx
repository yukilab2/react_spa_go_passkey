import React, { useState } from 'react';
import * as Toast from '@radix-ui/react-toast';
import * as Form from '@radix-ui/react-form';
import { startRegistration, startAuthentication } from '@simplewebauthn/browser';
import authService from '../services/authService';
import { User } from '../App';
import styles from './LoginPage.module.css';

interface LoginPageProps {
    onLogin: (userData: User) => void;
}

const LoginPage: React.FC<LoginPageProps> = ({ onLogin }) => {
    const [email, setEmail] = useState<string>('someone@example.com');
    const [message, setMessage] = useState<string>('');
    const [error, setError] = useState<string>('');
    const [showToast, setShowToast] = useState<boolean>(false);
    const [loading, setLoading] = useState<boolean>(false);

    // Passkey登録処理
    const handleRegister = async (e: React.FormEvent): Promise<void> => {
        e.preventDefault();
        if (!email || !validateEmail(email)) {
            setError('有効なメールアドレスを入力してください。');
            return;
        }

        setLoading(true);
        setError('');

        try {
            // バックエンドから登録オプションを取得
            const registrationOptions = await authService.getRegistrationOptions(email);
            // ブラウザのPasskey登録APIを呼び出し
            const attResp = await startRegistration(registrationOptions);

            // 検証のためにバックエンドに結果を送信
            const verificationResponse = await authService.verifyRegistration(email, attResp);

            if (verificationResponse.success) {
                setMessage('Passkey登録が完了しました。登録したPasskeyでログインできます。');
                setShowToast(true);
                setEmail('');
            } else {
                setError('Passkey登録に失敗しました。');
            }
        } catch (err: any) {
            console.error('登録エラー:', err);
            setError(err.message || 'Passkey登録中にエラーが発生しました。');
        } finally {
            setLoading(false);
        }
    };

    // Passkey認証処理
    const handleLogin = async (): Promise<void> => {
        setLoading(true);
        setError('');

        try {
            // バックエンドから認証オプションを取得
            const authOptions = await authService.getAuthenticationOptions();

            // ブラウザのPasskey認証APIを呼び出し
            const authResp = await startAuthentication(authOptions);

            // 検証のためにバックエンドに結果を送信
            const verificationResponse = await authService.verifyAuthentication(authResp);

            if (verificationResponse.success) {
                // 認証成功後の処理
                onLogin({
                    email: verificationResponse.email || '',
                    displayName: verificationResponse.displayName || verificationResponse.email
                });
            } else {
                setError('認証に失敗しました。');
            }
        } catch (err: any) {
            console.error('認証エラー:', err);
            setError(err.message || '認証中にエラーが発生しました。');
        } finally {
            setLoading(false);
        }
    };

    const validateEmail = (email: string): boolean => {
        const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        return re.test(email);
    };

    return (
        <Toast.Provider swipeDirection="right">
            <div className={styles.container}>
                <div className={styles.card}>
                    <div className={styles.header}>
                        <h1 className={styles.title}>Passkey認証サンプル</h1>
                    </div>

                    {error && (
                        <div className={styles.errorAlert}>
                            {error}
                        </div>
                    )}

                    <Form.Root className={styles.form} onSubmit={handleRegister}>
                        <h2 className={styles.sectionTitle}>新規登録</h2>

                        <Form.Field className={styles.inputGroup} name="email">
                            <Form.Label className={styles.label}>メールアドレス</Form.Label>
                            <Form.Control asChild>
                                <input
                                    className={styles.input}
                                    type="email"
                                    required
                                    placeholder="example@example.com"
                                    value={email}
                                    onChange={(e) => setEmail(e.target.value)}
                                />
                            </Form.Control>
                        </Form.Field>

                        <Form.Submit asChild>
                            <button
                                className={styles.buttonPrimary}
                                disabled={loading}
                            >
                                Passkey登録
                            </button>
                        </Form.Submit>
                    </Form.Root>

                    <div className={styles.divider}>または</div>

                    <div>
                        <h2 className={styles.sectionTitle}>ログイン</h2>

                        <button
                            className={styles.buttonOutline}
                            onClick={handleLogin}
                            disabled={loading}
                        >
                            Passkeyでログイン
                        </button>
                    </div>
                </div>
            </div>

            <Toast.Root
                className="ToastRoot"
                open={showToast}
                onOpenChange={setShowToast}
            >
                <Toast.Title>{message}</Toast.Title>
                <Toast.Action asChild altText="閉じる">
                    <button onClick={() => setShowToast(false)}>閉じる</button>
                </Toast.Action>
            </Toast.Root>

            <Toast.Viewport className="ToastViewport" />
        </Toast.Provider>
    );
};

export default LoginPage; 