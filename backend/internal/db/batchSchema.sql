create type severity_level as enum ('low', 'medium', 'high', 'critical');

create table public.ecosystem_scans (
  batch_id text not null,
  language text not null,
  status text check (status in ('pending', 'running', 'completed', 'failed')) default 'pending',

  repo_count int not null,
  completed_repos int default 0,

  created_at timestamp without time zone default now(),
  completed_at timestamp without time zone null,

  constraint ecosystem_scans_pkey primary key (batch_id)
);

create index idx_ecosystem_scans_created_at 
on public.ecosystem_scans (created_at desc);


create table public.ecosystem_repos (
  id bigserial primary key,
  batch_id text not null,
  repo text not null,

  stars int,
  rank int,

  created_at timestamp without time zone default now(),

  constraint fk_ecosystem_repos_batch
  foreign key (batch_id) references ecosystem_scans(batch_id) on delete cascade
);

create index idx_ecosystem_repos_batch 
on public.ecosystem_repos (batch_id);

create index idx_ecosystem_repos_repo 
on public.ecosystem_repos (repo);

create table public.vuln_catalog (
  normalized_id text primary key,  -- CVE/GHSA canonical ID

  severity severity_level,
  summary text,
  urls text[],
  cwe_ids text[],

  created_at timestamp default now()
);

create table public.ecosystem_vulnerabilities (
  id bigserial primary key,

  batch_id text not null,
  repo text not null,

  package text not null,
  version text not null,

  normalized_id text not null,

  fix_version text[],
  fix_state text,
  risk double precision,

  namespace text,
  match_type text,
  version_constraint text,
  data_source text,
  source text,
  ecosystem text,

  created_at timestamp without time zone default now(),

  constraint fk_ecosystem_vulns_batch
  foreign key (batch_id) references ecosystem_scans(batch_id) on delete cascade,

  constraint fk_vuln_catalog
  foreign key (normalized_id) references vuln_catalog(normalized_id)
);

create unique index uniq_ecosystem_vuln_entry
on public.ecosystem_vulnerabilities (
  batch_id, repo, package, version, normalized_id
);

-- batch filtering
create index idx_ecosystem_vulns_batch 
on public.ecosystem_vulnerabilities (batch_id);

-- aggregation (top CVEs)
create index idx_ecosystem_vulns_normalized 
on public.ecosystem_vulnerabilities (normalized_id);

-- package stats
create index idx_ecosystem_vulns_package 
on public.ecosystem_vulnerabilities (package);

-- severity via join (optional optimization)
create index idx_vuln_catalog_severity 
on public.vuln_catalog (severity);

-- time sorting
create index idx_ecosystem_vulns_created_at 
on public.ecosystem_vulnerabilities (created_at desc);

create table public.ecosystem_stats (
  id bigserial primary key,
  batch_id text not null,

  stat_type text check (stat_type in ('cve', 'package', 'severity')),
  key text not null,
  count int not null,

  created_at timestamp default now(),

  constraint fk_ecosystem_stats_batch
  foreign key (batch_id) references ecosystem_scans(batch_id) on delete cascade
);

create index idx_ecosystem_stats_batch 
on public.ecosystem_stats (batch_id);

create index idx_ecosystem_stats_type 
on public.ecosystem_stats (stat_type);