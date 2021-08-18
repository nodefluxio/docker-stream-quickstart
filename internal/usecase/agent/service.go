package agent

import (
	"context"

	"gitlab.com/nodefluxio/vanilla-dashboard/internal/presenter"
)

// Service interface abstracts the controller layer and should be implemented in controller directory. Controller contains business logics and is independent of any database connection.
type Service interface {
	Sync(ctx context.Context) error
	CronjobSyncEnrollment(ctx context.Context) error
	Create(ctx context.Context, postData *presenter.CoordinatorRequest) error
	PingCoordinator(ctx context.Context) error
	GetStatus(ctx context.Context) (*presenter.AgentStatus, error)
}
