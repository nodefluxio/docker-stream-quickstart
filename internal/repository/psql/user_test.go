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

func Test_NewUserRepository(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		s := helper.Setup()
		defer s.Conn.Close()

		repo := psqlrepo.NewUserRepository(s.Db)

		t.Run("initialized", func(t *testing.T) {
			require.NotNil(t, repo)
		})

		t.Run("implements valid interface", func(t *testing.T) {
			repositoryInterface := reflect.TypeOf((*repository.User)(nil)).Elem()
			repositoryStruct := reflect.TypeOf(repo)

			require.True(t, repositoryStruct.Implements(repositoryInterface))
		})
	})
}

func Test_CreateUser(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewUserRepository(s.Db)

	query := `INSERT INTO "user_access" ("email","username","password","fullname","avatar","role","site_id","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING "user_access"."id"`

	id := uint64(1)
	now := time.Now()
	data := entity.User{
		Email:     "testing@nodeflux.io",
		Username:  "testing",
		Password:  "super_secure_password",
		Fullname:  "unit testing",
		Avatar:    byte(0),
		Role:      string(entity.UserRoleSuperAdmin),
		SiteID:    []int64{1, 2, 3},
		CreatedAt: now,
		UpdatedAt: &now,
		DeletedAt: &now,
	}

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.
			NewRows([]string{"id"}).
			AddRow(id)
		s.Mock.ExpectBegin()
		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(
				data.Email,
				data.Username,
				data.Password,
				data.Fullname,
				data.Avatar,
				data.Role,
				data.SiteID,
				now,
				data.UpdatedAt,
				data.DeletedAt,
			).WillReturnRows(rows)
		s.Mock.ExpectCommit()

		res, err := repo.Create(context.TODO(), &data)

		var expectedRes = data
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

func Test_GetByUsernameUser(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewUserRepository(s.Db)

	query := `SELECT * FROM "user_access" WHERE "user_access"."deleted_at" IS NULL AND ((username = $1)) ORDER BY "user_access"."id" ASC LIMIT 1`

	id := uint64(1)
	now := time.Now()
	username := "testing"
	data := entity.User{
		Email:     "testing@nodeflux.io",
		Username:  "testing",
		Password:  "super_secure_password",
		Fullname:  "unit testing",
		Avatar:    byte(0),
		Role:      string(entity.UserRoleSuperAdmin),
		SiteID:    []int64{1, 2, 3},
		CreatedAt: now,
		UpdatedAt: &now,
		DeletedAt: &now,
	}

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.
			NewRows([]string{"id", "email", "username", "password", "fullname", "avatar", "role", "site_id", "created_at", "updated_at", "deleted_at"}).
			AddRow(
				id,
				data.Email,
				data.Username,
				data.Password,
				data.Fullname,
				data.Avatar,
				data.Role,
				data.SiteID,
				now,
				data.UpdatedAt,
				data.DeletedAt,
			)
		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(username).
			WillReturnRows(rows)

		res, err := repo.GetByUsername(context.TODO(), username)

		var expectedRes = data
		expectedRes.ID = id

		assert.Equal(t, nil, s.Mock.ExpectationsWereMet())
		assert.Equal(t, nil, err)
		assert.Nil(t, deep.Equal(&expectedRes, res))

	})

	t.Run("error connection", func(t *testing.T) {
		res, err := repo.GetByUsername(context.TODO(), username)

		assert.Error(t, err)
		assert.Nil(t, res)

	})
}

