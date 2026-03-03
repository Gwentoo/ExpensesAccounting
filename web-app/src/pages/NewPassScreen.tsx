import React, {useState} from "react";
import "../css/NewPassScreen.css";
import {Link, useNavigate} from 'react-router-dom';
import {newPass} from "../services/api";

export default function NewPassScreen() {
    const [errors, setErrors] = useState<{ email?: string }>({});
    const [email, setEmail] = useState('');
    const [loading, setLoading] = useState(false);
    const [isSent, setIsSent] = useState(false);
    const [apiError, setApiError] = useState('');

    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

    const isFormValid = () => {
        return emailRegex.test(email.trim());
    };

    const validateForm = () => {
        const newErrors: { email?: string } = {};

        if (!email.trim()) {
            newErrors.email = 'Email обязателен';
        } else if (!emailRegex.test(email.trim())) {
            newErrors.email = 'Неверный формат email';
        }

        setErrors(newErrors);
        setApiError('');
        return Object.keys(newErrors).length === 0;
    };

    const onSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        if (!validateForm()) return;

        setLoading(true);
        setApiError('');

        try {

            const res = await newPass(email);

            if (!res.success) {

                if (res.message?.includes('не существует') || res.message?.includes('не найден')) {
                    setApiError('Пользователя с такой почтой не существует');
                } else {
                    setApiError(res.message || 'Ошибка отправки');
                }
                return;
            }

            setIsSent(true);

        } catch (err: any) {
            if (err.message?.includes('не найден')) {
                setApiError('Пользователя с такой почтой не существует');
            } else {
                setApiError(err?.message || 'Произошла ошибка при отправке');
            }
        } finally {
            setLoading(false);
        }
    };

    const handleSendAgain = () => {
        setIsSent(false);
        setApiError('');
        setEmail('');
    };

    return (
        <div className="newpass-container">
            <div>
                <form onSubmit={onSubmit} className="newpass-form">

                    {isSent ? (
                        <div className="success-message">
                            <h2>Ссылка отправлена</h2>
                            <p>Проверьте вашу почту <strong>{email}</strong> и перейдите по ссылке для восстановления пароля</p>
                            <button
                                type="button"
                                onClick={handleSendAgain}
                                className="submit-btn"
                            >
                                Отправить ещё раз
                            </button>
                        </div>
                    ) : (

                        <>
                            <h1 className="title">Восстановление пароля</h1>

                            <div className="field">
                                <label className="label">Email</label>
                                <input
                                    value={email}
                                    onChange={(e) => {
                                        setEmail(e.target.value);
                                        setApiError('');
                                    }}
                                    type="email"
                                    autoCapitalize="none"
                                    placeholder="you@example.com"
                                    className={`input ${errors.email || apiError ? 'input-error' : ''}`}
                                />
                                {errors.email && <span className="error-text">{errors.email}</span>}
                            </div>


                            {apiError && (
                                <div className="api-error">
                                    <span className="error-icon">⚠</span>
                                    <span className="error-text">{apiError}</span>
                                </div>
                            )}

                            <button
                                type="submit"
                                className={`submit-btn ${!isFormValid() || loading ? 'submit-disabled' : ''}`}
                                disabled={!isFormValid() || loading}
                            >
                                {loading ? "Отправка..." : "Отправить ссылку"}
                            </button>
                        </>
                    )}

                    <div className="footer">
                        <Link to="/register" className="link">Вернуться к регистрации</Link>
                    </div>

                </form>
            </div>
        </div>
    );
}