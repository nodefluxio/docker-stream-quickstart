package psql_test

import (
	"context"
	"encoding/json"
	"log"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/pkg/util"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/repository"
	psqlrepo "gitlab.com/nodefluxio/vanilla-dashboard/internal/repository/psql"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/repository/psql/helper"
)

func Test_NewEventRepository(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		s := helper.Setup()
		defer s.Conn.Close()

		repo := psqlrepo.NewEventRepository(s.Db)

		t.Run("initialized", func(t *testing.T) {
			require.NotNil(t, repo)
		})

		t.Run("implements valid interface", func(t *testing.T) {
			repositoryInterface := reflect.TypeOf((*repository.Event)(nil)).Elem()
			repositoryStruct := reflect.TypeOf(repo)

			require.True(t, repositoryStruct.Implements(repositoryInterface))
		})
	})
}

func Test_CreateEvent(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewEventRepository(s.Db)

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

	now := time.Now()
	eventResult := entity.EventResult{
		Label:     "recognize",
		Result:    "ade yusup permana",
		Location:  "camera depan",
		Timestamp: now,
	}
	eventResultByte, err := json.Marshal(eventResult)
	if err != nil {
		log.Println("error marshal struct with error :", err)
	}

	detectionByte, err := json.Marshal(entity.Message{})
	if err != nil {
		log.Println("error marshal struct with error :", err)
	}

	data := entity.Event{
		EventType:      "NFV4-FR",
		StreamID:       "this-is-stream-id",
		Detection:      entity.Message{},
		PrimaryImage:   []byte{},
		SecondaryImage: []byte{},
		Result:         eventResultByte,
		Status:         "UNKNOWN",
		EventTime:      now,
	}
	t.Run("success", func(t *testing.T) {
		result := sqlmock.NewResult(1, 1)

		s.Mock.
			ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(
				data.EventType,
				data.StreamID,
				detectionByte,
				data.PrimaryImage,
				data.SecondaryImage,
				data.Result,
				data.Status,
				data.EventTime,
				eventResult.Label,
				eventResult.Result,
				eventResult.Location).
			WillReturnResult(result)

		err := repo.Create(context.TODO(), &data)

		assert.Equal(t, nil, s.Mock.ExpectationsWereMet())
		assert.Equal(t, nil, err)

	})

	t.Run("error connection", func(t *testing.T) {
		err := repo.Create(context.TODO(), &data)

		assert.Error(t, err)

	})
}

