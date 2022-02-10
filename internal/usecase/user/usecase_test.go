package user_test

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
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/usecase/user"
)

func TestCreate(t *testing.T) {
	logutil.Init("info")
	userRepo := &mocks.User{}
	siteRepo := &mocks.Site{}
	userService := user.ServiceImpl{
		UserRepo: userRepo,
		SiteRepo: siteRepo,
	}

	t.Run("Failed create user, error validate password confirmation", func(t *testing.T) {
		postData := &presenter.UserRequest{
			Password:   "password",
			RePassword: "repassword",
		}
		_, err := userService.Create(context.Background(), postData)
		assert.NotNil(t, err)
	})

	t.Run("Failed create user, error validate password length", func(t *testing.T) {
		postData := &presenter.UserRequest{
			Password:   "1",
			RePassword: "1",
		}
		_, err := userService.Create(context.Background(), postData)
		assert.NotNil(t, err)
	})

	t.Run("Failed create user, error username already use", func(t *testing.T) {
		postData := &presenter.UserRequest{
			Password:   "password",
			RePassword: "password",
			Username:   "username",
		}
		userRepo.On("GetByUsername", mock.Anything, mock.Anything).Once().Return(&entity.User{}, nil)

		_, err := userService.Create(context.Background(), postData)
		assert.NotNil(t, err)
	})

	t.Run("Failed create user, error email already use", func(t *testing.T) {
		postData := &presenter.UserRequest{
			Password:   "password",
			RePassword: "password",
			Username:   "username",
			Email:      "email@nodeflux.io",
		}
		userRepo.On("GetByUsername", mock.Anything, mock.Anything).Once().Return(nil, errors.New("record not found"))
		userRepo.On("GetByEmail", mock.Anything, mock.Anything).Once().Return(&entity.User{}, nil)

		_, err := userService.Create(context.Background(), postData)
		assert.NotNil(t, err)
	})

	t.Run("Failed create user, error role not avail", func(t *testing.T) {
		postData := &presenter.UserRequest{
			Password:   "password",
			RePassword: "password",
			Username:   "username",
			Email:      "email@nodeflux.io",
			Role:       "pengawas",
		}
		userRepo.On("GetByUsername", mock.Anything, mock.Anything).Once().Return(nil, errors.New("record not found"))
		userRepo.On("GetByEmail", mock.Anything, mock.Anything).Once().Return(nil, errors.New("record not found"))

		_, err := userService.Create(context.Background(), postData)
		assert.NotNil(t, err)
	})

	t.Run("Failed create user, error site not avail", func(t *testing.T) {
		postData := &presenter.UserRequest{
			Password:   "password",
			RePassword: "password",
			Username:   "username",
			Email:      "email@nodeflux.io",
			Role:       string(entity.UserRoleSuperAdmin),
			SiteID:     []int64{9, 8},
		}
		userRepo.On("GetByUsername", mock.Anything, mock.Anything).Once().Return(nil, errors.New("record not found"))
		userRepo.On("GetByEmail", mock.Anything, mock.Anything).Once().Return(nil, errors.New("record not found"))
		siteRepo.On("GetSiteByIDs", mock.Anything, mock.Anything).Once().Return(nil, errors.New("Site doesn't existsd"))

		_, err := userService.Create(context.Background(), postData)
		assert.NotNil(t, err)
	})

	t.Run("Failed create user, error create", func(t *testing.T) {
		postData := &presenter.UserRequest{
			Password:   "password",
			RePassword: "password",
			Username:   "username",
			Email:      "email@nodeflux.io",
			Role:       string(entity.UserRoleSuperAdmin),
			SiteID:     []int64{9, 8},
		}
		fakeListSite := []*entity.Site{}
		userRepo.On("GetByUsername", mock.Anything, mock.Anything).Once().Return(nil, errors.New("record not found"))
		userRepo.On("GetByEmail", mock.Anything, mock.Anything).Once().Return(nil, errors.New("record not found"))
		siteRepo.On("GetSiteByIDs", mock.Anything, mock.Anything).Once().Return(fakeListSite, nil)
		userRepo.On("Create", mock.Anything, mock.Anything).Once().Return(nil, errors.New("error when create user"))

		_, err := userService.Create(context.Background(), postData)
		assert.NotNil(t, err)
	})

	t.Run("Success create user", func(t *testing.T) {
		postData := &presenter.UserRequest{
			Password:   "password",
			RePassword: "password",
			Username:   "username",
			Email:      "email@nodeflux.io",
			Role:       string(entity.UserRoleSuperAdmin),
			SiteID:     []int64{9, 8},
		}
		fakeListSite := []*entity.Site{}
		userRepo.On("GetByUsername", mock.Anything, mock.Anything).Once().Return(nil, errors.New("record not found"))
		userRepo.On("GetByEmail", mock.Anything, mock.Anything).Once().Return(nil, errors.New("record not found"))
		siteRepo.On("GetSiteByIDs", mock.Anything, mock.Anything).Once().Return(fakeListSite, nil)
		userRepo.On("Create", mock.Anything, mock.Anything).Once().Return(&entity.User{}, nil)

		_, err := userService.Create(context.Background(), postData)
		assert.Nil(t, err)
	})
}
func TestUpdate(t *testing.T) {
	logutil.Init("info")
	userRepo := &mocks.User{}
	siteRepo := &mocks.Site{}
	userService := user.ServiceImpl{
		UserRepo: userRepo,
		SiteRepo: siteRepo,
	}

	t.Run("Failed update, error site not found", func(t *testing.T) {
		postData := &presenter.UserRequest{
			SiteID: []int64{1},
		}
		siteRepo.On("GetSiteByIDs", mock.Anything, mock.Anything).Once().Return(nil, errors.New("Site doesn't exists"))
		err := userService.Update(context.Background(), postData)
		assert.NotNil(t, err)
	})

	t.Run("Failed update, error site not found", func(t *testing.T) {
		postData := &presenter.UserRequest{
			SiteID: []int64{1},
			Role:   "pengawas",
		}
		fakeListSite := []*entity.Site{}
		siteRepo.On("GetSiteByIDs", mock.Anything, mock.Anything).Once().Return(fakeListSite, nil)
		err := userService.Update(context.Background(), postData)
		assert.NotNil(t, err)
	})

	t.Run("Failed update, error update basic data failed", func(t *testing.T) {
		postData := &presenter.UserRequest{
			SiteID: []int64{1},
			Role:   string(entity.UserRoleSuperAdmin),
		}
		fakeListSite := []*entity.Site{}
		siteRepo.On("GetSiteByIDs", mock.Anything, mock.Anything).Once().Return(fakeListSite, nil)
		userRepo.On("UpdateBasicData", mock.Anything, mock.Anything).Once().Return(errors.New("error when update user"))
		err := userService.Update(context.Background(), postData)
		assert.NotNil(t, err)
	})
	t.Run("Success", func(t *testing.T) {
		postData := &presenter.UserRequest{
			SiteID: []int64{1},
			Role:   string(entity.UserRoleSuperAdmin),
		}
		fakeListSite := []*entity.Site{}
		siteRepo.On("GetSiteByIDs", mock.Anything, mock.Anything).Once().Return(fakeListSite, nil)
		userRepo.On("UpdateBasicData", mock.Anything, mock.Anything).Once().Return(nil)
		err := userService.Update(context.Background(), postData)
		assert.Nil(t, err)
	})
}

