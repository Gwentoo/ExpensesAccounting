
import React, { useState } from "react";
import { useNavigate, useLocation } from "react-router-dom";
import OtpInput from "../components/OtpInput";
import "../css/VerifyScreen.css"
import { verify } from "../services/api";

export default function VerifyScreen() {
    const navigate = useNavigate();
    const location = useLocation();
    const email = location.state?.email;

    const [code, setCode] = useState(["", "", "", "", ""]);
    const [loading, setLoading] = useState(false);

    const onVerify = async () => {
        const joined = code.join("");
        if (joined.length < 5) {
            alert("Введите 5 цифр");
            return;
        }
        try {
            setLoading(true);
            const res = await verify(email, joined);
            if (!res.success) throw new Error("Неверный код");
            localStorage.setItem("user_token", res.token);
            navigate("/dashboard");
        } catch (err: any) {
            alert(err.message);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="verify-container">
            <div className="verify-form">
                <h1 className="verify-title">Подтвердите регистрацию</h1>
                <p className="verify-subtitle">Код отправлен на {email}</p>
                <OtpInput value={code} setValue={setCode} />
                <button
                    onClick={onVerify}
                    className={`verify-btn ${loading ? 'verify-btn-disabled' : ''}`}
                    disabled={loading}
                >
                    {loading ? "Проверка..." : "Подтвердить"}
                </button>
            </div>
        </div>
    );
}