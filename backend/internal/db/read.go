package db

import (
	"encoding/json"
	"errors"
	"math"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/fullstacksam23/GitSecure/internal/models"
	"github.com/supabase-community/postgrest-go"
)

type VulnerabilityFilter struct {
	Page      int
	PageSize  int
	Severity  string
	Search    string
	JobID     string
	Ecosystem string
	FixState  string
	SortBy    string
	SortOrder string
}

type ScanListFilter struct {
	Page     int
	PageSize int
	Repo     string
	Status   string
}

type scanRow struct {
	JobID      string  `json:"job_id"`
	Repo       string  `json:"repo"`
	Status     string  `json:"status"`
	CommitHash *string `json:"commit_hash"`
	CreatedAt  string  `json:"created_at"`
}

type scanDetailRow struct {
	JobID      string  `json:"job_id"`
	Repo       string  `json:"repo"`
	Status     string  `json:"status"`
	CommitHash *string `json:"commit_hash"`
	CreatedAt  string  `json:"created_at"`
}

var dashboardSummaryCache = struct {
	mu        sync.RWMutex
	value     models.DashboardSummary
	expiresAt time.Time
}{}

func GetDashboardSummary() (models.DashboardSummary, error) {
	if Client == nil {
		return models.DashboardSummary{}, errors.New("client not initialized")
	}

	if cached, ok := getCachedDashboardSummary(); ok {
		return cached, nil
	}

	summary := models.DashboardSummary{}

	totalScans, err := countRecords("scan_jobs", "")
	if err != nil {
		return summary, err
	}
	summary.TotalScans = totalScans

	totalVulns, err := countRecords("vulnerabilities", "")
	if err != nil {
		return summary, err
	}
	summary.TotalVulnerabilities = totalVulns

	breakdown, err := getSeverityBreakdown("")
	if err != nil {
		return summary, err
	}
	summary.SeverityDistribution = breakdown
	summary.Critical = breakdown.Critical
	summary.High = breakdown.High
	summary.Medium = breakdown.Medium
	summary.Low = breakdown.Low

	packagesFixed, err := countFixedPackages()
	if err != nil {
		return summary, err
	}
	summary.PackagesFixed = packagesFixed

	recent, _, err := ListScans(ScanListFilter{Page: 1, PageSize: 6})
	if err != nil {
		return summary, err
	}
	summary.RecentScans = recent

	topPackages, err := listTopRiskPackages("", 6)
	if err != nil {
		return summary, err
	}
	summary.TopRiskPackages = topPackages

	repos, err := listRepoSummaries(6)
	if err != nil {
		return summary, err
	}
	summary.RepoSummaries = repos

	trend, err := getRiskTrend(7)
	if err != nil {
		return summary, err
	}
	summary.RiskTrend = trend

	cacheDashboardSummary(summary)

	return summary, nil
}

func ListScans(filter ScanListFilter) ([]models.ScanListItem, int64, error) {
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
		From("scan_jobs").
		Select("job_id,repo,status,commit_hash,created_at", "exact", false)

	if repo := strings.TrimSpace(filter.Repo); repo != "" {
		query = query.Ilike("repo", "%"+escapeLike(repo)+"%")
	}
	if status := strings.TrimSpace(filter.Status); status != "" {
		query = query.Eq("status", status)
	}

	from := (filter.Page - 1) * filter.PageSize
	to := from + filter.PageSize - 1

	data, count, err := query.
		Order("created_at", &postgrest.OrderOpts{Ascending: false}).
		Range(from, to, "").
		Execute()
	if err != nil {
		return nil, 0, err
	}
	var rawRows []scanRow
	if err := json.Unmarshal(data, &rawRows); err != nil {
		return nil, 0, err
	}
	rows := make([]models.ScanListItem, 0, len(rawRows))
	for _, row := range rawRows {
		rows = append(rows, models.ScanListItem{
			JobID:      row.JobID,
			Repo:       row.Repo,
			Status:     row.Status,
			CommitHash: stringOrEmpty(row.CommitHash),
			CreatedAt:  row.CreatedAt,
		})
	}

	if len(rows) == 0 {
		return []models.ScanListItem{}, count, nil
	}

	jobIDs := make([]string, 0, len(rows))
	for _, row := range rows {
		jobIDs = append(jobIDs, row.JobID)
	}

	aggregates, err := aggregateScanVulnerabilities(jobIDs)
	if err != nil {
		return nil, 0, err
	}

	for index := range rows {
		rows[index].TopSeverity = "unknown"
		if aggregate, ok := aggregates[rows[index].JobID]; ok {
			rows[index].TopSeverity = aggregate.TopSeverity
			rows[index].VulnerabilityCount = aggregate.Count
			rows[index].SeverityCounts = aggregate.SeverityCounts
		}
	}

	return rows, count, nil
}

