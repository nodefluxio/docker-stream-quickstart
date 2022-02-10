package agent

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"image/jpeg"
	"time"

	"github.com/robfig/cron"
	logutil "gitlab.com/nodefluxio/goutils/pkg/log"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/presenter"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/repository"
)

// ServiceImpl struct to represent quota transaction service
type ServiceImpl struct {
	LatestTimestampRepo   repository.LatestTimestamp
	CoordinatorRepo       repository.Coordinator
	EnrollmentVanillaRepo repository.EnrollmentVanilla
	SyncPeriod            string
	AgentName             string
	TotalEventSync        string
}

var status string = "starting"
var LastSyncTimestamp time.Time

func decodeImage(image string) (*bytes.Buffer, error) {
	unbased, err := base64.StdEncoding.DecodeString(image)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"error on when unbased base64 string ")
		return nil, err
	}

	newImage, err := jpeg.Decode(bytes.NewReader(unbased))
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"error on when decode image")
		return nil, err
	}

	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, newImage, nil)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"error on when encode image")
		return nil, err
	}
	return buf, nil
}

// Create is function for add enrollment
func (s *ServiceImpl) Sync(ctx context.Context) error {
	if status == "running" {
		logutil.LogObj.SetInfoLog(map[string]interface{}{"status": status}, "skipping sync, because 1 job is already running")
		return nil
	}

	status = "running"
	LastSyncTimestamp = time.Now()
	var latestTimestamp string
	dataLastTimestamp, err := s.LatestTimestampRepo.Get(ctx)
	if err != nil && err.Error() != "record not found" {
		status = "done"
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"error when get latest timestamp data")
		return err
	}
	logutil.LogObj.SetInfoLog(map[string]interface{}{"latest_timestamp": dataLastTimestamp.Timestamp}, "Starting Sync Enrollment")
	latestTimestamp = dataLastTimestamp.Timestamp
	data, err := s.CoordinatorRepo.GetEnrollmentEvent(ctx, s.TotalEventSync, latestTimestamp)
	if err != nil {
		status = "done"
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"error when get data from CES Coordinator")
		return err
	}
	if len(data.Results) > 0 {
		for _, v := range data.Results {
			if s.AgentName == v.Agent {
				// assume job success because event come from same agent
				s.LatestTimestampRepo.CreateOrUpdate(ctx, &entity.LatestTimestamp{
					Timestamp: v.CreatedAt,
				})
				status = "done"
				continue
			}
			logutil.LogObj.SetInfoLog(map[string]interface{}{
				"action": v.EventAction,
			},
				"start process event enrollment")
			switch v.EventAction {
			case "create":
				decodedImage := make([]*entity.EnrollmentImage, 0)
				if len(v.Images) > 0 {
					for _, image := range v.Images {
						buf, err := decodeImage(image)
						if err != nil {
							status = "done"
							logutil.LogObj.SetErrorLog(map[string]interface{}{
								"err": err,
							},
								"error when decode image")
							return err
						}
						decodedImage = append(decodedImage, &entity.EnrollmentImage{Image: buf})
					}
				}

				dataEnroll := &entity.VanillaEnrollmentPayload{
					FaceID:         v.Payload.FaceID,
					IdentityNumber: v.Payload.IdentityNumber,
					Name:           v.Payload.Name,
					Gender:         v.Payload.Gender,
					BirthPlace:     v.Payload.BirthPlace,
					BirthDate:      v.Payload.BirthDate,
					Status:         v.Payload.Status,
				}
				err := s.EnrollmentVanillaRepo.CreateFaceEnrollment(ctx, dataEnroll, decodedImage)
				if err != nil {
					status = "done"
					logutil.LogObj.SetErrorLog(map[string]interface{}{
						"err": err,
					},
						"error when hit api vanilla dashboard CREATE face enrollment")
					return err
				}
			case "update":
				dataVanilla, err := s.EnrollmentVanillaRepo.GetByFaceID(ctx, v.Payload.FaceID)
				if err != nil {
					status = "done"
					logutil.LogObj.SetErrorLog(map[string]interface{}{
						"err":     err,
						"face_id": v.Payload.FaceID,
					},
						"error when hit api vanilla dashboard GET face enrollment by face id")
					return err
				}
				if len(dataVanilla.Results.Enrollments) > 0 {
					decodedImage := make([]*entity.EnrollmentImage, 0)
					if len(v.Images) > 0 {
						for _, image := range v.Images {
							buf, err := decodeImage(image)
							if err != nil {
								status = "done"
								logutil.LogObj.SetErrorLog(map[string]interface{}{
									"err": err,
								},
									"error when decode image")
								return err
							}
							decodedImage = append(decodedImage, &entity.EnrollmentImage{Image: buf})
						}
					}

					dataEnroll := &entity.VanillaEnrollmentPayload{
						FaceID:         v.Payload.FaceID,
						IdentityNumber: v.Payload.IdentityNumber,
						Name:           v.Payload.Name,
						Gender:         v.Payload.Gender,
						BirthPlace:     v.Payload.BirthPlace,
						BirthDate:      v.Payload.BirthDate,
						Status:         v.Payload.Status,
					}
					err := s.EnrollmentVanillaRepo.UpdateFaceEnrollment(ctx, dataVanilla.Results.Enrollments[0].ID, dataEnroll, decodedImage)
					if err != nil {
						status = "done"
						logutil.LogObj.SetErrorLog(map[string]interface{}{
							"err": err,
						},
							"error when hit api vanilla dashboard UPDATE face enrollment")
						return err
					}
				}
			case "delete":
				dataVanilla, err := s.EnrollmentVanillaRepo.GetByFaceID(ctx, v.Payload.FaceID)
				if err != nil {
					status = "done"
					logutil.LogObj.SetErrorLog(map[string]interface{}{
						"err":     err,
						"face_id": v.Payload.FaceID,
					},
						"error when hit api vanilla dashboard GET face enrollment by face id")
					return err
				}
				if len(dataVanilla.Results.Enrollments) > 0 {
					err := s.EnrollmentVanillaRepo.DeleteFaceEnrollment(ctx, dataVanilla.Results.Enrollments[0].ID)
					if err != nil {
						status = "done"
						logutil.LogObj.SetErrorLog(map[string]interface{}{
							"err": err,
						},
							"error when hit api vanilla dashboard DELETE face enrollment")
						return err
					}
				}
			}
			s.LatestTimestampRepo.CreateOrUpdate(ctx, &entity.LatestTimestamp{
				Timestamp: v.CreatedAt,
			})
			logutil.LogObj.SetInfoLog(map[string]interface{}{
				"event_id":        v.EventID,
				"action":          v.EventAction,
				"face_id":         v.Payload.FaceID,
				"event_timestamp": v.CreatedAt,
			},
				"success process event enrollment")
		}
	}
	logutil.LogObj.SetInfoLog(map[string]interface{}{"total_event_executed": len(data.Results)}, "Finished Sync Enrollment")
	status = "done"
	return nil
}

