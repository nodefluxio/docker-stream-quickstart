package vehicle

import (
	"context"

	"gitlab.com/nodefluxio/vanilla-dashboard/internal/pkg/util"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/presenter"
)

// Service interface abstracts the controller layer and should be implemented in controller directory. Controller contains business logics and is independent of any database connection.
type Service interface {
	GetList(ctx context.Context, paging *util.Pagination) (*presenter.VehiclePaging, error)
	Create(ctx context.Context, postData *presenter.VehicleEnrollmentRequest) (*presenter.VehicleEnrollmentResponse, error)
	Delete(ctx context.Context, ID uint64) error
	DeleteAll(ctx context.Context) error
	GetDetail(ctx context.Context, ID uint64) (*presenter.VehicleEnrollmentResponse, error)
	Update(ctx context.Context, enrollmentID uint64, postData *presenter.VehicleEnrollmentRequest) error
}
