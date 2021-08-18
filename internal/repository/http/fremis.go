package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	logutil "gitlab.com/nodefluxio/goutils/pkg/log"

	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/repository"
)

// fremisRepo connection FRemis service
type fremisRepo struct {
	BaseURL  string
	Keyspace string
}

type postData struct {
	Image    string `json:"image"`
	Keyspace string `json:"keyspace"`
}

type responsePostData struct {
	FaceID    string `json:"face_id"`
	Variation string `json:"variation"`
}

type deleteData struct {
	Keyspace string   `json:"keyspace"`
	FaceIds  []string `json:"face_ids"`
}

type responseDeleteData struct {
	DeletedFaceIds []string `json:"deleted_face_ids"`
}
type postDataFR struct {
	Image            string                 `json:"image"`
	Keyspace         string                 `json:"keyspace"`
	AdditionalParams map[string]interface{} `json:"additional_params"`
}

type postDeleteVariation struct {
	Keyspace   string   `json:"keyspace"`
	FaceID     string   `json:"face_id"`
	Variations []string `json:"variations"`
}

// NewFremisRepository will create an object that represent the analyticsetting.Repository interface
func NewFremisRepository(baseURL, keyspace string) repository.FRemis {
	return &fremisRepo{
		BaseURL:  baseURL,
		Keyspace: keyspace,
	}
}

func (r *fremisRepo) FaceEnrollment(ctx context.Context, image string) (*entity.FaceEnrollment, error) {
	data := postData{
		Image:    image,
		Keyspace: r.Keyspace,
	}
	logutil.LogObj.SetInfoLog(map[string]interface{}{"keyspace": r.Keyspace}, "keyspace value")
	logutil.LogObj.SetDebugLog(map[string]interface{}{"data": data}, "data to post")

	postData, _ := json.Marshal(data)
	apiURL := r.BaseURL + "/enrollment"
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(postData))
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"fremisn_url": apiURL,
			"error":       err,
			"data":        postData,
		}, "Failed create new request to fremisn service")
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "Error Create client "+apiURL)
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "Error when trying to read body")
		return nil, err
	}
	logutil.LogObj.SetInfoLog(map[string]interface{}{
		"http_code": res.StatusCode,
		"body":      string(body),
	}, "Log body FaceEnrollment")

	var resp map[string]interface{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "Error when trying unmarshal data")
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		err = errors.New(resp["description"].(string))
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error":         err,
			"response_body": string(body),
		}, "Error enroll face")
		return nil, err
	}

	logutil.LogObj.SetDebugLog(map[string]interface{}{
		"response": resp,
	}, "Log response FaceEnrollment")
	return &entity.FaceEnrollment{
		FaceID:    resp["face_id"].(string),
		Variation: resp["variation"].(string),
	}, nil
}

func (r *fremisRepo) FaceDeleteEnrollment(ctx context.Context, faceIDs []string) error {
	data := deleteData{
		Keyspace: r.Keyspace,
		FaceIds:  faceIDs,
	}
	postData, _ := json.Marshal(data)
	logutil.LogObj.SetInfoLog(map[string]interface{}{
		"data": data,
	}, "info deleted data faceIDs")
	req, err := http.NewRequest("POST", r.BaseURL+"/delete-enrollment", bytes.NewBuffer(postData))
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
			"data":  postData,
		}, "Create New Request")
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "Error fetch "+r.BaseURL+"/delete-enrollment")
		return err
	}
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "Error when trying to read body")
		return err
	}
	logutil.LogObj.SetInfoLog(map[string]interface{}{
		"body": body,
	}, "information of body")
	var resp map[string]interface{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "Error when trying unmarshal data")
		return err
	}
	logutil.LogObj.SetInfoLog(map[string]interface{}{
		"response": resp,
	}, "Response data from Face Delete Enrollment")
	// if len(resp["deleted_face_ids"].([]interface{})) < 1 {
	// 	logutil.LogObj.SetErrorLog(map[string]interface{}{
	// 		"response": resp,
	// 	}, "Error response deleted face_id not return anything")
	// 	return errors.New("deleted face_id not return in response")
	// }

	return nil
}

