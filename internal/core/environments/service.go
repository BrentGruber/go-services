package environments

import (
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type service struct {
}

func NewService() *service {
	log.Info().Msg("Creating new environment service...")
	return &service{}
}

func (s *service) GetEnvironment(id string) (*Environment, error) {
	log.Info().Msg("Getting environment...")
	return &Environment{
		ID:          "123",
		Name:        "test",
		Host:        "test",
		AwsAccount:  "test",
		AwsRegion:   "test",
		ClusterName: "test",
		Domain:      "test",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

func (s *service) CreateEnvironment(env *Environment) (*Environment, error) {
	log.Info().Msg("Creating environment...")
	env.ID = uuid.New().String()
	env.CreatedAt = time.Now()
	env.UpdatedAt = time.Now()

	return env, nil
}

func (s *service) UpdateEnvironment(env *Environment) (*Environment, error) {
	log.Info().Msg("Updating environment...")
	env.UpdatedAt = time.Now()

	return env, nil
}

func (s *service) DeleteEnvironment(env *Environment) error {
	log.Info().Msg("Deleting environment...")
	return nil
}

func (s *service) ListEnvironments() ([]*Environment, error) {
	log.Info().Msg("Listing environments...")
	return []*Environment{}, nil
}
