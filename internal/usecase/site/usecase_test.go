package site_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	logutil "gitlab.com/nodefluxio/goutils/pkg/log"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/pkg/util"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/presenter"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/repository/mocks"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/usecase/site"
)

func TestGetList(t *testing.T) {
	logutil.Init("info")
	siteRepo := &mocks.Site{}
	siteService := site.ServiceImpl{
		SiteRepo: siteRepo,
	}

	var fakeListSite []*entity.Site

	t.Run("Failed GetList Site", func(t *testing.T) {
		userInfo := presenter.AuthInfoResponse{}
		paging := util.Pagination{}
		siteRepo.On("GetList", mock.Anything, mock.Anything).Once().Return(nil, errors.New("error when get list site"))
		_, err := siteService.GetList(context.Background(), &paging, &userInfo)
		assert.NotNil(t, err)
	})

	t.Run("Success GetList Site with role admin", func(t *testing.T) {
		userInfo := presenter.AuthInfoResponse{
			Role: string(entity.UserRoleSuperAdmin),
		}
		paging := util.Pagination{
			Filter: map[string]string{
				"site_id": "1",
			},
		}
		siteRepo.On("GetList", mock.Anything, mock.Anything).Once().Return(fakeListSite, nil)
		_, err := siteService.GetList(context.Background(), &paging, &userInfo)
		assert.Nil(t, err)
	})

	t.Run("Success GetList Site with role operator and allowed site id empty", func(t *testing.T) {
		userInfo := presenter.AuthInfoResponse{
			Role: string(entity.UserRoleOperator),
		}
		paging := util.Pagination{
			Filter: map[string]string{
				"site_id": "1",
			},
		}
		_, err := siteService.GetList(context.Background(), &paging, &userInfo)
		assert.Nil(t, err)
	})

	t.Run("Success GetList Site with role operator", func(t *testing.T) {
		userInfo := presenter.AuthInfoResponse{
			Role:   string(entity.UserRoleOperator),
			SiteID: []int64{1, 2, 3},
		}
		paging := util.Pagination{
			Filter: map[string]string{
				"site_id": "1",
			},
		}
		siteRepo.On("GetList", mock.Anything, mock.Anything).Once().Return(fakeListSite, nil)
		_, err := siteService.GetList(context.Background(), &paging, &userInfo)
		assert.Nil(t, err)
	})

}

func TestCreate(t *testing.T) {
	logutil.Init("info")
	siteRepo := &mocks.Site{}
	siteService := site.ServiceImpl{
		SiteRepo: siteRepo,
	}

	t.Run("Success Create Site", func(t *testing.T) {
		siteRepo.On("Create", mock.Anything, mock.Anything).Once().Return(&entity.Site{}, nil)
		_, err := siteService.Create(context.Background(), &presenter.SiteRequest{
			Name: "Site 1",
		})
		assert.Nil(t, err)
	})

	t.Run("Failed Create Site", func(t *testing.T) {
		siteRepo.On("Create", mock.Anything, mock.Anything).Once().Return(nil, errors.New("error on create site"))
		_, err := siteService.Create(context.Background(), &presenter.SiteRequest{
			Name: "Site 1",
		})
		assert.NotNil(t, err)
	})
}

func TestUpdate(t *testing.T) {
	logutil.Init("info")
	siteRepo := &mocks.Site{}
	siteService := site.ServiceImpl{
		SiteRepo: siteRepo,
	}

	t.Run("Success Update Site", func(t *testing.T) {
		siteRepo.On("GetDetail", mock.Anything, mock.Anything).Once().Return(&entity.Site{}, nil)
		siteRepo.On("Update", mock.Anything, mock.Anything).Once().Return(nil)
		err := siteService.Update(context.Background(), &presenter.SiteRequest{
			ID:   1,
			Name: "Site 1",
		})
		assert.Nil(t, err)
	})

	t.Run("Failed GetDetail Site", func(t *testing.T) {
		siteRepo.On("GetDetail", mock.Anything, mock.Anything).Once().Return(nil, errors.New("error on get detail site"))
		err := siteService.Update(context.Background(), &presenter.SiteRequest{
			ID:   1,
			Name: "Site 1",
		})
		assert.NotNil(t, err)
	})

	t.Run("Failed Update Site", func(t *testing.T) {
		siteRepo.On("GetDetail", mock.Anything, mock.Anything).Once().Return(&entity.Site{}, nil)
		siteRepo.On("Update", mock.Anything, mock.Anything).Once().Return(errors.New("error on update site"))
		err := siteService.Update(context.Background(), &presenter.SiteRequest{
			ID:   1,
			Name: "Site 1",
		})
		assert.NotNil(t, err)
	})
}

