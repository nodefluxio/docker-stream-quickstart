package enrollment

import (
	"context"
	"os"

	"gitlab.com/nodefluxio/vanilla-dashboard/internal/pkg/util"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/presenter"
)

// Service interface abstracts the controller layer and should be implemented in controller directory. Controller contains business logics and is independent of any database connection.
type Service interface {
	GetList(ctx context.Context, paging *util.Pagination) (*presenter.EnrollmentPaging, error)
	Create(ctx context.Context, postData *presenter.EnrollmentRequest, isAgent string) (*presenter.EnrollmentResponse, error)
	Update(ctx context.Context, enrollmentID uint64, postData *presenter.EnrollmentRequest, isAgent string) error
	Delete(ctx context.Context, ID uint64, isAgent string) error
	DeleteAll(ctx context.Context) error
	GetDetail(ctx context.Context, ID uint64) (*presenter.EnrollmentResponse, error)
	GetFaceImage(ctx context.Context, ID uint64) ([]byte, error)
	Backup(ctx context.Context) (*os.File, error)
}
