package vehicle

import (
	"context"

	"github.com/jinzhu/copier"
	logutil "gitlab.com/nodefluxio/goutils/pkg/log"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/pkg/util"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/presenter"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/repository"
)

// ServiceImpl struct to represent quota transaction service
type ServiceImpl struct {
	VehicleRepo repository.Vehicle
}

// GetList for get all data enrollment with paging
func (s *ServiceImpl) GetList(ctx context.Context, paging *util.Pagination) (*presenter.VehiclePaging, error) {
	totalData, err := s.VehicleRepo.Count(ctx, paging)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"error on count event history from repository")
		return nil, err
	}
	pgDetail := paging.CreateProperties(totalData)

	paging.Offset = pgDetail.Offset
	list, err := s.VehicleRepo.GetList(ctx, paging)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"error on get event history from repository")
		return nil, err
	}
	output := make([]*presenter.VehicleEnrollmentResponse, 0)
	copier.Copy(&output, &list)
	var result presenter.VehiclePaging
	result.Limit = paging.Limit
	result.TotalPage = pgDetail.TotalPage
	result.TotalData = totalData
	result.CurrentPage = pgDetail.CurrentPage
	result.Vehicles = output

	return &result, nil
}

// Create is function for add enrollment
func (s *ServiceImpl) Create(ctx context.Context, postData *presenter.VehicleEnrollmentRequest) (*presenter.VehicleEnrollmentResponse, error) {
	data := entity.Vehicle{}
	copier.Copy(&data, &postData)
	result, err := s.VehicleRepo.Create(ctx, &data)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"error on when insert to vehicle repo")
		return nil, err
	}
	response := presenter.VehicleEnrollmentResponse{}
	copier.Copy(&response, &data)
	response.ID = result.ID
	return &response, nil
}

// GetDetail is for get detail vehicle by id
func (s *ServiceImpl) GetDetail(ctx context.Context, ID uint64) (*presenter.VehicleEnrollmentResponse, error) {
	detail, err := s.VehicleRepo.GetDetail(ctx, ID)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"id":  ID,
			"err": err,
		},
			"error on when GetDetail repo")
		return nil, err
	}
	response := presenter.VehicleEnrollmentResponse{}
	copier.Copy(&response, &detail)
	return &response, nil
}

// Delete Vehicle
func (s *ServiceImpl) Delete(ctx context.Context, ID uint64) error {
	objDetail, errGetDetail := s.VehicleRepo.GetDetail(ctx, ID)
	if errGetDetail != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": errGetDetail,
		},
			"error on get detail base one ID")
		return errGetDetail
	}
	logutil.LogObj.SetDebugLog(map[string]interface{}{
		"obj": objDetail,
	},
		"log obj detail ")
	errDelete := s.VehicleRepo.Delete(ctx, objDetail.ID)
	if errDelete != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": errDelete,
		},
			"error on delete from database")
		return errDelete
	}
	return nil
}

// Delete Face ID
func (s *ServiceImpl) DeleteAll(ctx context.Context) error {
	errGetDetail := s.VehicleRepo.DeleteAll(ctx)
	if errGetDetail != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": errGetDetail,
		},
			"error on delete all data")
		return errGetDetail
	}
	return nil
}

// Update is function for add enrollment
func (s *ServiceImpl) Update(ctx context.Context, ID uint64, postData *presenter.VehicleEnrollmentRequest) error {
	detailEnroll, err := s.VehicleRepo.GetDetail(ctx, ID)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"error on get detail vehicle")
		return err
	}
	copier.Copy(&detailEnroll, &postData)
	err = s.VehicleRepo.Update(ctx, detailEnroll)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"error on update vehicle")
		return err
	}
	return nil
}
