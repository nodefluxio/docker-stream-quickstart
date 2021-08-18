package enrollment

import (
	"bytes"
	"context"
	b64 "encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"strconv"
	"time"

	"github.com/nfnt/resize"
	logutil "gitlab.com/nodefluxio/goutils/pkg/log"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/pkg/imageprocessing"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/pkg/util"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/presenter"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/repository"
)

// ServiceImpl struct to represent quota transaction service
type ServiceImpl struct {
	EnrolledFaceRepo       repository.EnrolledFace
	FRemisRepo             repository.FRemis
	FaceImageRepo          repository.FaceImage
	UseCES                 string
	AgentRepo              repository.Agent
	PsqlTransactionRepo    repository.PsqlTransaction
	MaxSizeImageEnrollment int64
}

const (
	maxImageWidth     = 2000
	maxImageHeight    = 2000
	size              = 1000
	sizeForAspecRatio = 0 // 0 for maintain aspec ratio
	resolutionDivider = 2
)

// GetList for get all data enrollment with paging
func (s *ServiceImpl) GetList(ctx context.Context, paging *util.Pagination) (*presenter.EnrollmentPaging, error) {
	totalData, err := s.EnrolledFaceRepo.Count(ctx, paging)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"error on count event history from repository")
		return nil, err
	}
	pgDetail := paging.CreateProperties(totalData)

	paging.Offset = pgDetail.Offset
	enrolledFace, err := s.EnrolledFaceRepo.GetList(ctx, paging)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"error on get event history from repository")
		return nil, err
	}
	output := make([]*presenter.EnrollmentResponse, 0)
	for _, enrollment := range enrolledFace {
		faces, err := s.FaceImageRepo.GetDetailByEnrollID(ctx, enrollment.ID)
		if err != nil {
			logutil.LogObj.SetErrorLog(map[string]interface{}{
				"enrollment_id": enrollment.ID,
				"err":           err,
			},
				"error on when GetDetailByEnrollID repo")
			return nil, err
		}
		output = append(output, &presenter.EnrollmentResponse{
			ID:             enrollment.ID,
			Name:           enrollment.Name,
			IdentityNumber: enrollment.IdentityNumber,
			Status:         enrollment.Status,
			FaceID:         enrollment.FaceID,
			CreatedAt:      enrollment.CreatedAt,
			UpdatedAt:      enrollment.UpdatedAt,
			DeletedAt:      enrollment.DeletedAt,
			Faces:          faces,
		})
	}

	var result presenter.EnrollmentPaging
	result.Limit = paging.Limit
	result.TotalPage = pgDetail.TotalPage
	result.TotalData = totalData
	result.CurrentPage = pgDetail.CurrentPage
	result.Enrollments = output

	return &result, nil
}

func (s *ServiceImpl) deleteFremis(ctx context.Context, faceID string) error {
	faceIDs := []string{faceID}
	logutil.LogObj.SetInfoLog(map[string]interface{}{"face_ids": faceID}, "delete fremis enrollment, because error happen at enrollment process")
	return s.FRemisRepo.FaceDeleteEnrollment(ctx, faceIDs)
}

func (s *ServiceImpl) deleteFremisFaceVariation(ctx context.Context, faceID string, variations []string) error {
	logutil.LogObj.SetInfoLog(map[string]interface{}{"face_ids": faceID, "variations": variations}, "delete fremis face variations enrollment, because error happen at enrollment process")
	return s.FRemisRepo.DeleteFaceVariation(ctx, faceID, variations)
}

