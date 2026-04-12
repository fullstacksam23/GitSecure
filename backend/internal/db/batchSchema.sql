CREATE TABLE ecosystem_batches (
    batch_id TEXT PRIMARY KEY,

    language TEXT NOT NULL,

    repo_count INT NOT NULL,         -- requested number (user input)
    total_repos INT DEFAULT 0,       -- actual fetched
    completed_repos INT DEFAULT 0,   -- progress tracking

    status TEXT NOT NULL,            -- pending, running, completed, failed

    created_at TIMESTAMPTZ DEFAULT now(),
    completed_at TIMESTAMPTZ
);

CREATE TABLE ecosystem_repos (
    id BIGSERIAL PRIMARY KEY,

    batch_id TEXT NOT NULL,

    repo_name TEXT NOT NULL, -- full_name (owner/repo)
    stars INT,
    forks INT,
    repo_rank INT, -- position in top X

    CONSTRAINT fk_batch
        FOREIGN KEY(batch_id)
        REFERENCES ecosystem_batches(batch_id)
        ON DELETE CASCADE
);

CREATE OR REPLACE FUNCTION increment_batch_progress(batch_id TEXT)
RETURNS void AS $$
BEGIN
    UPDATE ecosystem_batches
    SET 
        completed_repos = completed_repos + 1,
        status = CASE 
            WHEN completed_repos + 1 >= total_repos THEN 'completed'
            ELSE 'running'
        END,
        completed_at = CASE 
            WHEN completed_repos + 1 >= total_repos THEN now()
            ELSE completed_at
        END
    WHERE ecosystem_batches.batch_id = increment_batch_progress.batch_id;
END;
$$ LANGUAGE plpgsql;

CREATE INDEX idx_batches_status ON ecosystem_batches(status);