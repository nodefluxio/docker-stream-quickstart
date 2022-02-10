package repository

import (
	"context"

	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/pkg/util"
)

// User interface abstracts the repository layer and should be implemented in repository
type User interface {
	Create(ctx context.Context, data *entity.User) (*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	UpdateBasicData(ctx context.Context, data *entity.User) error
	UpdatePassword(ctx context.Context, password string, ID uint64) error
	GetDetail(ctx context.Context, ID uint64) (*entity.User, error)
	Delete(ctx context.Context, ID uint64) error
	GetList(ctx context.Context, paging *util.Pagination) ([]*entity.User, error)
	Count(ctx context.Context, paging *util.Pagination) (int, error)
}