func (s *ServiceImpl) reEnrollFremis(ctx context.Context, enrolledFaceID uint64, faceID string) error {
	faceImages, err := s.FaceImageRepo.GetDetailByEnrollID(ctx, enrolledFaceID)
	if err != nil {
		return err
	}
	for _, faceImage := range faceImages {
		sEnc := b64.StdEncoding.EncodeToString([]byte(faceImage.Image))
		_, err := s.FRemisRepo.AddFaceVariation(ctx, faceID, sEnc)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *ServiceImpl) prepareImage(ctx context.Context, faceImage []byte) ([]byte, error) {
	imgConvert, err := imageprocessing.ConvertToJPG(faceImage)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"error on when convert image to jpeg")
		return nil, err
	}

	decodedImage, _, err := imageprocessing.Decode(bytes.NewReader(imgConvert))
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"error on when remove exif at image")
		return nil, err
	}
	buf := new(bytes.Buffer)
	jpeg.Encode(buf, decodedImage, nil)

	resImg := buf.Bytes()
	imgReader := bytes.NewReader(resImg)
	imgSize, err := imgReader.Seek(0, os.SEEK_END)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"error on when get size image")
		return nil, err
	}

	img, _, err := image.DecodeConfig(bytes.NewReader(resImg))
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"error on when get image resolution")
		return nil, err
	}
	logutil.LogObj.SetInfoLog(map[string]interface{}{
		"width":  img.Width,
		"height": img.Height,
		"size":   imgSize,
	},
		"image information")

	if img.Width > maxImageWidth || img.Height > maxImageHeight || imgSize > s.MaxSizeImageEnrollment {
		var initWidth uint
		var initHeight uint

		if img.Width > img.Height {
			initHeight = size
			initWidth = sizeForAspecRatio
		} else {
			initHeight = sizeForAspecRatio
			initWidth = size
		}
		for {
			decodedImage, _, err := image.Decode(bytes.NewReader(resImg))
			if err != nil {
				logutil.LogObj.SetErrorLog(map[string]interface{}{
					"err": err,
				},
					"error on when decode image")
				return nil, err
			}

			resizedImg := resize.Resize(initWidth, initHeight, decodedImage, resize.Lanczos3)
			buf := new(bytes.Buffer)
			err = jpeg.Encode(buf, resizedImg, nil)
			if err != nil {
				logutil.LogObj.SetErrorLog(map[string]interface{}{
					"err": err,
				},
					"error on when encode image")
				return nil, err
			}
			resImg = buf.Bytes()

			r := bytes.NewReader(resImg)
			newImgSize, err := r.Seek(0, os.SEEK_END)
			if err != nil {
				logutil.LogObj.SetErrorLog(map[string]interface{}{
					"err": err,
				},
					"error on when get resized image size")
				return nil, err
			}

			// get new resolution
			nImg, _, err := image.DecodeConfig(bytes.NewReader(resImg))
			if err != nil {
				logutil.LogObj.SetErrorLog(map[string]interface{}{
					"err": err,
				},
					"error on when get resized image resolution")
				return nil, err
			}

			logutil.LogObj.SetInfoLog(map[string]interface{}{
				"width":  nImg.Width,
				"height": nImg.Height,
				"size":   newImgSize,
			},
				"new image information")

			if newImgSize < s.MaxSizeImageEnrollment {
				break
			} else {
				if nImg.Width > nImg.Height {
					initHeight = size
					initWidth = uint(nImg.Width) / resolutionDivider
				} else {
					initHeight = uint(nImg.Height) / resolutionDivider
					initWidth = size
				}
			}
		}
	}
	return resImg, nil
}

