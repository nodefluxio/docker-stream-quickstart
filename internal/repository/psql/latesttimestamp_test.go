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
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/repository"
	psqlrepo "gitlab.com/nodefluxio/vanilla-dashboard/internal/repository/psql"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/repository/psql/helper"
)

func Test_NewLatestTimestampRepository(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		s := helper.Setup()
		defer s.Conn.Close()

		repo := psqlrepo.NewLatestTimestampRepository(s.Db)

		t.Run("initialized", func(t *testing.T) {
			require.NotNil(t, repo)
		})

		t.Run("implements valid interface", func(t *testing.T) {
			repositoryInterface := reflect.TypeOf((*repository.LatestTimestamp)(nil)).Elem()
			repositoryStruct := reflect.TypeOf(repo)

			require.True(t, repositoryStruct.Implements(repositoryInterface))
		})
	})
}

func Test_CreateOrUpdateLatestTimestamp(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewLatestTimestampRepository(s.Db)

	now := time.Now()
	data := entity.LatestTimestamp{
		Timestamp: now.String(),
	}

	t.Run("success update", func(t *testing.T) {
		query := `UPDATE "latest_timestamp" SET "timestamp" = $1`
		result := sqlmock.NewResult(1, 1)

		s.Mock.ExpectBegin()
		s.Mock.
			ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(data.Timestamp).
			WillReturnResult(result)
		s.Mock.ExpectCommit()

		err := repo.CreateOrUpdate(context.TODO(), &data)

		assert.Equal(t, nil, s.Mock.ExpectationsWereMet())
		assert.Equal(t, nil, err)
	})

	t.Run("error connection", func(t *testing.T) {
		s.Conn.Close()
		err := repo.CreateOrUpdate(context.TODO(), &data)
		assert.Error(t, err)
	})
}

func Test_GetLatestTimestamp(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewLatestTimestampRepository(s.Db)

	query := `SELECT * FROM "latest_timestamps" LIMIT 1`
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.
			NewRows([]string{"timestamp"}).
			AddRow(now.String())
		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnRows(rows)

		expectedRes := entity.LatestTimestamp{
			Timestamp: now.String(),
		}

		res, err := repo.Get(context.TODO())

		assert.Equal(t, nil, s.Mock.ExpectationsWereMet())
		assert.Equal(t, nil, err)
		assert.Nil(t, deep.Equal(&expectedRes, res))

	})

	t.Run("error connection", func(t *testing.T) {
		s.Conn.Close()
		expectedRes := entity.LatestTimestamp{}
		res, err := repo.Get(context.TODO())

		assert.Error(t, err)
		assert.Nil(t, deep.Equal(&expectedRes, res))
	})
}
