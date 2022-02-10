package psql

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/pkg/util"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/repository"
)

type psqlSiteRepo struct {
	Conn *gorm.DB
}

// NewSiteRepository is method to initiate Site repo
func NewSiteRepository(conn *gorm.DB) repository.Site {
	return &psqlSiteRepo{
		Conn: conn,
	}
}

func (p *psqlSiteRepo) generateFilter(ctx context.Context, filter map[string]string, tx *gorm.DB) *gorm.DB {
	// create query for filter
	for key, val := range filter {
		if val != "" {
			switch key {
			case "stream_id":
				tx = tx.Where("map_site_stream.stream_id = ?", val)

			case "site_id":
				newval := strings.Split(val, ",")
				if len(newval) > 1 {
					tx = tx.Where("site.id IN (?)", newval)
				} else {
					tx = tx.Where("site.id = ?", newval)
				}
			}
		}
	}
	return tx
}

func (p *psqlSiteRepo) generateSearch(ctx context.Context, search string, tx *gorm.DB) *gorm.DB {
	if search != "" {
		tx = tx.Where("site.name ~* ?", search)
	}
	return tx
}

func (p *psqlSiteRepo) GetList(ctx context.Context, paging *util.Pagination) ([]*entity.Site, error) {
	var data []*entity.Site
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
		tx = tx.Order("site.created_at DESC")
	}

	err := tx.Select("site.*").Joins("LEFT join map_site_stream on site.id = map_site_stream.site_id").Group("map_site_stream.site_id, site.name, site.id").Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, err
}

func (p *psqlSiteRepo) GetDetail(ctx context.Context, ID uint64) (*entity.Site, error) {
	data := &entity.Site{}
	err := p.Conn.Select("*").Limit("1").Find(data, ID).Error
	if err != nil {
		return nil, err
	}
	return data, err
}

func (p *psqlSiteRepo) Create(ctx context.Context, data *entity.Site) (*entity.Site, error) {
	err := p.Conn.Save(data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (p *psqlSiteRepo) Update(ctx context.Context, data *entity.Site) error {
	err := p.Conn.Save(data).Error
	return err
}

func (p *psqlSiteRepo) Delete(ctx context.Context, ID uint64) error {
	timeNow := time.Now()
	query := `
	UPDATE "site" SET deleted_at = $1
	WHERE id = $2
	`
	err := p.Conn.Exec(query, &timeNow, ID).Error
	return err
}

func (p *psqlSiteRepo) AddStreamToSite(ctx context.Context, data *entity.MapSiteStream) error {
	return p.Conn.Save(data).Error
}

func (p *psqlSiteRepo) GetSiteByIDs(ctx context.Context, IDs []int64) ([]*entity.Site, error) {
	var data []*entity.Site
	err := p.Conn.Table("site").Select("*").Where("id IN (?)", IDs).Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (p *psqlSiteRepo) GetSiteWithStream(ctx context.Context, IDs []int64) ([]*entity.SiteWithStream, error) {
	var data []*entity.SiteWithStream
	tx := p.Conn.Table("site").Select("site.id, site.name, map_site_stream.stream_id, site.created_at, site.updated_at, site.deleted_at")
	if len(IDs) > 0 {
		tx = tx.Where("site.id IN (?)", IDs)
	}
	err := tx.Joins("INNER JOIN map_site_stream on site.id = map_site_stream.site_id").Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetDetailByStreamID is function for get detail site by key stream id
func (p *psqlSiteRepo) GetDetailByStreamID(ctx context.Context, streamID string) (*entity.Site, error) {
	data := &entity.Site{}
	err := p.Conn.Select("site.*").Joins("INNER join map_site_stream on site.id = map_site_stream.site_id").Where("map_site_stream.stream_id=?", streamID).Limit("1").Find(data).Error
	if err != nil {
		return nil, err
	}
	return data, err
}

// GetMapStreamSiteByStreamID is function for get map stream site by sreamID
func (p *psqlSiteRepo) GetMapStreamSiteByStreamID(ctx context.Context, streamID string) (*entity.MapSiteStream, error) {
	data := &entity.MapSiteStream{}
	err := p.Conn.Select("*").Where("stream_id = ?", streamID).First(data).Error
	if err != nil {
		return nil, err
	}
	return data, err
}