func GetScanByID(jobID string) (*models.ScanDetails, error) {
	if Client == nil {
		return nil, errors.New("client not initialized")
	}

	var baseRows []scanDetailRow
	_, err := Client.
		From("scan_jobs").
		Select("job_id,repo,status,commit_hash,created_at", "", false).
		Eq("job_id", jobID).
		Limit(1, "").
		ExecuteTo(&baseRows)
	if err != nil {
		return nil, err
	}
	if len(baseRows) == 0 {
		return nil, nil
	}

	scan := models.ScanDetails{
		JobID:      baseRows[0].JobID,
		Repo:       baseRows[0].Repo,
		Status:     baseRows[0].Status,
		CommitHash: stringOrEmpty(baseRows[0].CommitHash),
		CreatedAt:  baseRows[0].CreatedAt,
	}
	vulns, err := listVulnerabilitiesForJob(jobID)
	if err != nil {
		return nil, err
	}

	scan.VulnerabilityCount = int64(len(vulns))
	scan.SeverityCounts = summarizeSeverity(vulns)
	scan.Ecosystems = summarizeFacet(vulns, func(item models.VulnerabilityRecord) string { return item.Ecosystem })
	scan.FixStates = summarizeFacet(vulns, func(item models.VulnerabilityRecord) string { return item.FixState })
	scan.TopPackages = summarizePackages(vulns, 5)

	return &scan, nil
}

func ListVulnerabilities(filter VulnerabilityFilter) ([]models.VulnerabilityRecord, int64, models.VulnerabilityFacets, error) {
	if Client == nil {
		return nil, 0, models.VulnerabilityFacets{}, errors.New("client not initialized")
	}

	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 50
	}

	query := Client.
		From("vulnerability_records").
		Select("id,job_id,package,version,severity,normalized_severity,summary,urls,fix_version,fix_state,risk,namespace,match_type,version_constraint,data_source,source,cwe_ids,ecosystem,created_at", "exact", false)

	query = applyVulnerabilityFilters(query, filter)
	query = applyVulnerabilitySorting(query, filter)

	from := (filter.Page - 1) * filter.PageSize
	to := from + filter.PageSize - 1

	data, count, err := query.Range(from, to, "").Execute()
	if err != nil {
		return nil, 0, models.VulnerabilityFacets{}, err
	}

	var items []models.VulnerabilityRecord
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, 0, models.VulnerabilityFacets{}, err
	}

	facets, err := listVulnerabilityFacets(filter)
	if err != nil {
		return nil, 0, models.VulnerabilityFacets{}, err
	}

	return items, count, facets, nil
}

func GetVulnerabilityByID(id, jobID string) (*models.VulnerabilityRecord, error) {
	if Client == nil {
		return nil, errors.New("client not initialized")
	}

	query := Client.
		From("vulnerability_records").
		Select("id,job_id,package,version,severity,normalized_severity,summary,urls,fix_version,fix_state,risk,namespace,match_type,version_constraint,data_source,source,cwe_ids,ecosystem,created_at", "", false).
		Eq("id", id)

	if strings.TrimSpace(jobID) != "" {
		query = query.Eq("job_id", jobID)
	}

	var items []models.VulnerabilityRecord
	_, err := query.Limit(1, "").ExecuteTo(&items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, nil
	}

	return &items[0], nil
}