// Create is function for add enrollment
func (s *ServiceImpl) Create(ctx context.Context, postData *presenter.EnrollmentRequest, isAgent string) (*presenter.EnrollmentResponse, error) {
	var faceID string
	var enrollmentID uint64
	var variation string
	var enrolledData *entity.EnrolledFace
	var images []string
	faces := make([]*entity.FaceImage, 0)

	// check if this system use CES or not, if use we must check ces service avaibility
	if s.UseCES == "true" && isAgent != "true" {
		err := s.AgentRepo.Ping(ctx)
		if err != nil {
			logutil.LogObj.SetErrorLog(map[string]interface{}{
				"err": err,
			},
				"failed ping CES agent CUD opration disabled")
			return nil, errors.New("failed ping CES agent CUD opration disabled")
		}
	}
	// init database transaction
	tx := s.PsqlTransactionRepo.BeginTransaction(ctx)

	// iterate image for enroll to FREMIS and save to database
	for i, v := range postData.Images {
		newImg, err := s.prepareImage(ctx, v.Image)
		if err != nil {
			s.PsqlTransactionRepo.RollbackTransaction(ctx, tx)
			if faceID != "" {
				s.deleteFremis(ctx, faceID)
			}
			logutil.LogObj.SetErrorLog(map[string]interface{}{
				"err": err,
			},
				"error on when preparing image")
			return nil, err
		}

		sEnc := b64.StdEncoding.EncodeToString(newImg)
		// fmt.Println(sEnc)
		images = append(images, sEnc)
		// first image is will be to enroll first for get face id and save
		// data image to enrolled face table for get enrollment id
		if i == 0 {
			var startFaceEnroll time.Time = time.Now()
			var enrollmentData *entity.FaceEnrollment
			var err error

			// check user custom face enrollment or not for enroll face to fremis
			if postData.FaceID != "" {
				enrollmentData, err = s.FRemisRepo.AddFaceVariation(ctx, postData.FaceID, sEnc)
			} else {
				enrollmentData, err = s.FRemisRepo.FaceEnrollment(ctx, sEnc)
			}

			if err != nil {
				logutil.LogObj.SetErrorLog(map[string]interface{}{
					"err": err,
				},
					"error on when insert to Fremis Repo")
				return nil, err
			}
			logutil.LogObj.SetInfoLog(map[string]interface{}{"time_elapsed": fmt.Sprintf("%f Second", time.Since(startFaceEnroll).Seconds())}, "fremis face erollment time elapsed")
			faceIDFormat, _ := strconv.ParseUint(enrollmentData.FaceID, 10, 64)
			logutil.LogObj.SetInfoLog(map[string]interface{}{
				"face_id": enrollmentData.FaceID,
			},
				"Info faceIDFormat")

			// add data face to database
			tx, enrolledData, err = s.EnrolledFaceRepo.Create(ctx, tx, &entity.EnrolledFace{
				FaceID:         faceIDFormat,
				Name:           postData.Name,
				IdentityNumber: postData.IdentityNumber,
				Status:         postData.Status,
			})
			if err != nil {
				s.PsqlTransactionRepo.RollbackTransaction(ctx, tx)
				s.deleteFremis(ctx, faceID)
				logutil.LogObj.SetErrorLog(map[string]interface{}{
					"err": err,
				},
					"error on when insert to enrolledFaceRepo")
				return nil, err
			}
			enrollmentID = enrolledData.ID
			faceID = enrollmentData.FaceID
			variation = enrollmentData.Variation
		} else {
			// next and rest of image is face variation
			// add image variation to FREMIS
			startFaceEnroll := time.Now()
			variationData, err := s.FRemisRepo.AddFaceVariation(ctx, faceID, sEnc)
			if err != nil {
				s.PsqlTransactionRepo.RollbackTransaction(ctx, tx)
				s.deleteFremis(ctx, faceID)
				logutil.LogObj.SetErrorLog(map[string]interface{}{
					"err": err,
				},
					"error on when add face variation on FRemis")
				return nil, err
			}
			logutil.LogObj.SetInfoLog(map[string]interface{}{"time_elapsed": fmt.Sprintf("%f Second", time.Since(startFaceEnroll).Seconds())}, "fremis custom face erollment time elapsed")
			variation = variationData.Variation
		}

		// pre proccess image to create thumbnail before save image to database face image
		image, _, err := image.Decode(bytes.NewReader(newImg))
		if err != nil {
			s.PsqlTransactionRepo.RollbackTransaction(ctx, tx)
			s.deleteFremis(ctx, faceID)
			logutil.LogObj.SetErrorLog(map[string]interface{}{
				"err": err,
			},
				"error on when decode image")
			return nil, err
		}

		newImage := resize.Resize(200, 200, image, resize.Lanczos3)
		buf := new(bytes.Buffer)
		err = jpeg.Encode(buf, newImage, nil)
		if err != nil {
			s.PsqlTransactionRepo.RollbackTransaction(ctx, tx)
			s.deleteFremis(ctx, faceID)
			logutil.LogObj.SetErrorLog(map[string]interface{}{
				"err": err,
			},
				"error on when encode image")
			return nil, err
		}

		// save image to table face image
		tx, faceImage, err := s.FaceImageRepo.Create(ctx, tx, &entity.FaceImage{
			EnrolledFaceID: enrollmentID,
			Variation:      variation,
			Image:          newImg,
			ImageThumbnail: buf.Bytes(),
		})
		if err != nil {
			s.PsqlTransactionRepo.RollbackTransaction(ctx, tx)
			s.deleteFremis(ctx, faceID)
			logutil.LogObj.SetErrorLog(map[string]interface{}{
				"err": err,
			},
				"error on when insert to FaceImageRepo")
			return nil, err
		}
		// cast data to list of face
		faces = append(faces, &entity.FaceImage{
			ID:             faceImage.ID,
			Variation:      faceImage.Variation,
			ImageThumbnail: buf.Bytes(),
			CreatedAt:      faceImage.CreatedAt,
		})
	}

	// push event to agent if use CES
	if s.UseCES == "true" && isAgent != "true" {
		faceIDFormat, _ := strconv.ParseUint(faceID, 10, 64)
		dataAgent := entity.CreateEnrollmentEventCoordinator{
			EventAction: "create",
			Images:      images,
			Payload: entity.PayloadEnrollmentEventCoordinator{
				FaceID:         faceIDFormat,
				Name:           postData.Name,
				IdentityNumber: postData.IdentityNumber,
				Status:         postData.Status,
			},
		}
		err := s.AgentRepo.CreateEnrollmentEvent(ctx, &dataAgent)
		if err != nil {
			s.PsqlTransactionRepo.RollbackTransaction(ctx, tx)
			s.deleteFremis(ctx, faceID)
			logutil.LogObj.SetErrorLog(map[string]interface{}{
				"err": err,
			},
				"failed push event create enrollment to CES agent, rollback all event")
			return nil, err
		}
		logutil.LogObj.SetInfoLog(map[string]interface{}{}, "successfully push event create enrollment to CES agent")
		s.PsqlTransactionRepo.CommitTransaction(ctx, tx)
	} else {
		s.PsqlTransactionRepo.CommitTransaction(ctx, tx)
	}

	return &presenter.EnrollmentResponse{
		ID:             enrolledData.ID,
		Name:           enrolledData.Name,
		IdentityNumber: enrolledData.IdentityNumber,
		Status:         enrolledData.Status,
		FaceID:         enrolledData.FaceID,
		CreatedAt:      enrolledData.CreatedAt,
		UpdatedAt:      enrolledData.UpdatedAt,
		DeletedAt:      enrolledData.DeletedAt,
		Faces:          faces,
	}, nil
}

