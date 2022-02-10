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

func Test_NewGlobalSettingRepository(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		s := helper.Setup()
		defer s.Conn.Close()

		repo := psqlrepo.NewGlobalSettingRepository(s.Db)

		t.Run("initialized", func(t *testing.T) {
			require.NotNil(t, repo)
		})

		t.Run("implements valid interface", func(t *testing.T) {
			repositoryInterface := reflect.TypeOf((*repository.GlobalSetting)(nil)).Elem()
			repositoryStruct := reflect.TypeOf(repo)

			require.True(t, repositoryStruct.Implements(repositoryInterface))
		})
	})
}

func Test_CreateGlobalSetting(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewGlobalSettingRepository(s.Db)

	query := `INSERT INTO "global_settings" ("similarity","created_at") VALUES ($1,$2) RETURNING "global_settings"."id"`

	id := uint64(1)
	now := time.Now()
	data := entity.GlobalSetting{
		Similarity: 0.9,
		CreatedAt:  now,
	}

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.
			NewRows([]string{"id"}).
			AddRow(id)
		s.Mock.ExpectBegin()
		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(data.Similarity, now).
			WillReturnRows(rows)
		s.Mock.ExpectCommit()

		err := repo.Create(context.TODO(), &data)

		assert.Equal(t, nil, s.Mock.ExpectationsWereMet())
		assert.Equal(t, nil, err)
		assert.Equal(t, id, data.ID)

	})

	t.Run("error connection", func(t *testing.T) {
		err := repo.Create(context.TODO(), &data)

		assert.Error(t, err)

	})
}

func Test_GetCurrentGlobalSetting(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewGlobalSettingRepository(s.Db)

	query := `SELECT * FROM "global_settings" ORDER BY "global_settings"."id" ASC LIMIT `

	t.Run("success", func(t *testing.T) {
		now := time.Now()
		rows := sqlmock.
			NewRows([]string{"id", "similarity", "created_at"}).
			AddRow(1, 0.9, now)
		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnRows(rows)

		expectedRes := entity.GlobalSetting{
			ID:         1,
			Similarity: 0.9,
			CreatedAt:  now,
		}

		res, err := repo.GetCurrent(context.TODO())

		assert.Equal(t, nil, s.Mock.ExpectationsWereMet())
		assert.Equal(t, nil, err)
		assert.Nil(t, deep.Equal(&expectedRes, res))

	})

	t.Run("error connection", func(t *testing.T) {
		s.Conn.Close()
		res, err := repo.GetCurrent(context.TODO())

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func Test_UpdateGlobalSetting(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewGlobalSettingRepository(s.Db)

	query := `UPDATE "global_settings" SET "similarity" = $1  WHERE "global_settings"."id" = $2`

	data := entity.GlobalSetting{
		ID:         2,
		Similarity: 0.9,
	}

	t.Run("success", func(t *testing.T) {
		result := sqlmock.NewResult(1, 1)

		s.Mock.ExpectBegin()
		s.Mock.
			ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(data.Similarity, data.ID).
			WillReturnResult(result)
		s.Mock.ExpectCommit()

		err := repo.Update(context.TODO(), &data)

		assert.Equal(t, nil, s.Mock.ExpectationsWereMet())
		assert.Equal(t, nil, err)
	})

	t.Run("error connection", func(t *testing.T) {
		err := repo.Update(context.TODO(), &data)
		assert.Error(t, err)
	})
}