func CompareScans(baseJobID, targetJobID string) (*models.ScanCompareResponse, error) {
	baseScan, err := GetScanByID(baseJobID)
	if err != nil || baseScan == nil {
		return nil, err
	}
	targetScan, err := GetScanByID(targetJobID)
	if err != nil || targetScan == nil {
		return nil, err
	}

	baseVulns, err := listVulnerabilitiesForJob(baseJobID)
	if err != nil {
		return nil, err
	}
	targetVulns, err := listVulnerabilitiesForJob(targetJobID)
	if err != nil {
		return nil, err
	}

	baseMap := make(map[string]models.VulnerabilityRecord, len(baseVulns))
	for _, item := range baseVulns {
		baseMap[compareKey(item)] = item
	}

	targetMap := make(map[string]models.VulnerabilityRecord, len(targetVulns))
	for _, item := range targetVulns {
		targetMap[compareKey(item)] = item
	}

	newItems := make([]models.VulnerabilityRecord, 0)
	fixedItems := make([]models.VulnerabilityRecord, 0)
	persistingItems := make([]models.VulnerabilityRecord, 0)

	for key, item := range targetMap {
		if _, ok := baseMap[key]; ok {
			persistingItems = append(persistingItems, item)
			continue
		}
		newItems = append(newItems, item)
	}
	for key, item := range baseMap {
		if _, ok := targetMap[key]; ok {
			continue
		}
		fixedItems = append(fixedItems, item)
	}

	sortVulnerabilities(newItems)
	sortVulnerabilities(fixedItems)
	sortVulnerabilities(persistingItems)

	response := &models.ScanCompareResponse{
		BaseScan:   baseScan,
		TargetScan: targetScan,
		New: models.CompareBucket{
			Count: int64(len(newItems)),
			Items: truncateVulnerabilities(newItems, 25),
		},
		Fixed: models.CompareBucket{
			Count: int64(len(fixedItems)),
			Items: truncateVulnerabilities(fixedItems, 25),
		},
		Persisting: models.CompareBucket{
			Count: int64(len(persistingItems)),
			Items: truncateVulnerabilities(persistingItems, 25),
		},
		NewSeverity:   summarizeSeverity(newItems),
		FixedSeverity: summarizeSeverity(fixedItems),
	}

	return response, nil
}

func applyVulnerabilityFilters(query *postgrest.FilterBuilder, filter VulnerabilityFilter) *postgrest.FilterBuilder {
	if jobID := strings.TrimSpace(filter.JobID); jobID != "" {
		query = query.Eq("job_id", jobID)
	}

	if severities := normalizeSeverityFilters(filter.Severity); len(severities) > 0 {
		parts := make([]string, 0, len(severities))
		for _, severity := range severities {
			parts = append(parts, "normalized_severity.eq."+severity)
		}
		query = query.Or(strings.Join(parts, ","), "")
	}

	if search := strings.TrimSpace(filter.Search); search != "" {
		pattern := "%" + escapeLike(search) + "%"
		query = query.Or(
			strings.Join([]string{
				"id.ilike." + pattern,
				"package.ilike." + pattern,
				"summary.ilike." + pattern,
			}, ","),
			"",
		)
	}

	if ecosystem := strings.TrimSpace(filter.Ecosystem); ecosystem != "" {
		query = query.Eq("ecosystem", ecosystem)
	}
	if fixState := strings.TrimSpace(filter.FixState); fixState != "" {
		query = query.Eq("fix_state", fixState)
	}

	return query
}