func TestChangePassword(t *testing.T) {
	logutil.Init("info")
	userRepo := &mocks.User{}
	siteRepo := &mocks.Site{}
	userService := user.ServiceImpl{
		UserRepo: userRepo,
		SiteRepo: siteRepo,
	}

	t.Run("Failed change password, error validate password confirmation", func(t *testing.T) {
		postData := &presenter.UserChangePassRequest{
			Password:   "password",
			RePassword: "repassword",
		}
		err := userService.ChangePassword(context.Background(), postData)
		assert.NotNil(t, err)
	})

	t.Run("Failed change password, error validate password to short", func(t *testing.T) {
		postData := &presenter.UserChangePassRequest{
			Password:   "1",
			RePassword: "1",
		}
		err := userService.ChangePassword(context.Background(), postData)
		assert.NotNil(t, err)
	})

	t.Run("Failed change password, error get detail user", func(t *testing.T) {
		postData := &presenter.UserChangePassRequest{
			ID:         1,
			Password:   "NewPassword",
			RePassword: "NewPassword",
		}
		userRepo.On("GetDetail", mock.Anything, mock.Anything).Once().Return(nil, errors.New("record not found"))
		err := userService.ChangePassword(context.Background(), postData)
		assert.NotNil(t, err)
	})

	t.Run("Failed change password, error update password", func(t *testing.T) {
		postData := &presenter.UserChangePassRequest{
			ID:         1,
			Password:   "NewPassword",
			RePassword: "NewPassword",
		}
		userRepo.On("GetDetail", mock.Anything, mock.Anything).Once().Return(&entity.User{}, nil)
		userRepo.On("UpdatePassword", mock.Anything, mock.Anything, mock.Anything).Once().Return(errors.New("error when update user password"))
		err := userService.ChangePassword(context.Background(), postData)
		assert.NotNil(t, err)
	})

	t.Run("Success change password", func(t *testing.T) {
		postData := &presenter.UserChangePassRequest{
			ID:         1,
			Password:   "NewPassword",
			RePassword: "NewPassword",
		}
		userRepo.On("GetDetail", mock.Anything, mock.Anything).Once().Return(&entity.User{}, nil)
		userRepo.On("UpdatePassword", mock.Anything, mock.Anything, mock.Anything).Once().Return(nil)
		err := userService.ChangePassword(context.Background(), postData)
		assert.Nil(t, err)
	})
}

