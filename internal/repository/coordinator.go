package repository

import (
	"context"

	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
)

// Coordinator interface abstracts the repository layer and should be implemented in repository
type Coordinator interface {
	Ping(ctx context.Context) error
	GetEnrollmentEvent(ctx context.Context, limit, latestTimestamp string) (*entity.EnrollmentEventCoordinator, error)
	CreateEnrollmentEvent(ctx context.Context, data *entity.CreateEnrollmentEventCoordinator) error
}
