package psql_test

import (
	"context"
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

func Test_NewEventEnrollmentRepository(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		s := helper.Setup()
		defer s.Conn.Close()

		repo := psqlrepo.NewEventEnrollmentRepository(s.Db)

		t.Run("initialized", func(t *testing.T) {
			require.NotNil(t, repo)
		})

		t.Run("implements valid interface", func(t *testing.T) {
			repositoryInterface := reflect.TypeOf((*repository.EventEnrollment)(nil)).Elem()
			repositoryStruct := reflect.TypeOf(repo)

			require.True(t, repositoryStruct.Implements(repositoryInterface))
		})
	})
}

func Test_CreateEventEnrollment(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewEventEnrollmentRepository(s.Db)
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		id := uint64(1)
		data := &entity.EventEnrollment{
			EventID:     "event-id",
			Agent:       "nodeflux",
			EventAction: "GET",
			Payload:     []byte{},
			CreatedAt:   now,
		}
		image := []*entity.EventEnrollmentImage{}
		query := `INSERT INTO "event_enrollments" ("event_id","agent","event_action","payload","created_at") VALUES ($1,$2,$3,$4,$5) RETURNING "event_enrollments"."id"`

		rows := sqlmock.
			NewRows([]string{"id"}).
			AddRow(id)

		s.Mock.ExpectBegin()
		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(data.EventID, data.Agent, data.EventAction, data.Payload, now).
			WillReturnRows(rows)
		s.Mock.ExpectCommit()

		res, err := repo.Create(context.TODO(), data, image)

		expectedRes := entity.EventEnrollment{
			ID:          1,
			EventID:     "event-id",
			Agent:       "nodeflux",
			EventAction: "GET",
			Payload:     []byte{},
			CreatedAt:   now,
		}

		assert.Equal(t, nil, s.Mock.ExpectationsWereMet())
		assert.Equal(t, nil, err)
		assert.Nil(t, deep.Equal(&expectedRes, res))
	})

	t.Run("error connection", func(t *testing.T) {
		s.Conn.Close()
		data := &entity.EventEnrollment{
			EventID:     "event-id",
			Agent:       "nodeflux",
			EventAction: "GET",
			Payload:     []byte{},
			CreatedAt:   now,
		}
		image := []*entity.EventEnrollmentImage{}

		s.Mock.ExpectBegin()
		s.Mock.ExpectRollback()

		res, err := repo.Create(context.TODO(), data, image)

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func Test_GetEventEnrollment(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewEventEnrollmentRepository(s.Db)
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		query := `SELECT * FROM "event_enrollments" ORDER BY created_at ASC`

		rows := sqlmock.
			NewRows([]string{"id", "event_id", "agent", "event_action", "payload", "created_at"}).
			AddRow(1, "event-id", "nodeflux", "GET", []byte{}, now)

		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnRows(rows)

		filter := &util.Pagination{}
		res, err := repo.Get(context.TODO(), filter)

		expectedRes := []*entity.EventEnrollment{
			{
				ID:          1,
				EventID:     "event-id",
				Agent:       "nodeflux",
				EventAction: "GET",
				Payload:     []byte{},
				CreatedAt:   now,
			},
		}

		assert.Equal(t, nil, s.Mock.ExpectationsWereMet())
		assert.Equal(t, nil, err)
		assert.Nil(t, deep.Equal(expectedRes, res))
	})

	t.Run("success wtih parameter pagination (limit and offset)", func(t *testing.T) {
		query := `SELECT * FROM "event_enrollments" ORDER BY created_at ASC LIMIT 2 OFFSET 0`

		rows := sqlmock.
			NewRows([]string{"id", "event_id", "agent", "event_action", "payload", "created_at"}).
			AddRow(1, "event-id", "nodeflux", "GET", []byte{}, now)

		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnRows(rows)

		filter := &util.Pagination{
			Limit:  2,
			Offset: 0,
		}
		res, err := repo.Get(context.TODO(), filter)

		expectedRes := []*entity.EventEnrollment{
			{
				ID:          1,
				EventID:     "event-id",
				Agent:       "nodeflux",
				EventAction: "GET",
				Payload:     []byte{},
				CreatedAt:   now,
			},
		}

		assert.Equal(t, nil, s.Mock.ExpectationsWereMet())
		assert.Equal(t, nil, err)
		assert.Nil(t, deep.Equal(expectedRes, res))
	})

	t.Run("success wtih parameter filter latest timestamp", func(t *testing.T) {
		query := `SELECT * FROM "event_enrollments" WHERE (created_at > $1::timestamp with time zone) ORDER BY created_at ASC`

		rows := sqlmock.
			NewRows([]string{"id", "event_id", "agent", "event_action", "payload", "created_at"}).
			AddRow(1, "event-id", "nodeflux", "GET", []byte{}, now)

		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(now.String()).
			WillReturnRows(rows)

		filter := &util.Pagination{
			Filter: map[string]string{
				"latest_timestamp": now.String(),
			},
		}
		res, err := repo.Get(context.TODO(), filter)

		expectedRes := []*entity.EventEnrollment{
			{
				ID:          1,
				EventID:     "event-id",
				Agent:       "nodeflux",
				EventAction: "GET",
				Payload:     []byte{},
				CreatedAt:   now,
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

func Test_PartitionEventEnrollment(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewEventEnrollmentRepository(s.Db)

	query := `SELECT create_daily_event_enrollment($1::timestamp);`

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