func Test_GetByEmail(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewUserRepository(s.Db)

	query := `SELECT * FROM "user_access" WHERE "user_access"."deleted_at" IS NULL AND ((email = $1)) ORDER BY "user_access"."id" ASC LIMIT 1`

	id := uint64(1)
	now := time.Now()
	email := "testing@nodeflux.io"
	data := entity.User{
		Email:     "testing@nodeflux.io",
		Username:  "testing",
		Password:  "super_secure_password",
		Fullname:  "unit testing",
		Avatar:    byte(0),
		Role:      string(entity.UserRoleSuperAdmin),
		SiteID:    []int64{1, 2, 3},
		CreatedAt: now,
		UpdatedAt: &now,
		DeletedAt: &now,
	}

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.
			NewRows([]string{"id", "email", "username", "password", "fullname", "avatar", "role", "site_id", "created_at", "updated_at", "deleted_at"}).
			AddRow(
				id,
				data.Email,
				data.Username,
				data.Password,
				data.Fullname,
				data.Avatar,
				data.Role,
				data.SiteID,
				now,
				data.UpdatedAt,
				data.DeletedAt,
			)
		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(email).
			WillReturnRows(rows)

		res, err := repo.GetByEmail(context.TODO(), email)

		var expectedRes = data
		expectedRes.ID = id

		assert.Equal(t, nil, s.Mock.ExpectationsWereMet())
		assert.Equal(t, nil, err)
		assert.Nil(t, deep.Equal(&expectedRes, res))

	})

	t.Run("error connection", func(t *testing.T) {
		res, err := repo.GetByEmail(context.TODO(), email)

		assert.Error(t, err)
		assert.Nil(t, res)

	})
}

func Test_UpdateBasicData(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewUserRepository(s.Db)

	query := `UPDATE "user_access" SET "email" = $1, "fullname" = $2, "id" = $3, "role" = $4, "site_id" = $5, "updated_at" = $6, "username" = $7  WHERE "user_access"."deleted_at" IS NULL AND "user_access"."id" = $8`

	now := time.Now()
	data := entity.User{
		ID:        1,
		Email:     "testing@nodeflux.io",
		Username:  "testing",
		Fullname:  "unit testing",
		Role:      string(entity.UserRoleSuperAdmin),
		SiteID:    []int64{1, 2, 3},
		CreatedAt: now,
		UpdatedAt: &now,
	}

	t.Run("success", func(t *testing.T) {
		s.Mock.ExpectBegin()
		s.Mock.
			ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(
				data.Email,
				data.Fullname,
				data.ID,
				data.Role,
				data.SiteID,
				data.UpdatedAt,
				data.Username,
				data.ID,
			).WillReturnResult(sqlmock.NewResult(1, 1))
		s.Mock.ExpectCommit()

		err := repo.UpdateBasicData(context.TODO(), &data)

		assert.Equal(t, nil, s.Mock.ExpectationsWereMet())
		assert.Equal(t, nil, err)

	})

	t.Run("error connection", func(t *testing.T) {
		err := repo.UpdateBasicData(context.TODO(), &data)
		assert.Error(t, err)
	})
}

func Test_UpdatePassword(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewUserRepository(s.Db)

	query := `UPDATE "user_access" SET "password" = $1, "updated_at" = $2 WHERE "user_access"."deleted_at" IS NULL AND ((id=$3))`

	id := uint64(1)
	password := "super-duper-strong-password"

	t.Run("success", func(t *testing.T) {
		s.Mock.ExpectBegin()
		s.Mock.
			ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(password, helper.AnyTime{}, id).WillReturnResult(sqlmock.NewResult(1, 1))
		s.Mock.ExpectCommit()

		err := repo.UpdatePassword(context.TODO(), password, id)

		assert.Equal(t, nil, s.Mock.ExpectationsWereMet())
		assert.Equal(t, nil, err)

	})

	t.Run("error connection", func(t *testing.T) {
		err := repo.UpdatePassword(context.TODO(), password, id)
		assert.Error(t, err)
	})
}

