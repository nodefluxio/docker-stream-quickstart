package repository

import (
	"context"

	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
)

// Stream interface abstracts the repository layer and should be implemented in repository
type Stream interface {
	GetDetail(ctx context.Context, nodeNumber int64, streamID string) (*entity.StreamDetail, error)
}
