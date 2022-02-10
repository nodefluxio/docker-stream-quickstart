package stream

import (
	"context"

	"gitlab.com/nodefluxio/vanilla-dashboard/internal/pkg/util"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/presenter"
)

// Service interface abstracts the controller layer and should be implemented in controller directory. Controller contains business logics and is independent of any database connection.
type Service interface {
	GetList(ctx context.Context, paging *util.Pagination, userInfo *presenter.AuthInfoResponse) (*presenter.StreamResponse, error)
	GetDetail(ctx context.Context, streamRequest *presenter.StreamRequest, userInfo *presenter.AuthInfoResponse) (*presenter.StreamDetailWithSite, error)
}
