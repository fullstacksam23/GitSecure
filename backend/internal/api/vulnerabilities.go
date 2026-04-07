package api

import (
	"net/http"

	"github.com/fullstacksam23/GitSecure/internal/db"
	"github.com/fullstacksam23/GitSecure/internal/models"
	"github.com/go-chi/chi/v5"
)

func ListVulnerabilitiesHandler(w http.ResponseWriter, r *http.Request) {
	page := parsePositiveInt(r, "page", 1, 1, 10000)
	pageSize := parsePositiveInt(r, "page_size", 50, 1, 100)

	filter := db.VulnerabilityFilter{
		Page:      page,
		PageSize:  pageSize,
		Severity:  r.URL.Query().Get("severity"),
		Search:    r.URL.Query().Get("search"),
		JobID:     r.URL.Query().Get("job_id"),
		Ecosystem: r.URL.Query().Get("ecosystem"),
		FixState:  r.URL.Query().Get("fix_state"),
		SortBy:    r.URL.Query().Get("sort_by"),
		SortOrder: r.URL.Query().Get("sort_order"),
	}

	items, total, facets, err := db.ListVulnerabilities(filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load vulnerabilities")
		return
	}

	response := models.VulnerabilityListResponse{
		Items: items,
		Pagination: models.Pagination{
			Page:       page,
			PageSize:   pageSize,
			TotalItems: total,
			TotalPages: calculateTotalPages(total, pageSize),
		},
		Facets: facets,
	}

	writeJSON(w, http.StatusOK, response)
}

func GetVulnerabilityHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "vulnerability id is required")
		return
	}

	record, err := db.GetVulnerabilityByID(id, r.URL.Query().Get("job_id"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load vulnerability")
		return
	}
	if record == nil {
		writeError(w, http.StatusNotFound, "vulnerability not found")
		return
	}

	writeJSON(w, http.StatusOK, record)
}
