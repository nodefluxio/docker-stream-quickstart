package repository

import (
	"context"

	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
)

// LatestTimestamp is repository for table latest_timestamp
type LatestTimestamp interface {
	CreateOrUpdate(ctx context.Context, data *entity.LatestTimestamp) error
	Get(ctx context.Context) (*entity.LatestTimestamp, error)
}
