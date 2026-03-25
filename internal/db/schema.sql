CREATE TABLE scan_jobs (
    job_id TEXT PRIMARY KEY,
    repo TEXT,
    status TEXT,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE vulnerabilities (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    job_id TEXT REFERENCES scan_jobs(job_id),

    vuln_id TEXT,
    package TEXT,
    version TEXT,

    severity TEXT,
    summary TEXT,

    urls TEXT[],
    fix_versions TEXT[],

    source TEXT
);

CREATE INDEX idx_job ON vulnerabilities(job_id);
CREATE INDEX idx_vuln_id ON vulnerabilities(vuln_id);