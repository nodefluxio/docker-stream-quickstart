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

func Test_NewEventEnrollmentFaceImageRepository(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		s := helper.Setup()
		defer s.Conn.Close()

		repo := psqlrepo.NewEventEnrollmentFaceImageRepository(s.Db)

		t.Run("initialized", func(t *testing.T) {
			require.NotNil(t, repo)
		})

		t.Run("implements valid interface", func(t *testing.T) {
			repositoryInterface := reflect.TypeOf((*repository.EventEnrollmentFaceImage)(nil)).Elem()
			repositoryStruct := reflect.TypeOf(repo)

			require.True(t, repositoryStruct.Implements(repositoryInterface))
		})
	})
}

func Test_GetByEventEnrollmendID(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewEventEnrollmentFaceImageRepository(s.Db)
	id := uint64(1)

	query := `SELECT * FROM "face_image" WHERE (event_enrollment_id = $1)`

	t.Run("success", func(t *testing.T) {
		now := time.Now()
		rows := sqlmock.
			NewRows([]string{"id", "event_enrollment_id", "image", "created_at"}).
			AddRow(1, 2, []byte{}, now)
		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnRows(rows)

		expectedRes := []*entity.EventEnrollmentFaceImage{
			{
				ID:                1,
				EventEnrollmentID: 2,
				Image:             []byte{},
				CreatedAt:         now,
			},
		}

		res, err := repo.GetByEventEnrollmendID(context.TODO(), id)

		assert.Equal(t, nil, s.Mock.ExpectationsWereMet())
		assert.Equal(t, nil, err)
		assert.Nil(t, deep.Equal(expectedRes, res))

	})

	t.Run("error connection", func(t *testing.T) {
		s.Conn.Close()
		res, err := repo.GetByEventEnrollmendID(context.TODO(), id)

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}
