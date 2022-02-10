package repository

import (
	"context"

	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/pkg/util"
)

// Plate interface for access with datastructure platea
type Vehicle interface {
	Create(ctx context.Context, data *entity.Vehicle) (*entity.Vehicle, error)
	Delete(ctx context.Context, ID uint64) error
	GetAll(ctx context.Context) ([]*entity.Vehicle, error)
	DeleteAll(ctx context.Context) error
	GetList(ctx context.Context, paging *util.Pagination) ([]*entity.Vehicle, error)
	Count(ctx context.Context, paging *util.Pagination) (int, error)
	GetDetail(ctx context.Context, ID uint64) (*entity.Vehicle, error)
	Update(ctx context.Context, data *entity.Vehicle) error
	GetByPlateNumber(ctx context.Context, plate string) (*entity.Vehicle, error)
}
