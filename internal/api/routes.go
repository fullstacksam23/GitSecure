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
	router.Post("/scan", handler.StartScan)

	return router
}
