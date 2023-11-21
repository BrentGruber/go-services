// The api package handles configuring routing for http and websocket requests into the server
// it validates the requests and passes them to the appropriate handler
package environments

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"github.com/rs/zerolog/log"

	a "cordy/internal/core/api"
)

const (
	CtxKeyEnvironment a.CtxKey = "environment"
)

type EnvironmentService interface {
	GetEnvironment(id string) (*Environment, error)
	CreateEnvironment(env *Environment) (*Environment, error)
	UpdateEnvironment(env *Environment) (*Environment, error)
	DeleteEnvironment(env *Environment) error
	ListEnvironments() ([]*Environment, error)
}

type EnvironmentApi struct {
	service EnvironmentService
}

func NewEnvironmentApi(service EnvironmentService) *EnvironmentApi {
	return &EnvironmentApi{
		service: service,
	}
}

func (api *EnvironmentApi) ConfigureRouter(r chi.Router) {
	r.Route("/", func(r chi.Router) {
		r.With(a.Paginate).Get("/", api.ListEnvironments)
		r.Post("/", api.CreateEnvironment)

		r.Route("/{id}", func(r chi.Router) {
			r.Use(api.EnvironmentCtx)
			r.Get("/", api.GetEnvironment)
			r.Put("/", api.UpdateEnvironment)
			r.Delete("/", api.DeleteEnvironment)
		})
	})
}

// List all environments
func (api *EnvironmentApi) ListEnvironments(w http.ResponseWriter, r *http.Request) {
	envs, err := api.service.ListEnvironments()
	if err != nil {
		render.Render(w, r, a.ErrInternalServer)
		return
	}

	render.JSON(w, r, envs)
}

// Get a single environment
func (api *EnvironmentApi) GetEnvironment(w http.ResponseWriter, r *http.Request) {
	env := r.Context().Value(CtxKeyEnvironment).(*Environment)

	render.JSON(w, r, env)
}

// Create a new environment
func (api *EnvironmentApi) CreateEnvironment(w http.ResponseWriter, r *http.Request) {
	data := &CreateEnvironmentRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, a.ErrInvalidRequest(err))
		return
	}

	env, err := api.service.CreateEnvironment(&data.Environment)

	if err != nil {
		if errors.Is(err, ErrNotFound) {
			render.Render(w, r, a.ErrNotFound)
		} else {
			log.Error().Err(err).Interface("data", data).Msg("Failed to create environment")
		}
		return
	}

	resp := &EnvironmentResponse{Environment: *env}
	render.Status(r, http.StatusCreated)
	a.Render(w, r, resp)
}

// Update an existing environment
func (api *EnvironmentApi) UpdateEnvironment(w http.ResponseWriter, r *http.Request) {
	env := r.Context().Value(CtxKeyEnvironment).(*Environment)

	data := &UpdateEnvironmentRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, a.ErrInvalidRequest(err))
		return
	}

	env.Name = data.Name
	env.Host = data.Host
	env.AwsAccount = data.AwsAccount
	env.AwsRegion = data.AwsRegion
	env.ClusterName = data.ClusterName
	env.Domain = data.Domain

	env, err := api.service.UpdateEnvironment(env)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			render.Render(w, r, a.ErrNotFound)
		} else {
			log.Error().Err(err).Interface("data", data).Msg("Failed to update environment")
		}
		return
	}

	resp := &EnvironmentResponse{Environment: *env}
	a.Render(w, r, resp)
}

// Delete an existing environment
func (api *EnvironmentApi) DeleteEnvironment(w http.ResponseWriter, r *http.Request) {
	env := r.Context().Value(CtxKeyEnvironment).(*Environment)

	err := api.service.DeleteEnvironment(env)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			render.Render(w, r, a.ErrNotFound)
		} else {
			log.Error().Err(err).Interface("data", env).Msg("Failed to delete environment")
		}
		return
	}

	render.Status(r, http.StatusNoContent)
}

// EnvironmentCtx middleware is used to load an Environment object from
// the URL parameters passed through as the request. In case
// the Environment could not be found, we stop here and return a 404.
func (api *EnvironmentApi) EnvironmentCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var environment *Environment
		var err error

		id := chi.URLParam(r, "id")
		if id == "" {
			render.Render(w, r, a.ErrInvalidRequest(errors.New("environment id is required")))
			return
		}

		environment, err = api.service.GetEnvironment(id)

		if err != nil {
			if errors.Is(err, ErrNotFound) {
				render.Render(w, r, a.ErrNotFound)
			} else {
				log.Error().Err(err).Str("id", id).Msg("Failed to get environment")
			}
			return
		}

		ctx := context.WithValue(r.Context(), CtxKeyEnvironment, environment)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