func applyVulnerabilitySorting(query *postgrest.FilterBuilder, filter VulnerabilityFilter) *postgrest.FilterBuilder {
	ascending := strings.ToLower(strings.TrimSpace(filter.SortOrder)) == "asc"

	switch strings.ToLower(strings.TrimSpace(filter.SortBy)) {
	case "package":
		return query.Order("package", &postgrest.OrderOpts{Ascending: ascending})
	case "risk":
		return query.Order("risk", &postgrest.OrderOpts{Ascending: ascending, NullsFirst: !ascending})
	case "severity":
		return query.Order("normalized_severity", &postgrest.OrderOpts{Ascending: ascending})
	case "fix_state":
		return query.Order("fix_state", &postgrest.OrderOpts{Ascending: ascending})
	case "ecosystem":
		return query.Order("ecosystem", &postgrest.OrderOpts{Ascending: ascending})
	default:
		return query.Order("created_at", &postgrest.OrderOpts{Ascending: false})
	}
}

func listVulnerabilityFacets(filter VulnerabilityFilter) (models.VulnerabilityFacets, error) {
	query := Client.
		From("vulnerability_records").
		Select("ecosystem,fix_state", "", false)

	query = applyVulnerabilityFilters(query, VulnerabilityFilter{
		JobID:     filter.JobID,
		Severity:  filter.Severity,
		Search:    filter.Search,
		Ecosystem: "",
		FixState:  "",
	})

	var rows []struct {
		Ecosystem string `json:"ecosystem"`
		FixState  string `json:"fix_state"`
	}

	_, err := query.ExecuteTo(&rows)
	if err != nil {
		return models.VulnerabilityFacets{}, err
	}

	ecosystems := map[string]int64{}
	fixStates := map[string]int64{}
	for _, row := range rows {
		if value := strings.TrimSpace(row.Ecosystem); value != "" {
			ecosystems[value]++
		}
		if value := strings.TrimSpace(row.FixState); value != "" {
			fixStates[value]++
		}
	}

	return models.VulnerabilityFacets{
		Ecosystems: mapToFacets(ecosystems),
		FixStates:  mapToFacets(fixStates),
	}, nil
}

func listVulnerabilitiesForJob(jobID string) ([]models.VulnerabilityRecord, error) {
	var items []models.VulnerabilityRecord
	_, err := Client.
		From("vulnerability_records").
		Select("id,job_id,package,version,severity,normalized_severity,summary,urls,fix_version,fix_state,risk,namespace,match_type,version_constraint,data_source,source,cwe_ids,ecosystem,created_at", "", false).
		Eq("job_id", jobID).
		ExecuteTo(&items)
	if err != nil {
		return nil, err
	}
	return items, nil
}

type scanAggregate struct {
	TopSeverity    string
	Count          int64
	SeverityCounts models.SeverityBreakdown
}

func aggregateScanVulnerabilities(jobIDs []string) (map[string]scanAggregate, error) {
	if len(jobIDs) == 0 {
		return map[string]scanAggregate{}, nil
	}

	var rows []struct {
		JobID              string `json:"job_id"`
		NormalizedSeverity string `json:"normalized_severity"`
	}

	_, err := Client.
		From("vulnerability_records").
		Select("job_id,normalized_severity", "", false).
		In("job_id", jobIDs).
		ExecuteTo(&rows)
	if err != nil {
		return nil, err
	}

	result := make(map[string]scanAggregate, len(jobIDs))
	for _, jobID := range jobIDs {
		result[jobID] = scanAggregate{TopSeverity: "unknown"}
	}

	for _, row := range rows {
		current := result[row.JobID]
		current.Count++
		current.SeverityCounts = incrementSeverity(current.SeverityCounts, row.NormalizedSeverity)
		if severityRank(row.NormalizedSeverity) < severityRank(current.TopSeverity) {
			current.TopSeverity = row.NormalizedSeverity
		}
		result[row.JobID] = current
	}

	return result, nil
}

