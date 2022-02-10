package http

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	logutil "gitlab.com/nodefluxio/goutils/pkg/log"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/repository"
)

// polriServiceRepo connection polri http REST API
type polriServiceRepo struct {
	URL             string
	Username        string
	Password        string
	SearchPlatePath string
	SearchNikPath   string
}

type searchPlateRequest struct {
	Nopol string `json:"nopol"`
}

type searchNikRequest struct {
	Nik string `json:"nik"`
}

// NewPolriServiceRepo will create an object that represent the analyticsetting.Repository interface
func NewPolriServiceRepo(conn, username, password, searchPlatePath, searchNikPath string) repository.Polri {
	return &polriServiceRepo{
		URL:             conn,
		Username:        username,
		Password:        password,
		SearchPlatePath: searchPlatePath,
		SearchNikPath:   searchNikPath,
	}
}

func (r *polriServiceRepo) SearchPlateNumber(ctx context.Context, plateNumber string) (*entity.PolriPlateNuberInfo, error) {
	// prepare endpoint
	newURL, err := url.Parse(r.URL)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "error parsing url")
		return nil, err
	}
	newURL.Path = r.SearchPlatePath

	// prepare data request body
	requestData := searchPlateRequest{
		Nopol: plateNumber,
	}
	postData, _ := json.Marshal(requestData)
	req, err := http.NewRequest("POST", newURL.String(), bytes.NewBuffer(postData))
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"coordinator": newURL.String(),
			"error":       err,
		}, "Failed create new request to coordinator service")
		return nil, err
	}

	// prepare header
	b64BasicAuth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", r.Username, r.Password)))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", b64BasicAuth))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// make a request
	res, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "Error Create client "+newURL.String())
		return nil, err
	}

	// read response body
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "Error when trying to read body")
		return nil, err
	}
	logutil.LogObj.SetDebugLog(map[string]interface{}{
		"http_code": res.StatusCode,
		"body":      string(body),
	}, "Log body")

	// http code error handling
	if res.StatusCode >= http.StatusBadRequest {
		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		if err != nil {
			logutil.LogObj.SetErrorLog(map[string]interface{}{
				"error": err,
			}, "Error when trying unmarshal data")
			return nil, fmt.Errorf("error %d response from polri API, with body : %s", res.StatusCode, string(body))
		}

		return nil, errors.New(response["message"].(string))
	}

	var respBody []map[string]interface{}
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "Error when trying unmarshal data")
		return nil, err
	}

	if respBody[0]["STATUS"] == "0" {
		return nil, errors.New(respBody[0]["DESKRIPSI"].(string))
	}

	var resp []*entity.PolriPlateNuberInfo
	err = json.Unmarshal(body, &resp)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "Error when trying unmarshal data")
		return nil, err
	}
	logutil.LogObj.SetDebugLog(map[string]interface{}{
		"response": resp,
	}, "Log response")
	return resp[0], nil
}

func (r *polriServiceRepo) SearchNIK(ctx context.Context, NIK string) (*entity.PolriCitizenData, error) {
	// prepare endpoint
	newURL, err := url.Parse(r.URL)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "error parsing url")
		return nil, err
	}
	newURL.Path = r.SearchNikPath

	// prepare data request body
	requestData := searchNikRequest{
		Nik: NIK,
	}
	postData, _ := json.Marshal(requestData)
	req, err := http.NewRequest("POST", newURL.String(), bytes.NewBuffer(postData))
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"coordinator": newURL.String(),
			"error":       err,
		}, "Failed create new request to coordinator service")
		return nil, err
	}

	// prepare header
	b64BasicAuth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", r.Username, r.Password)))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", b64BasicAuth))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// make a request
	res, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "Error Create client "+newURL.String())
		return nil, err
	}

	// read response body
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "Error when trying to read body")
		return nil, err
	}
	logutil.LogObj.SetDebugLog(map[string]interface{}{
		"http_code": res.StatusCode,
		"body":      string(body),
	}, "Log body")

	// http code error handling
	if res.StatusCode >= http.StatusBadRequest {
		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		if err != nil {
			logutil.LogObj.SetErrorLog(map[string]interface{}{
				"error": err,
			}, "Error when trying unmarshal data")
			return nil, fmt.Errorf("error %d response from polri API, with body : %s", res.StatusCode, string(body))
		}

		return nil, errors.New(response["message"].(string))
	}

	var resp *entity.PolriCitizenResponseData
	err = json.Unmarshal(body, &resp)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "Error when trying unmarshal data")
		return nil, err
	}

	if resp.Payload.Content[0].Respon != nil {
		return nil, errors.New(resp.Payload.Content[0].Respon.(string))
	}

	logutil.LogObj.SetDebugLog(map[string]interface{}{
		"response": resp,
	}, "Log response")
	return &resp.Payload.Content[0], nil
}
