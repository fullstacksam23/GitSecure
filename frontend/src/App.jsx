import { lazy, Suspense } from "react";
import { Navigate, Route, Routes } from "react-router-dom";
import AppShell from "./components/layout/AppShell";
import PageSkeleton from "./components/shared/PageSkeleton";

const DashboardPage = lazy(() => import("./pages/DashboardPage"));
const ScansPage = lazy(() => import("./pages/ScansPage"));
const HistoryPage = lazy(() => import("./pages/HistoryPage"));
const ComparePage = lazy(() => import("./pages/ComparePage"));
const ScanDetailPage = lazy(() => import("./pages/ScanDetailPage"));
const NewScanPage = lazy(() => import("./pages/NewScanPage"));
const EcosystemBatchesPage = lazy(() => import("./pages/EcosystemBatchesPage"));
const BatchDetailPage = lazy(() => import("./pages/BatchDetailPage"));
const NewEcosystemScanPage = lazy(() => import("./pages/NewEcosystemScanPage"));

function App() {
  return (
    <Suspense fallback={<AppShell><PageSkeleton cards={4} rows={6} /></AppShell>}>
      <Routes>
        <Route element={<AppShell />}>
          <Route path="/" element={<Navigate to="/dashboard" replace />} />
          <Route path="/dashboard" element={<DashboardPage />} />
          <Route path="/scans" element={<ScansPage />} />
          <Route path="/history" element={<HistoryPage />} />
          <Route path="/compare" element={<ComparePage />} />
          <Route path="/ecosystem/batches" element={<EcosystemBatchesPage />} />
          <Route path="/ecosystem/new-scan" element={<NewEcosystemScanPage />} />
          <Route path="/ecosystem/batches/:batchId" element={<BatchDetailPage />} />
          <Route path="/scans/:jobId" element={<ScanDetailPage />} />
          <Route path="/new-scan" element={<NewScanPage />} />
        </Route>
      </Routes>
    </Suspense>
  );
}

export default App;
