import React, { useState, useEffect, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import { embedDashboard } from "@superset-ui/embedded-sdk";
import { getSupersetToken, uploadStatement } from "../services/api";
import '../css/DashboardScreen.css';



export default function DashboardScreen() {
    const navigate = useNavigate();
    const dashboardRef = useRef<HTMLDivElement>(null);
    const fileInputRef = useRef<HTMLInputElement>(null);

    const [startDate, setStartDate] = useState('2024-01-01');
    const [endDate, setEndDate] = useState('2026-03-01');
    const [loading, setLoading] = useState(false);
    const [uploading, setUploading] = useState(false);
    const [refreshKey, setRefreshKey] = useState(0);

    useEffect(() => {
        const renderChart = async () => {
            if (!dashboardRef.current) return;
            setLoading(true);

            try {

                const res = await getSupersetToken(startDate, endDate);

                if (!res.token || !res.dashboardId) {
                    console.error("Бэкенд не вернул необходимые данные:", res);
                    return;
                }

                dashboardRef.current.innerHTML = "";

                await embedDashboard({
                    id: res.dashboardId,
                    supersetDomain: "http://localhost:8088",
                    mountPoint: dashboardRef.current,
                    fetchGuestToken: async () => res.token,
                    dashboardUiConfig: {
                        hideTitle: true,
                        hideChartControls: true,
                        filters: { visible: false, expanded: false }
                    },
                });
            } catch (err) {
                console.error("Ошибка встраивания дашборда:", err);
            } finally {
                setLoading(false);
            }
        };

        renderChart();
    }, [startDate, endDate, refreshKey]);

    const handleFileUpload = async (event: React.ChangeEvent<HTMLInputElement>) => {
        const file = event.target.files?.[0];
        if (!file) return;

        setUploading(true);
        try {
            await uploadStatement(file);
            setRefreshKey(prev => prev + 1);
            alert("Выписка успешно загружена!");
            setStartDate(startDate);
        } catch (err) {
            alert("Ошибка при загрузке файла");
            console.error(err);
        } finally {
            setUploading(false);
            if (fileInputRef.current) fileInputRef.current.value = "";
        }
    };

    const handleLogout = () => {
        localStorage.removeItem("user_token");
        navigate("/login");
    };

    return (
        <div className="dashboard-page">
            <header className="dashboard-header">
                <div className="header-content">
                    <h1>Аналитика расходов</h1>
                    <button onClick={handleLogout} className="logout-btn">Выйти</button>
                </div>
            </header>

            <main className="dashboard-main">
                <section className="filter-section">
                    <div className="filters">
                        <div className="input-group">
                            <label>Период с</label>
                            <input type="date" value={startDate} onChange={(e) => setStartDate(e.target.value)} />
                        </div>
                        <div className="input-group">
                            <label>по</label>
                            <input type="date" value={endDate} onChange={(e) => setEndDate(e.target.value)} />
                        </div>
                    </div>

                    <div className="upload-block">
                        <input
                            type="file"
                            ref={fileInputRef}
                            onChange={handleFileUpload}
                            style={{ display: 'none' }}
                            accept=".csv"
                        />
                        <button
                            className="upload-btn"
                            onClick={() => fileInputRef.current?.click()}
                            disabled={uploading}
                        >
                            {uploading ? "Загрузка..." : "Загрузить выписку"}
                        </button>
                    </div>

                    {(loading || uploading) && <div className="loading-indicator">Обновление данных...</div>}
                </section>

                <section className="chart-container">
                    <div ref={dashboardRef} className="superset-embed" />
                </section>
            </main>
        </div>
    );
}