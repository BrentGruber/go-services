package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"

	a "cordy/internal/core/api"
	"cordy/internal/core/environments"
)

const (
	HealthEndpoint = "/health"

	ApiPath         = "/api/v1"
	EnvironmentPath = "/environments"
)

// ConfigureRouter instnatiates a new chi router with middleware and routes for the server
func ConfigureRouter(envSvc environments.EnvironmentService) *chi.Mux {
	log.Info().Msg("Configuring router")
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value cannot be ignored by any major browsers
	}))
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(a.Logging)

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("UP"))
	})

	r.Route(ApiPath, func(r chi.Router) {
		r.Route(EnvironmentPath, environments.NewEnvironmentApi(envSvc).ConfigureRouter)
	})

	return r
}
