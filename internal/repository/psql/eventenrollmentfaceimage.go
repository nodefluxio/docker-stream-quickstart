package psql

import (
	"context"

	"github.com/jinzhu/gorm"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/repository"
)

type psqlEventEnrollmentFaceImageRepo struct {
	Conn *gorm.DB
}

// NewFaceImageRepository is method to initiate FaceImage repo
func NewEventEnrollmentFaceImageRepository(conn *gorm.DB) repository.EventEnrollmentFaceImage {
	return &psqlEventEnrollmentFaceImageRepo{
		Conn: conn,
	}
}

func (p *psqlEventEnrollmentFaceImageRepo) GetByEventEnrollmendID(ctx context.Context, id uint64) ([]*entity.EventEnrollmentFaceImage, error) {
	var data []*entity.EventEnrollmentFaceImage
	err := p.Conn.Where("event_enrollment_id = ?", id).Find(&data).Error
	return data, err
}
