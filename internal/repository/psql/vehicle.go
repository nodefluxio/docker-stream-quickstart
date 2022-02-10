package psql

import (
	"context"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/pkg/util"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/repository"
)

type psqlVehicleRepository struct {
	Conn *gorm.DB
}

// NewVehicleRepository is method to initiate vehicle repository
func NewVehicleRepository(conn *gorm.DB) repository.Vehicle {
	return &psqlVehicleRepository{
		Conn: conn,
	}
}
func (p *psqlVehicleRepository) Create(ctx context.Context, data *entity.Vehicle) (*entity.Vehicle, error) {
	err := p.Conn.Save(data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (p *psqlVehicleRepository) Delete(ctx context.Context, ID uint64) error {
	timeNow := time.Now()
	query := `
		UPDATE "vehicle" SET deleted_at = $1
		WHERE id = $2
	`
	return p.Conn.Exec(query, &timeNow, ID).Error
}
func (p *psqlVehicleRepository) GetAll(ctx context.Context) ([]*entity.Vehicle, error) {
	var data []*entity.Vehicle

	err := p.Conn.Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, err
}
func (p *psqlVehicleRepository) DeleteAll(ctx context.Context) error {
	query := `
		DELETE from "vehicle"
	`
	return p.Conn.Exec(query).Error
}
func (p *psqlVehicleRepository) generateSearch(ctx context.Context, search string, tx *gorm.DB) *gorm.DB {
	if search != "" {
		tx = tx.Where("name ~* ? OR plate_number ~* ? OR type ~* ? OR brand ~* ? OR name ~* ? OR status ~* ? OR color  ~* ? OR unique_id ~* ? ",
			search, search, search, search, search, search, search, search)
	}
	return tx
}

func (p *psqlVehicleRepository) GetList(ctx context.Context, paging *util.Pagination) ([]*entity.Vehicle, error) {
	var data []*entity.Vehicle
	tx := p.Conn
	// create query for search
	tx = p.generateSearch(ctx, paging.Search, tx)
	// sort
	querySort := paging.Sort
	if len(querySort) != 0 {
		for key, val := range querySort {
			sort := fmt.Sprintf("%s %s", key, val)
			tx = tx.Order(sort)
		}
	} else {
		tx = tx.Order("created_at DESC")
	}

	// Limit and offset
	if paging.Limit != 0 {
		tx = tx.Limit(paging.Limit).Offset(paging.Offset)
	}

	err := tx.Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, err
}

func (p *psqlVehicleRepository) Count(ctx context.Context, paging *util.Pagination) (int, error) {
	var count int
	tx := p.Conn.Where("deleted_at IS NULL")
	// create query for search
	tx = p.generateSearch(ctx, paging.Search, tx)
	err := tx.Table("vehicle").Count(&count).Error
	return count, err
}

func (p *psqlVehicleRepository) GetDetail(ctx context.Context, ID uint64) (*entity.Vehicle, error) {
	object := &entity.Vehicle{}
	err := p.Conn.First(object, ID).Error
	if err != nil {
		return nil, err
	}
	return object, err
}
func (p *psqlVehicleRepository) GetByPlateNumber(ctx context.Context, plate string) (*entity.Vehicle, error) {
	object := &entity.Vehicle{}
	err := p.Conn.Where("plate_number = ?", plate).First(object).Error
	if err != nil {
		return nil, err
	}
	return object, err
}
func (p *psqlVehicleRepository) Update(ctx context.Context, data *entity.Vehicle) error {
	return p.Conn.Save(data).Error
}
