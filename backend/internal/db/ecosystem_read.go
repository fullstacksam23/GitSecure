package db

import (
	"encoding/json"
	"errors"
	"sort"
	"strings"

	"github.com/fullstacksam23/GitSecure/internal/models"
	"github.com/supabase-community/postgrest-go"
)

type EcosystemBatchFilter struct {
	Page     int
	PageSize int
	Status   string
	Language string
	Search   string
}

type EcosystemRepoFilter struct {
	Page      int
	PageSize  int
	Search    string
	Status    string
	Severity  string
	SortBy    string
	SortOrder string
}

type ecosystemBatchRow struct {
	BatchID        string  `json:"batch_id"`
	Language       string  `json:"language"`
	Status         string  `json:"status"`
	RepoCount      int64   `json:"repo_count"`
	TotalRepos     int64   `json:"total_repos"`
	CompletedRepos int64   `json:"completed_repos"`
	CreatedAt      string  `json:"created_at"`
	CompletedAt    *string `json:"completed_at"`
}

type ecosystemRepoRow struct {
	ID       int64  `json:"id"`
	BatchID  string `json:"batch_id"`
	RepoName string `json:"repo_name"`
	Stars    int64  `json:"stars"`
	RepoRank int64  `json:"repo_rank"`
}

type ecosystemScanRow struct {
	JobID      string  `json:"job_id"`
	BatchID    string  `json:"batch_id"`
	RepoID     int64   `json:"repo_id"`
	Repo       string  `json:"repo"`
	Status     string  `json:"status"`
	CreatedAt  string  `json:"created_at"`
	CommitHash *string `json:"commit_hash"`
}

func ListEcosystemBatches(filter EcosystemBatchFilter) ([]models.EcosystemBatchListItem, int64, error) {
	if Client == nil {
		return nil, 0, errors.New("client not initialized")
	}

	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}

	query := Client.
		From("ecosystem_batches").
		Select("batch_id,language,status,repo_count,total_repos,completed_repos,created_at,completed_at", "exact", false)

	if value := strings.TrimSpace(filter.Status); value != "" {
		query = query.Eq("status", value)
	}
	if value := strings.TrimSpace(filter.Language); value != "" {
		query = query.Eq("language", value)
	}
	if value := strings.TrimSpace(filter.Search); value != "" {
		query = query.Or(strings.Join([]string{
			"batch_id.ilike.%" + escapeLike(value) + "%",
			"language.ilike.%" + escapeLike(value) + "%",
		}, ","), "")
	}

	from := (filter.Page - 1) * filter.PageSize
	to := from + filter.PageSize - 1

	var rows []ecosystemBatchRow
	count, err := executeRangeTo(query, from, to, &rows)
	if err != nil {
		return nil, 0, err
	}

	items := make([]models.EcosystemBatchListItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, models.EcosystemBatchListItem{
			BatchID:        row.BatchID,
			Language:       row.Language,
			Status:         NormalizeJobStatus(row.Status),
			RepoCount:      row.RepoCount,
			TotalRepos:     row.TotalRepos,
			CompletedRepos: row.CompletedRepos,
			CreatedAt:      row.CreatedAt,
			CompletedAt:    stringOrEmpty(row.CompletedAt),
		})
	}

	return items, count, nil
}

func GetEcosystemBatchByID(batchID string) (*models.EcosystemBatchDetail, error) {
	if Client == nil {
		return nil, errors.New("client not initialized")
	}

	var rows []ecosystemBatchRow
	_, err := Client.
		From("ecosystem_batches").
		Select("batch_id,language,status,repo_count,total_repos,completed_repos,created_at,completed_at", "", false).
		Eq("batch_id", batchID).
		Limit(1, "").
		ExecuteTo(&rows)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil
	}

	row := rows[0]
	return &models.EcosystemBatchDetail{
		BatchID:        row.BatchID,
		Language:       row.Language,
		Status:         NormalizeJobStatus(row.Status),
		RepoCount:      row.RepoCount,
		TotalRepos:     row.TotalRepos,
		CompletedRepos: row.CompletedRepos,
		CreatedAt:      row.CreatedAt,
		CompletedAt:    stringOrEmpty(row.CompletedAt),
	}, nil
}

