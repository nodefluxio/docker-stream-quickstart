package site

import (
	"context"
	"strconv"
	"strings"

	logutil "gitlab.com/nodefluxio/goutils/pkg/log"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/pkg/util"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/presenter"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/repository"
)

// ServiceImpl struct to represent site service
type ServiceImpl struct {
	SiteRepo repository.Site
}

// GetList for get all data site
func (s *ServiceImpl) GetList(ctx context.Context, paging *util.Pagination, userInfo *presenter.AuthInfoResponse) ([]*entity.Site, error) {
	var filterSiteID string
	switch userInfo.Role {
	case string(entity.UserRoleOperator):
		if len(userInfo.SiteID) == 0 {
			return make([]*entity.Site, 0), nil
		}

		var IDs []string
		for _, i := range userInfo.SiteID {
			IDs = append(IDs, strconv.FormatInt(i, 10))
		}

		filterSiteID = strings.Join(IDs, ", ")
	}

	if filterSiteID != "" {
		paging.Filter["site_id"] = filterSiteID
	}
	dataSites, err := s.SiteRepo.GetList(ctx, paging)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"paging": paging,
			"err":    err,
		},
			"error when get list site")
		return nil, err
	}

	return dataSites, nil
}

// Create is function for create site
func (s *ServiceImpl) Create(ctx context.Context, data *presenter.SiteRequest) (*entity.Site, error) {
	dataSite, err := s.SiteRepo.Create(ctx, &entity.Site{Name: data.Name})
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"data": data,
			"err":  err,
		},
			"error when create site")
		return nil, err
	}
	return dataSite, nil
}

// Update is function for update one site data
func (s *ServiceImpl) Update(ctx context.Context, data *presenter.SiteRequest) error {
	_, err := s.SiteRepo.GetDetail(ctx, data.ID)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"data": data,
			"err":  err,
		},
			"error when get detail site for check avaibility")
		return err
	}
	err = s.SiteRepo.Update(ctx, &entity.Site{Name: data.Name, ID: data.ID})
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"data": data,
			"err":  err,
		},
			"error when update site")
		return err
	}
	return nil
}

// Delete is function for delete one site data
func (s *ServiceImpl) Delete(ctx context.Context, ID uint64) error {
	_, err := s.SiteRepo.GetDetail(ctx, ID)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"id":  ID,
			"err": err,
		},
			"error when get detail site for check avaibility")
		return err
	}
	err = s.SiteRepo.Delete(ctx, ID)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"id":  ID,
			"err": err,
		},
			"error when delete site")
		return err
	}
	return nil
}

func (s *ServiceImpl) AssignStreamToSite(ctx context.Context, data *presenter.AssignStreamRequest) error {
	_, err := s.SiteRepo.GetDetail(ctx, data.SiteID)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"data": data,
			"err":  err,
		},
			"error when get detail site for check avaibility")
		return err
	}

	checkMapSite, err := s.SiteRepo.GetMapStreamSiteByStreamID(ctx, data.StreamID)
	if err != nil && err.Error() != "record not found" {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"stream_id": data.StreamID,
			"site_id":   data.SiteID,
			"err":       err,
		},
			"error when get check maps site and stream id")
		return err
	}

	dataAssignSite := entity.MapSiteStream{
		SiteID:   data.SiteID,
		StreamID: data.StreamID,
	}

	if checkMapSite != nil {
		dataAssignSite.ID = checkMapSite.ID
	}

	err = s.SiteRepo.AddStreamToSite(ctx, &dataAssignSite)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"data": data,
			"err":  err,
		},
			"error when assign stream id to site")
		return err
	}
	return nil
}