func TestDelete(t *testing.T) {
	logutil.Init("info")
	userRepo := &mocks.User{}
	siteRepo := &mocks.Site{}
	userService := user.ServiceImpl{
		UserRepo: userRepo,
		SiteRepo: siteRepo,
	}

	t.Run("Success Delete User", func(t *testing.T) {
		userRepo.On("Delete", mock.Anything, mock.Anything).Once().Return(nil)
		err := userService.Delete(context.Background(), 1)
		assert.Nil(t, err)
	})

	t.Run("Failed Delete User", func(t *testing.T) {
		userRepo.On("Delete", mock.Anything, mock.Anything).Once().Return(errors.New("error on Delete user"))
		err := userService.Delete(context.Background(), 1)
		assert.NotNil(t, err)
	})
}

func TestGetDetail(t *testing.T) {
	logutil.Init("info")
	userRepo := &mocks.User{}
	siteRepo := &mocks.Site{}
	userService := user.ServiceImpl{
		UserRepo: userRepo,
		SiteRepo: siteRepo,
	}

	t.Run("Success GetDetail User", func(t *testing.T) {
		userRepo.On("GetDetail", mock.Anything, mock.Anything).Once().Return(&entity.User{}, nil)
		_, err := userService.GetDetail(context.Background(), 1)
		assert.Nil(t, err)
	})

	t.Run("Failed GetDetail User", func(t *testing.T) {
		userRepo.On("GetDetail", mock.Anything, mock.Anything).Once().Return(nil, errors.New("error on Delete user"))
		_, err := userService.GetDetail(context.Background(), 1)
		assert.NotNil(t, err)
	})
}

func TestGetList(t *testing.T) {
	logutil.Init("info")
	userRepo := &mocks.User{}
	siteRepo := &mocks.Site{}
	userService := user.ServiceImpl{
		UserRepo: userRepo,
		SiteRepo: siteRepo,
	}

	t.Run("Failed count user", func(t *testing.T) {
		paging := util.Pagination{}
		userRepo.On("Count", mock.Anything, mock.Anything).Once().Return(0, errors.New("error on count user from repository"))
		_, err := userService.GetList(context.Background(), &paging)
		assert.NotNil(t, err)
	})

	t.Run("Failed count user", func(t *testing.T) {
		paging := util.Pagination{}
		userRepo.On("Count", mock.Anything, mock.Anything).Once().Return(10, nil)
		userRepo.On("GetList", mock.Anything, mock.Anything).Once().Return(nil, errors.New("error on get user from repository"))
		_, err := userService.GetList(context.Background(), &paging)
		assert.NotNil(t, err)
	})
	t.Run("Failed count user", func(t *testing.T) {
		paging := util.Pagination{}
		fakeListUser := []*entity.User{}
		userRepo.On("Count", mock.Anything, mock.Anything).Once().Return(10, nil)
		userRepo.On("GetList", mock.Anything, mock.Anything).Once().Return(fakeListUser, nil)
		_, err := userService.GetList(context.Background(), &paging)
		assert.Nil(t, err)
	})

}
