import React, { useState, useEffect, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import { embedDashboard } from "@superset-ui/embedded-sdk";
import { getSupersetToken, uploadStatement } from "../services/api";
import '../css/DashboardScreen.css';



export default function DashboardScreen() {
    const navigate = useNavigate();
    const dashboardRef = useRef<HTMLDivElement>(null);
    const fileInputRef = useRef<HTMLInputElement>(null);

    const getDefaultDates = () => {
        const now = new Date();
        const firstDay = new Date(now.getFullYear(), now.getMonth(), 1);

        return {
            start: firstDay.toLocaleDateString('en-CA'),
            end: now.toLocaleDateString('en-CA')
        };
    };

    const initialDates = getDefaultDates();

    const [startDate, setStartDate] = useState(initialDates.start);
    const [endDate, setEndDate] = useState(initialDates.end);
    const [loading, setLoading] = useState(false);
    const [uploading, setUploading] = useState(false);
    const [refreshKey, setRefreshKey] = useState(0);
    const dashboard_id = process.env.REACT_APP_DASHBOARD_OVERVIEW;


    useEffect(() => {
        const renderChart = async () => {
            if (!dashboardRef.current) return;
            setLoading(true);

            try {

                const res = await getSupersetToken(startDate, endDate, String(dashboard_id));

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
                        hideTab: true,
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
            <header className="compact-header">
                <div className="header-left">
                    <h1>Финансы</h1>
                    <div className="date-picker-group">
                        <input type="date" value={startDate} onChange={(e) => setStartDate(e.target.value)} title="Начало" />
                        <span className="date-separator">—</span>
                        <input type="date" value={endDate} onChange={(e) => setEndDate(e.target.value)} title="Конец" />
                    </div>
                </div>

                <div className="header-right">
                    <input type="file" ref={fileInputRef} onChange={handleFileUpload} style={{ display: 'none' }} accept=".csv" />
                    <button className="action-btn upload" onClick={() => fileInputRef.current?.click()} disabled={uploading}>
                        {uploading ? "..." : "Загрузить CSV"}
                    </button>
                    <button className="action-btn logout" onClick={() => navigate('/login')}>Выход</button>
                </div>
            </header>

            <main className="dashboard-content">
                {loading && <div className="overlay-loader">Обновление графиков...</div>}
                <div ref={dashboardRef} className="superset-container" />
            </main>
        </div>
    );
}
