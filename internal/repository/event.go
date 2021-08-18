package repository

import (
	"context"
	"time"

	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/pkg/util"
)

// Event interface abstracts the repository layer and should be implemented in repository
type Event interface {
	Create(ctx context.Context, data *entity.Event) error
	Get(ctx context.Context, paging *util.Pagination) ([]*entity.Event, error)
	GetWithoutImage(ctx context.Context, paging *util.Pagination) ([]*entity.EventWithoutImage, error)
	Count(ctx context.Context, paging *util.Pagination) (int, error)
	Partition(ctx context.Context, date time.Time) error
}
