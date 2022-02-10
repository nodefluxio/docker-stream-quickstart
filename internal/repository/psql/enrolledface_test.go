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

func Test_NewEnrollFaceRepository(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		s := helper.Setup()
		defer s.Conn.Close()

		repo := psqlrepo.NewEnrollFaceRepository(s.Db)

		t.Run("initialized", func(t *testing.T) {
			require.NotNil(t, repo)
		})

		t.Run("implements valid interface", func(t *testing.T) {
			repositoryInterface := reflect.TypeOf((*repository.EnrolledFace)(nil)).Elem()
			repositoryStruct := reflect.TypeOf(repo)

			require.True(t, repositoryStruct.Implements(repositoryInterface))
		})
	})
}

func Test_GetAllEnrolledFace(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewEnrollFaceRepository(s.Db)

	query := `SELECT * FROM "enrolled_faces" WHERE "enrolled_faces"."deleted_at" IS NULL`

	t.Run("success", func(t *testing.T) {
		now := time.Now()
		rows := sqlmock.
			NewRows([]string{"id", "face_id", "name", "identity_number", "status", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, 123456789, "testing", "99887766554433211", "active", now, now, nil)
		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnRows(rows)

		expectedRes := []*entity.EnrolledFace{
			{
				ID:             1,
				FaceID:         123456789,
				Name:           "testing",
				IdentityNumber: "99887766554433211",
				Status:         "active",
				CreatedAt:      now,
				UpdatedAt:      &now,
				DeletedAt:      nil,
			},
		}

		res, err := repo.GetAll(context.TODO())

		assert.Equal(t, nil, s.Mock.ExpectationsWereMet())
		assert.Equal(t, nil, err)
		assert.Nil(t, deep.Equal(expectedRes, res))

	})

	t.Run("error connection", func(t *testing.T) {
		s.Conn.Close()
		res, err := repo.GetAll(context.TODO())

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func Test_DeleteAllEnrollFace(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewEnrollFaceRepository(s.Db)

	query := `DELETE from "enrolled_face"`

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

func Test_GetList(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewEnrollFaceRepository(s.Db)
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		query := `SELECT * FROM "enrolled_faces" WHERE "enrolled_faces"."deleted_at" IS NULL ORDER BY created_at DESC`

		rows := sqlmock.
			NewRows([]string{"id", "face_id", "name", "identity_number", "status", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, 123456789, "testing", "99887766554433211", "active", now, now, nil)

		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnRows(rows)

		filter := &util.Pagination{}
		res, err := repo.GetList(context.TODO(), filter)

		expectedRes := []*entity.EnrolledFace{
			{
				ID:             1,
				FaceID:         123456789,
				Name:           "testing",
				IdentityNumber: "99887766554433211",
				Status:         "active",
				CreatedAt:      now,
				UpdatedAt:      &now,
				DeletedAt:      nil,
			},
		}

		assert.Equal(t, nil, s.Mock.ExpectationsWereMet())
		assert.Equal(t, nil, err)
		assert.Nil(t, deep.Equal(expectedRes, res))
	})

	t.Run("success with parameter sort", func(t *testing.T) {
		query := `SELECT * FROM "enrolled_faces" WHERE "enrolled_faces"."deleted_at" IS NULL ORDER BY name ASC`

		filter := &util.Pagination{
			Sort: map[string]string{
				"name": "ASC",
			},
		}

		now := time.Now()
		rows := sqlmock.
			NewRows([]string{"id", "face_id", "name", "identity_number", "status", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, 123456789, "testing", "99887766554433211", "active", now, now, nil)

		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnRows(rows)

		expectedRes := []*entity.EnrolledFace{
			{
				ID:             1,
				FaceID:         123456789,
				Name:           "testing",
				IdentityNumber: "99887766554433211",
				Status:         "active",
				CreatedAt:      now,
				UpdatedAt:      &now,
				DeletedAt:      nil,
			},
		}
		res, err := repo.GetList(context.TODO(), filter)

		assert.Equal(t, nil, s.Mock.ExpectationsWereMet())
		assert.Equal(t, nil, err)
		assert.Nil(t, deep.Equal(expectedRes, res))
	})

	t.Run("success wtih parameter pagination (limit and offset)", func(t *testing.T) {
		query := `SELECT * FROM "enrolled_faces" WHERE "enrolled_faces"."deleted_at" IS NULL ORDER BY created_at DESC LIMIT 2 OFFSET 0`

		rows := sqlmock.
			NewRows([]string{"id", "face_id", "name", "identity_number", "status", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, 123456789, "testing", "99887766554433211", "active", now, now, nil)

		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnRows(rows)

		filter := &util.Pagination{
			Limit:  2,
			Offset: 0,
		}
		res, err := repo.GetList(context.TODO(), filter)

		expectedRes := []*entity.EnrolledFace{
			{
				ID:             1,
				FaceID:         123456789,
				Name:           "testing",
				IdentityNumber: "99887766554433211",
				Status:         "active",
				CreatedAt:      now,
				UpdatedAt:      &now,
				DeletedAt:      nil,
			},
		}

		assert.Equal(t, nil, s.Mock.ExpectationsWereMet())
		assert.Equal(t, nil, err)
		assert.Nil(t, deep.Equal(expectedRes, res))
	})

	t.Run("error connection", func(t *testing.T) {
		s.Conn.Close()
		filter := &util.Pagination{}
		res, err := repo.GetList(context.TODO(), filter)

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func Test_CountEnrolledFace(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewEnrollFaceRepository(s.Db)

	query := `SELECT count(*) FROM "enrolled_face" WHERE (deleted_at IS NULL) AND (name ~* $1)`

	filter := &util.Pagination{
		Search: "testing2",
	}
	t.Run("success", func(t *testing.T) {
		totalData := 10
		rows := sqlmock.
			NewRows([]string{"count"}).
			AddRow(totalData)
		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(filter.Search).
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

func Test_GetDetailEnrolledFace(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewEnrollFaceRepository(s.Db)
	id := uint64(123456789)

	query := `SELECT * FROM "enrolled_faces" WHERE "enrolled_faces"."deleted_at" IS NULL AND (("enrolled_faces"."id" = 123456789)) ORDER BY "enrolled_faces"."id" ASC LIMIT 1`

	t.Run("success", func(t *testing.T) {
		now := time.Now()
		rows := sqlmock.
			NewRows([]string{"id", "face_id", "name", "identity_number", "status", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, 123456789, "testing", "99887766554433211", "active", now, now, nil)
		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnRows(rows)

		expectedRes := entity.EnrolledFace{
			ID:             1,
			FaceID:         123456789,
			Name:           "testing",
			IdentityNumber: "99887766554433211",
			Status:         "active",
			CreatedAt:      now,
			UpdatedAt:      &now,
			DeletedAt:      nil,
		}

		res, err := repo.GetDetail(context.TODO(), id)

		assert.Equal(t, nil, s.Mock.ExpectationsWereMet())
		assert.Equal(t, nil, err)
		assert.Nil(t, deep.Equal(&expectedRes, res))

	})

	t.Run("error connection", func(t *testing.T) {
		s.Conn.Close()
		res, err := repo.GetDetail(context.TODO(), id)

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func Test_GetDetailwFaceIDEnrolledFace(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewEnrollFaceRepository(s.Db)
	faceID := uint64(123456789)

	query := `SELECT enrolled_face.*, face_image.image_thumbnail as image FROM "enrolled_face" INNER JOIN face_image on face_image.enrolled_face_id = enrolled_face.id WHERE (enrolled_face.deleted_at is NULL AND enrolled_face.face_id= $1) ORDER BY enrolled_face ASC LIMIT 1`

	t.Run("success", func(t *testing.T) {
		now := time.Now()
		rows := sqlmock.
			NewRows([]string{"id", "face_id", "name", "identity_number", "status", "image", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, 123456789, "testing", "99887766554433211", "active", []byte{}, now, now, nil)
		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnRows(rows)

		expectedRes := entity.EnrolledFaceWithImage{
			ID:             1,
			FaceID:         123456789,
			Name:           "testing",
			IdentityNumber: "99887766554433211",
			Status:         "active",
			Image:          []byte{},
			CreatedAt:      now,
			UpdatedAt:      &now,
			DeletedAt:      nil,
		}

		res, err := repo.GetDetailwFaceID(context.TODO(), faceID)

		assert.Equal(t, nil, s.Mock.ExpectationsWereMet())
		assert.Equal(t, nil, err)
		assert.Nil(t, deep.Equal(&expectedRes, res))

	})

	t.Run("error connection", func(t *testing.T) {
		s.Conn.Close()
		res, err := repo.GetDetailwFaceID(context.TODO(), faceID)

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}
