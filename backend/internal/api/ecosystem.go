package api

import (
	"net/http"

	"github.com/fullstacksam23/GitSecure/internal/db"
	"github.com/fullstacksam23/GitSecure/internal/models"
	"github.com/go-chi/chi/v5"
)

func ListEcosystemBatchesHandler(w http.ResponseWriter, r *http.Request) {
	page := parsePositiveInt(r, "page", 1, 1, 10000)
	pageSize := parsePositiveInt(r, "page_size", 20, 1, 100)

	items, total, err := db.ListEcosystemBatches(db.EcosystemBatchFilter{
		Page:     page,
		PageSize: pageSize,
		Status:   r.URL.Query().Get("status"),
		Language: r.URL.Query().Get("language"),
		Search:   r.URL.Query().Get("search"),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load ecosystem batches")
		return
	}

	writeJSON(w, http.StatusOK, models.EcosystemBatchListResponse{
		Items: items,
		Pagination: models.Pagination{
			Page:       page,
			PageSize:   pageSize,
			TotalItems: total,
			TotalPages: calculateTotalPages(total, pageSize),
		},
	})
}

func GetEcosystemBatchHandler(w http.ResponseWriter, r *http.Request) {
	batchID := chi.URLParam(r, "batchId")
	if batchID == "" {
		writeError(w, http.StatusBadRequest, "batchId is required")
		return
	}

	batch, err := db.GetEcosystemBatchByID(batchID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load ecosystem batch")
		return
	}
	if batch == nil {
		writeError(w, http.StatusNotFound, "ecosystem batch not found")
		return
	}

	writeJSON(w, http.StatusOK, batch)
}

func GetEcosystemBatchSummaryHandler(w http.ResponseWriter, r *http.Request) {
	batchID := chi.URLParam(r, "batchId")
	if batchID == "" {
		writeError(w, http.StatusBadRequest, "batchId is required")
		return
	}

	summary, err := db.GetEcosystemBatchSummary(batchID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load batch summary")
		return
	}

	writeJSON(w, http.StatusOK, summary)
}

func ListEcosystemBatchReposHandler(w http.ResponseWriter, r *http.Request) {
	batchID := chi.URLParam(r, "batchId")
	if batchID == "" {
		writeError(w, http.StatusBadRequest, "batchId is required")
		return
	}

	page := parsePositiveInt(r, "page", 1, 1, 10000)
	pageSize := parsePositiveInt(r, "page_size", 20, 1, 100)

	items, total, err := db.ListEcosystemRepos(batchID, db.EcosystemRepoFilter{
		Page:      page,
		PageSize:  pageSize,
		Search:    r.URL.Query().Get("search"),
		Status:    r.URL.Query().Get("status"),
		Severity:  r.URL.Query().Get("severity"),
		SortBy:    r.URL.Query().Get("sort_by"),
		SortOrder: r.URL.Query().Get("sort_order"),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load batch repositories")
		return
	}

	writeJSON(w, http.StatusOK, models.EcosystemRepoListResponse{
		Items: items,
		Pagination: models.Pagination{
			Page:       page,
			PageSize:   pageSize,
			TotalItems: total,
			TotalPages: calculateTotalPages(total, pageSize),
		},
	})
}
