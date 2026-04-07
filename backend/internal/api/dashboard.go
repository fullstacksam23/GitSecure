package api

import (
	"fmt"
	"net/http"

	"github.com/fullstacksam23/GitSecure/internal/db"
)

func DashboardSummaryHandler(w http.ResponseWriter, r *http.Request) {
	summary, err := db.GetDashboardSummary()
	if err != nil {
		errorString := fmt.Sprintf("failed to load dashboard summary - Error: %s", err)
		writeError(w, http.StatusInternalServerError, errorString)
		return
	}

	writeJSON(w, http.StatusOK, summary)
}
