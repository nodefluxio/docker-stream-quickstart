package polrisearching

import (
	"bytes"
	"context"
	"errors"
	"image"
	"image/jpeg"
	"os"

	"encoding/base64"
	"encoding/json"

	"github.com/nfnt/resize"
	logutil "gitlab.com/nodefluxio/goutils/pkg/log"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/pkg/imageprocessing"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/presenter"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/repository"
)

// ServiceImpl struct to represent quota transaction service
type ServiceImpl struct {
	PolriRepo          repository.Polri
	SeagateRepo        repository.Seagate
	FRemisRepo         repository.FRemis
	MaxSizeImageUpload int64
}

const (
	defaultLimitFaceSearch = 10
	maxImageWidth          = 2000
	maxImageHeight         = 2000
	size                   = 1000
	sizeForAspecRatio      = 0 // 0 for maintain aspec ratio
	resolutionDivider      = 2
	dateLayout             = "2006-01-02"
)

// CreateOrUpdate is function for create or update global setting
func (s *ServiceImpl) SearchPlateNumber(ctx context.Context, nopol string) (*presenter.PolriSearchPlateResponse, error) {
	data, err := s.PolriRepo.SearchPlateNumber(ctx, nopol)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"error on when get data from polri API")
		return nil, err
	}
	newData := presenter.PolriSearchPlateResponse(*data)
	return &newData, nil
}

// SearchNik is function for create or update global setting
func (s *ServiceImpl) SearchNik(ctx context.Context, nik string) (*presenter.PolriCitizenResponse, error) {
	data, err := s.PolriRepo.SearchNIK(ctx, nik)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"error on when get data from polri API")
		return nil, err
	}
	newData := presenter.PolriCitizenResponse(*data)
	return &newData, nil
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

	if img.Width > maxImageWidth || img.Height > maxImageHeight || imgSize > s.MaxSizeImageUpload {
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

			if newImgSize < s.MaxSizeImageUpload {
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

// GetFaceSearchToken is function for get token face search from seagate service
func (s *ServiceImpl) GetFaceSearchToken(ctx context.Context, Image []byte, limit uint64) (string, error) {
	// prepare face embedding data
	newImg, err := s.prepareImage(ctx, Image)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"error on when preparing image")
		return "", err
	}
	encodedImg := base64.StdEncoding.EncodeToString(newImg)
	embeddingData, err := s.FRemisRepo.GetFaceEmbedings(ctx, encodedImg)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"error on when get data embedding")
		return "", err
	}
	if len(embeddingData.Embeddings) < 1 {
		logutil.LogObj.SetErrorLog(map[string]interface{}{}, "unrecognized Image Data")
		return "", errors.New("unrecognized Image Data")
	}
	stringifyEmbeddingData, err := json.MarshalIndent(embeddingData.Embeddings[0], "", " ")
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"error on when stringify embedding")
		return "", err
	}

	// get token from seagate
	if limit == 0 {
		limit = defaultLimitFaceSearch
	}
	requestData := &entity.SeagateGetTokenRequest{
		Embedding: string(stringifyEmbeddingData),
		Limit:     limit,
	}
	token, err := s.SeagateRepo.GetToken(ctx, requestData)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"error on when get token from seagate API")
		return "", err
	}
	return token.Token, nil
}

func (s *ServiceImpl) GetFaceSearchResult(ctx context.Context, token string) ([]*presenter.PolriFaceResultResponse, error) {
	dataFaceSearch, err := s.SeagateRepo.GetFaceSearchResult(ctx, token)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"error on when get result from seagate API")
		return nil, err
	}

	result := make([]*presenter.PolriFaceResultResponse, 0)
	for _, v := range dataFaceSearch.Data {
		newData := &presenter.PolriFaceResultResponse{
			Similiarity:  v.Similiarity,
			Token:        v.Token,
			Probability:  v.Probability,
			DukcapilData: &v.Dukcapil,
			DukcapilStatus: entity.SeagateDukcapilStatus{
				Ok:     true,
				Respon: "Data Ditemukan",
			},
		}
		if v.Dukcapil.Ok && v.Dukcapil.Respon == "Data Tidak Ditemukan" {
			newData.DukcapilStatus = entity.SeagateDukcapilStatus{
				Ok:     false,
				Respon: v.Dukcapil.Respon,
			}
			newData.DukcapilData = nil
		}
		result = append(result, newData)
	}
	return result, nil
}
