ALTER TABLE scan_jobs
    ADD COLUMN IF NOT EXISTS commit_hash TEXT;

ALTER TABLE vulnerabilities
    ALTER COLUMN package SET NOT NULL,
    ALTER COLUMN version SET NOT NULL;

ALTER TABLE vulnerabilities
    DROP CONSTRAINT IF EXISTS vulnerabilities_pkey;

ALTER TABLE vulnerabilities
    ADD CONSTRAINT vulnerabilities_pkey PRIMARY KEY (job_id, id, package, version);

DROP INDEX IF EXISTS idx_job;
DROP INDEX IF EXISTS idx_vuln_id;

CREATE INDEX IF NOT EXISTS idx_vulnerabilities_job_id ON vulnerabilities(job_id);
CREATE INDEX IF NOT EXISTS idx_vulnerabilities_id ON vulnerabilities(id);
CREATE INDEX IF NOT EXISTS idx_vulnerabilities_created_at ON vulnerabilities(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_vulnerabilities_severity ON vulnerabilities(lower(severity));
CREATE INDEX IF NOT EXISTS idx_scan_jobs_created_at ON scan_jobs(created_at DESC);

CREATE OR REPLACE VIEW vulnerability_records AS
SELECT
    v.job_id,
    v.id,
    v.package,
    v.version,
    v.severity,
    v.summary,
    v.urls,
    v.fix_version,
    v.fix_state,
    v.risk,
    v.namespace,
    v.match_type,
    v.version_constraint,
    v.data_source,
    v.source,
    v.cwe_ids,
    v.ecosystem,
    v.created_at,
    CASE
        WHEN lower(coalesce(v.severity, '')) LIKE '%critical%' THEN 'critical'
        WHEN lower(coalesce(v.severity, '')) LIKE '%high%' THEN 'high'
        WHEN lower(coalesce(v.severity, '')) LIKE '%medium%' THEN 'medium'
        WHEN lower(coalesce(v.severity, '')) LIKE '%moderate%' THEN 'medium'
        WHEN lower(coalesce(v.severity, '')) LIKE '%low%' THEN 'low'
        WHEN lower(coalesce(v.severity, '')) LIKE '%negligible%' THEN 'low'
        ELSE 'unknown'
    END AS normalized_severity
FROM vulnerabilities v;

CREATE OR REPLACE VIEW dashboard_summary AS
SELECT
    COUNT(*)::BIGINT AS total_vulnerabilities,
    COUNT(*) FILTER (WHERE normalized_severity = 'critical')::BIGINT AS critical,
    COUNT(*) FILTER (WHERE normalized_severity = 'high')::BIGINT AS high,
    COUNT(*) FILTER (WHERE normalized_severity = 'medium')::BIGINT AS medium,
    COUNT(*) FILTER (WHERE normalized_severity = 'low')::BIGINT AS low
FROM vulnerability_records;

CREATE OR REPLACE VIEW scan_job_summaries AS
SELECT
    sj.job_id,
    sj.repo,
    sj.status,
    sj.commit_hash,
    sj.created_at,
    COALESCE(
        CASE MIN(
            CASE vr.normalized_severity
                WHEN 'critical' THEN 1
                WHEN 'high' THEN 2
                WHEN 'medium' THEN 3
                WHEN 'low' THEN 4
                ELSE 5
            END
        )
            WHEN 1 THEN 'critical'
            WHEN 2 THEN 'high'
            WHEN 3 THEN 'medium'
            WHEN 4 THEN 'low'
            ELSE 'unknown'
        END,
        'unknown'
    ) AS top_severity,
    COUNT(vr.id)::BIGINT AS vulnerability_count
FROM scan_jobs sj
LEFT JOIN vulnerability_records vr ON vr.job_id = sj.job_id
GROUP BY sj.job_id, sj.repo, sj.status, sj.commit_hash, sj.created_at;
