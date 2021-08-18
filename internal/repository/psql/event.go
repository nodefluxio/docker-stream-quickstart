package psql

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/pkg/util"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/repository"
)

type psqlEventRepo struct {
	Conn *gorm.DB
}

// NewEventRepository is method to initiate Event repo
func NewEventRepository(conn *gorm.DB) repository.Event {
	return &psqlEventRepo{
		Conn: conn,
	}
}
func (p *psqlEventRepo) Create(ctx context.Context, data *entity.Event) error {
	query := `
	INSERT INTO event(
		type,
		stream_id, 
		detection, 
		primary_image,
		secondary_image,
		result, 
		status, 
		event_time,
		keyword
	) VALUES( 
		$1, $2, $3, $4, $5, $6, $7, $8, 
		to_tsvector('indonesian', coalesce($9, '') || ' ' || coalesce($10, '') || ' ' || coalesce($11, ''))
	)`

	detection, err := json.Marshal(data.Detection)
	if err != nil {
		return err
	}
	var nr entity.EventResult
	err = json.Unmarshal(data.Result, &nr)
	if err != nil {
		return err
	}
	err = p.Conn.Exec(
		query,
		data.EventType,
		data.StreamID,
		detection,
		data.PrimaryImage,
		data.SecondaryImage,
		data.Result,
		data.Status,
		data.EventTime,
		nr.Label,
		nr.Result,
		nr.Location,
	).Error
	return err
}

func (p *psqlEventRepo) mappingField(key string) string {
	switch key {
	case "type":
		key = "status"
	case "timestamp":
		key = "event_time"
	}
	return key
}

func (p *psqlEventRepo) mappingValue(value []string) []string {
	for i := 0; i < len(value); i++ {
		strings.ToLower(value[i])
		switch strings.ToLower(value[i]) {
		case "unrecognized":
			value[i] = "UNKNOWN"
		case "recognized":
			value[i] = "KNOWN"
		}
	}
	return value
}

func (p *psqlEventRepo) generateFilter(ctx context.Context, filter map[string]string, tx *gorm.DB) *gorm.DB {
	// create query for filter
	for key, val := range filter {
		switch key {
		case "timestamp_from":
		case "timestamp_to":
			dateFrom := filter["timestamp_from"]
			dateTo := val
			tx = tx.Where("event_time BETWEEN ?::timestamp with time zone AND ?::timestamp with time zone", dateFrom, dateTo)
		default:
			if val != "" {
				newval := p.mappingValue(strings.Split(val, ","))
				key = p.mappingField(key)
				tx = tx.Where(fmt.Sprintf("%s IN (?)", key), newval)
			}
		}
	}
	return tx
}

func (p *psqlEventRepo) generateSearch(ctx context.Context, search string, tx *gorm.DB) *gorm.DB {
	if search != "" {
		constructor := util.Constructor{
			FTSearch: search,
		}
		search = constructor.FTSQuery()
		tx = tx.Where("keyword @@ to_tsquery('indonesian',?)", search)
	}
	return tx
}

func (p *psqlEventRepo) Get(ctx context.Context, paging *util.Pagination) ([]*entity.Event, error) {
	var data []*entity.Event
	tx := p.Conn.Select("id, type, stream_id, primary_image, secondary_image, result, status, event_time")

	// create query for filter
	tx = p.generateFilter(ctx, paging.Filter, tx)

	// create query for search
	tx = p.generateSearch(ctx, paging.Search, tx)

	// sort
	querySort := paging.Sort
	if len(querySort) != 0 {
		for key, val := range querySort {
			key = p.mappingField(key)
			sort := fmt.Sprintf("%s %s", key, val)
			tx = tx.Order(sort)
		}
	} else {
		tx = tx.Order("event_time DESC")
	}

	// Limit and offset
	if paging.Limit != 0 {
		tx = tx.Limit(paging.Limit).Offset(paging.Offset)
	}

	err := tx.Find(&data).Error
	return data, err
}

func (p *psqlEventRepo) GetWithoutImage(ctx context.Context, paging *util.Pagination) ([]*entity.EventWithoutImage, error) {
	var data []*entity.EventWithoutImage
	tx := p.Conn.Table("event").Select("id, type, detection, stream_id, result, status, event_time, created_at")

	// create query for filter
	tx = p.generateFilter(ctx, paging.Filter, tx)

	// create query for search
	tx = p.generateSearch(ctx, paging.Search, tx)

	// sort
	querySort := paging.Sort
	if len(querySort) != 0 {
		for key, val := range querySort {
			key = p.mappingField(key)
			sort := fmt.Sprintf("%s %s", key, val)
			tx = tx.Order(sort)
		}
	} else {
		tx = tx.Order("event_time DESC")
	}

	// Limit and offset
	if paging.Limit != 0 {
		tx = tx.Limit(paging.Limit).Offset(paging.Offset)
	}

	err := tx.Find(&data).Error
	return data, err
}

func (p *psqlEventRepo) Count(ctx context.Context, paging *util.Pagination) (int, error) {
	var count int
	tx := p.Conn
	// create query for filter
	tx = p.generateFilter(ctx, paging.Filter, tx)
	// create query for search
	tx = p.generateSearch(ctx, paging.Search, tx)
	err := tx.Table("event").Count(&count).Error
	return count, err
}

func (p *psqlEventRepo) Partition(ctx context.Context, date time.Time) error {
	query := "SELECT create_daily_event(?::timestamp);"
	return p.Conn.Exec(query, date).Error
}
