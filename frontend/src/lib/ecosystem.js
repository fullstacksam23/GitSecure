import { normalizeRisk, normalizeSeverity, severityOrder } from "./utils";

export const NO_KNOWN_VULNERABILITIES_STATUS = "no_known_vulnerabilities";
export const NO_KNOWN_VULNERABILITIES_LABEL = "No Known Vulnerabilities";

export function isActiveBatchStatus(status) {
  return ["queued", "running"].includes(String(status || "").toLowerCase());
}

export function getBatchId(batch) {
  return batch?.batch_id || batch?.id || "";
}

export function getBatchLanguage(batch) {
  return batch?.language || batch?.ecosystem || batch?.lang || "Unknown";
}

export function getBatchProgress(batch) {
  const completed = Number(batch?.completed_repos ?? batch?.progress?.completed_repos ?? batch?.completed ?? 0);
  const total = Number(batch?.total_repos ?? batch?.progress?.total_repos ?? batch?.total ?? 0);

  return {
    completed,
    total,
    percent: total > 0 ? Math.min(100, Math.round((completed / total) * 100)) : 0,
  };
}

export function getSummaryTotalRepos(summary, batch) {
  return Number(summary?.total_repositories ?? summary?.total_repos ?? batch?.total_repos ?? 0);
}

export function getSummaryTotalVulnerabilities(summary) {
  return Number(summary?.total_vulnerabilities ?? summary?.vulnerability_count ?? summary?.total_findings ?? 0);
}

export function getSummarySeverityBreakdown(summary) {
  const counts = summary?.severity_breakdown || summary?.severity_counts || {};
  return {
    critical: Number(counts.critical || 0),
    high: Number(counts.high || 0),
    medium: Number(counts.medium || 0),
    low: Number(counts.low || 0),
  };
}

export function getSummaryAverageRisk(summary) {
  return normalizeRisk(summary?.average_risk_score ?? summary?.avg_risk_score ?? summary?.average_risk ?? 0);
}

export function getMostVulnerableRepo(summary) {
  return summary?.most_vulnerable_repo || null;
}

export function getRepoName(repo) {
  return repo?.repo_name || repo?.full_name || repo?.repo || repo?.name || "Unknown repository";
}

export function getRepoRank(repo) {
  return Number(repo?.rank ?? repo?.repo_rank ?? repo?.position ?? 0);
}

export function getRepoStars(repo) {
  return Number(repo?.stars ?? repo?.stargazers_count ?? 0);
}

export function getRepoStatus(repo) {
  const status = String(repo?.scan_status || repo?.status || "").trim().toLowerCase();
  if (status) return status;

  return getRepoVulnerabilityCount(repo) === 0 ? NO_KNOWN_VULNERABILITIES_STATUS : "unknown";
}

export function getRepoVulnerabilityCount(repo) {
  return Number(repo?.vulnerability_count ?? repo?.total_vulnerabilities ?? repo?.findings ?? 0);
}

export function getRepoTopSeverity(repo) {
  if (repo?.top_severity) return normalizeSeverity(repo.top_severity);

  const counts = repo?.severity_counts || {};
  const severity = ["critical", "high", "medium", "low"].find((key) => Number(counts[key] || 0) > 0);
  return severity || "unknown";
}

export function getRepoRiskScore(repo) {
  return normalizeRisk(repo?.risk_score ?? repo?.average_risk_score ?? repo?.risk ?? 0);
}

export function getRepoJobId(repo) {
  return repo?.job_id || repo?.scan_job_id || repo?.latest_job_id || "";
}

export function getRepoVulnerabilityPreview(repoOrSummary) {
  return (
    repoOrSummary?.top_vulnerabilities ||
    repoOrSummary?.vulnerabilities ||
    repoOrSummary?.recent_vulnerabilities ||
    []
  );
}

export function sortReposForHighlight(items = []) {
  return [...items].sort((left, right) => {
    const bySeverity = severityOrder(getRepoTopSeverity(left)) - severityOrder(getRepoTopSeverity(right));
    if (bySeverity !== 0) return bySeverity;

    const byCount = getRepoVulnerabilityCount(right) - getRepoVulnerabilityCount(left);
    if (byCount !== 0) return byCount;

    return getRepoRiskScore(right) - getRepoRiskScore(left);
  });
}
