import React from 'react';
import * as Avatar from '@radix-ui/react-avatar';
import { User } from '../App';
import styles from './WelcomePage.module.css';

interface WelcomePageProps {
    user: User | null;
    onLogout: () => void;
}

const WelcomePage: React.FC<WelcomePageProps> = ({ user, onLogout }) => {
    return (
        <div className={styles.container}>
            <div className={styles.card}>
                <Avatar.Root className={styles.avatar}>
                    <Avatar.Fallback>
                        {user?.displayName?.charAt(0) || user?.email?.charAt(0) || '?'}
                    </Avatar.Fallback>
                </Avatar.Root>

                <h1 className={styles.title}>ようこそ！</h1>

                <p className={styles.text}>
                    {user?.displayName || user?.email || 'ユーザー'} さん、正常にログインしました。
                </p>

                <p className={styles.emailText}>
                    メールアドレス: {user?.email || 'N/A'}
                </p>

                <button
                    className={styles.logoutButton}
                    onClick={onLogout}
                >
                    ログアウト
                </button>
            </div>
        </div>
    );
};

export default WelcomePage; 