func GetEcosystemBatchSummary(batchID string) (*models.EcosystemBatchSummary, error) {
	repos, err := buildEcosystemRepoItems(batchID)
	if err != nil {
		return nil, err
	}

	summary := &models.EcosystemBatchSummary{
		BatchID:           batchID,
		TotalRepositories: int64(len(repos)),
	}

	if len(repos) == 0 {
		return summary, nil
	}

	var totalRisk float64
	var mostVulnerable *models.EcosystemRepoSummary
	var mostVulnerableVulns []models.VulnerabilityRecord

	for _, repo := range repos {
		summary.TotalVulnerabilities += repo.VulnerabilityCount
		summary.SeverityBreakdown.Critical += repo.SeverityCounts.Critical
		summary.SeverityBreakdown.High += repo.SeverityCounts.High
		summary.SeverityBreakdown.Medium += repo.SeverityCounts.Medium
		summary.SeverityBreakdown.Low += repo.SeverityCounts.Low
		summary.SeverityBreakdown.Unknown += repo.SeverityCounts.Unknown
		totalRisk += repo.RiskScore

		if mostVulnerable == nil || repo.VulnerabilityCount > mostVulnerable.VulnerabilityCount || (repo.VulnerabilityCount == mostVulnerable.VulnerabilityCount && severityRank(repo.TopSeverity) < severityRank(mostVulnerable.TopSeverity)) {
			mostVulnerable = &models.EcosystemRepoSummary{
				RepoName:           repo.RepoName,
				JobID:              repo.JobID,
				ScanStatus:         repo.ScanStatus,
				VulnerabilityCount: repo.VulnerabilityCount,
				TopSeverity:        repo.TopSeverity,
				RiskScore:          repo.RiskScore,
				SeverityCounts:     repo.SeverityCounts,
			}

			if strings.TrimSpace(repo.JobID) != "" {
				vulns, vulnErr := listVulnerabilitiesForJob(repo.JobID)
				if vulnErr != nil {
					return nil, vulnErr
				}
				sortVulnerabilities(vulns)
				mostVulnerableVulns = truncateVulnerabilities(vulns, 5)
			} else {
				mostVulnerableVulns = nil
			}
		}
	}

	summary.AverageRiskScore = totalRisk / float64(len(repos))
	if mostVulnerable != nil {
		mostVulnerable.TopVulnerabilities = mostVulnerableVulns
		summary.MostVulnerableRepo = mostVulnerable
	}

	return summary, nil
}

func ListEcosystemRepos(batchID string, filter EcosystemRepoFilter) ([]models.EcosystemRepoListItem, int64, error) {
	items, err := buildEcosystemRepoItems(batchID)
	if err != nil {
		return nil, 0, err
	}

	filtered := make([]models.EcosystemRepoListItem, 0, len(items))
	for _, item := range items {
		if value := strings.TrimSpace(filter.Search); value != "" && !strings.Contains(strings.ToLower(item.RepoName), strings.ToLower(value)) {
			continue
		}
		if value := strings.TrimSpace(filter.Status); value != "" && strings.ToLower(item.ScanStatus) != strings.ToLower(value) {
			continue
		}
		if value := strings.TrimSpace(filter.Severity); value != "" && strings.ToLower(item.TopSeverity) != strings.ToLower(value) {
			continue
		}
		filtered = append(filtered, item)
	}

	sortEcosystemRepoItems(filtered, filter)

	total := int64(len(filtered))
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}

	from := (filter.Page - 1) * filter.PageSize
	if from >= len(filtered) {
		return []models.EcosystemRepoListItem{}, total, nil
	}

	to := from + filter.PageSize
	if to > len(filtered) {
		to = len(filtered)
	}

	return filtered[from:to], total, nil
}

