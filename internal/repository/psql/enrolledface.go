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

type psqlEnrolledFaceRepo struct {
	Conn *gorm.DB
}

// NewEnrollFaceRepository is method to initiate EnrolledFace repo
func NewEnrollFaceRepository(conn *gorm.DB) repository.EnrolledFace {
	return &psqlEnrolledFaceRepo{
		Conn: conn,
	}
}

func (p *psqlEnrolledFaceRepo) Create(ctx context.Context, tx *gorm.DB, data *entity.EnrolledFace) (*gorm.DB, *entity.EnrolledFace, error) {
	err := tx.Save(data).Error
	if err != nil {
		return tx, nil, err
	}
	return tx, data, nil
}

func (p *psqlEnrolledFaceRepo) Update(ctx context.Context, tx *gorm.DB, data *entity.EnrolledFace) (*gorm.DB, error) {
	err := tx.Save(data).Error
	return tx, err
}

func (p *psqlEnrolledFaceRepo) Delete(ctx context.Context, tx *gorm.DB, ID uint64) (*gorm.DB, error) {
	timeNow := time.Now()
	query := `
		UPDATE "enrolled_face" SET deleted_at = $1
		WHERE id = $2
	`
	err := tx.Exec(query, &timeNow, ID).Error
	return tx, err
}
func (p *psqlEnrolledFaceRepo) GetAll(ctx context.Context) ([]*entity.EnrolledFace, error) {
	var data []*entity.EnrolledFace

	err := p.Conn.Find(&data).Error
	return data, err
}
func (p *psqlEnrolledFaceRepo) DeleteAll(ctx context.Context) error {
	query := `
		DELETE from "enrolled_face"
	`
	return p.Conn.Exec(query).Error
}
func (p *psqlEnrolledFaceRepo) generateSearch(ctx context.Context, search string, tx *gorm.DB) *gorm.DB {
	if search != "" {
		tx = tx.Where("name ~* ?", search)
	}
	return tx
}

func (p *psqlEnrolledFaceRepo) generateFilter(ctx context.Context, filter map[string]string, tx *gorm.DB) *gorm.DB {
	// create query for filter
	for key, val := range filter {
		switch key {
		case "face_id":
			tx = tx.Where("face_id = ?", val)
		}
	}
	return tx
}

func (p *psqlEnrolledFaceRepo) GetList(ctx context.Context, paging *util.Pagination) ([]*entity.EnrolledFace, error) {
	var data []*entity.EnrolledFace
	tx := p.Conn

	// create query for filter
	tx = p.generateFilter(ctx, paging.Filter, tx)
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
	return data, err
}

func (p *psqlEnrolledFaceRepo) Count(ctx context.Context, paging *util.Pagination) (int, error) {
	var count int
	tx := p.Conn.Where("deleted_at IS NULL")
	// create query for search
	tx = p.generateSearch(ctx, paging.Search, tx)
	err := tx.Table("enrolled_face").Count(&count).Error
	return count, err
}

func (p *psqlEnrolledFaceRepo) GetDetail(ctx context.Context, ID uint64) (*entity.EnrolledFace, error) {
	object := &entity.EnrolledFace{}
	err := p.Conn.First(object, ID).Error
	return object, err
}

func (p *psqlEnrolledFaceRepo) GetDetailwFaceID(ctx context.Context, faceID uint64) (*entity.EnrolledFaceWithImage, error) {
	object := &entity.EnrolledFaceWithImage{}
	err := p.Conn.Select("enrolled_face.*, face_image.image_thumbnail as image").Table("enrolled_face").Joins("INNER JOIN face_image on face_image.enrolled_face_id = enrolled_face.id").Where("enrolled_face.deleted_at is NULL AND enrolled_face.face_id= ?", faceID).Order("enrolled_face ASC").Limit(1).Scan(&object).Error
	return object, err
}
