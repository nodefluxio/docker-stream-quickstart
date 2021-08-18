package repository

import (
	"context"

	"github.com/jinzhu/gorm"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/pkg/util"
)

// EnrolledFace ,,,
type EnrolledFace interface {
	Create(ctx context.Context, tx *gorm.DB, data *entity.EnrolledFace) (*gorm.DB, *entity.EnrolledFace, error)
	Update(ctx context.Context, tx *gorm.DB, data *entity.EnrolledFace) (*gorm.DB, error)
	Delete(ctx context.Context, tx *gorm.DB, ID uint64) (*gorm.DB, error)
	GetAll(ctx context.Context) ([]*entity.EnrolledFace, error)
	DeleteAll(ctx context.Context) error
	GetList(ctx context.Context, paging *util.Pagination) ([]*entity.EnrolledFace, error)
	Count(ctx context.Context, paging *util.Pagination) (int, error)
	GetDetail(ctx context.Context, ID uint64) (*entity.EnrolledFace, error)
	GetDetailwFaceID(ctx context.Context, faceID uint64) (*entity.EnrolledFaceWithImage, error)
}