func buildEcosystemRepoItems(batchID string) ([]models.EcosystemRepoListItem, error) {
	repos, err := listEcosystemRepos(batchID)
	if err != nil {
		return nil, err
	}
	scans, err := listEcosystemScans(batchID)
	if err != nil {
		return nil, err
	}

	scansByRepoID := map[int64]ecosystemScanRow{}
	jobIDs := make([]string, 0, len(scans))
	for _, scan := range scans {
		jobIDs = append(jobIDs, scan.JobID)
		current, ok := scansByRepoID[scan.RepoID]
		if !ok || scan.CreatedAt > current.CreatedAt {
			scansByRepoID[scan.RepoID] = scan
		}
	}

	aggregates, err := aggregateScanVulnerabilities(jobIDs)
	if err != nil {
		return nil, err
	}

	items := make([]models.EcosystemRepoListItem, 0, len(repos))
	for _, repo := range repos {
		item := models.EcosystemRepoListItem{
			ID:         repo.ID,
			BatchID:    repo.BatchID,
			RepoName:   repo.RepoName,
			Stars:      repo.Stars,
			Rank:       repo.RepoRank,
			ScanStatus: "queued",
			TopSeverity: "unknown",
		}

		if scan, ok := scansByRepoID[repo.ID]; ok {
			item.JobID = scan.JobID
			item.ScanStatus = NormalizeJobStatus(scan.Status)
			if aggregate, ok := aggregates[scan.JobID]; ok {
				item.VulnerabilityCount = aggregate.Count
				item.TopSeverity = aggregate.TopSeverity
				item.SeverityCounts = aggregate.SeverityCounts
			}
			item.RiskScore = estimateRepoRisk(item.SeverityCounts, item.VulnerabilityCount)
		}

		items = append(items, item)
	}

	return items, nil
}

func listEcosystemRepos(batchID string) ([]ecosystemRepoRow, error) {
	var rows []ecosystemRepoRow
	_, err := Client.
		From("ecosystem_repos").
		Select("id,batch_id,repo_name,stars,repo_rank", "", false).
		Eq("batch_id", batchID).
		Order("repo_rank", &postgrest.OrderOpts{Ascending: true}).
		ExecuteTo(&rows)
	return rows, err
}

func listEcosystemScans(batchID string) ([]ecosystemScanRow, error) {
	var rows []ecosystemScanRow
	_, err := Client.
		From("scan_jobs").
		Select("job_id,batch_id,repo_id,repo,status,created_at,commit_hash", "", false).
		Eq("batch_id", batchID).
		Order("created_at", &postgrest.OrderOpts{Ascending: false}).
		ExecuteTo(&rows)
	return rows, err
}

func sortEcosystemRepoItems(items []models.EcosystemRepoListItem, filter EcosystemRepoFilter) {
	ascending := strings.ToLower(strings.TrimSpace(filter.SortOrder)) == "asc"

	sort.Slice(items, func(i, j int) bool {
		left := items[i]
		right := items[j]

		var result bool
		switch strings.ToLower(strings.TrimSpace(filter.SortBy)) {
		case "repo":
			result = left.RepoName < right.RepoName
		case "stars":
			result = left.Stars < right.Stars
		case "status":
			result = left.ScanStatus < right.ScanStatus
		case "vulnerability_count":
			result = left.VulnerabilityCount < right.VulnerabilityCount
		case "top_severity":
			result = severityRank(left.TopSeverity) > severityRank(right.TopSeverity)
		default:
			result = left.Rank < right.Rank
		}

		if ascending {
			return result
		}
		return !result
	})
}

func estimateRepoRisk(breakdown models.SeverityBreakdown, total int64) float64 {
	if total == 0 {
		return 0
	}

	score := float64(
		breakdown.Critical*100 +
			breakdown.High*75 +
			breakdown.Medium*45 +
			breakdown.Low*20,
	) / float64(total)

	if score < 0 {
		return 0
	}
	if score > 100 {
		return 100
	}
	return score
}

func executeRangeTo(query *postgrest.FilterBuilder, from, to int, target interface{}) (int64, error) {
	data, count, err := query.
		Order("created_at", &postgrest.OrderOpts{Ascending: false}).
		Range(from, to, "").
		Execute()
	if err != nil {
		return 0, err
	}

	if err := json.Unmarshal(data, target); err != nil {
		return 0, err
	}

	return count, nil
}
