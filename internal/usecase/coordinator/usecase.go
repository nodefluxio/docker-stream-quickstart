package coordinator

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"math/rand"
	"time"

	"github.com/robfig/cron"
	logutil "gitlab.com/nodefluxio/goutils/pkg/log"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/pkg/util"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/presenter"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/repository"
)

// ServiceImpl struct to represent quota transaction service
type ServiceImpl struct {
	EventEnrollmentRepo          repository.EventEnrollment
	EventEnrollmentFaceImageRepo repository.EventEnrollmentFaceImage
}

// Create is function for add enrollment
func (s *ServiceImpl) Create(ctx context.Context, postData *presenter.CoordinatorRequest) error {
	var faceImages []*entity.EventEnrollmentImage
	if len(postData.Images) != 0 {
		// decode b64 to blob
		for _, image := range postData.Images {
			unbased, err := base64.StdEncoding.DecodeString(image)
			if err != nil {
				logutil.LogObj.SetErrorLog(map[string]interface{}{
					"err": err,
				},
					"error on when unbased base64 string ")
				return err
			}

			newImage, err := jpeg.Decode(bytes.NewReader(unbased))
			if err != nil {
				logutil.LogObj.SetErrorLog(map[string]interface{}{
					"err": err,
				},
					"error on when decode image")
				return err
			}

			buf := new(bytes.Buffer)
			err = jpeg.Encode(buf, newImage, nil)
			if err != nil {
				logutil.LogObj.SetErrorLog(map[string]interface{}{
					"err": err,
				},
					"error on when encode image")
				return err
			}

			faceImages = append(faceImages, &entity.EventEnrollmentImage{
				Image: buf.Bytes(),
			})
		}
	}

	// generate event id
	rand.Seed(time.Now().UnixNano())
	eventId := fmt.Sprintf("%d%d", time.Now().Unix(), rand.Intn(999))
	eventEnrolment := entity.EventEnrollment{
		EventID:     eventId,
		Agent:       postData.Agent,
		EventAction: postData.EventAction,
		Payload:     postData.Payload,
	}
	_, err := s.EventEnrollmentRepo.Create(ctx, &eventEnrolment, faceImages)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"error on when save data to database")
		return err
	}
	return nil
}

func (s *ServiceImpl) Get(ctx context.Context, paging *util.Pagination) ([]*presenter.CoordinatorResponse, error) {
	var result []*presenter.CoordinatorResponse = make([]*presenter.CoordinatorResponse, 0)
	data, err := s.EventEnrollmentRepo.Get(ctx, paging)

	for _, v := range data {
		var dataPayload map[string]interface{}
		err = json.Unmarshal(v.Payload, &dataPayload)
		if err != nil {
			return result, err
		}
		faces, err := s.EventEnrollmentFaceImageRepo.GetByEventEnrollmendID(ctx, v.ID)
		if err != nil {
			logutil.LogObj.SetErrorLog(map[string]interface{}{
				"err": err,
			},
				"error on when GetByEventEnrollmendID repo")
			return nil, err
		}

		images := []string{}
		if len(faces) != 0 {
			for _, img := range faces {
				nImg := base64.StdEncoding.EncodeToString(img.Image)
				images = append(images, nImg)
			}
		}
		coordResponse := &presenter.CoordinatorResponse{
			EventID:     v.EventID,
			Agent:       v.Agent,
			EventAction: v.EventAction,
			Payload:     dataPayload,
			CreatedAt:   v.CreatedAt,
			Images:      images,
		}
		result = append(result, coordResponse)
	}
	return result, err
}

// CronjobPartition is function for scheduling running partition
func (s *ServiceImpl) CronjobPartition(ctx context.Context) error {
	logutil.LogObj.SetInfoLog(map[string]interface{}{}, "setup cronjob partition")
	c := cron.New()
	cronjobSpec := " 0 0 * * *"
	c.AddFunc(cronjobSpec, func() { s.Partition(ctx) })

	// Start cron with one scheduled job
	logutil.LogObj.SetInfoLog(map[string]interface{}{}, "start cron")
	c.Start()
	return nil
}

// Partition is function for running event partition
func (s *ServiceImpl) Partition(ctx context.Context) error {
	logutil.LogObj.SetInfoLog(map[string]interface{}{}, "start running event partition")
	err := s.EventEnrollmentRepo.Partition(ctx, time.Now())
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{"err": err}, "failed running event partition")
		return err
	}

	err = s.EventEnrollmentRepo.Partition(ctx, time.Now().AddDate(0, 0, 1))
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{"err": err}, "failed running event partition")
		return err
	}

	err = s.EventEnrollmentRepo.Partition(ctx, time.Now().AddDate(0, 0, 2))
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{"err": err}, "failed running event partition")
		return err
	}
	logutil.LogObj.SetInfoLog(map[string]interface{}{}, "finish running event partition")
	return nil
}
