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

func Test_NewFaceImageRepository(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		s := helper.Setup()
		defer s.Conn.Close()

		repo := psqlrepo.NewFaceImageRepository(s.Db)

		t.Run("initialized", func(t *testing.T) {
			require.NotNil(t, repo)
		})

		t.Run("implements valid interface", func(t *testing.T) {
			repositoryInterface := reflect.TypeOf((*repository.FaceImage)(nil)).Elem()
			repositoryStruct := reflect.TypeOf(repo)

			require.True(t, repositoryStruct.Implements(repositoryInterface))
		})
	})
}

func Test_GetDetailByEnrollIDFaceImage(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewFaceImageRepository(s.Db)
	id := uint64(1)

	query := `SELECT id, enrolled_face_id, variation, image, created_at, deleted_at FROM "face_images" WHERE "face_images"."deleted_at" IS NULL AND ((enrolled_face_id = $1))`

	t.Run("success", func(t *testing.T) {
		now := time.Now()
		rows := sqlmock.
			NewRows([]string{"id", "enrolled_face_id", "variation", "image", "created_at", "deleted_at"}).
			AddRow(id, 1, "this-is-variation", []byte{}, now, nil)
		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(id).
			WillReturnRows(rows)

		expectedRes := []*entity.FaceImage{
			{
				ID:             id,
				EnrolledFaceID: 1,
				Variation:      "this-is-variation",
				Image:          []byte{},
				CreatedAt:      now,
				DeletedAt:      nil,
			},
		}

		res, err := repo.GetDetailByEnrollID(context.TODO(), id)

		assert.Equal(t, nil, s.Mock.ExpectationsWereMet())
		assert.Equal(t, nil, err)
		assert.Nil(t, deep.Equal(expectedRes, res))

	})

	t.Run("error connection", func(t *testing.T) {
		s.Conn.Close()
		res, err := repo.GetDetailByEnrollID(context.TODO(), id)

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func Test_GetDetailWithoutImgByEnrollIDFaceImage(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewFaceImageRepository(s.Db)
	id := uint64(1)

	query := `SELECT id, enrolled_face_id, variation, created_at, deleted_at FROM "face_images" WHERE "face_images"."deleted_at" IS NULL AND ((enrolled_face_id = $1))`

	t.Run("success", func(t *testing.T) {
		now := time.Now()
		rows := sqlmock.
			NewRows([]string{"id", "enrolled_face_id", "variation", "created_at", "deleted_at"}).
			AddRow(id, 1, "this-is-variation", now, nil)
		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(id).
			WillReturnRows(rows)

		expectedRes := []*entity.FaceImage{
			{
				ID:             id,
				EnrolledFaceID: 1,
				Variation:      "this-is-variation",
				CreatedAt:      now,
				DeletedAt:      nil,
			},
		}

		res, err := repo.GetDetailWithoutImgByEnrollID(context.TODO(), id)

		assert.Equal(t, nil, s.Mock.ExpectationsWereMet())
		assert.Equal(t, nil, err)
		assert.Nil(t, deep.Equal(expectedRes, res))

	})

	t.Run("error connection", func(t *testing.T) {
		s.Conn.Close()
		res, err := repo.GetDetailWithoutImgByEnrollID(context.TODO(), id)

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func Test_DeleteAllFaceImage(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewFaceImageRepository(s.Db)

	query := `DELETE from "face_image`

	t.Run("success", func(t *testing.T) {
		result := sqlmock.NewResult(1, 1)
		s.Mock.
			ExpectExec(regexp.QuoteMeta(query)).
			WillReturnResult(result)

		err := repo.DeleteAll(context.TODO())

		assert.Equal(t, nil, s.Mock.ExpectationsWereMet())
		assert.Equal(t, nil, err)

	})

	t.Run("error connection", func(t *testing.T) {
		s.Conn.Close()
		err := repo.DeleteAll(context.TODO())
		assert.Error(t, err)
	})
}

func Test_GetImageByIDFaceImage(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewFaceImageRepository(s.Db)
	id := uint64(1)

	query := `SELECT * FROM "face_images"  WHERE "face_images"."deleted_at" IS NULL AND (("face_images"."id" = 1)) ORDER BY "face_images"."id" ASC LIMIT 1`

	t.Run("success", func(t *testing.T) {
		now := time.Now()
		rows := sqlmock.
			NewRows([]string{"id", "enrolled_face_id", "variation", "image", "image_thumbnail", "created_at", "deleted_at"}).
			AddRow(id, 1, "this-is-variation", []byte{}, []byte{}, now, nil)
		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnRows(rows)

		expectedRes := entity.FaceImage{
			ID:             id,
			EnrolledFaceID: 1,
			Variation:      "this-is-variation",
			Image:          []byte{},
			ImageThumbnail: []byte{},
			CreatedAt:      now,
			DeletedAt:      nil,
		}

		res, err := repo.GetImageByID(context.TODO(), id)

		assert.Equal(t, nil, s.Mock.ExpectationsWereMet())
		assert.Equal(t, nil, err)
		assert.Nil(t, deep.Equal(&expectedRes, res))

	})

	t.Run("error connection", func(t *testing.T) {
		s.Conn.Close()
		res, err := repo.GetImageByID(context.TODO(), id)

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func Test_GetDetailByVariationsImage(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewFaceImageRepository(s.Db)
	variation := []string{"variation-1", "variation-2"}

	query := `SELECT id, enrolled_face_id, variation, image, created_at, deleted_at FROM "face_images" WHERE "face_images"."deleted_at" IS NULL AND ((variation IN ($1,$2)))`

	t.Run("success", func(t *testing.T) {
		now := time.Now()
		rows := sqlmock.
			NewRows([]string{"id", "enrolled_face_id", "variation", "image", "created_at", "deleted_at"}).
			AddRow(1, 1, "variation-1", []byte{}, now, nil).
			AddRow(2, 1, "variation-2", []byte{}, now, nil)
		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnRows(rows)

		expectedRes := []*entity.FaceImage{
			{
				ID:             1,
				EnrolledFaceID: 1,
				Variation:      "variation-1",
				Image:          []byte{},
				CreatedAt:      now,
				DeletedAt:      nil,
			},
			{
				ID:             2,
				EnrolledFaceID: 1,
				Variation:      "variation-2",
				Image:          []byte{},
				CreatedAt:      now,
				DeletedAt:      nil,
			},
		}

		res, err := repo.GetDetailByVariations(context.TODO(), variation)

		assert.Equal(t, nil, s.Mock.ExpectationsWereMet())
		assert.Equal(t, nil, err)
		assert.Nil(t, deep.Equal(expectedRes, res))

	})

	t.Run("error connection", func(t *testing.T) {
		s.Conn.Close()
		res, err := repo.GetDetailByVariations(context.TODO(), variation)

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}
