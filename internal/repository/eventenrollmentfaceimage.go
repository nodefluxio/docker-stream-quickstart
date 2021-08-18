package repository

import (
	"context"

	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
)

// EventEnrollmentFaceImage interface abstracts the repository layer and should be implemented in repository
type EventEnrollmentFaceImage interface {
	GetByEventEnrollmendID(ctx context.Context, id uint64) ([]*entity.EventEnrollmentFaceImage, error)
}
