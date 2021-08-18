package repository

import (
	"context"

	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
)

// Agent interface abstracts the repository layer and should be implemented in repository
type Agent interface {
	Ping(ctx context.Context) error
	CreateEnrollmentEvent(ctx context.Context, data *entity.CreateEnrollmentEventCoordinator) error
}
