import RiskBar from "../shared/RiskBar";

export default function TopRiskPackages({ items = [] }) {
  return (
    <div className="panel p-6">
      <div className="mb-5">
        <h3 className="text-lg font-semibold text-white">Top Risk Packages</h3>
        <p className="text-sm text-slate-400">Packages with the highest observed risk score.</p>
      </div>
      <div className="space-y-3">
        {items.map((item) => (
          <div key={`${item.package}-${item.ecosystem}`} className="surface p-4">
            <div className="mb-3 flex items-center justify-between gap-3">
              <div>
                <p className="font-mono text-sm text-white">{item.package}</p>
                <p className="text-xs text-slate-500">{item.ecosystem || "Unknown ecosystem"}</p>
              </div>
              <span className="text-xs text-slate-400">{item.vulnerability_count} findings</span>
            </div>
            <RiskBar value={item.risk} />
          </div>
        ))}
      </div>
    </div>
  );
}
