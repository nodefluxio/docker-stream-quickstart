package repository

import (
	"context"

	"github.com/jinzhu/gorm"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
)

// FaceImage interface abstracts the repository layer and should be implemented in repository
type FaceImage interface {
	Create(ctx context.Context, tx *gorm.DB, data *entity.FaceImage) (*gorm.DB, *entity.FaceImage, error)
	GetDetailByEnrollID(ctx context.Context, EnrolledFaceID uint64) ([]*entity.FaceImage, error)
	GetDetailWithoutImgByEnrollID(ctx context.Context, EnrolledFaceID uint64) ([]*entity.FaceImage, error)
	GetImageByID(ctx context.Context, ID uint64) (*entity.FaceImage, error)
	DeleteByEnrollID(ctx context.Context, tx *gorm.DB, EnrolledFaceID uint64) (*gorm.DB, error)
	DeleteAll(ctx context.Context) error
	DeleteVariation(ctx context.Context, tx *gorm.DB, deletedVariations []string) (*gorm.DB, error)
	GetDetailByVariations(ctx context.Context, variations []string) ([]*entity.FaceImage, error)
}
