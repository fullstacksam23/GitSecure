package services

import (
	"strings"

	"github.com/fullstacksam23/GitSecure/internal/models"
)

type VulnGraph struct {
	adj map[string]map[string]struct{}
}

func NewVulnGraph() *VulnGraph {
	return &VulnGraph{
		adj: make(map[string]map[string]struct{}),
	}
}

func (g *VulnGraph) AddEdge(a, b string) {

	if g.adj[a] == nil {
		g.adj[a] = map[string]struct{}{}
	}

	if g.adj[b] == nil {
		g.adj[b] = map[string]struct{}{}
	}

	g.adj[a][b] = struct{}{}
	g.adj[b][a] = struct{}{}
}

func BuildVulnGraph(advisories map[string]models.OSVAdvisory) *VulnGraph {

	graph := NewVulnGraph()

	for id, adv := range advisories {

		if graph.adj[id] == nil {
			graph.adj[id] = map[string]struct{}{}
		}

		for _, alias := range adv.Aliases {
			graph.AddEdge(id, alias)
		}

	}

	return graph
}
func chooseCanonical(ids []string) string {

	for _, id := range ids {
		if strings.HasPrefix(id, "CVE-") {
			return id
		}
	}

	for _, id := range ids {
		if strings.HasPrefix(id, "GHSA-") {
			return id
		}
	}

	return ids[0]
}

func (g *VulnGraph) CanonicalMap() map[string]string {

	visited := map[string]bool{}
	result := map[string]string{}

	for node := range g.adj {

		if visited[node] {
			continue
		}

		stack := []string{node}
		component := []string{}

		for len(stack) > 0 {

			n := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if visited[n] {
				continue
			}

			visited[n] = true
			component = append(component, n)

			for neighbor := range g.adj[n] {
				stack = append(stack, neighbor)
			}
		}

		canonical := chooseCanonical(component)

		for _, id := range component {
			result[id] = canonical
		}
	}

	return result
}

func CanonicalizeAdvisories(advisories map[string]models.OSVAdvisory, canonical map[string]string) map[string]models.OSVAdvisory {

	newMap := map[string]models.OSVAdvisory{}

	for id, adv := range advisories {

		canonicalID := id

		if c, ok := canonical[id]; ok {
			canonicalID = c
		}

		newMap[canonicalID] = adv
	}

	return newMap
}

func NormalizeGrype(grype models.GrypeResponse, canonical map[string]string, jobID string) []models.UnifiedVuln {

	var vulns []models.UnifiedVuln

	for _, match := range grype.Matches {

		id := match.Vulnerability.ID
		if canonicalID, ok := canonical[id]; ok && canonicalID != "" {
			id = canonicalID
		}

		// Extract match details safely
		matchType, constraint := pickBestMatch(match.MatchDetails)

		v := models.UnifiedVuln{
			ID:       id,
			JobID:    jobID,
			Package:  match.Artifact.Name,
			Version:  match.Artifact.Version,
			Severity: match.Vulnerability.Severity,
			Summary:  match.Vulnerability.Description,
			Urls:     match.Vulnerability.Urls,

			FixVersion: match.Vulnerability.Fix.Versions,
			FixState:   match.Vulnerability.Fix.State,

			Risk:      match.Vulnerability.Risk,
			Namespace: match.Vulnerability.Namespace,

			MatchType:  matchType,
			Constraint: constraint,

			DataSource: match.Vulnerability.DataSource,
			Source:     "grype",
		}
		vulns = append(vulns, v)
	}

	return vulns
}

func pickBestMatch(details []models.MatchDetail) (string, string) {
	bestType := ""
	constraint := ""

	for _, d := range details {
		if d.Type == "exact-direct-match" {
			return d.Type, d.Found.VersionConstraint
		}
		if d.Type == "exact-indirect-match" {
			bestType = d.Type
			constraint = d.Found.VersionConstraint
		}
	}

	if bestType != "" {
		return bestType, constraint
	}

	if len(details) > 0 {
		return details[0].Type, details[0].Found.VersionConstraint
	}

	return "", ""
}
