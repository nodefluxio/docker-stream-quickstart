package repository

import (
	"context"

	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
)

// FRemis interface abstracts the repository layer and should be implemented in repository
type FRemis interface {
	FaceEnrollment(ctx context.Context, image string) (*entity.FaceEnrollment, error)
	FaceRecognition(image string) ([]*entity.FRemisCandidate, error)
	FaceDeleteEnrollment(ctx context.Context, faceIDs []string) error
	AddFaceVariation(ctx context.Context, faceID, image string) (*entity.FaceEnrollment, error)
	DeleteFaceVariation(ctx context.Context, faceID string, variations []string) error
	GetFaceEmbedings(ctx context.Context, image string) (*entity.FaceEmbedings, error)
}
