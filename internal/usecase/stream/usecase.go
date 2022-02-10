package stream

import (
	"context"
	"errors"
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
	SiteRepo   repository.Site
	StreamRepo repository.Stream
	UserRepo   repository.User
}

// GetList for get all data site
func (s *ServiceImpl) GetList(ctx context.Context, paging *util.Pagination, userInfo *presenter.AuthInfoResponse) (*presenter.StreamResponse, error) {
	response := presenter.StreamResponse{}
	streamList := make([]presenter.StreamDetailWithSite, 0)

	// check filter site_id
	isSiteIDFiltered := false
	if len(paging.Filter["site_id"]) > 0 {
		isSiteIDFiltered = true
	}

	// check site id and role, if user role operator and site is empty will be return empty value
	if userInfo.Role == string(entity.UserRoleOperator) && len(userInfo.SiteID) == 0 {
		response.Streams = streamList
		return &response, nil
	}

	// checking site id
	var desiredSiteID []int64
	if isSiteIDFiltered {
		filterSiteID := strings.Split(paging.Filter["site_id"], ",")
		for _, id := range filterSiteID {
			// checking is filtered site id is available or not at user site id
			siteID, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				logutil.LogObj.SetErrorLog(map[string]interface{}{
					"site_id": id,
					"err":     err,
				},
					"failed parse site id")
				return nil, err
			}
			if userInfo.Role == string(entity.UserRoleOperator) {
				for _, v := range userInfo.SiteID {
					if siteID == v {
						desiredSiteID = append(desiredSiteID, siteID)
					}
				}
			} else if userInfo.Role == string(entity.UserRoleSuperAdmin) {
				desiredSiteID = append(desiredSiteID, siteID)
			}
		}
	} else {
		desiredSiteID = userInfo.SiteID
	}

	// get list site
	siteWithStream, err := s.SiteRepo.GetSiteWithStream(ctx, desiredSiteID)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"user_site_id": userInfo.SiteID,
			"err":          err,
		},
			"error when get list site with stream")
		return nil, err
	}

	// iterate data site
	var streamIDs []string
	listSiteName := make(map[string]string)
	listSiteID := make(map[string]uint64)
	for _, site := range siteWithStream {
		streamIDs = append(streamIDs, site.StreamID)
		listSiteName[site.StreamID] = site.Name
		listSiteID[site.StreamID] = site.ID
	}

	// get list stream
	streams, err := s.StreamRepo.GetList(ctx)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"error when get list stream")
		return nil, err
	}

	// cast data site to stream list
	for _, streamData := range streams.Streams {
		checkStreamID := util.ArrayStringAvailability(streamIDs, streamData.StreamID)
		dataWithSite := presenter.StreamDetailWithSite{
			VisionaireStreamDetail: streamData,
			StreamSiteID:           listSiteID[streamData.StreamID],
			StreamSiteName:         listSiteName[streamData.StreamID],
		}

		switch userInfo.Role {
		case string(entity.UserRoleOperator):
			if checkStreamID {
				streamList = append(streamList, dataWithSite)
			}
		case string(entity.UserRoleSuperAdmin):
			if isSiteIDFiltered {
				if checkStreamID {
					streamList = append(streamList, dataWithSite)
				}
			} else {
				streamList = append(streamList, dataWithSite)
			}
		}
	}
	response.StreamNumber = len(streamList)
	response.Streams = streamList

	return &response, nil
}

// GetDetail is usecase for get detail stream with tempered site and role base checking
func (s *ServiceImpl) GetDetail(ctx context.Context, streamRequest *presenter.StreamRequest, userInfo *presenter.AuthInfoResponse) (*presenter.StreamDetailWithSite, error) {
	var response presenter.StreamDetailWithSite

	// get site detail
	site, err := s.SiteRepo.GetDetailByStreamID(ctx, streamRequest.StreamID)
	if err != nil && err.Error() != "record not found" {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"stream_id": streamRequest.StreamID,
			"err":       err,
		},
			"error when get site data")
		return nil, err
	}

	if site == nil {
		site = &entity.Site{}
	}

	// get stream detail
	streamData, err := s.StreamRepo.GetDetail(ctx, streamRequest.NodeNumber, streamRequest.StreamID)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"stream_request": streamRequest,
			"err":            err,
		},
			"error when get stream detail from visionaire api")
		return nil, err
	}

	// check role
	if site == nil && userInfo.Role == string(entity.UserRoleOperator) {
		err = errors.New("stream not found")
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"site":      site,
			"user_role": userInfo.Role,
			"err":       err,
		},
			"stream not found for this user because either site or stream not available")
		return nil, err
	}

	switch userInfo.Role {
	case string(entity.UserRoleOperator):
		// check site id with assigned site to this role
		streamSiteFound := false
		for _, siteID := range userInfo.SiteID {
			if siteID == int64(site.ID) {
				streamSiteFound = true
				response.VisionaireStreamDetail = *streamData
				response.StreamSiteID = site.ID
				response.StreamSiteName = site.Name
				break
			}
		}
		if !streamSiteFound {
			logutil.LogObj.SetErrorLog(map[string]interface{}{
				"user_role":      userInfo.Role,
				"user_site_id":   userInfo.SiteID,
				"stream_site_id": site.ID,
				"err":            err,
			},
				"stream not found for this user because either site or stream not available")
			return nil, errors.New("stream not found")
		}
	case string(entity.UserRoleSuperAdmin):
		response.VisionaireStreamDetail = *streamData
		response.StreamSiteID = site.ID
		response.StreamSiteName = site.Name
	}

	return &response, nil
}