func Test_GetEvent(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewEventRepository(s.Db)
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		query := `SELECT id, type, stream_id, primary_image, secondary_image, result, status, event_time FROM "events" ORDER BY event_time DESC`

		rows := sqlmock.
			NewRows([]string{"id", "type", "stream_id", "primary_image", "secondary_image", "result", "status", "event_time"}).
			AddRow(1, "NFV4-FR", "this-is-stream-id", []byte{}, []byte{}, []byte{}, "UNKNOWN", now)

		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnRows(rows)

		filter := &util.Pagination{}
		res, err := repo.Get(context.TODO(), filter)

		expectedRes := []*entity.Event{
			{
				ID:             1,
				EventType:      "NFV4-FR",
				StreamID:       "this-is-stream-id",
				Detection:      entity.Message{},
				PrimaryImage:   []byte{},
				SecondaryImage: []byte{},
				Result:         []byte{},
				Status:         "UNKNOWN",
				EventTime:      now,
			},
		}

		assert.Equal(t, nil, s.Mock.ExpectationsWereMet())
		assert.Equal(t, nil, err)
		assert.Nil(t, deep.Equal(expectedRes, res))
	})

	t.Run("success with parameter search", func(t *testing.T) {
		query := `SELECT id, type, stream_id, primary_image, secondary_image, result, status, event_time FROM "events" WHERE (keyword @@ to_tsquery('indonesian',$1)) ORDER BY event_time DESC`

		filter := &util.Pagination{
			Search: "nodeflux",
		}
		ftsFormat := "nodeflux:*"

		rows := sqlmock.
			NewRows([]string{"id", "type", "stream_id", "primary_image", "secondary_image", "result", "status", "event_time"}).
			AddRow(1, "NFV4-FR", "this-is-stream-id", []byte{}, []byte{}, []byte{}, "UNKNOWN", now)
		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(ftsFormat).
			WillReturnRows(rows)

		res, err := repo.Get(context.TODO(), filter)

		expectedRes := []*entity.Event{
			{
				ID:             1,
				EventType:      "NFV4-FR",
				StreamID:       "this-is-stream-id",
				Detection:      entity.Message{},
				PrimaryImage:   []byte{},
				SecondaryImage: []byte{},
				Result:         []byte{},
				Status:         "UNKNOWN",
				EventTime:      now,
			},
		}

		assert.Equal(t, nil, s.Mock.ExpectationsWereMet())
		assert.Equal(t, nil, err)
		assert.Nil(t, deep.Equal(expectedRes, res))

	})

	t.Run("error connection", func(t *testing.T) {
		s.Conn.Close()
		filter := &util.Pagination{}
		res, err := repo.Get(context.TODO(), filter)

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func Test_GetWithoutImageEvent(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewEventRepository(s.Db)
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		query := `SELECT id, type, detection, stream_id, result, status, event_time, created_at FROM "event" ORDER BY event_time DESC`

		rows := sqlmock.
			NewRows([]string{"id", "type", "stream_id", "result", "status", "event_time"}).
			AddRow(1, "NFV4-FR", "this-is-stream-id", []byte{}, "UNKNOWN", now)

		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnRows(rows)

		filter := &util.Pagination{}
		res, err := repo.GetWithoutImage(context.TODO(), filter)

		expectedRes := []*entity.EventWithoutImage{
			{
				ID:        1,
				EventType: "NFV4-FR",
				StreamID:  "this-is-stream-id",
				Result:    []byte{},
				Status:    "UNKNOWN",
				EventTime: now,
			},
		}

		assert.Equal(t, nil, s.Mock.ExpectationsWereMet())
		assert.Equal(t, nil, err)
		assert.Nil(t, deep.Equal(expectedRes, res))
	})

	t.Run("success with parameter search", func(t *testing.T) {
		query := `SELECT id, type, detection, stream_id, result, status, event_time, created_at FROM "event" WHERE (keyword @@ to_tsquery('indonesian',$1)) ORDER BY event_time DESC`

		filter := &util.Pagination{
			Search: "nodeflux",
		}
		ftsFormat := "nodeflux:*"

		rows := sqlmock.
			NewRows([]string{"id", "type", "stream_id", "result", "status", "event_time"}).
			AddRow(1, "NFV4-FR", "this-is-stream-id", []byte{}, "UNKNOWN", now)

		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(ftsFormat).
			WillReturnRows(rows)

		res, err := repo.GetWithoutImage(context.TODO(), filter)

		expectedRes := []*entity.EventWithoutImage{
			{
				ID:        1,
				EventType: "NFV4-FR",
				StreamID:  "this-is-stream-id",
				Result:    []byte{},
				Status:    "UNKNOWN",
				EventTime: now,
			},
		}

		assert.Equal(t, nil, s.Mock.ExpectationsWereMet())
		assert.Equal(t, nil, err)
		assert.Nil(t, deep.Equal(expectedRes, res))

	})

	t.Run("error connection", func(t *testing.T) {
		s.Conn.Close()
		filter := &util.Pagination{}
		res, err := repo.GetWithoutImage(context.TODO(), filter)

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}
func Test_CountEvent(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewEventRepository(s.Db)

	t.Run("success", func(t *testing.T) {
		query := `SELECT count(*) FROM "event" WHERE (stream_id = $1) AND (keyword @@ to_tsquery('indonesian',$2))`

		filter := &util.Pagination{
			Search: "nodeflux",
			Filter: map[string]string{
				"stream_id": "this-is-stream-id",
			},
		}
		totalData := 10
		ftsFormat := "nodeflux:*"
		rows := sqlmock.
			NewRows([]string{"count"}).
			AddRow(totalData)
		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(filter.Filter["stream_id"], ftsFormat).
			WillReturnRows(rows)

		res, err := repo.Count(context.TODO(), filter)

		assert.Equal(t, nil, s.Mock.ExpectationsWereMet())
		assert.Equal(t, nil, err)
		assert.Equal(t, totalData, res)

	})

	t.Run("error connection", func(t *testing.T) {
		s.Conn.Close()
		filter := &util.Pagination{
			Search: "testing2",
		}
		totalData := 0
		res, err := repo.Count(context.TODO(), filter)
		assert.Error(t, err)
		assert.Equal(t, totalData, res)

	})
}

func Test_PartitionEvent(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewEventRepository(s.Db)

	query := `SELECT create_daily_event($1::timestamp);`

	now := time.Now()

	t.Run("success", func(t *testing.T) {
		result := sqlmock.NewResult(1, 1)

		s.Mock.
			ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(now).
			WillReturnResult(result)

		err := repo.Partition(context.TODO(), now)

		assert.Equal(t, nil, s.Mock.ExpectationsWereMet())
		assert.Equal(t, nil, err)
	})

	t.Run("error connection", func(t *testing.T) {
		err := repo.Partition(context.TODO(), now)
		assert.Error(t, err)
	})
}
