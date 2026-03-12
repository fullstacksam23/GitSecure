package models

type UnifiedVuln struct {
	ID         string
	Package    string
	Version    string
	Severity   string
	Summary    string
	Urls       []string
	FixVersion []string
	Source     string
}