// Update is function for add enrollment
func (s *ServiceImpl) Update(ctx context.Context, enrollmentID uint64, postData *presenter.EnrollmentRequest, isAgent string) error {
	if s.UseCES == "true" && isAgent != "true" {
		err := s.AgentRepo.Ping(ctx)
		if err != nil {
			return errors.New("failed ping CES agent CUD opration disabled")
		}
	}
	detailEnroll, err := s.EnrolledFaceRepo.GetDetail(ctx, enrollmentID)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"error on get detail enrolled face by ID")
		return err
	}

	// init database transaction
	tx := s.PsqlTransactionRepo.BeginTransaction(ctx)
	faceID := strconv.FormatUint(detailEnroll.FaceID, 10)
	tx, err = s.EnrolledFaceRepo.Update(ctx, tx, &entity.EnrolledFace{
		ID:             enrollmentID,
		Name:           postData.Name,
		IdentityNumber: postData.IdentityNumber,
		Status:         postData.Status,
		FaceID:         detailEnroll.FaceID,
	})
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"error on update enrolled face")
		return err
	}

	var images []string
	var variations []string
	for _, v := range postData.Images {
		newImg, err := s.prepareImage(ctx, v.Image)
		if err != nil {
			s.PsqlTransactionRepo.RollbackTransaction(ctx, tx)
			if len(variations) > 0 {
				s.deleteFremisFaceVariation(ctx, faceID, variations)
			}
			logutil.LogObj.SetErrorLog(map[string]interface{}{
				"err": err,
			},
				"error on when preparing image")
			return err
		}

		sEnc := b64.StdEncoding.EncodeToString(newImg)
		images = append(images, sEnc)
		var startAddFaceVariation time.Time = time.Now()
		// next and rest of image is face variation
		variationData, err := s.FRemisRepo.AddFaceVariation(ctx, faceID, sEnc)
		logutil.LogObj.SetInfoLog(map[string]interface{}{"time_elapsed": fmt.Sprintf("%f Second", time.Since(startAddFaceVariation).Seconds())}, "fremis face erollment time elapsed")
		if err != nil {
			s.PsqlTransactionRepo.RollbackTransaction(ctx, tx)
			if len(variations) > 0 {
				s.deleteFremisFaceVariation(ctx, faceID, variations)
			}
			logutil.LogObj.SetErrorLog(map[string]interface{}{
				"err": err,
			},
				"error on when add face variation on FRemis")
			return err
		}
		variations = append(variations, variationData.Variation)
		variation := variationData.Variation

		// create image thumbnail
		image, _, err := image.Decode(bytes.NewReader(newImg))
		if err != nil {
			s.PsqlTransactionRepo.RollbackTransaction(ctx, tx)
			s.deleteFremisFaceVariation(ctx, faceID, variations)
			logutil.LogObj.SetErrorLog(map[string]interface{}{
				"err": err,
			},
				"error on when decode image")
			return err
		}
		newImage := resize.Resize(200, 200, image, resize.Lanczos3)

		buf := new(bytes.Buffer)
		err = jpeg.Encode(buf, newImage, nil)
		if err != nil {
			s.PsqlTransactionRepo.RollbackTransaction(ctx, tx)
			s.deleteFremisFaceVariation(ctx, faceID, variations)
			logutil.LogObj.SetErrorLog(map[string]interface{}{
				"err": err,
			},
				"error on when encode image")
			return err
		}

		// save image to table face image
		tx, _, err = s.FaceImageRepo.Create(ctx, tx, &entity.FaceImage{
			EnrolledFaceID: enrollmentID,
			Variation:      variation,
			Image:          newImg,
			ImageThumbnail: buf.Bytes(),
		})
		if err != nil {
			s.PsqlTransactionRepo.RollbackTransaction(ctx, tx)
			s.deleteFremisFaceVariation(ctx, faceID, variations)
			logutil.LogObj.SetErrorLog(map[string]interface{}{
				"err": err,
			},
				"error on when insert to FaceImageRepo")
			return err
		}
	}

	if s.UseCES == "true" && isAgent != "true" {
		dataAgent := entity.CreateEnrollmentEventCoordinator{
			EventAction: "update",
			Images:      images,
			Payload: entity.PayloadEnrollmentEventCoordinator{
				FaceID:         detailEnroll.FaceID,
				Name:           postData.Name,
				IdentityNumber: postData.IdentityNumber,
				Status:         postData.Status,
			},
		}
		err := s.AgentRepo.CreateEnrollmentEvent(ctx, &dataAgent)
		if err != nil {
			s.PsqlTransactionRepo.RollbackTransaction(ctx, tx)
			if len(variations) != 0 {
				s.deleteFremisFaceVariation(ctx, faceID, variations)
			}
			logutil.LogObj.SetErrorLog(map[string]interface{}{
				"err": err,
			},
				"failed push event update enrollment to CES agent, rollback all event")
			return err
		}
		logutil.LogObj.SetInfoLog(map[string]interface{}{}, "successfully push event update enrollment to CES agent")
		s.PsqlTransactionRepo.CommitTransaction(ctx, tx)
	} else {
		s.PsqlTransactionRepo.CommitTransaction(ctx, tx)
	}

	return nil
}

