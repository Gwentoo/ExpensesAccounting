import React, {useState} from 'react';
import {Link, useNavigate} from 'react-router-dom';
import {register} from '../services/api';
import '../css/RegisterScreen.css';

const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

export default function RegisterScreen() {
    const navigate = useNavigate();
    const [email, setEmail] = useState('');
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [showPassword, setShowPassword] = useState(false);
    const [loading, setLoading] = useState(false);
    const [errors, setErrors] = useState<{ email?: string; username?: string; password?: string }>({});


    const onSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        setLoading(true);
        try {
            const res = await register(email, username, password);
            if (!res.success) {
                new Error(res.message || "Регистрация не удалась");
            }
            alert("Успех! Теперь подтвердите email");
            navigate("/verify", {state: {email}});
        } catch (err: any) {
            alert(err?.message || "Ошибка регистрации");
        } finally {
            setLoading(false);
        }
    };

    const isFormValid = () => {
        return (
            emailRegex.test(email.trim()) &&
            username.trim().length >= 3 &&
            password.length >= 6 &&
            password.length < 50
        );
    };

    return (
        <div className="register-container">
            <form onSubmit={onSubmit} className="register-form">
                <h1 className="register-title">Регистрация</h1>

                <div className="field">
                    <label className="label">Email</label>
                    <input
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                        type="email"
                        autoCapitalize="none"
                        placeholder="you@example.com"
                        className={`input ${errors.email ? 'input-error' : ''}`}
                    />
                    {errors.email && <span className="error-text">{errors.email}</span>}
                </div>

                <div className="field">
                    <label className="label">Имя пользователя</label>
                    <input
                        value={username}
                        onChange={(e) => setUsername(e.target.value)}
                        type="text"
                        autoCapitalize="none"
                        placeholder="username"
                        className={`input ${errors.username ? 'input-error' : ''}`}
                    />
                    {errors.username && <span className="error-text">{errors.username}</span>}
                </div>

                <div className="field">
                    <label className="label">Пароль</label>
                    <div className="password-row">
                        <input
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            type={showPassword ? "text" : "password"}
                            placeholder="Пароль"
                            className={`input flex-input ${errors.password ? 'input-error' : ''}`}
                        />
                        <button
                            type="button"
                            onClick={() => setShowPassword((s) => !s)}
                            className="show-btn"
                        >
                            {showPassword ? 'Скрыть' : 'Показать'}
                        </button>
                    </div>
                    {errors.password && <span className="error-text">{errors.password}</span>}
                </div>

                <button
                    type="submit"
                    className={`submit-btn ${!isFormValid() || loading ? 'submit-disabled' : ''}`}
                    disabled={!isFormValid() || loading}
                >
                    {loading ? "Регистрация..." : "Зарегистрироваться"}
                </button>

                <div className="footer">
                    <span className="small">Уже есть аккаунт?</span>
                    <Link to="/login" className="link">Войти</Link>
                </div>
                <div className="footer">
                    <Link to="/newpass" className="link">Забыли пароль?</Link>
                </div>
            </form>
        </div>
    );
}