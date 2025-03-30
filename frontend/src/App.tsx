import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import LoginPage from './components/LoginPage';
import WelcomePage from './components/WelcomePage';
import styles from './App.module.css';

// 型定義
export interface User {
    email: string;
    displayName?: string;
}

const App: React.FC = () => {
    const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);
    const [user, setUser] = useState<User | null>(null);

    // ローカルストレージからログイン状態を復元
    useEffect(() => {
        const storedAuth = localStorage.getItem('isAuthenticated');
        const storedUser = localStorage.getItem('user');

        if (storedAuth === 'true' && storedUser) {
            setIsAuthenticated(true);
            setUser(JSON.parse(storedUser));
        }
    }, []);

    // ログイン処理
    const handleLogin = (userData: User): void => {
        setIsAuthenticated(true);
        setUser(userData);
        localStorage.setItem('isAuthenticated', 'true');
        localStorage.setItem('user', JSON.stringify(userData));
    };

    // ログアウト処理
    const handleLogout = (): void => {
        setIsAuthenticated(false);
        setUser(null);
        localStorage.removeItem('isAuthenticated');
        localStorage.removeItem('user');
    };

    return (
        <div className={styles.app}>
            <Router>
                <Routes>
                    <Route
                        path="/"
                        element={
                            isAuthenticated ?
                                <Navigate to="/welcome" replace /> :
                                <LoginPage onLogin={handleLogin} />
                        }
                    />
                    <Route
                        path="/welcome"
                        element={
                            isAuthenticated ?
                                <WelcomePage user={user} onLogout={handleLogout} /> :
                                <Navigate to="/" replace />
                        }
                    />
                </Routes>
            </Router>
        </div>
    );
};

export default App; 