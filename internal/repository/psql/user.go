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

type psqlUserRepo struct {
	Conn *gorm.DB
}

// NewUserRepository is method to initiate User repo
func NewUserRepository(conn *gorm.DB) repository.User {
	return &psqlUserRepo{
		Conn: conn,
	}
}

func (p *psqlUserRepo) Create(ctx context.Context, data *entity.User) (*entity.User, error) {
	err := p.Conn.Save(data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (p *psqlUserRepo) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	err := p.Conn.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (p *psqlUserRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := p.Conn.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (p *psqlUserRepo) UpdateBasicData(ctx context.Context, data *entity.User) error {
	var user entity.User
	err := p.Conn.Model(&user).Omit("created_at", "password", "deleted_at").Update(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (p *psqlUserRepo) UpdatePassword(ctx context.Context, password string, ID uint64) error {
	var user entity.User
	err := p.Conn.Model(&user).Where("id=?", ID).Update("password", password).Error
	if err != nil {
		return err
	}
	return nil
}

func (p *psqlUserRepo) GetDetail(ctx context.Context, ID uint64) (*entity.User, error) {
	var user entity.User
	err := p.Conn.First(&user, ID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (p *psqlUserRepo) Delete(ctx context.Context, ID uint64) error {
	var user entity.User
	timeNow := time.Now()
	err := p.Conn.Model(&user).Where("id=?", ID).Update("deleted_at", timeNow).Error
	if err != nil {
		return err
	}
	return nil
}

func (p *psqlUserRepo) generateSearch(ctx context.Context, search string, tx *gorm.DB) *gorm.DB {
	if search != "" {
		tx = tx.Where("fullname ~* ? OR email ~* ?", search, search)
	}
	return tx
}

func (p *psqlUserRepo) generateFilter(ctx context.Context, filter map[string]string, tx *gorm.DB) *gorm.DB {
	// create query for filter
	for key, val := range filter {
		switch key {
		case "role":
			tx = tx.Where("role = ?", val)
		}
	}
	return tx
}

func (p *psqlUserRepo) GetList(ctx context.Context, paging *util.Pagination) ([]*entity.User, error) {
	var user []*entity.User
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

	err := tx.Find(&user).Error
	if err != nil {
		return nil, err
	}
	return user, err
}

func (p *psqlUserRepo) Count(ctx context.Context, paging *util.Pagination) (int, error) {
	var count int
	tx := p.Conn.Where("deleted_at IS NULL")
	// create query for search
	tx = p.generateSearch(ctx, paging.Search, tx)
	err := tx.Table("user_access").Count(&count).Error
	return count, err
}
