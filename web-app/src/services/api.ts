import axios from 'axios';

const GW_PORT = process.env.REACT_APP_GATEWAY_PORT;

const API_URL = `http://localhost:${GW_PORT}`;



const api = axios.create({
    baseURL: API_URL,
});

api.interceptors.request.use((config) => {
    const token = localStorage.getItem('user_token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

export async function register(email: string, username: string, password: string) {
    const res = await fetch(`${API_URL}/auth/register`, {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify({
            email,
            username,
            password
        }),
    });
    if (!res.ok) throw new Error("Ошибка регистрации");
    return res.json();
}

export async function newPass(email: string) {
    const res = await fetch(`${API_URL}/auth/newpass`, {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify({
            email
        }),
    });

    if (!res.ok) throw new Error("Ошибка восстановления");
    return res.json();
}

export async function verify(email: string, code: string) {
    const res = await fetch(`${API_URL}/auth/verify`, {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify({
            email,
            code
        }),
    });
    if (!res.ok) throw new Error("Ошибка подтверждения");
    return res.json();
}

export async function login(email: string, password: string) {
    const res = await fetch(`${API_URL}/auth/login`, {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify({
            email,
            password
        }),
    });
    if (!res.ok) throw new Error("Не удалось войти");
    return res.json();
}

export async function validateResetToken(token: string){
    const res = await fetch(`${API_URL}/auth/validate-reset-token`, {
       method: "POST",
       headers: {"Content-Type": "application/json"},
       body: JSON.stringify({
          token
       }),
    });
    if (!res.ok) throw new Error("Не удалось проверить ссылку");
    return res.json();
}

export async function setNewPass(token: string, newPass: string) {
    const res = await fetch(`${API_URL}/auth/reset-password`, {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify({
            token,
            newPass,
        }),
    })
    if (!res.ok) throw new Error("Не удалось проверить ссылку");
    return res.json();
}

export async function getSupersetToken(startDate: string, endDate: string) {
    const token = localStorage.getItem("user_token");
    const res = await fetch(`${API_URL}/v1/expenses/superset-token`, {
        method: "POST",
        headers: {
            "Accept": "application/json",
            "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`
        },
        body: JSON.stringify({
            start_date: startDate,
            end_date: endDate
        }),
    });

    return await res.json();
}

export async function uploadStatement(file: File) {
    const token = localStorage.getItem("user_token");

    const arrayBuffer = await file.arrayBuffer();
    const uint8Array = new Uint8Array(arrayBuffer);
    let binary = '';
    for (let i = 0; i < uint8Array.byteLength; i++) {
        binary += String.fromCharCode(uint8Array[i]);
    }
    const base64Content = btoa(binary);
    const res = await fetch(`${API_URL}/v1/expenses/upload`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`
        },
        body: JSON.stringify({
            file_name: file.name,
            content: base64Content
        }),
    });

    if (!res.ok) {
        const errorData = await res.json();
        throw new Error(errorData.message || "Ошибка при загрузке");
    }
    
    return res.json();
}