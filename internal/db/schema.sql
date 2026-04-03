CREATE TABLE scan_jobs (
    job_id TEXT PRIMARY KEY,
    repo TEXT,
    status TEXT,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE vulnerabilities (
    id TEXT NOT NULL,
    job_id TEXT NOT NULL,

    package TEXT,
    version TEXT,
    severity TEXT,
    summary TEXT,

    urls TEXT[],
    fix_version TEXT[],
    fix_state TEXT,

    risk DOUBLE PRECISION,
    namespace TEXT,

    match_type TEXT,
    version_constraint TEXT,

    data_source TEXT,
    source TEXT,

    cwe_ids TEXT[],
    ecosystem TEXT,

    created_at TIMESTAMP DEFAULT now(),

    -- relationships
    CONSTRAINT fk_job
        FOREIGN KEY(job_id)
        REFERENCES scan_jobs(job_id)
        ON DELETE CASCADE,

    -- composite primary key
    PRIMARY KEY (id, job_id)
);

CREATE INDEX idx_job ON vulnerabilities(job_id);
CREATE INDEX idx_vuln_id ON vulnerabilities(vuln_id);