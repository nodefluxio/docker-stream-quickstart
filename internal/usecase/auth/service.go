package auth

import (
	"context"

	"gitlab.com/nodefluxio/vanilla-dashboard/internal/presenter"
)

// Service interface abstracts the controller layer and should be implemented in controller directory. Controller contains business logics and is independent of any database connection.
type Service interface {
	Login(ctx context.Context, input *presenter.LoginInput) (*presenter.LoginResponse, error)
	RefreshToken(ctx context.Context, postData *presenter.LoginResponse) (*presenter.LoginResponse, error)
}
