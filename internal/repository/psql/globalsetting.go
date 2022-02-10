package psql

import (
	"context"

	"github.com/jinzhu/gorm"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/repository"
)

type psqlGlobalSettingRepo struct {
	Conn *gorm.DB
}

// NewGlobalSettingRepository is method to initiate gloabl setting repo
func NewGlobalSettingRepository(conn *gorm.DB) repository.GlobalSetting {
	return &psqlGlobalSettingRepo{
		Conn: conn,
	}
}
func (p *psqlGlobalSettingRepo) Create(ctx context.Context, data *entity.GlobalSetting) error {
	err := p.Conn.Save(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (p *psqlGlobalSettingRepo) GetCurrent(ctx context.Context) (*entity.GlobalSetting, error) {
	object := entity.GlobalSetting{}
	err := p.Conn.First(&object).Error
	if err != nil {
		return nil, err
	}
	return &object, err
}

func (p *psqlGlobalSettingRepo) Update(ctx context.Context, data *entity.GlobalSetting) error {
	return p.Conn.Save(data).Error
}
