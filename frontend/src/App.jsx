import { Navigate, Route, Routes } from "react-router-dom";
import AppShell from "./components/layout/AppShell";
import ComparePage from "./pages/ComparePage";
import DashboardPage from "./pages/DashboardPage";
import HistoryPage from "./pages/HistoryPage";
import NewScanPage from "./pages/NewScanPage";
import ScanDetailPage from "./pages/ScanDetailPage";
import ScansPage from "./pages/ScansPage";

function App() {
  return (
    <Routes>
      <Route element={<AppShell />}>
        <Route path="/" element={<Navigate to="/dashboard" replace />} />
        <Route path="/dashboard" element={<DashboardPage />} />
        <Route path="/scans" element={<ScansPage />} />
        <Route path="/history" element={<HistoryPage />} />
        <Route path="/compare" element={<ComparePage />} />
        <Route path="/scans/:jobId" element={<ScanDetailPage />} />
        <Route path="/new-scan" element={<NewScanPage />} />
      </Route>
    </Routes>
  );
}

export default App;