// GetDetail is for get detail enrollment by id
func (s *ServiceImpl) GetDetail(ctx context.Context, ID uint64) (*presenter.EnrollmentResponse, error) {
	enrolledData, err := s.EnrolledFaceRepo.GetDetail(ctx, ID)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"id":  ID,
			"err": err,
		},
			"error on when GetDetail repo")
		return nil, err
	}
	faces, err := s.FaceImageRepo.GetDetailByEnrollID(ctx, ID)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"id":  ID,
			"err": err,
		},
			"error on when GetDetailByEnrollID repo")
		return nil, err
	}

	return &presenter.EnrollmentResponse{
		ID:             enrolledData.ID,
		Name:           enrolledData.Name,
		FaceID:         enrolledData.FaceID,
		IdentityNumber: enrolledData.IdentityNumber,
		Status:         enrolledData.Status,
		CreatedAt:      enrolledData.CreatedAt,
		UpdatedAt:      enrolledData.UpdatedAt,
		DeletedAt:      enrolledData.DeletedAt,
		Faces:          faces,
	}, nil
}

// Delete Face ID
func (s *ServiceImpl) Delete(ctx context.Context, ID uint64, isAgent string) error {
	if s.UseCES == "true" && isAgent != "true" {
		err := s.AgentRepo.Ping(ctx)
		if err != nil {
			logutil.LogObj.SetErrorLog(map[string]interface{}{
				"err": err,
			},
				"failed ping CES agent CUD opration disabled")
			return errors.New("failed ping CES agent CUD opration disabled")
		}
	}
	detailFace, errGetDetail := s.EnrolledFaceRepo.GetDetail(ctx, ID)
	// for idempotent, if a record not found delete enrollment will assume success
	if errGetDetail != nil && errGetDetail.Error() == "record not found" {
		return nil
	}
	if errGetDetail != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": errGetDetail,
		},
			"error on get detail base one ID")
		return errGetDetail
	}
	logutil.LogObj.SetDebugLog(map[string]interface{}{
		"detail_face": detailFace,
	},
		"get detail enrolled face data")

	tx := s.PsqlTransactionRepo.BeginTransaction(ctx)
	tx, errDelete := s.EnrolledFaceRepo.Delete(ctx, tx, detailFace.ID)
	if errDelete != nil {
		s.PsqlTransactionRepo.RollbackTransaction(ctx, tx)
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": errDelete,
		},
			"error on delete from database")
		return errDelete
	}

	tx, errDeleteFaceImg := s.FaceImageRepo.DeleteByEnrollID(ctx, tx, detailFace.ID)
	if errDeleteFaceImg != nil {
		s.PsqlTransactionRepo.RollbackTransaction(ctx, tx)
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": errDelete,
		},
			"error on delete face image from database")
		return errDeleteFaceImg
	}

	errFremis := s.FRemisRepo.FaceDeleteEnrollment(ctx, []string{strconv.FormatUint(detailFace.FaceID, 10)})
	if errFremis != nil {
		s.PsqlTransactionRepo.RollbackTransaction(ctx, tx)
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": errFremis,
		},
			"error on delete faceId")
		return errFremis
	}

	if s.UseCES == "true" && isAgent != "true" {
		dataAgent := entity.CreateEnrollmentEventCoordinator{
			EventAction: "delete",
			Images:      []string{},
			Payload: entity.PayloadEnrollmentEventCoordinator{
				FaceID: detailFace.FaceID,
			},
		}
		err := s.AgentRepo.CreateEnrollmentEvent(ctx, &dataAgent)
		if err != nil {
			s.PsqlTransactionRepo.RollbackTransaction(ctx, tx)
			s.reEnrollFremis(ctx, detailFace.ID, strconv.FormatUint(detailFace.FaceID, 10))
			logutil.LogObj.SetErrorLog(map[string]interface{}{
				"err": err,
			},
				"failed push event delete enrollment to CES agent, rollback all event")
			return err
		}
		logutil.LogObj.SetInfoLog(map[string]interface{}{}, "successfully push event delete enrollment to CES agent")
		s.PsqlTransactionRepo.CommitTransaction(ctx, tx)
	} else {
		s.PsqlTransactionRepo.CommitTransaction(ctx, tx)
	}

	return nil
}

