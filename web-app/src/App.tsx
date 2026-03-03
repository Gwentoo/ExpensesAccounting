
import React from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import LoginScreen from "./pages/LoginScreen";
import RegisterScreen from "./pages/RegisterScreen";
import VerifyScreen from "./pages/VerifyScreen";
import NewPassScreen from "./pages/NewPassScreen";
import RestorePassScreen from "./pages/RestorePassScreen";
import DashboardScreen from "./pages/DashboardScreen";

export default function App() {
  return (
      <Router>
        <div className="App">
          <Routes>
            <Route path="/" element={<RegisterScreen />} />
            <Route path="/login" element={<LoginScreen />} />
            <Route path="/register" element={<RegisterScreen />} />
            <Route path="/verify" element={<VerifyScreen />} />
            <Route path="/newpass" element={<NewPassScreen />} /><Route path="/dashboard" element={<DashboardScreen />} />
            <Route path="/restore-pass/:token" element={<RestorePassScreen />} />
          </Routes>
        </div>
      </Router>
  );
}