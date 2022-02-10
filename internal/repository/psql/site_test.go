package psql_test

import (
	"context"
	"fmt"
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

func Test_NewSiteRepository(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		s := helper.Setup()
		defer s.Conn.Close()

		repo := psqlrepo.NewSiteRepository(s.Db)

		t.Run("initialized", func(t *testing.T) {
			require.NotNil(t, repo)
		})

		t.Run("implements valid interface", func(t *testing.T) {
			repositoryInterface := reflect.TypeOf((*repository.Site)(nil)).Elem()
			repositoryStruct := reflect.TypeOf(repo)

			require.True(t, repositoryStruct.Implements(repositoryInterface))
		})
	})
}

func Test_GetListSite(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewSiteRepository(s.Db)

	t.Run("success", func(t *testing.T) {
		query := `SELECT site.* FROM "sites" LEFT join map_site_stream on site.id = map_site_stream.site_id WHERE "sites"."deleted_at" IS NULL GROUP BY map_site_stream.site_id, site.name, site.id ORDER BY site.created_at DESC`

		now := time.Now()
		rows := sqlmock.
			NewRows([]string{"id", "name", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, "site 1", now, now, nil).
			AddRow(2, "site 2", now, now, nil)

		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnRows(rows)

		var expectedRes = []*entity.Site{
			{
				ID:        1,
				Name:      "site 1",
				CreatedAt: now,
				UpdatedAt: &now,
				DeletedAt: nil,
			},
			{
				ID:        2,
				Name:      "site 2",
				CreatedAt: now,
				UpdatedAt: &now,
				DeletedAt: nil,
			},
		}
		filter := &util.Pagination{}
		res, err := repo.GetList(context.TODO(), filter)

		assert.Equal(t, nil, s.Mock.ExpectationsWereMet())
		assert.Equal(t, nil, err)
		assert.Nil(t, deep.Equal(expectedRes, res))
	})

	t.Run("success with parameter search", func(t *testing.T) {
		query := `SELECT site.* FROM "sites" LEFT join map_site_stream on site.id = map_site_stream.site_id WHERE "sites"."deleted_at" IS NULL AND ((site.name ~* $1)) GROUP BY map_site_stream.site_id, site.name, site.id ORDER BY site.created_at DESC`

		filter := &util.Pagination{
			Search: "site 1",
		}
		now := time.Now()
		rows := sqlmock.
			NewRows([]string{"id", "name", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, "site 1", now, now, nil)

		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(filter.Search).
			WillReturnRows(rows)

		var expectedRes = []*entity.Site{
			{
				ID:        1,
				Name:      "site 1",
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

	t.Run("success with parameter filter stream_id", func(t *testing.T) {
		query := `SELECT site.* FROM "sites" LEFT join map_site_stream on site.id = map_site_stream.site_id WHERE "sites"."deleted_at" IS NULL AND ((map_site_stream.stream_id = $1)) GROUP BY map_site_stream.site_id, site.name, site.id ORDER BY site.created_at DESC`

		filter := &util.Pagination{
			Filter: map[string]string{
				"stream_id": "this-is-stream-id",
			},
		}
		now := time.Now()
		rows := sqlmock.
			NewRows([]string{"id", "name", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, "site 1", now, now, nil)

		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(filter.Filter["stream_id"]).
			WillReturnRows(rows)

		var expectedRes = []*entity.Site{
			{
				ID:        1,
				Name:      "site 1",
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

	t.Run("success with parameter filter one site_id", func(t *testing.T) {
		query := `SELECT site.* FROM "sites" LEFT join map_site_stream on site.id = map_site_stream.site_id WHERE "sites"."deleted_at" IS NULL AND ((site.id = $1)) GROUP BY map_site_stream.site_id, site.name, site.id ORDER BY site.created_at DESC`

		filter := &util.Pagination{
			Filter: map[string]string{
				"site_id": "1",
			},
		}
		now := time.Now()
		rows := sqlmock.
			NewRows([]string{"id", "name", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, "site 1", now, now, nil)

		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(filter.Filter["site_id"]).
			WillReturnRows(rows)

		var expectedRes = []*entity.Site{
			{
				ID:        1,
				Name:      "site 1",
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

	t.Run("success with parameter filter multiple site_id", func(t *testing.T) {
		query := `SELECT site.* FROM "sites" LEFT join map_site_stream on site.id = map_site_stream.site_id WHERE "sites"."deleted_at" IS NULL AND ((site.id IN ($1,$2))) GROUP BY map_site_stream.site_id, site.name, site.id ORDER BY site.created_at DESC`

		filter := &util.Pagination{
			Filter: map[string]string{
				"site_id": "1,2",
			},
		}
		now := time.Now()
		rows := sqlmock.
			NewRows([]string{"id", "name", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, "site 1", now, now, nil).
			AddRow(2, "site 2", now, now, nil)

		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs("1", "2").
			WillReturnRows(rows)

		var expectedRes = []*entity.Site{
			{
				ID:        1,
				Name:      "site 1",
				CreatedAt: now,
				UpdatedAt: &now,
				DeletedAt: nil,
			},
			{
				ID:        2,
				Name:      "site 2",
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

	t.Run("success with parameter sort", func(t *testing.T) {
		query := `SELECT site.* FROM "sites" LEFT join map_site_stream on site.id = map_site_stream.site_id WHERE "sites"."deleted_at" IS NULL GROUP BY map_site_stream.site_id, site.name, site.id ORDER BY name ASC`

		filter := &util.Pagination{
			Sort: map[string]string{
				"name": "ASC",
			},
		}
		now := time.Now()
		rows := sqlmock.
			NewRows([]string{"id", "name", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, "site 1", now, now, nil)

		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnRows(rows)

		var expectedRes = []*entity.Site{
			{
				ID:        1,
				Name:      "site 1",
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

	t.Run("error connection", func(t *testing.T) {
		s.Conn.Close()
		filter := &util.Pagination{}
		res, err := repo.GetList(context.TODO(), filter)

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func Test_GetDetailSite(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewSiteRepository(s.Db)

	query := `SELECT * FROM "sites" WHERE "sites"."deleted_at" IS NULL AND (("sites"."id" = 1)) LIMIT 1`

	id := uint64(1)
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.
			NewRows([]string{"id", "name", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, "site 1", now, now, nil)

		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnRows(rows)

		var expectedRes = &entity.Site{
			ID:        1,
			Name:      "site 1",
			CreatedAt: now,
			UpdatedAt: &now,
			DeletedAt: nil,
		}

		res, err := repo.GetDetail(context.TODO(), id)

		assert.Equal(t, nil, s.Mock.ExpectationsWereMet())
		assert.Equal(t, nil, err)
		assert.Nil(t, deep.Equal(expectedRes, res))

	})

	t.Run("error connection", func(t *testing.T) {
		s.Conn.Close()
		res, err := repo.GetDetail(context.TODO(), id)

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func Test_CreateSite(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewSiteRepository(s.Db)

	query := `INSERT INTO "sites" ("name","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4) RETURNING "sites"."id"`

	id := uint64(1)
	now := time.Now()
	data := entity.Site{
		Name:      "Site Name 1",
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
			WithArgs(data.Name, now, data.UpdatedAt, data.DeletedAt).
			WillReturnRows(rows)
		s.Mock.ExpectCommit()

		res, err := repo.Create(context.TODO(), &data)

		var expectedRes = &entity.Site{
			ID:        id,
			Name:      data.Name,
			CreatedAt: now,
			UpdatedAt: &now,
			DeletedAt: &now,
		}

		assert.Equal(t, nil, s.Mock.ExpectationsWereMet())
		assert.Equal(t, nil, err)
		assert.Equal(t, id, data.ID)
		assert.Nil(t, deep.Equal(expectedRes, res))

	})

	t.Run("error connection", func(t *testing.T) {
		res, err := repo.Create(context.TODO(), &data)

		assert.Error(t, err)
		assert.Nil(t, res)

	})
}

func Test_UpdateSite(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewSiteRepository(s.Db)

	query := `INSERT INTO "sites" ("name","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4) RETURNING "sites"."id"`

	id := uint64(1)
	now := time.Now()
	data := entity.Site{
		Name:      "Site Name 2",
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
			WithArgs(data.Name, now, data.UpdatedAt, data.DeletedAt).
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

func Test_DeteleteSite(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewSiteRepository(s.Db)

	query := `UPDATE "site" SET deleted_at = $1 WHERE id = $2`
	id := uint64(1)

	t.Run("success", func(t *testing.T) {
		result := sqlmock.NewResult(1, 1)
		s.Mock.
			ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(helper.AnyTime{}, id).
			WillReturnResult(result)

		err := repo.Delete(context.TODO(), id)

		assert.Equal(t, nil, err)

	})

	t.Run("error connection", func(t *testing.T) {
		s.Conn.Close()
		err := repo.Delete(context.TODO(), id)
		assert.Error(t, err)
	})
}

func Test_AddStreamToSite(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewSiteRepository(s.Db)

	query := `INSERT INTO "map_site_streams" ("site_id","stream_id","created_at") VALUES ($1,$2,$3) RETURNING "map_site_streams"."id"`

	id := uint64(1)
	now := time.Now()
	data := entity.MapSiteStream{
		SiteID:    1,
		StreamID:  "this-is-stream-id",
		CreatedAt: now,
	}

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.
			NewRows([]string{"id"}).
			AddRow(id)
		s.Mock.ExpectBegin()
		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(data.SiteID, data.StreamID, now).
			WillReturnRows(rows)
		s.Mock.ExpectCommit()

		err := repo.AddStreamToSite(context.TODO(), &data)

		assert.Equal(t, nil, s.Mock.ExpectationsWereMet())
		assert.Equal(t, nil, err)
		assert.Equal(t, id, data.ID)

	})

	t.Run("error connection", func(t *testing.T) {
		s.Conn.Close()
		err := repo.AddStreamToSite(context.TODO(), &data)
		assert.Error(t, err)
	})
}

func Test_GetSiteByIDs(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewSiteRepository(s.Db)

	query := `SELECT * FROM "site" WHERE "site"."deleted_at" IS NULL AND ((id IN ($1,$2)))`

	IDs := []int64{1, 2}
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.
			NewRows([]string{"id", "name", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, "site 1", now, now, nil).
			AddRow(2, "site 2", now, now, nil)

		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(IDs[0], IDs[1]).
			WillReturnRows(rows)

		var expectedRes = []*entity.Site{
			{
				ID:        1,
				Name:      "site 1",
				CreatedAt: now,
				UpdatedAt: &now,
				DeletedAt: nil,
			},
			{
				ID:        2,
				Name:      "site 2",
				CreatedAt: now,
				UpdatedAt: &now,
				DeletedAt: nil,
			},
		}

		res, err := repo.GetSiteByIDs(context.TODO(), IDs)

		fmt.Println("res :", res)
		assert.Equal(t, nil, s.Mock.ExpectationsWereMet())
		assert.Equal(t, nil, err)
		assert.Nil(t, deep.Equal(expectedRes, res))

	})

	t.Run("error connection", func(t *testing.T) {
		s.Conn.Close()
		res, err := repo.GetSiteByIDs(context.TODO(), IDs)

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func Test_GetSiteWithStreams(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewSiteRepository(s.Db)

	query := `SELECT site.id, site.name, map_site_stream.stream_id, site.created_at, site.updated_at, site.deleted_at FROM "site" INNER JOIN map_site_stream on site.id = map_site_stream.site_id WHERE "site"."deleted_at" IS NULL AND ((site.id IN ($1,$2)))`

	IDs := []int64{1, 2}
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.
			NewRows([]string{"id", "name", "stream_id", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, "site 1", "stream-id-site-1", now, now, nil).
			AddRow(2, "site 2", "stream-id-site-2", now, now, nil)

		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(IDs[0], IDs[1]).
			WillReturnRows(rows)

		var expectedRes = []*entity.SiteWithStream{
			{
				ID:        1,
				Name:      "site 1",
				StreamID:  "stream-id-site-1",
				CreatedAt: now,
				UpdatedAt: &now,
				DeletedAt: nil,
			},
			{
				ID:        2,
				Name:      "site 2",
				StreamID:  "stream-id-site-2",
				CreatedAt: now,
				UpdatedAt: &now,
				DeletedAt: nil,
			},
		}

		res, err := repo.GetSiteWithStream(context.TODO(), IDs)

		fmt.Println("res :", res)
		assert.Equal(t, nil, s.Mock.ExpectationsWereMet())
		assert.Equal(t, nil, err)
		assert.Nil(t, deep.Equal(expectedRes, res))

	})

	t.Run("error connection", func(t *testing.T) {
		s.Conn.Close()
		res, err := repo.GetSiteWithStream(context.TODO(), IDs)

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func Test_GetDetailByStreamID(t *testing.T) {
	s := helper.Setup()
	defer s.Conn.Close()

	repo := psqlrepo.NewSiteRepository(s.Db)

	query := `SELECT site.* FROM "sites" INNER join map_site_stream on site.id = map_site_stream.site_id WHERE "sites"."deleted_at" IS NULL AND ((map_site_stream.stream_id=$1)) LIMIT 1`

	StreamID := "stream-id"
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.
			NewRows([]string{"id", "name", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, "site 1", now, now, nil)

		s.Mock.
			ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(StreamID).
			WillReturnRows(rows)

		var expectedRes = &entity.Site{
			ID:        1,
			Name:      "site 1",
			CreatedAt: now,
			UpdatedAt: &now,
			DeletedAt: nil,
		}

		res, err := repo.GetDetailByStreamID(context.TODO(), StreamID)

		assert.Equal(t, nil, s.Mock.ExpectationsWereMet())
		assert.Equal(t, nil, err)
		assert.Nil(t, deep.Equal(expectedRes, res))

	})

	t.Run("error connection", func(t *testing.T) {
		s.Conn.Close()
		res, err := repo.GetDetailByStreamID(context.TODO(), StreamID)

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}
