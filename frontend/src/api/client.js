const API_BASE_URL = import.meta.env.VITE_API_URL || "http://localhost:3000";

async function request(path, options = {}) {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    cache: options.method && options.method !== "GET" ? "default" : "no-store",
    headers: {
      "Content-Type": "application/json",
      ...(options.headers || {}),
    },
    ...options,
  });

  if (!response.ok) {
    let message = `Request failed with status ${response.status}`;
    try {
      const payload = await response.json();
      message = payload.error || message;
    } catch {
      try {
        message = await response.text();
      } catch {
        // ignore fallback errors
      }
    }
    throw new Error(message);
  }

  const type = response.headers.get("content-type") || "";
  if (!type.includes("application/json")) return null;
  return response.json();
}

function toParams(values) {
  const params = new URLSearchParams();
  Object.entries(values).forEach(([key, value]) => {
    if (value === undefined || value === null || value === "" || value === "all") return;
    params.set(key, String(value));
  });
  const text = params.toString();
  return text ? `?${text}` : "";
}

export const api = {
  getDashboardSummary() {
    return request("/dashboard/summary");
  },

  getScans({ page = 1, pageSize = 20, repo = "", status = "" } = {}) {
    return request(
      `/scans${toParams({
        page,
        page_size: pageSize,
        repo,
        status,
      })}`
    );
  },

  getScan(jobId) {
    return request(`/scans/${jobId}`);
  },

  compareScans({ base, target }) {
    return request(`/scans/compare${toParams({ base, target })}`);
  },

  getVulnerabilities({
    page = 1,
    pageSize = 50,
    search = "",
    severity = "",
    jobId = "",
    ecosystem = "",
    fixState = "",
    sortBy = "created_at",
    sortOrder = "desc",
  } = {}) {
    return request(
      `/vulnerabilities${toParams({
        page,
        page_size: pageSize,
        search,
        severity,
        job_id: jobId,
        ecosystem,
        fix_state: fixState,
        sort_by: sortBy,
        sort_order: sortOrder,
      })}`
    );
  },

  async getAllVulnerabilities({
    pageSize = 100,
    search = "",
    severity = "",
    jobId = "",
    ecosystem = "",
    fixState = "",
    sortBy = "created_at",
    sortOrder = "desc",
  } = {}) {
    const firstPage = await api.getVulnerabilities({
      page: 1,
      pageSize,
      search,
      severity,
      jobId,
      ecosystem,
      fixState,
      sortBy,
      sortOrder,
    });

    const totalPages = firstPage?.pagination?.total_pages || 1;
    if (totalPages <= 1) {
      return firstPage;
    }

    const rest = await Promise.all(
      Array.from({ length: totalPages - 1 }, (_, index) =>
        api.getVulnerabilities({
          page: index + 2,
          pageSize,
          search,
          severity,
          jobId,
          ecosystem,
          fixState,
          sortBy,
          sortOrder,
        })
      )
    );

    const items = [firstPage, ...rest].flatMap((page) => page?.items || []);

    return {
      ...firstPage,
      items,
      pagination: {
        ...firstPage.pagination,
        page: 1,
        page_size: items.length,
        total_items: items.length,
        total_pages: 1,
      },
    };
  },

  getVulnerabilityById(vulnerabilityId, jobId, packageName) {
    return request(`/vulnerabilities/${vulnerabilityId}${toParams({ job_id: jobId, package: packageName })}`);
  },

  startScan({ owner, repo }) {
    return request("/scan", {
      method: "POST",
      body: JSON.stringify({ owner, repo }),
    });
  },
};
