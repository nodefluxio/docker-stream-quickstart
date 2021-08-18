package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"

	logutil "gitlab.com/nodefluxio/goutils/pkg/log"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/repository"
)

// enrollmentVanillaServiceRepo connection visionaire v4 stream service
type enrollmentVanillaServiceRepo struct {
	URL string
}

// NewEnrollmentVanillaServiceRepo will create an object that represent the analyticsetting.Repository interface
func NewEnrollmentVanillaServiceRepo(conn string) repository.EnrollmentVanilla {
	return &enrollmentVanillaServiceRepo{
		URL: conn,
	}
}

func (r *enrollmentVanillaServiceRepo) CreateFaceEnrollment(ctx context.Context, data *entity.VanillaEnrollmentPayload, images []*entity.EnrollmentImage) error {
	url := fmt.Sprintf("%s/api/enrollment", r.URL)
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	if len(images) > 0 {
		for i, v := range images {
			part, err := writer.CreateFormFile("images", fmt.Sprintf("image-%d", i))
			if err != nil {
				logutil.LogObj.SetErrorLog(map[string]interface{}{
					"error": err,
				}, "Error create form file")
				return err
			}
			_, err = io.Copy(part, v.Image)
			if err != nil {
				logutil.LogObj.SetErrorLog(map[string]interface{}{
					"error": err,
				}, "Error copy file to form file")
				return err
			}
		}
	}

	writer.WriteField("identity_number", data.IdentityNumber)
	writer.WriteField("name", data.Name)
	writer.WriteField("status", data.Status)
	writer.WriteField("face_id", strconv.FormatUint(data.FaceID, 10))
	err := writer.Close()
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "Error create writer")
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, method, url, payload)
	req.Header.Set("Is-Agent", "true")

	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "Error create new request")
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "Error do request")
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "Error read body")
		return err
	}

	if res.StatusCode != http.StatusCreated {
		err = errors.New("create new enrollment failed")
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error":         err,
			"response_body": string(body),
			"url":           url,
		}, "Error when trying create enrollment vanilla")
		return err
	}
	return err
}

func (r *enrollmentVanillaServiceRepo) UpdateFaceEnrollment(ctx context.Context, ID uint64, data *entity.VanillaEnrollmentPayload, images []*entity.EnrollmentImage) error {
	url := fmt.Sprintf("%s/api/enrollment/%d", r.URL, ID)
	method := "PUT"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	if len(images) > 0 {
		for i, v := range images {
			part, err := writer.CreateFormFile("images", fmt.Sprintf("image-%d", i))
			if err != nil {
				logutil.LogObj.SetErrorLog(map[string]interface{}{
					"error": err,
				}, "Error create form file")
				return err
			}
			_, err = io.Copy(part, v.Image)
			if err != nil {
				logutil.LogObj.SetErrorLog(map[string]interface{}{
					"error": err,
				}, "Error copy file to form file")
				return err
			}
		}
	}

	writer.WriteField("identity_number", data.IdentityNumber)
	writer.WriteField("name", data.Name)
	writer.WriteField("status", data.Status)
	writer.WriteField("face_id", strconv.FormatUint(data.FaceID, 10))
	err := writer.Close()
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "Error create writer")
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, method, url, payload)
	req.Header.Set("Is-Agent", "true")
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "Error create new request")
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "Error do request")
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "Error read body")
		return err
	}

	if res.StatusCode != http.StatusOK {
		err = errors.New("update enrollment failed")
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error":         err,
			"response_body": string(body),
			"url":           url,
		}, "Error when trying update enrollment vanilla")
		return err
	}
	return err
}

func (r *enrollmentVanillaServiceRepo) DeleteFaceEnrollment(ctx context.Context, ID uint64) error {
	url := fmt.Sprintf("%s/api/enrollment/%d", r.URL, ID)
	method := "DELETE"

	payload := &bytes.Buffer{}
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, method, url, payload)
	req.Header.Set("Is-Agent", "true")

	if err != nil {
		fmt.Println(err)
		return err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if res.StatusCode != http.StatusOK {
		err = errors.New("delete enrollment failed")
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error":         err,
			"response_body": string(body),
			"url":           url,
		}, "Error when trying delete enrollment vanilla")
		return err
	}

	return err
}

func (r *enrollmentVanillaServiceRepo) GetByFaceID(ctx context.Context, faceID uint64) (*entity.VanillaEnrollmentData, error) {
	var result entity.VanillaEnrollmentData
	url := fmt.Sprintf("%s/api/enrollment?filter[face_id]=%d", r.URL, faceID)

	resp, err := http.Get(url)
	logutil.LogObj.SetDebugLog(map[string]interface{}{
		"resp": resp,
	}, "Data response")

	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "Request get data vanille enrollment failed")
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	logutil.LogObj.SetDebugLog(map[string]interface{}{
		"body": body,
	}, "Data body")
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "Error when trying to read body")
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "Error when trying unmarshal data")
		return nil, err
	}

	logutil.LogObj.SetDebugLog(map[string]interface{}{
		"response": &result,
	}, "Info get data vanille enrollment")

	return &result, nil
}
