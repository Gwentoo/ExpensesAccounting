// src/pages/RestorePassScreen.tsx
import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import '../css/RestorePassScreen.css';
import {setNewPass, validateResetToken} from "../services/api";

export default function RestorePassScreen() {
    const { token } = useParams<{ token: string }>();
    const navigate = useNavigate();

    const [newPassword, setNewPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [loading, setLoading] = useState(false);
    const [message, setMessage] = useState('');
    const [isValidToken, setIsValidToken] = useState(true);
    const [userEmail, setUserEmail] = useState('');

    useEffect(() => {
        if (token) {
            validateToken(token);
        } else {
            setIsValidToken(false);
            setMessage('Неверная ссылка для восстановления пароля');
        }
    }, [token]);

    const validateToken = async (token: string) => {
        try {
            const res = await validateResetToken(token);
            if (res.success) {
                setUserEmail(res.email);
                setIsValidToken(true);
            } else {
                setIsValidToken(false);
            }
        } catch (error) {
            setIsValidToken(false);
            setMessage('Ошибка проверки ссылки');
        }
    };

    const handleResetPassword = async (e: React.FormEvent) => {
        e.preventDefault();

        if (!token) {
            setMessage('Неверная ссылка для восстановления пароля');
            return;
        }

        if (newPassword !== confirmPassword) {
            setMessage('Пароли не совпадают');
            return;
        }

        if (newPassword.length < 6) {
            setMessage('Пароль должен содержать минимум 6 символов');
            return;
        }

        setLoading(true);
        setMessage('');

        try {

            console.log(token);
            console.log(newPassword);
            const res = await setNewPass(token, newPassword)

            if (res.success) {
                setMessage('Пароль успешно изменен!');
                setTimeout(() => navigate('/login'));
            }
        } catch (error) {
            setMessage('Произошла ошибка при подключении к серверу');
        } finally {
            setLoading(false);
        }
    };

    if (!isValidToken) {
        return (
            <div className="restore-container">
                <div className="restore-form error-message">
                    <h2>Ошибка</h2>
                    <p>{message}</p>
                    <button
                        onClick={() => navigate('/newpass')}
                        className="submit-btn"
                    >
                        Запросить новую ссылку
                    </button>
                </div>
            </div>
        );
    }

    return (
        <div className="restore-container">
            <form onSubmit={handleResetPassword} className="restore-form">
                <h2>Восстановление пароля</h2>
                {userEmail && <p className="user-email">Для: {userEmail}</p>}

                <div className="field">
                    <label className="label">Новый пароль</label>
                    <input
                        type="password"
                        value={newPassword}
                        onChange={(e) => setNewPassword(e.target.value)}
                        placeholder="Введите новый пароль"
                        className="input"
                        required
                    />
                </div>

                <div className="field">
                    <label className="label">Подтвердите пароль</label>
                    <input
                        type="password"
                        value={confirmPassword}
                        onChange={(e) => setConfirmPassword(e.target.value)}
                        placeholder="Повторите новый пароль"
                        className="input"
                        required
                    />
                </div>

                <button
                    type="submit"
                    className={`submit-btn ${loading ? 'submit-disabled' : ''}`}
                    disabled={loading}
                >
                    {loading ? 'Изменение...' : 'Изменить пароль'}
                </button>

                {message && (
                    <div className={`message ${message.includes('успешно') ? 'success' : 'error'}`}>
                        {message}
                    </div>
                )}
            </form>
        </div>
    );
}