package environments

import (
	"time"
)

type Environment struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Host        string    `json:"host"`
	AwsAccount  string    `json:"aws_account"`
	AwsRegion   string    `json:"aws_region"`
	ClusterName string    `json:"cluster_name"`
	Domain      string    `json:"domain"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
