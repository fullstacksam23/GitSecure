package api

import (
	"github.com/fullstacksam23/GitSecure/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func loadRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/health", handler.HealthHandler)
	router.Route("/sbom", loadSbomRoutes)

	return router
}

func loadSbomRoutes(r chi.Router) {
	// - /dependencies?repo=username/repoName
	r.Get("/dependencies", handler.ExtractDependencies)
}