func countFixedPackages() (int64, error) {
	var rows []struct {
		Package    string   `json:"package"`
		FixVersion []string `json:"fix_version"`
		FixState   string   `json:"fix_state"`
	}

	_, err := Client.
		From("vulnerability_records").
		Select("package,fix_version,fix_state", "", false).
		ExecuteTo(&rows)
	if err != nil {
		return 0, err
	}

	set := map[string]struct{}{}
	for _, row := range rows {
		if len(row.FixVersion) > 0 || strings.Contains(strings.ToLower(row.FixState), "fixed") {
			if pkg := strings.TrimSpace(row.Package); pkg != "" {
				set[pkg] = struct{}{}
			}
		}
	}

	return int64(len(set)), nil
}

func listTopRiskPackages(jobID string, limit int) ([]models.PackageRiskItem, error) {
	query := Client.
		From("vulnerability_records").
		Select("package,ecosystem,risk", "", false)
	if strings.TrimSpace(jobID) != "" {
		query = query.Eq("job_id", jobID)
	}

	var rows []struct {
		Package   string  `json:"package"`
		Ecosystem string  `json:"ecosystem"`
		Risk      float64 `json:"risk"`
	}

	_, err := query.ExecuteTo(&rows)
	if err != nil {
		return nil, err
	}

	grouped := map[string]models.PackageRiskItem{}
	for _, row := range rows {
		key := strings.TrimSpace(row.Package) + "|" + strings.TrimSpace(row.Ecosystem)
		item := grouped[key]
		item.Package = row.Package
		item.Ecosystem = row.Ecosystem
		item.VulnerabilityCount++
		if row.Risk > item.Risk {
			item.Risk = row.Risk
		}
		grouped[key] = item
	}

	items := make([]models.PackageRiskItem, 0, len(grouped))
	for _, item := range grouped {
		items = append(items, item)
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].Risk == items[j].Risk {
			return items[i].VulnerabilityCount > items[j].VulnerabilityCount
		}
		return items[i].Risk > items[j].Risk
	})

	if limit > 0 && len(items) > limit {
		items = items[:limit]
	}

	return items, nil
}

func listRepoSummaries(limit int) ([]models.RepoSummary, error) {
	var scans []scanRow
	_, err := Client.
		From("scan_jobs").
		Select("job_id,repo,status,commit_hash,created_at", "", false).
		Order("created_at", &postgrest.OrderOpts{Ascending: false}).
		Limit(max(limit*6, 24), "").
		ExecuteTo(&scans)
	if err != nil {
		return nil, err
	}

	repoLatest := map[string]models.RepoSummary{}
	jobIDs := make([]string, 0, len(scans))
	for _, scan := range scans {
		jobIDs = append(jobIDs, scan.JobID)
		if _, ok := repoLatest[scan.Repo]; ok {
			continue
		}
		repoLatest[scan.Repo] = models.RepoSummary{
			Repo:        scan.Repo,
			Status:      scan.Status,
			LastJobID:   scan.JobID,
			LastScanned: scan.CreatedAt,
			TopSeverity: "unknown",
		}
	}

	aggregates, err := aggregateScanVulnerabilities(jobIDs)
	if err != nil {
		return nil, err
	}

	result := make([]models.RepoSummary, 0, len(repoLatest))
	for _, item := range repoLatest {
		if aggregate, ok := aggregates[item.LastJobID]; ok {
			item.TopSeverity = aggregate.TopSeverity
		}
		result = append(result, item)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].LastScanned > result[j].LastScanned
	})
	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}

	return result, nil
}

