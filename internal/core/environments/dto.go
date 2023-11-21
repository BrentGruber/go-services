package environments

import (
	"net/http"

	"github.com/go-chi/render"
)

type CreateEnvironmentRequest struct {
	Environment
}

func (r *CreateEnvironmentRequest) Bind(_ *http.Request) error {
	// Put any validation logic here before response is marshalled

	return nil
}

type UpdateEnvironmentRequest struct {
	Name        string `json:"name"`
	Host        string `json:"host"`
	AwsAccount  string `json:"aws_account"`
	AwsRegion   string `json:"aws_region"`
	ClusterName string `json:"cluster_name"`
	Domain      string `json:"domain"`
}

func (r *UpdateEnvironmentRequest) Bind(_ *http.Request) error {
	// Put any validation logic here before response is marshalled

	return nil
}

type EnvironmentResponse struct {
	Environment
}

func NewEnvironmentResponse(env *Environment) *EnvironmentResponse {
	return &EnvironmentResponse{
		Environment: *env,
	}
}

func (rd *EnvironmentResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	// Put any pre-processing logic here before response is marshalled

	return nil
}

func NewEnvironmentListResponse(envs []*Environment) []render.Renderer {
	list := make([]render.Renderer, 0)
	for _, env := range envs {
		list = append(list, NewEnvironmentResponse(env))
	}
	return list
}