func TestDelete(t *testing.T) {
	logutil.Init("info")
	siteRepo := &mocks.Site{}
	siteService := site.ServiceImpl{
		SiteRepo: siteRepo,
	}

	t.Run("Success Delete Site", func(t *testing.T) {
		siteRepo.On("GetDetail", mock.Anything, mock.Anything).Once().Return(&entity.Site{}, nil)
		siteRepo.On("Delete", mock.Anything, mock.Anything).Once().Return(nil)
		err := siteService.Delete(context.Background(), 1)
		assert.Nil(t, err)
	})

	t.Run("Failed GetDetail Site", func(t *testing.T) {
		siteRepo.On("GetDetail", mock.Anything, mock.Anything).Once().Return(nil, errors.New("error on get detail site"))
		err := siteService.Delete(context.Background(), 1)
		assert.NotNil(t, err)
	})

	t.Run("Failed Delete Site", func(t *testing.T) {
		siteRepo.On("GetDetail", mock.Anything, mock.Anything).Once().Return(&entity.Site{}, nil)
		siteRepo.On("Delete", mock.Anything, mock.Anything).Once().Return(errors.New("error on update site"))
		err := siteService.Delete(context.Background(), 1)
		assert.NotNil(t, err)
	})
}

func TestAssignStreamToSite(t *testing.T) {
	logutil.Init("info")
	siteRepo := &mocks.Site{}
	siteService := site.ServiceImpl{
		SiteRepo: siteRepo,
	}

	t.Run("Success Delete Site", func(t *testing.T) {
		siteRepo.On("GetDetail", mock.Anything, mock.Anything).Once().Return(&entity.Site{}, nil)
		siteRepo.On("GetMapStreamSiteByStreamID", mock.Anything, mock.Anything).Once().Return(&entity.MapSiteStream{ID: 1}, nil)
		siteRepo.On("AddStreamToSite", mock.Anything, mock.Anything).Once().Return(nil)
		err := siteService.AssignStreamToSite(context.Background(), &presenter.AssignStreamRequest{
			SiteID:   1,
			StreamID: "stream-id",
		})
		assert.Nil(t, err)
	})

	t.Run("Failed GetDetail Site", func(t *testing.T) {
		siteRepo.On("GetDetail", mock.Anything, mock.Anything).Once().Return(nil, errors.New("error on get detail site"))
		err := siteService.AssignStreamToSite(context.Background(), &presenter.AssignStreamRequest{
			SiteID:   1,
			StreamID: "stream-id",
		})
		assert.NotNil(t, err)
	})

	t.Run("Failed get map stream site by stream id", func(t *testing.T) {
		siteRepo.On("GetDetail", mock.Anything, mock.Anything).Once().Return(&entity.Site{}, nil)
		siteRepo.On("GetMapStreamSiteByStreamID", mock.Anything, mock.Anything).Once().Return(nil, errors.New("error on get detail map stream by stream id"))

		err := siteService.AssignStreamToSite(context.Background(), &presenter.AssignStreamRequest{
			SiteID:   1,
			StreamID: "stream-id",
		})
		assert.NotNil(t, err)
	})

	t.Run("Failed save to map stream site", func(t *testing.T) {
		siteRepo.On("GetDetail", mock.Anything, mock.Anything).Once().Return(&entity.Site{}, nil)
		siteRepo.On("GetMapStreamSiteByStreamID", mock.Anything, mock.Anything).Once().Return(nil, nil)
		siteRepo.On("AddStreamToSite", mock.Anything, mock.Anything).Once().Return(errors.New("error on update site"))
		err := siteService.AssignStreamToSite(context.Background(), &presenter.AssignStreamRequest{
			SiteID:   1,
			StreamID: "stream-id",
		})
		assert.NotNil(t, err)
	})
}