// Delete Face ID
func (s *ServiceImpl) DeleteAll(ctx context.Context) error {
	objs, errGetDetail := s.EnrolledFaceRepo.GetAll(ctx)
	if errGetDetail != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": errGetDetail,
		},
			"error on get detail base one ID")
		return errGetDetail
	}
	logutil.LogObj.SetDebugLog(map[string]interface{}{
		"obj": objs,
	},
		"log obj detail ")
	faceIds := []string{}
	for _, s := range objs {
		faceIds = append(faceIds, strconv.FormatUint(s.FaceID, 10))
	}
	errFremis := s.FRemisRepo.FaceDeleteEnrollment(ctx, faceIds)
	if errFremis != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": errFremis,
		},
			"error on delete faceId")
		return errFremis
	}
	errDeleteFaceImg := s.FaceImageRepo.DeleteAll(ctx)
	if errDeleteFaceImg != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": errDeleteFaceImg,
		},
			"error on delete face image from database")
		return errDeleteFaceImg
	}
	errDeleteFaceRepo := s.EnrolledFaceRepo.DeleteAll(ctx)
	if errDeleteFaceRepo != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": errDeleteFaceRepo,
		},
			"error on delete from database")
		return errDeleteFaceRepo
	}

	return nil
}

// GetFaceImage is usecase for get image face by id
func (s *ServiceImpl) GetFaceImage(ctx context.Context, ID uint64) ([]byte, error) {
	getFace, err := s.FaceImageRepo.GetImageByID(ctx, ID)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"error on get face by ID")
		return nil, err
	}
	return getFace.ImageThumbnail, nil
}
