package user

import (
	"context"

	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/pkg/util"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/presenter"
)

// Service interface abstracts the controller layer and should be implemented in controller directory. Controller contains business logics and is independent of any database connection.
type Service interface {
	Create(ctx context.Context, postData *presenter.UserRequest) (*entity.User, error)
	Update(ctx context.Context, postData *presenter.UserRequest) error
	ChangePassword(ctx context.Context, postData *presenter.UserChangePassRequest) error
	Delete(ctx context.Context, ID uint64) error
	GetDetail(ctx context.Context, ID uint64) (*entity.User, error)
	GetList(ctx context.Context, paging *util.Pagination) (*presenter.UserPaging, error)
}
