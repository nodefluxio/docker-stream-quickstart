package repository

import (
	"context"

	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
)

// Event interface abstracts the repository layer and should be implemented in repository
type GlobalSetting interface {
	Create(ctx context.Context, data *entity.GlobalSetting) error
	GetCurrent(ctx context.Context) (*entity.GlobalSetting, error)
	Update(ctx context.Context, data *entity.GlobalSetting) error
}
