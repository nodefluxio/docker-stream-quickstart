package repository

import (
	"context"

	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
)

// Seagate interface abstracts the repository layer and should be implemented in repository
type Seagate interface {
	GetToken(ctx context.Context, requestData *entity.SeagateGetTokenRequest) (*entity.SeagateGetTokenResult, error)
	GetFaceSearchResult(ctx context.Context, token string) (*entity.SeagateFaceSearchResult, error)
}
