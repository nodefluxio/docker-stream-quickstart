package psql

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/pkg/util"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/repository"
)

type psqlEventEnrollmentRepo struct {
	Conn *gorm.DB
}

// NewEventEnrollmentRepository is method to initiate EnrolledFace repo
func NewEventEnrollmentRepository(conn *gorm.DB) repository.EventEnrollment {
	return &psqlEventEnrollmentRepo{
		Conn: conn,
	}
}
func (p *psqlEventEnrollmentRepo) Create(ctx context.Context, data *entity.EventEnrollment, images []*entity.EventEnrollmentImage) (*entity.EventEnrollment, error) {
	tx := p.Conn.Begin()
	err := tx.Save(data).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if len(images) != 0 {
		for _, v := range images {
			dataImages := entity.EventEnrollmentFaceImage{
				EventEnrollmentID: data.ID,
				Image:             v.Image,
			}
			err := tx.Save(&dataImages).Error
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	tx.Commit()
	return data, nil
}

func (p *psqlEventEnrollmentRepo) generateFilter(ctx context.Context, filter map[string]string, tx *gorm.DB) *gorm.DB {
	// create query for filter
	for key, val := range filter {
		switch key {
		case "latest_timestamp":
			tx = tx.Where("created_at > ?::timestamp with time zone", val)
		}

	}
	return tx
}

func (p *psqlEventEnrollmentRepo) Get(ctx context.Context, paging *util.Pagination) ([]*entity.EventEnrollment, error) {
	var data []*entity.EventEnrollment
	tx := p.Conn.Select("*")

	// create query for filter
	tx = p.generateFilter(ctx, paging.Filter, tx)

	// Limit and offset
	if paging.Limit != 0 {
		tx = tx.Limit(paging.Limit).Offset(paging.Offset)
	}
	tx = tx.Order("created_at ASC")
	err := tx.Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, err
}

func (p *psqlEventEnrollmentRepo) Partition(ctx context.Context, date time.Time) error {
	query := "SELECT create_daily_event_enrollment(?::timestamp);"
	return p.Conn.Exec(query, date).Error
}
