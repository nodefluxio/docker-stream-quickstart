package globalsetting

import (
	"context"
	"fmt"
	logutil "gitlab.com/nodefluxio/goutils/pkg/log"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/presenter"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/repository"
)

// ServiceImpl struct to represent quota transaction service
type ServiceImpl struct {
	GlobalSettingRepo repository.GlobalSetting
}


// CreateOrUpdate is function for create or update global setting
func (s *ServiceImpl) CreateOrUpdate(ctx context.Context, postData *presenter.GlobalSettingRequest) (*entity.GlobalSetting, error) {

	result, err := s.GlobalSettingRepo.GetCurrent(ctx)
	if err!=nil{
		if err.Error()=="record not found"{
			data:=&entity.GlobalSetting{
				Similarity:postData.Similarity,
			}
			errCreate:=s.GlobalSettingRepo.Create(ctx,data)
			return data,errCreate
		}
		fmt.Println(result,"resujlt",err)
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"error on when Create Update Global Setting")
		return nil, err
	}
	if result!=nil{
		result.Similarity=postData.Similarity
		errUpdate:=s.GlobalSettingRepo.Update(ctx,result)
		return result,errUpdate
	}
	return &entity.GlobalSetting{}, nil
}

// GetDetail is for get detail enrollment by id
func (s *ServiceImpl) GetDetail(ctx context.Context) (*entity.GlobalSetting, error) {
	result, err := s.GlobalSettingRepo.GetCurrent(ctx)
	
	if err != nil {
		
		if err.Error()=="record not found"{
			return &entity.GlobalSetting{
				Similarity: 0.7,
			}, nil

		}
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"error on when GetDetail repo")
		return nil, err
	}
	
	return result, nil
}
