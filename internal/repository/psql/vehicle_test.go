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

func Test_NewVehicleRepository(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		s := helper.Setup()
		defer s.Conn.Close()

		repo := psqlrepo.NewVehicleRepository(s.Db)

		t.Run("initialized", func(t *testing.T) {
			require.NotNil(t, repo)
		})

		t.Run("implements valid interface", func(t *testing.T) {
			repositoryInterface := reflect.TypeOf((*repository.Vehicle)(nil)).Elem()
			repositoryStruct := reflect.TypeOf(repo)

			require.True(t, repositoryStruct.Implements(repositoryInterface))
		})
	})
}

func Test_CreateVehicle(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewVehicleRepository(s.Db)

	query := `INSERT INTO "vehicles" ("plate_number","type","brand","color","name","unique_id","status","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING "vehicles"."id"`

	id := uint64(1)
	now := time.Now()
	data := entity.Vehicle{
		Plate:     "B5721JHF",
		Type:      "motorcycle",
		Brand:     "honda",
		Color:     "hitam",
		Name:      "motor pegawai nodeflux",
		UniqueID:  "this_is_unique_id",
		Status:    "active",
		CreatedAt: now,
		UpdatedAt: &now,
		DeletedAt: nil,
	}

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.
			NewRows([]string{"id"}).
			AddRow(id)
		s.Mock.ExpectBegin()
		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(data.Plate, data.Type, data.Brand, data.Color, data.Name, data.UniqueID, data.Status, data.CreatedAt, data.UpdatedAt, data.DeletedAt).
			WillReturnRows(rows)
		s.Mock.ExpectCommit()

		res, err := repo.Create(context.TODO(), &data)

		expectedRes := data
		expectedRes.ID = id

		assert.Equal(t, nil, s.Mock.ExpectationsWereMet())
		assert.Equal(t, nil, err)
		assert.Equal(t, id, data.ID)
		assert.Nil(t, deep.Equal(&expectedRes, res))

	})

	t.Run("error connection", func(t *testing.T) {
		res, err := repo.Create(context.TODO(), &data)

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func Test_DeleteVehicle(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewVehicleRepository(s.Db)

	query := `UPDATE "vehicle" SET deleted_at = $1 WHERE id = $2`
	id := uint64(1)

	t.Run("success", func(t *testing.T) {
		result := sqlmock.NewResult(1, 1)
		s.Mock.
			ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(helper.AnyTime{}, id).
			WillReturnResult(result)

		err := repo.Delete(context.TODO(), id)

		assert.Equal(t, nil, s.Mock.ExpectationsWereMet())
		assert.Equal(t, nil, err)

	})

	t.Run("error connection", func(t *testing.T) {
		s.Conn.Close()
		err := repo.Delete(context.TODO(), id)
		assert.Error(t, err)
	})
}

func Test_GetAllVehicle(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewVehicleRepository(s.Db)

	query := `SELECT * FROM "vehicles" WHERE "vehicles"."deleted_at" IS NULL`

	t.Run("success", func(t *testing.T) {
		now := time.Now()
		rows := sqlmock.
			NewRows([]string{"id", "plate_number", "type", "brand", "color", "name", "unique_id", "status", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, "B5721JHF", "motorcycle", "honda", "hitam", "motor pegawai nodeflux", "this_is_unique_id", "active", now, now, nil).
			AddRow(2, "D5721JHF", "car", "honda", "hitam", "mobil pegawai nodeflux", "this_is_unique_id_2", "active", now, now, nil)

		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnRows(rows)

		expectedRes := []*entity.Vehicle{
			{
				ID:        1,
				Plate:     "B5721JHF",
				Type:      "motorcycle",
				Brand:     "honda",
				Color:     "hitam",
				Name:      "motor pegawai nodeflux",
				UniqueID:  "this_is_unique_id",
				Status:    "active",
				CreatedAt: now,
				UpdatedAt: &now,
				DeletedAt: nil,
			},
			{
				ID:        2,
				Plate:     "D5721JHF",
				Type:      "car",
				Brand:     "honda",
				Color:     "hitam",
				Name:      "mobil pegawai nodeflux",
				UniqueID:  "this_is_unique_id_2",
				Status:    "active",
				CreatedAt: now,
				UpdatedAt: &now,
				DeletedAt: nil,
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

func Test_DeleteAllVehicle(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewVehicleRepository(s.Db)

	query := `DELETE from "vehicle`

	t.Run("success", func(t *testing.T) {
		result := sqlmock.NewResult(1, 1)
		s.Mock.
			ExpectExec(regexp.QuoteMeta(query)).
			WillReturnResult(result)

		err := repo.DeleteAll(context.TODO())

		assert.Equal(t, nil, err)

	})

	t.Run("error connection", func(t *testing.T) {
		s.Conn.Close()
		err := repo.DeleteAll(context.TODO())
		assert.Error(t, err)
	})
}

func Test_GetListVehicle(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewVehicleRepository(s.Db)
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		query := `SELECT * FROM "vehicles" WHERE "vehicles"."deleted_at" IS NULL ORDER BY created_at DESC`

		rows := sqlmock.
			NewRows([]string{"id", "plate_number", "type", "brand", "color", "name", "unique_id", "status", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, "B5721JHF", "motorcycle", "honda", "hitam", "motor pegawai nodeflux", "this_is_unique_id", "active", now, now, nil)

		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnRows(rows)

		filter := &util.Pagination{}
		res, err := repo.GetList(context.TODO(), filter)

		expectedRes := []*entity.Vehicle{
			{
				ID:        1,
				Plate:     "B5721JHF",
				Type:      "motorcycle",
				Brand:     "honda",
				Color:     "hitam",
				Name:      "motor pegawai nodeflux",
				UniqueID:  "this_is_unique_id",
				Status:    "active",
				CreatedAt: now,
				UpdatedAt: &now,
				DeletedAt: nil,
			},
		}

		assert.Equal(t, nil, s.Mock.ExpectationsWereMet())
		assert.Equal(t, nil, err)
		assert.Nil(t, deep.Equal(expectedRes, res))
	})

	t.Run("success with parameter sort", func(t *testing.T) {
		query := `SELECT * FROM "vehicles" WHERE "vehicles"."deleted_at" IS NULL ORDER BY plate_number ASC`

		filter := &util.Pagination{
			Sort: map[string]string{
				"plate_number": "ASC",
			},
		}
		now := time.Now()
		rows := sqlmock.
			NewRows([]string{"id", "plate_number", "type", "brand", "color", "name", "unique_id", "status", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, "B5721JHF", "motorcycle", "honda", "hitam", "motor pegawai nodeflux", "this_is_unique_id", "active", now, now, nil)

		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnRows(rows)

		expectedRes := []*entity.Vehicle{
			{
				ID:        1,
				Plate:     "B5721JHF",
				Type:      "motorcycle",
				Brand:     "honda",
				Color:     "hitam",
				Name:      "motor pegawai nodeflux",
				UniqueID:  "this_is_unique_id",
				Status:    "active",
				CreatedAt: now,
				UpdatedAt: &now,
				DeletedAt: nil,
			},
		}
		res, err := repo.GetList(context.TODO(), filter)

		assert.Equal(t, nil, s.Mock.ExpectationsWereMet())
		assert.Equal(t, nil, err)
		assert.Nil(t, deep.Equal(expectedRes, res))

	})

	t.Run("success wtih parameter pagination (limit and offset)", func(t *testing.T) {
		query := `SELECT * FROM "vehicles" WHERE "vehicles"."deleted_at" IS NULL ORDER BY created_at DESC LIMIT 2 OFFSET 0`

		rows := sqlmock.
			NewRows([]string{"id", "plate_number", "type", "brand", "color", "name", "unique_id", "status", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, "B5721JHF", "motorcycle", "honda", "hitam", "motor pegawai nodeflux", "this_is_unique_id", "active", now, now, nil)

		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnRows(rows)

		filter := &util.Pagination{
			Limit:  2,
			Offset: 0,
		}
		res, err := repo.GetList(context.TODO(), filter)

		expectedRes := []*entity.Vehicle{
			{
				ID:        1,
				Plate:     "B5721JHF",
				Type:      "motorcycle",
				Brand:     "honda",
				Color:     "hitam",
				Name:      "motor pegawai nodeflux",
				UniqueID:  "this_is_unique_id",
				Status:    "active",
				CreatedAt: now,
				UpdatedAt: &now,
				DeletedAt: nil,
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

func Test_CountVehicle(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewVehicleRepository(s.Db)

	query := `SELECT count(*) FROM "vehicle" WHERE (deleted_at IS NULL) AND (name ~* $1 OR plate_number ~* $2 OR type ~* $3 OR brand ~* $4 OR name ~* $5 OR status ~* $6 OR color ~* $7 OR unique_id ~* $8 )`

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
			WithArgs(filter.Search, filter.Search, filter.Search, filter.Search, filter.Search, filter.Search, filter.Search, filter.Search).
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

func Test_GetDetailVehicle(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewVehicleRepository(s.Db)
	id := uint64(1)

	query := `SELECT * FROM "vehicles"  WHERE "vehicles"."deleted_at" IS NULL AND (("vehicles"."id" = 1)) ORDER BY "vehicles"."id" ASC LIMIT 1`

	t.Run("success", func(t *testing.T) {
		now := time.Now()
		rows := sqlmock.
			NewRows([]string{"id", "plate_number", "type", "brand", "color", "name", "unique_id", "status", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, "B5721JHF", "motorcycle", "honda", "hitam", "motor pegawai nodeflux", "this_is_unique_id", "active", now, now, nil)
		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnRows(rows)

		expectedRes := entity.Vehicle{
			ID:        1,
			Plate:     "B5721JHF",
			Type:      "motorcycle",
			Brand:     "honda",
			Color:     "hitam",
			Name:      "motor pegawai nodeflux",
			UniqueID:  "this_is_unique_id",
			Status:    "active",
			CreatedAt: now,
			UpdatedAt: &now,
			DeletedAt: nil,
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

func Test_GetByPlateNumberlVehicle(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewVehicleRepository(s.Db)
	plate := "B5721JHF"

	query := `SELECT * FROM "vehicles" WHERE "vehicles"."deleted_at" IS NULL AND ((plate_number = $1)) ORDER BY "vehicles"."id" ASC LIMIT 1`

	t.Run("success", func(t *testing.T) {
		now := time.Now()
		rows := sqlmock.
			NewRows([]string{"id", "plate_number", "type", "brand", "color", "name", "unique_id", "status", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, "B5721JHF", "motorcycle", "honda", "hitam", "motor pegawai nodeflux", "this_is_unique_id", "active", now, now, nil)
		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(plate).
			WillReturnRows(rows)

		expectedRes := entity.Vehicle{
			ID:        1,
			Plate:     "B5721JHF",
			Type:      "motorcycle",
			Brand:     "honda",
			Color:     "hitam",
			Name:      "motor pegawai nodeflux",
			UniqueID:  "this_is_unique_id",
			Status:    "active",
			CreatedAt: now,
			UpdatedAt: &now,
			DeletedAt: nil,
		}

		res, err := repo.GetByPlateNumber(context.TODO(), plate)

		assert.Equal(t, nil, s.Mock.ExpectationsWereMet())
		assert.Equal(t, nil, err)
		assert.Nil(t, deep.Equal(&expectedRes, res))

	})

	t.Run("error connection", func(t *testing.T) {
		s.Conn.Close()
		res, err := repo.GetByPlateNumber(context.TODO(), plate)

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func Test_UpdateVehicle(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewVehicleRepository(s.Db)

	query := `INSERT INTO "vehicles" ("plate_number","type","brand","color","name","unique_id","status","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING "vehicles"."id"`

	id := uint64(1)
	now := time.Now()
	data := entity.Vehicle{
		Plate:     "B5721JHF",
		Type:      "motorcycle",
		Brand:     "honda",
		Color:     "hitam",
		Name:      "motor pegawai nodeflux",
		UniqueID:  "this_is_unique_id",
		Status:    "active",
		CreatedAt: now,
		UpdatedAt: &now,
		DeletedAt: nil,
	}

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.
			NewRows([]string{"id"}).
			AddRow(id)
		s.Mock.ExpectBegin()
		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(data.Plate, data.Type, data.Brand, data.Color, data.Name, data.UniqueID, data.Status, data.CreatedAt, data.UpdatedAt, data.DeletedAt).
			WillReturnRows(rows)
		s.Mock.ExpectCommit()

		err := repo.Update(context.TODO(), &data)

		assert.Equal(t, nil, s.Mock.ExpectationsWereMet())
		assert.Equal(t, nil, err)
		assert.Equal(t, id, data.ID)

	})

	t.Run("error connection", func(t *testing.T) {
		s.Conn.Close()
		err := repo.Update(context.TODO(), &data)
		assert.Error(t, err)
	})
}
