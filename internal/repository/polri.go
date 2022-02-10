package repository

import (
	"context"

	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
)

// Polri interface abstracts the repository layer and should be implemented in repository
type Polri interface {
	SearchPlateNumber(ctx context.Context, plateNumber string) (*entity.PolriPlateNuberInfo, error)
	SearchNIK(ctx context.Context, NIK string) (*entity.PolriCitizenData, error)
}
