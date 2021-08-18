package repository

import (
	"context"
	"time"

	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/pkg/util"
)

// EventEnrollment ,,,
type EventEnrollment interface {
	Create(ctx context.Context, data *entity.EventEnrollment, images []*entity.EventEnrollmentImage) (*entity.EventEnrollment, error)
	Get(ctx context.Context, paging *util.Pagination) ([]*entity.EventEnrollment, error)
	Partition(ctx context.Context, date time.Time) error
}
