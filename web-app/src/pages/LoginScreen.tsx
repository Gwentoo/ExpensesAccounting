import React, { useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import { login } from "../services/api";
import '../css/LoginScreen.css';

const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

export default function LoginScreen() {
    const navigate = useNavigate();
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [showPassword, setShowPassword] = useState(false);
    const [loading, setLoading] = useState(false);
    const [errors, setErrors] = useState<{ email?: string; password?: string }>({});

    const validate = () => {
        const e: typeof errors = {};
        if (!email.trim()) e.email = "Email обязателен";
        else if (!emailRegex.test(email.trim())) e.email = "Неверный формат email";

        if (!password) e.password = "Пароль обязателен";
        else if (password.length < 6) e.password = "Минимум 6 символов";

        setErrors(e);
        return Object.keys(e).length === 0;
    };

    const onSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!validate()) return;
        setLoading(true);
        try {
            const res = await login(email, password);
            if (res.token === "") {
                throw new Error(res.message);
            }

            localStorage.setItem("user_token", res.token);
            navigate("/dashboard");

        } catch (err: any) {
            alert(err?.message || "Не удалось войти");
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="login-container">
            <form onSubmit={onSubmit} className="login-form">
                <h1 className="login-title">Вход</h1>

                <div className="field">
                    <label className="label">Email</label>
                    <input
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                        type="email"
                        placeholder="you@example.com"
                        className={`input ${errors.email ? 'input-error' : ''}`}
                    />
                    {errors.email && <span className="error-text">{errors.email}</span>}
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
                            {showPassword ? "Скрыть" : "Показать"}
                        </button>
                    </div>
                    {errors.password && <span className="error-text">{errors.password}</span>}
                </div>

                <button
                    type="submit"
                    className={`submit-btn ${loading ? 'submit-disabled' : ''}`}
                    disabled={loading}
                >
                    {loading ? "Загрузка..." : "Войти"}
                </button>

                <div className="footer">
                    <span className="small">Нет аккаунта?</span>
                    <Link to="/register" className="link">Зарегистрироваться</Link>
                </div>
                <div className="footer">
                    <Link to="/newpass" className="link">Забыли пароль?</Link>
                </div>
            </form>
        </div>
    );
}