func (r *fremisRepo) FaceRecognition(image string) ([]*entity.FRemisCandidate, error) {

	data := postDataFR{
		Keyspace: r.Keyspace,
		Image:    image,
		AdditionalParams: map[string]interface{}{
			"candidateCount": 1,
		},
	}
	postData, _ := json.Marshal(data)

	req, err := http.NewRequest("POST", r.BaseURL+"/recognition", bytes.NewBuffer(postData))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
			"data":  postData,
		}, "Create New Request")
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "Error fetch "+r.BaseURL+"/recognition")
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "Error when trying to read body")
		return nil, err
	}
	var resp map[string]interface{}
	err = json.Unmarshal(body, &resp)
	result := resp["result"].(map[string]interface{})
	frResult := result["face_recognition"].(map[string]interface{})
	listCandidate := frResult["candidates"].([]entity.FRemisCandidate)
	logutil.LogObj.SetDebugLog(map[string]interface{}{
		"response":       resp,
		"result":         result,
		"fr_result":      frResult,
		"list_candidate": listCandidate,
	}, "Debug log for response")
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "Error when trying unmarshal data")
		return nil, err
	}

	return nil, err
}

func (r *fremisRepo) AddFaceVariation(ctx context.Context, faceID, image string) (*entity.FaceEnrollment, error) {
	data := postDataFR{
		Image:    image,
		Keyspace: r.Keyspace,
		AdditionalParams: map[string]interface{}{
			"face_id": faceID,
		},
	}
	logutil.LogObj.SetInfoLog(map[string]interface{}{"keyspace": r.Keyspace}, "keyspace value")
	logutil.LogObj.SetDebugLog(map[string]interface{}{"data": data}, "data to post")

	postData, _ := json.Marshal(data)
	apiURL := r.BaseURL + "/enrollment"
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(postData))
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"fremisn_url": apiURL,
			"error":       err,
			"data":        postData,
		}, "Failed create new request to fremisn service")
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
			"data":  postData,
		}, "Request face enrollment ")
		return nil, err
	}

	res, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "Error Create client "+apiURL)
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "Error when trying to read body")
		return nil, err
	}
	logutil.LogObj.SetInfoLog(map[string]interface{}{
		"http_code": res.StatusCode,
		"body":      string(body),
	}, "Log body AddFaceVariation")

	var resp map[string]interface{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "Error when trying unmarshal data")
		return nil, err
	}
	logutil.LogObj.SetDebugLog(map[string]interface{}{
		"response": resp,
	}, "Log response AddFaceVariation")

	if res.StatusCode != http.StatusOK {
		err = errors.New(resp["description"].(string))
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error":         err,
			"response_body": string(body),
		}, "Error add face variation")
		return nil, err
	}

	return &entity.FaceEnrollment{
		FaceID:    resp["face_id"].(string),
		Variation: resp["variation"].(string),
	}, nil
}

func (r *fremisRepo) DeleteFaceVariation(ctx context.Context, faceID string, variations []string) error {
	data := postDeleteVariation{
		FaceID:     faceID,
		Keyspace:   r.Keyspace,
		Variations: variations,
	}
	logutil.LogObj.SetInfoLog(map[string]interface{}{"keyspace": r.Keyspace}, "keyspace value")
	logutil.LogObj.SetDebugLog(map[string]interface{}{"data": data}, "data to post")

	postData, _ := json.Marshal(data)
	apiURL := r.BaseURL + "/delete-enrollment-variation"
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(postData))
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"fremisn_url": apiURL,
			"error":       err,
			"data":        postData,
		}, "Failed create new request to fremisn service")
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
			"data":  postData,
		}, "Request delete face variation ")
		return err
	}

	res, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "Error Create client "+apiURL)
		return err
	}

	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "Error when trying to read body")
		return err
	}
	logutil.LogObj.SetInfoLog(map[string]interface{}{
		"body": string(body),
	}, "Log body delete face variation")

	var resp map[string]interface{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "Error when trying unmarshal data")
		return err
	}
	logutil.LogObj.SetDebugLog(map[string]interface{}{
		"response": resp,
	}, "Log response delete face variation")

	return nil
}
