package repository

import (
	"context"

	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
)

// EnrollmentVanilla interface abstracts the repository layer and should be implemented in repository
type EnrollmentVanilla interface {
	CreateFaceEnrollment(ctx context.Context, data *entity.VanillaEnrollmentPayload, images []*entity.EnrollmentImage) error
	UpdateFaceEnrollment(ctx context.Context, ID uint64, data *entity.VanillaEnrollmentPayload, images []*entity.EnrollmentImage) error
	DeleteFaceEnrollment(ctx context.Context, ID uint64) error
	GetByFaceID(ctx context.Context, faceID uint64) (*entity.VanillaEnrollmentData, error)
}
