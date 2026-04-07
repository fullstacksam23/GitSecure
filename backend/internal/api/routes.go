package api

import (
	"github.com/fullstacksam23/GitSecure/internal/scanner"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func loadRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(corsMiddleware)

	router.Get("/health", HealthHandler)
	router.Post("/scan", scanner.StartScan)
	router.Get("/dashboard/summary", DashboardSummaryHandler)
	router.Get("/scans/compare", CompareScansHandler)
	router.Get("/scans", ListScansHandler)
	router.Get("/scans/{jobId}", GetScanHandler)
	router.Get("/vulnerabilities", ListVulnerabilitiesHandler)
	router.Get("/vulnerabilities/{id}", GetVulnerabilityHandler)

	return router
}