func getRiskTrend(days int) ([]models.TrendPoint, error) {
	now := time.Now().UTC()

	start := time.Date(
		now.Year(), now.Month(), now.Day(),
		0, 0, 0, 0, time.UTC,
	).AddDate(0, 0, -(days - 1))

	startTime := start.Format(time.RFC3339)

	var scanRows []struct {
		CreatedAt string `json:"created_at"`
	}

	_, err := Client.
		From("scan_jobs").
		Select("created_at", "", false).
		Gte("created_at", startTime).
		Order("created_at", &postgrest.OrderOpts{Ascending: true}).
		ExecuteTo(&scanRows)

	if err != nil {
		return nil, err
	}

	var vulnRows []struct {
		CreatedAt string `json:"created_at"`
	}

	_, err = Client.
		From("vulnerabilities").
		Select("created_at", "", false).
		Gte("created_at", startTime).
		Order("created_at", &postgrest.OrderOpts{Ascending: true}).
		ExecuteTo(&vulnRows)

	if err != nil {
		return nil, err
	}

	points := make([]models.TrendPoint, 0, days)
	indexByDate := make(map[string]int)

	for i := 0; i < days; i++ {
		day := start.AddDate(0, 0, i).Format("2006-01-02")

		indexByDate[day] = len(points)
		points = append(points, models.TrendPoint{
			Date: day,
		})
	}

	for _, row := range scanRows {
		createdAt, err := parseTimestamp(row.CreatedAt)
		if err != nil {
			continue
		}
		day := createdAt.UTC().Format("2006-01-02")

		if index, ok := indexByDate[day]; ok {
			points[index].Scans++
		}
	}

	for _, row := range vulnRows {
		createdAt, err := parseTimestamp(row.CreatedAt)
		if err != nil {
			continue
		}
		day := createdAt.UTC().Format("2006-01-02")

		if index, ok := indexByDate[day]; ok {
			points[index].Vulnerabilities++
		}
	}

	return points, nil
}

func getSeverityBreakdown(jobID string) (models.SeverityBreakdown, error) {
	filter := VulnerabilityFilter{JobID: jobID, Page: 1, PageSize: 1}
	query := Client.
		From("vulnerability_records").
		Select("normalized_severity", "", false)
	query = applyVulnerabilityFilters(query, filter)

	var rows []struct {
		NormalizedSeverity string `json:"normalized_severity"`
	}
	_, err := query.ExecuteTo(&rows)
	if err != nil {
		return models.SeverityBreakdown{}, err
	}

	breakdown := models.SeverityBreakdown{}
	for _, row := range rows {
		breakdown = incrementSeverity(breakdown, row.NormalizedSeverity)
	}
	return breakdown, nil
}

func compareKey(item models.VulnerabilityRecord) string {
	return strings.Join([]string{
		strings.TrimSpace(item.ID),
		strings.TrimSpace(item.Package),
		strings.TrimSpace(item.Ecosystem),
	}, "|")
}

func summarizeSeverity(items []models.VulnerabilityRecord) models.SeverityBreakdown {
	result := models.SeverityBreakdown{}
	for _, item := range items {
		result = incrementSeverity(result, item.NormalizedSeverity)
	}
	return result
}

func incrementSeverity(current models.SeverityBreakdown, value string) models.SeverityBreakdown {
	switch normalizeSeverityValue(value) {
	case "critical":
		current.Critical++
	case "high":
		current.High++
	case "medium":
		current.Medium++
	case "low":
		current.Low++
	default:
		current.Unknown++
	}
	return current
}

func summarizeFacet(items []models.VulnerabilityRecord, selector func(models.VulnerabilityRecord) string) []models.FacetCount {
	counts := map[string]int64{}
	for _, item := range items {
		if value := strings.TrimSpace(selector(item)); value != "" {
			counts[value]++
		}
	}
	return mapToFacets(counts)
}

func summarizePackages(items []models.VulnerabilityRecord, limit int) []models.PackageRiskItem {
	grouped := map[string]models.PackageRiskItem{}
	for _, item := range items {
		key := item.Package + "|" + item.Ecosystem
		current := grouped[key]
		current.Package = item.Package
		current.Ecosystem = item.Ecosystem
		current.VulnerabilityCount++
		if item.Risk > current.Risk {
			current.Risk = item.Risk
		}
		grouped[key] = current
	}

	result := make([]models.PackageRiskItem, 0, len(grouped))
	for _, item := range grouped {
		result = append(result, item)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Risk == result[j].Risk {
			return result[i].VulnerabilityCount > result[j].VulnerabilityCount
		}
		return result[i].Risk > result[j].Risk
	})
	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}
	return result
}

