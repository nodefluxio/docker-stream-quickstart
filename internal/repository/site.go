package repository

import (
	"context"

	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/pkg/util"
)

// Site interface abstracts the repository layer and should be implemented in repository
type Site interface {
	GetList(ctx context.Context, paging *util.Pagination) ([]*entity.Site, error)
	GetDetail(ctx context.Context, ID uint64) (*entity.Site, error)
	Create(ctx context.Context, data *entity.Site) (*entity.Site, error)
	Update(ctx context.Context, data *entity.Site) error
	Delete(ctx context.Context, ID uint64) error
	AddStreamToSite(ctx context.Context, data *entity.MapSiteStream) error
	GetSiteByIDs(ctx context.Context, IDs []int64) ([]*entity.Site, error)
	GetSiteWithStream(ctx context.Context, IDs []int64) ([]*entity.SiteWithStream, error)
	GetDetailByStreamID(ctx context.Context, streamID string) (*entity.Site, error)
	GetMapStreamSiteByStreamID(ctx context.Context, streamID string) (*entity.MapSiteStream, error)
}
