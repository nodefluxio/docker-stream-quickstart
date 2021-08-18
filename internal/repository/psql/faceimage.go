package psql

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/repository"
)

type psqlFaceImageRepo struct {
	Conn *gorm.DB
}

// NewFaceImageRepository is method to initiate FaceImage repo
func NewFaceImageRepository(conn *gorm.DB) repository.FaceImage {
	return &psqlFaceImageRepo{
		Conn: conn,
	}
}
func (p *psqlFaceImageRepo) Create(ctx context.Context, tx *gorm.DB, data *entity.FaceImage) (*gorm.DB, *entity.FaceImage, error) {
	err := tx.Save(data).Error
	if err != nil {
		return tx, nil, err
	}
	return tx, data, nil
}

func (p *psqlFaceImageRepo) GetDetailByEnrollID(ctx context.Context, EnrolledFaceID uint64) ([]*entity.FaceImage, error) {
	var data []*entity.FaceImage
	err := p.Conn.Select("id, enrolled_face_id, variation, image, created_at, deleted_at").Where("enrolled_face_id = ?", EnrolledFaceID).Find(&data).Error
	return data, err
}

func (p *psqlFaceImageRepo) DeleteByEnrollID(ctx context.Context, tx *gorm.DB, EnrolledFaceID uint64) (*gorm.DB, error) {
	timeNow := time.Now()
	query := `
	UPDATE "face_image" SET deleted_at = $1
	WHERE enrolled_face_id = $2
	`
	err := tx.Exec(query, &timeNow, EnrolledFaceID).Error
	return tx, err
}

func (p *psqlFaceImageRepo) DeleteAll(ctx context.Context) error {
	query := `
		DELETE from "face_image"
	`
	return p.Conn.Exec(query).Error
}

func (p *psqlFaceImageRepo) GetImageByID(ctx context.Context, ID uint64) (*entity.FaceImage, error) {
	object := &entity.FaceImage{}
	err := p.Conn.First(object, ID).Error
	return object, err
}
