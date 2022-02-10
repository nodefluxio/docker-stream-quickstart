package polrisearching

import (
	"context"

	"gitlab.com/nodefluxio/vanilla-dashboard/internal/presenter"
)

// Service interface abstracts the controller layer and should be implemented in controller directory. Controller contains business logics and is independent of any database connection.
type Service interface {
	SearchPlateNumber(ctx context.Context, nopol string) (*presenter.PolriSearchPlateResponse, error)
	SearchNik(ctx context.Context, nik string) (*presenter.PolriCitizenResponse, error)
	GetFaceSearchToken(ctx context.Context, Image []byte, limit uint64) (string, error)
	GetFaceSearchResult(ctx context.Context, token string) ([]*presenter.PolriFaceResultResponse, error)
}