// CronjobSyncEnrollment is function for scheduling running sync enrollemnt
func (s *ServiceImpl) CronjobSyncEnrollment(ctx context.Context) error {
	logutil.LogObj.SetInfoLog(map[string]interface{}{"sync_period": s.SyncPeriod}, "setup cronjob Sync")
	c := cron.New()
	c.AddFunc(s.SyncPeriod, func() { s.Sync(ctx) })

	// Start cron with one scheduled job
	logutil.LogObj.SetInfoLog(map[string]interface{}{}, "start cron")
	c.Start()
	return nil
}

func (s *ServiceImpl) Create(ctx context.Context, postData *presenter.CoordinatorRequest) error {
	var payload entity.PayloadEnrollmentEventCoordinator
	err := json.Unmarshal(postData.Payload, &payload)
	if err != nil {
		return err
	}
	data := entity.CreateEnrollmentEventCoordinator{
		Agent:       s.AgentName,
		EventAction: postData.EventAction,
		Images:      postData.Images,
		Payload:     payload,
	}
	err = s.CoordinatorRepo.CreateEnrollmentEvent(ctx, &data)
	if err != nil {
		return err
	}
	return nil
}

func (s *ServiceImpl) PingCoordinator(ctx context.Context) error {
	err := s.CoordinatorRepo.Ping(ctx)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"failed ping CES coordinator")
		return errors.New("failed ping CES coordinator")
	}
	return nil
}

func (s *ServiceImpl) GetStatus(ctx context.Context) (*presenter.AgentStatus, error) {
	err := s.CoordinatorRepo.Ping(ctx)
	newStatus := status
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"failed ping CES coordinator")
		newStatus = "error"
	}
	return &presenter.AgentStatus{
		Status:            newStatus,
		LastSyncTimestamp: LastSyncTimestamp,
	}, nil
}