func Test_GetDetailUser(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewUserRepository(s.Db)

	query := `SELECT * FROM "user_access" WHERE "user_access"."deleted_at" IS NULL AND (("user_access"."id" = 1)) ORDER BY "user_access"."id" ASC LIMIT 1`

	id := uint64(1)
	now := time.Now()
	data := entity.User{
		Email:     "testing@nodeflux.io",
		Username:  "testing",
		Password:  "super_secure_password",
		Fullname:  "unit testing",
		Avatar:    byte(0),
		Role:      string(entity.UserRoleSuperAdmin),
		SiteID:    []int64{1, 2, 3},
		CreatedAt: now,
		UpdatedAt: &now,
		DeletedAt: &now,
	}

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.
			NewRows([]string{"id", "email", "username", "password", "fullname", "avatar", "role", "site_id", "created_at", "updated_at", "deleted_at"}).
			AddRow(
				id,
				data.Email,
				data.Username,
				data.Password,
				data.Fullname,
				data.Avatar,
				data.Role,
				data.SiteID,
				now,
				data.UpdatedAt,
				data.DeletedAt,
			)

		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnRows(rows)

		var expectedRes = data
		expectedRes.ID = id

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

func Test_DeleteUser(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewUserRepository(s.Db)

	query := `UPDATE "user_access" SET "deleted_at" = $1, "updated_at" = $2 WHERE "user_access"."deleted_at" IS NULL AND ((id=$3))`

	id := uint64(1)

	t.Run("success", func(t *testing.T) {
		s.Mock.ExpectBegin()
		s.Mock.
			ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(
				helper.AnyTime{},
				helper.AnyTime{},
				id,
			).WillReturnResult(sqlmock.NewResult(1, 1))
		s.Mock.ExpectCommit()

		err := repo.Delete(context.TODO(), id)

		assert.Equal(t, nil, s.Mock.ExpectationsWereMet())
		assert.Equal(t, nil, err)

	})

	t.Run("error connection", func(t *testing.T) {
		err := repo.Delete(context.TODO(), id)
		assert.Error(t, err)
	})
}

func Test_GetListUser(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewUserRepository(s.Db)
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		query := `SELECT * FROM "user_access" WHERE "user_access"."deleted_at" IS NULL ORDER BY created_at DESC`

		rows := sqlmock.
			NewRows([]string{"id", "email", "username", "password", "fullname", "avatar", "role", "site_id", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, "testing@nodeflux.io", "testing", "super_secure_password", "unit testing", byte(0), string(entity.UserRoleOperator), "{1,2,3}", now, now, nil).
			AddRow(2, "testing2@nodeflux.io", "testing2", "super_secure_password2", "unit testing 2", byte(0), string(entity.UserRoleSuperAdmin), "{1,2,3}", now, now, nil)

		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnRows(rows)

		filter := &util.Pagination{}
		res, err := repo.GetList(context.TODO(), filter)

		expectedRes := []*entity.User{
			{
				ID:        1,
				Email:     "testing@nodeflux.io",
				Username:  "testing",
				Password:  "super_secure_password",
				Fullname:  "unit testing",
				Avatar:    byte(0),
				Role:      string(entity.UserRoleOperator),
				SiteID:    []int64{1, 2, 3},
				CreatedAt: now,
				UpdatedAt: &now,
				DeletedAt: nil,
			},
			{
				ID:        2,
				Email:     "testing2@nodeflux.io",
				Username:  "testing2",
				Password:  "super_secure_password2",
				Fullname:  "unit testing 2",
				Avatar:    byte(0),
				Role:      string(entity.UserRoleSuperAdmin),
				SiteID:    []int64{1, 2, 3},
				CreatedAt: now,
				UpdatedAt: &now,
				DeletedAt: nil,
			},
		}

		assert.Equal(t, nil, s.Mock.ExpectationsWereMet())
		assert.Equal(t, nil, err)
		assert.Nil(t, deep.Equal(expectedRes, res))
	})

	t.Run("success with parameter search", func(t *testing.T) {
		query := `SELECT * FROM "user_access" WHERE "user_access"."deleted_at" IS NULL AND ((fullname ~* $1 OR email ~* $2)) ORDER BY created_at DESC`

		filter := &util.Pagination{
			Search: "testing2",
		}
		rows := sqlmock.
			NewRows([]string{"id", "email", "username", "password", "fullname", "avatar", "role", "site_id", "created_at", "updated_at", "deleted_at"}).
			AddRow(2, "testing2@nodeflux.io", "testing2", "super_secure_password2", "unit testing 2", byte(0), string(entity.UserRoleSuperAdmin), "{1,2,3}", now, now, nil)
		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(filter.Search, filter.Search).
			WillReturnRows(rows)

		res, err := repo.GetList(context.TODO(), filter)

		expectedRes := []*entity.User{
			{
				ID:        2,
				Email:     "testing2@nodeflux.io",
				Username:  "testing2",
				Password:  "super_secure_password2",
				Fullname:  "unit testing 2",
				Avatar:    byte(0),
				Role:      string(entity.UserRoleSuperAdmin),
				SiteID:    []int64{1, 2, 3},
				CreatedAt: now,
				UpdatedAt: &now,
				DeletedAt: nil,
			},
		}

		assert.Equal(t, nil, s.Mock.ExpectationsWereMet())
		assert.Equal(t, nil, err)
		assert.Nil(t, deep.Equal(expectedRes, res))

	})

	t.Run("success with parameter filter role", func(t *testing.T) {
		query := `SELECT * FROM "user_access" WHERE "user_access"."deleted_at" IS NULL AND ((role = $1)) ORDER BY created_at DESC`

		filter := &util.Pagination{
			Filter: map[string]string{
				"role": string(entity.UserRoleSuperAdmin),
			},
		}

		rows := sqlmock.
			NewRows([]string{"id", "email", "username", "password", "fullname", "avatar", "role", "site_id", "created_at", "updated_at", "deleted_at"}).
			AddRow(2, "testing2@nodeflux.io", "testing2", "super_secure_password2", "unit testing 2", byte(0), string(entity.UserRoleSuperAdmin), "{1,2,3}", now, now, nil)

		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(filter.Filter["role"]).
			WillReturnRows(rows)

		res, err := repo.GetList(context.TODO(), filter)

		expectedRes := []*entity.User{
			{
				ID:        2,
				Email:     "testing2@nodeflux.io",
				Username:  "testing2",
				Password:  "super_secure_password2",
				Fullname:  "unit testing 2",
				Avatar:    byte(0),
				Role:      string(entity.UserRoleSuperAdmin),
				SiteID:    []int64{1, 2, 3},
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
		query := `SELECT * FROM "user_access" WHERE "user_access"."deleted_at" IS NULL ORDER BY username ASC`

		filter := &util.Pagination{
			Sort: map[string]string{
				"username": "ASC",
			},
		}
		now := time.Now()
		rows := sqlmock.
			NewRows([]string{"id", "email", "username", "password", "fullname", "avatar", "role", "site_id", "created_at", "updated_at", "deleted_at"}).
			AddRow(2, "testing2@nodeflux.io", "testing2", "super_secure_password2", "unit testing 2", byte(0), string(entity.UserRoleSuperAdmin), "{1,2,3}", now, now, nil)

		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnRows(rows)

		expectedRes := []*entity.User{
			{
				ID:        2,
				Email:     "testing2@nodeflux.io",
				Username:  "testing2",
				Password:  "super_secure_password2",
				Fullname:  "unit testing 2",
				Avatar:    byte(0),
				Role:      string(entity.UserRoleSuperAdmin),
				SiteID:    []int64{1, 2, 3},
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
		query := `SELECT * FROM "user_access" WHERE "user_access"."deleted_at" IS NULL ORDER BY created_at DESC LIMIT 2 OFFSET 0`

		rows := sqlmock.
			NewRows([]string{"id", "email", "username", "password", "fullname", "avatar", "role", "site_id", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, "testing@nodeflux.io", "testing", "super_secure_password", "unit testing", byte(0), string(entity.UserRoleOperator), "{1,2,3}", now, now, nil).
			AddRow(2, "testing2@nodeflux.io", "testing2", "super_secure_password2", "unit testing 2", byte(0), string(entity.UserRoleSuperAdmin), "{1,2,3}", now, now, nil)

		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnRows(rows)

		filter := &util.Pagination{
			Limit:  2,
			Offset: 0,
		}
		res, err := repo.GetList(context.TODO(), filter)

		expectedRes := []*entity.User{
			{
				ID:        1,
				Email:     "testing@nodeflux.io",
				Username:  "testing",
				Password:  "super_secure_password",
				Fullname:  "unit testing",
				Avatar:    byte(0),
				Role:      string(entity.UserRoleOperator),
				SiteID:    []int64{1, 2, 3},
				CreatedAt: now,
				UpdatedAt: &now,
				DeletedAt: nil,
			},
			{
				ID:        2,
				Email:     "testing2@nodeflux.io",
				Username:  "testing2",
				Password:  "super_secure_password2",
				Fullname:  "unit testing 2",
				Avatar:    byte(0),
				Role:      string(entity.UserRoleSuperAdmin),
				SiteID:    []int64{1, 2, 3},
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

func Test_CountUser(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewUserRepository(s.Db)

	query := `SELECT count(*) FROM "user_access"  WHERE (deleted_at IS NULL) AND (fullname ~* $1 OR email ~* $2)`

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
			WithArgs(filter.Search, filter.Search).
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
