package user

import (
	"context"
	"errors"
	"fmt"

	logutil "gitlab.com/nodefluxio/goutils/pkg/log"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/pkg/util"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/presenter"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/repository"
)

type ServiceImpl struct {
	UserRepo repository.User
	SiteRepo repository.Site
}

const minPasswordLength = 8

// Create is function for create site
func (s *ServiceImpl) Create(ctx context.Context, postData *presenter.UserRequest) (*entity.User, error) {
	if postData.Password != postData.RePassword {
		return nil, errors.New("Your password and confirmation password do not match")
	}

	if len(postData.Password) < 8 {
		return nil, fmt.Errorf("Your password to short, minimum length is %d", minPasswordLength)
	}

	_, err := s.UserRepo.GetByUsername(ctx, postData.Username)
	if err == nil {
		return nil, errors.New("Username is already used, please choose another username")
	}

	_, err = s.UserRepo.GetByEmail(ctx, postData.Email)
	if err == nil {
		return nil, errors.New("Email is already used, please choose another email")
	}

	if postData.Role != string(entity.UserRoleSuperAdmin) && postData.Role != string(entity.UserRoleOperator) {
		return nil, errors.New("Role doesn't exists")
	}

	_, err = s.SiteRepo.GetSiteByIDs(ctx, postData.SiteID)
	if err != nil {
		return nil, errors.New("Site doesn't exists")
	}

	// hash password
	pass, err := util.HashPassword(postData.Password)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"error when hash password")
		return nil, err
	}

	// cast data user to save
	dataUser := &entity.User{
		Email:    postData.Email,
		Username: postData.Username,
		Fullname: postData.Fullname,
		Role:     postData.Role,
		SiteID:   postData.SiteID,
		Password: pass,
	}
	savedUser, err := s.UserRepo.Create(ctx, dataUser)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"data_user": dataUser,
			"err":       err,
		},
			"error when create user")
		return nil, err
	}
	savedUser.Password = ""
	return savedUser, nil
}

func (s *ServiceImpl) Update(ctx context.Context, postData *presenter.UserRequest) error {
	_, err := s.SiteRepo.GetSiteByIDs(ctx, postData.SiteID)
	if err != nil {
		return errors.New("Site doesn't exists")
	}

	if postData.Role != string(entity.UserRoleSuperAdmin) && postData.Role != string(entity.UserRoleOperator) {
		return errors.New("Role doesn't exists")
	}

	dataUser := &entity.User{
		Fullname: postData.Fullname,
		Role:     postData.Role,
		SiteID:   postData.SiteID,
		ID:       postData.ID,
	}
	err = s.UserRepo.UpdateBasicData(ctx, dataUser)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"data_user": dataUser,
			"err":       err,
		},
			"error when update user")
		return err
	}
	return nil
}

func (s *ServiceImpl) ChangePassword(ctx context.Context, postData *presenter.UserChangePassRequest) error {
	if postData.Password != postData.RePassword {
		return errors.New("Your password and confirmation password do not match")
	}

	if len(postData.Password) < 8 {
		return fmt.Errorf("Your password to short, minimum length is %d", minPasswordLength)
	}

	// check user avaibility
	_, err := s.UserRepo.GetDetail(ctx, postData.ID)
	if err != nil {
		return errors.New("User not found")
	}
	// start hash password
	// note: please refactor this, cannot be tested
	pass, err := util.HashPassword(postData.Password)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"error when hash password")
		return err
	}
	err = s.UserRepo.UpdatePassword(ctx, pass, postData.ID)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"id":  postData.ID,
			"err": err,
		},
			"error when update user password")
		return err
	}
	return nil
}

func (s *ServiceImpl) Delete(ctx context.Context, ID uint64) error {
	err := s.UserRepo.Delete(ctx, ID)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"id":  ID,
			"err": err,
		},
			"error when delete user")
		return err
	}
	return nil
}

func (s *ServiceImpl) GetDetail(ctx context.Context, ID uint64) (*entity.User, error) {
	dataUser, err := s.UserRepo.GetDetail(ctx, ID)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"id":  ID,
			"err": err,
		},
			"error when get detail user")
		return nil, err
	}
	return dataUser, nil
}

// GetList for get all data user with paging
func (s *ServiceImpl) GetList(ctx context.Context, paging *util.Pagination) (*presenter.UserPaging, error) {
	totalData, err := s.UserRepo.Count(ctx, paging)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"error on count user from repository")
		return nil, err
	}
	pgDetail := paging.CreateProperties(totalData)

	paging.Offset = pgDetail.Offset
	dataUser, err := s.UserRepo.GetList(ctx, paging)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"error on get user from repository")
		return nil, err
	}

	if len(dataUser) == 0 {
		dataUser = make([]*entity.User, 0)
	}

	result := presenter.UserPaging{}
	result.Limit = paging.Limit
	result.TotalPage = pgDetail.TotalPage
	result.TotalData = totalData
	result.CurrentPage = pgDetail.CurrentPage
	result.Users = dataUser

	return &result, nil
}
