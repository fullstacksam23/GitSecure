const API_BASE_URL = import.meta.env.VITE_API_URL || "http://localhost:3000";

async function request(path, options = {}) {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    headers: {
      "Content-Type": "application/json",
      ...(options.headers || {}),
    },
    ...options,
  });

  if (!response.ok) {
    let message = `Request failed with status ${response.status}`;

    try {
      const text = await response.text();
      if (text) {
        message = text;
      }
    } catch {
      // Ignore parse errors and use fallback message.
    }

    throw new Error(message);
  }

  const contentType = response.headers.get("content-type") || "";
  if (!contentType.includes("application/json")) {
    return null;
  }

  return response.json();
}

export const api = {
  getHealth() {
    return request("/health", { headers: { Accept: "text/plain" } });
  },

  startScan({ owner, repo }) {
    return request("/scan", {
      method: "POST",
      body: JSON.stringify({ owner, repo }),
    });
  },

  getDashboardSummary() {
    return request("/dashboard/summary");
  },

  getRecentScans({ limit = 6 } = {}) {
    return request(`/scans?limit=${limit}`);
  },

  getScan(jobId) {
    return request(`/scans/${jobId}`);
  },

  getVulnerabilities({
    page = 1,
    pageSize = 10,
    search = "",
    severity = "all",
    jobId = "",
  } = {}) {
    const params = new URLSearchParams({
      page: String(page),
      page_size: String(pageSize),
    });

    if (search) params.set("search", search);
    if (severity && severity !== "all") params.set("severity", severity);
    if (jobId) params.set("job_id", jobId);

    return request(`/vulnerabilities?${params.toString()}`);
  },

  getVulnerabilityById(vulnerabilityId, jobId) {
    const params = new URLSearchParams();
    if (jobId) params.set("job_id", jobId);
    const suffix = params.toString() ? `?${params.toString()}` : "";
    return request(`/vulnerabilities/${vulnerabilityId}${suffix}`);
  },
};