func mapToFacets(values map[string]int64) []models.FacetCount {
	items := make([]models.FacetCount, 0, len(values))
	for key, count := range values {
		items = append(items, models.FacetCount{Value: key, Count: count})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Count == items[j].Count {
			return items[i].Value < items[j].Value
		}
		return items[i].Count > items[j].Count
	})
	return items
}

func normalizeSeverityFilters(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	seen := map[string]struct{}{}
	for _, part := range parts {
		normalized := normalizeSeverityValue(part)
		if normalized == "unknown" {
			continue
		}
		if _, ok := seen[normalized]; ok {
			continue
		}
		seen[normalized] = struct{}{}
		out = append(out, normalized)
	}
	return out
}

func normalizeSeverityValue(value string) string {
	normalized := strings.ToLower(strings.TrimSpace(value))
	switch {
	case strings.Contains(normalized, "critical"):
		return "critical"
	case strings.Contains(normalized, "high"):
		return "high"
	case strings.Contains(normalized, "medium"), strings.Contains(normalized, "moderate"):
		return "medium"
	case strings.Contains(normalized, "low"), strings.Contains(normalized, "negligible"):
		return "low"
	default:
		return "unknown"
	}
}

func severityRank(value string) int {
	switch normalizeSeverityValue(value) {
	case "critical":
		return 1
	case "high":
		return 2
	case "medium":
		return 3
	case "low":
		return 4
	default:
		return 5
	}
}

func countRecords(table string, jobID string) (int64, error) {
	query := Client.From(table).Select("id", "exact", true)
	if table == "scan_jobs" {
		query = Client.From(table).Select("job_id", "exact", true)
	}
	if strings.TrimSpace(jobID) != "" {
		query = query.Eq("job_id", jobID)
	}
	_, count, err := query.Execute()
	return count, err
}

func truncateVulnerabilities(items []models.VulnerabilityRecord, limit int) []models.VulnerabilityRecord {
	if limit <= 0 || len(items) <= limit {
		return items
	}
	return items[:limit]
}

func sortVulnerabilities(items []models.VulnerabilityRecord) {
	sort.Slice(items, func(i, j int) bool {
		leftRank := severityRank(items[i].NormalizedSeverity)
		rightRank := severityRank(items[j].NormalizedSeverity)
		if leftRank == rightRank {
			if items[i].Risk == items[j].Risk {
				return items[i].Package < items[j].Package
			}
			return items[i].Risk > items[j].Risk
		}
		return leftRank < rightRank
	})
}

func escapeLike(value string) string {
	replacer := strings.NewReplacer(",", `\,`, "(", `\(`, ")", `\)`)
	return replacer.Replace(value)
}

func totalPages(total int64, pageSize int) int {
	if total == 0 || pageSize <= 0 {
		return 0
	}
	return int(math.Ceil(float64(total) / float64(pageSize)))
}

func getCachedDashboardSummary() (models.DashboardSummary, bool) {
	dashboardSummaryCache.mu.RLock()
	defer dashboardSummaryCache.mu.RUnlock()

	if time.Now().Before(dashboardSummaryCache.expiresAt) {
		return dashboardSummaryCache.value, true
	}

	return models.DashboardSummary{}, false
}

func cacheDashboardSummary(summary models.DashboardSummary) {
	dashboardSummaryCache.mu.Lock()
	defer dashboardSummaryCache.mu.Unlock()

	dashboardSummaryCache.value = summary
	dashboardSummaryCache.expiresAt = time.Now().Add(15 * time.Second)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func stringOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func parseTimestamp(value string) (time.Time, error) {
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02T15:04:05.999999999",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05.999999999",
		"2006-01-02 15:04:05",
	}

	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, errors.New("unable to parse timestamp: " + value)
}
