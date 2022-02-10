package psql

import (
	"context"

	"github.com/jinzhu/gorm"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/repository"
)

type psqlLatestTimestampRepo struct {
	Conn *gorm.DB
}

// NewLatestTimestampRepository is method to initiate EnrolledFace repo
func NewLatestTimestampRepository(conn *gorm.DB) repository.LatestTimestamp {
	return &psqlLatestTimestampRepo{
		Conn: conn,
	}
}
func (p *psqlLatestTimestampRepo) CreateOrUpdate(ctx context.Context, data *entity.LatestTimestamp) error {
	update := p.Conn.Table("latest_timestamp").Update("timestamp", data.Timestamp)
	err := update.Error
	if err != nil {
		return err
	}
	if update.RowsAffected == 0 {
		err = p.Conn.Save(data).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *psqlLatestTimestampRepo) Get(ctx context.Context) (*entity.LatestTimestamp, error) {
	object := &entity.LatestTimestamp{}
	err := p.Conn.First(object).Error
	return object, err
}
