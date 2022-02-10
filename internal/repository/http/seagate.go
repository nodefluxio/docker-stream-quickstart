package http

import (
	"bytes"
	"context"
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

// seagateServiceRepo connection seagate http REST API
type seagateServiceRepo struct {
	URL string
}

// NewSeagateServiceRepo will create an object that represent the analyticsetting.Repository interface
func NewSeagateServiceRepo(conn string) repository.Seagate {
	return &seagateServiceRepo{
		URL: conn,
	}
}

func (r *seagateServiceRepo) GetToken(ctx context.Context, requestData *entity.SeagateGetTokenRequest) (*entity.SeagateGetTokenResult, error) {
	// prepare endpoint
	newURL, err := url.Parse(r.URL)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "error parsing url")
		return nil, err
	}
	newURL.Path = "/v1/search"

	// prepare data request body
	postData, _ := json.Marshal(requestData)
	req, err := http.NewRequest("POST", newURL.String(), bytes.NewBuffer(postData))
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"url":   newURL.String(),
			"error": err,
		}, "Failed create new request to seagate service")
		return nil, err
	}

	// prepare header
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
			return nil, fmt.Errorf("error %d response from seagate API, with body : %s", res.StatusCode, string(body))
		}

		return nil, errors.New(response["message"].(string))
	}

	var resp *entity.SeagateGetTokenResult
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

	return resp, nil
}

func (r *seagateServiceRepo) GetFaceSearchResult(ctx context.Context, token string) (*entity.SeagateFaceSearchResult, error) {
	// checking parameter value
	if token == "" {
		err := errors.New("value token is empty")
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "seagate token cannot be empty")
		return nil, err
	}
	// prepare endpoint
	newURL, err := url.Parse(r.URL)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "error parsing url")
		return nil, err
	}
	newURL.Path = fmt.Sprintf("/v1/result/%s", token)

	// prepare data request
	req, err := http.NewRequest("GET", newURL.String(), nil)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"url":   newURL.String(),
			"error": err,
		}, "Failed create new request to seagate service")
		return nil, err
	}

	// prepare header
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
			return nil, fmt.Errorf("error %d response from seagate API, with body : %s", res.StatusCode, string(body))
		}

		return nil, errors.New(response["message"].(string))
	}

	var resp *entity.SeagateFaceSearchResult
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

	return resp, nil
}
