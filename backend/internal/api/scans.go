package api

import (
	"net/http"
	"strconv"

	"github.com/fullstacksam23/GitSecure/internal/db"
	"github.com/fullstacksam23/GitSecure/internal/models"
	"github.com/go-chi/chi/v5"
)

func ListScansHandler(w http.ResponseWriter, r *http.Request) {
	page := parsePositiveInt(r, "page", 1, 1, 10000)
	pageSize := parsePositiveInt(r, "page_size", 20, 1, 100)

	items, total, err := db.ListScans(db.ScanListFilter{
		Page:     page,
		PageSize: pageSize,
		Repo:     r.URL.Query().Get("repo"),
		Status:   r.URL.Query().Get("status"),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load scans")
		return
	}

	writeJSON(w, http.StatusOK, models.ScanListResponse{
		Items: items,
		Pagination: models.Pagination{
			Page:       page,
			PageSize:   pageSize,
			TotalItems: total,
			TotalPages: calculateTotalPages(total, pageSize),
		},
	})
}

func GetScanHandler(w http.ResponseWriter, r *http.Request) {
	jobID := chi.URLParam(r, "jobId")
	if jobID == "" {
		writeError(w, http.StatusBadRequest, "jobId is required")
		return
	}

	scan, err := db.GetScanByID(jobID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load scan")
		return
	}
	if scan == nil {
		writeError(w, http.StatusNotFound, "scan not found")
		return
	}

	writeJSON(w, http.StatusOK, scan)
}

func CompareScansHandler(w http.ResponseWriter, r *http.Request) {
	baseJobID := r.URL.Query().Get("base")
	targetJobID := r.URL.Query().Get("target")

	if baseJobID == "" || targetJobID == "" {
		writeError(w, http.StatusBadRequest, "base and target query params are required")
		return
	}

	comparison, err := db.CompareScans(baseJobID, targetJobID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to compare scans")
		return
	}
	if comparison == nil {
		writeError(w, http.StatusNotFound, "one or both scans were not found")
		return
	}

	writeJSON(w, http.StatusOK, comparison)
}

func parsePositiveInt(r *http.Request, key string, fallback, min, max int) int {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(raw)
	if err != nil || parsed < min {
		return fallback
	}
	if parsed > max {
		return max
	}
	return parsed
}

func calculateTotalPages(total int64, pageSize int) int {
	if total == 0 || pageSize <= 0 {
		return 0
	}
	return int((total + int64(pageSize) - 1) / int64(pageSize))
}
