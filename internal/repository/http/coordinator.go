package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	logutil "gitlab.com/nodefluxio/goutils/pkg/log"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/repository"
)

// coordinatorerviceRepo connection CES coordinator
type coordinatorServiceRepo struct {
	URL string
}

// NewCoordinatorServiceRepo will create an object that represent the analyticsetting.Repository interface
func NewCoordinatorServiceRepo(conn string) repository.Coordinator {
	return &coordinatorServiceRepo{
		URL: conn,
	}
}

func (r *coordinatorServiceRepo) Ping(ctx context.Context) error {
	newURL, err := url.Parse(r.URL)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "error parsing url")
		return err
	}
	newURL.Path = "v1/coordinators/ping"

	resp, err := http.Get(newURL.String())
	logutil.LogObj.SetDebugLog(map[string]interface{}{
		"resp": resp,
	}, "Data response")

	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "ping coordinator failed")
		return err
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
		return err
	}

	return nil
}

func (r *coordinatorServiceRepo) GetEnrollmentEvent(ctx context.Context, limit, latestTimestamp string) (*entity.EnrollmentEventCoordinator, error) {
	var result entity.EnrollmentEventCoordinator
	newURL, err := url.Parse(r.URL)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "error parsing url")
		return nil, err
	}
	newURL.Path = "v1/coordinators"
	q := newURL.Query()
	if latestTimestamp != "" {
		q.Set("filter[latest_timestamp]", latestTimestamp)
	}
	q.Set("limit", limit)
	newURL.RawQuery = q.Encode()

	resp, err := http.Get(newURL.String())
	logutil.LogObj.SetDebugLog(map[string]interface{}{
		"resp": resp,
	}, "Data response")

	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "Request get event enrollment coordinator failed")
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
	}, "Info get event enrollment coordinator")

	return &result, nil
}

func (r *coordinatorServiceRepo) CreateEnrollmentEvent(ctx context.Context, data *entity.CreateEnrollmentEventCoordinator) error {
	newURL, err := url.Parse(r.URL)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "error parsing url")
		return err
	}
	newURL.Path = "v1/coordinators"

	postData, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", newURL.String(), bytes.NewBuffer(postData))
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"coordinator": newURL.String(),
			"error":       err,
		}, "Failed create new request to coordinator service")
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "Error Create client "+newURL.String())
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
	}, "Log body")
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
	}, "Log response")
	return nil
}
