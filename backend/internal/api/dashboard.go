package api

import (
	"net/http"

	"github.com/fullstacksam23/GitSecure/internal/db"
)

func DashboardSummaryHandler(w http.ResponseWriter, r *http.Request) {
	summary, err := db.GetDashboardSummary()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load dashboard summary")
		return
	}

	writeJSON(w, http.StatusOK, summary)
}
