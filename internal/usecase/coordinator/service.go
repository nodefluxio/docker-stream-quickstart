package coordinator

import (
	"context"

	"gitlab.com/nodefluxio/vanilla-dashboard/internal/pkg/util"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/presenter"
)

// Service interface abstracts the controller layer and should be implemented in controller directory. Controller contains business logics and is independent of any database connection.
type Service interface {
	Create(ctx context.Context, postData *presenter.CoordinatorRequest) error
	Get(ctx context.Context, paging *util.Pagination) ([]*presenter.CoordinatorResponse, error)
	Partition(ctx context.Context) error
	CronjobPartition(ctx context.Context) error
}
