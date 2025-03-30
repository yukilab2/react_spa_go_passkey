import axios, { AxiosInstance, AxiosResponse } from 'axios';

import {
    PublicKeyCredentialRequestOptionsJSON,
    PublicKeyCredentialCreationOptionsJSON,
    RegistrationResponseJSON,
    AuthenticationResponseJSON
} from '@simplewebauthn/typescript-types';

// APIのベースURL
const API_URL = 'http://localhost:8080/api';

// axios インスタンスの作成
const apiClient: AxiosInstance = axios.create({
    baseURL: API_URL,
    headers: {
        'Content-Type': 'application/json'
    },
    withCredentials: true // CSRF保護のためのクッキーを送信
});

// 認証サービスのレスポンス型
interface VerificationResponse {
    success: boolean;
    message?: string;
    email?: string;
    displayName?: string;
}

const authService = {
    // Passkey登録オプションを取得
    async getRegistrationOptions(email: string): Promise<PublicKeyCredentialCreationOptionsJSON> {
        try {
            const response: AxiosResponse = await apiClient.post('/register/options', { email });
            return response.data;
        } catch (error: any) {
            // エラーメッセージの抽出
            const errorMessage = error.response?.data?.message || 'サーバーとの通信に失敗しました';
            throw new Error(errorMessage);
        }
    },

    // Passkey登録の検証
    async verifyRegistration(
        email: string,
        attestationResponse: RegistrationResponseJSON
    ): Promise<VerificationResponse> {
        try {
            const response: AxiosResponse = await apiClient.post('/register/verify', {
                email,
                attestationResponse
            });
            return response.data;
        } catch (error: any) {
            const errorMessage = error.response?.data?.message || 'Passkey登録の検証に失敗しました';
            throw new Error(errorMessage);
        }
    },

    // Passkey認証オプションを取得
    async getAuthenticationOptions(): Promise<PublicKeyCredentialRequestOptionsJSON> {
        try {
            const response: AxiosResponse = await apiClient.post('/login/options');
            return response.data;
        } catch (error: any) {
            const errorMessage = error.response?.data?.message || 'サーバーとの通信に失敗しました';
            throw new Error(errorMessage);
        }
    },

    // Passkey認証の検証
    async verifyAuthentication(
        assertionResponse: AuthenticationResponseJSON
    ): Promise<VerificationResponse> {
        try {
            const response: AxiosResponse = await apiClient.post('/login/verify', {
                assertionResponse
            });
            return response.data;
        } catch (error: any) {
            const errorMessage = error.response?.data?.message || '認証に失敗しました';
            throw new Error(errorMessage);
        }
    }
};

export default authService; 