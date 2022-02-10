package site

import (
	"context"

	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/pkg/util"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/presenter"
)

// Service interface abstracts the controller layer and should be implemented in controller directory. Controller contains business logics and is independent of any database connection.
type Service interface {
	GetList(ctx context.Context, paging *util.Pagination, userInfo *presenter.AuthInfoResponse) ([]*entity.Site, error)
	Create(ctx context.Context, data *presenter.SiteRequest) (*entity.Site, error)
	Update(ctx context.Context, data *presenter.SiteRequest) error
	Delete(ctx context.Context, ID uint64) error
	AssignStreamToSite(ctx context.Context, data *presenter.AssignStreamRequest) error
}
