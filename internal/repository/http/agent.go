package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	logutil "gitlab.com/nodefluxio/goutils/pkg/log"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/repository"
)

// agentServiceRepo connection CES coordinator
type agentServiceRepo struct {
	URL string
}

// NewAgentServiceRepo will create an object that represent the analyticsetting.Repository interface
func NewAgentServiceRepo(conn string) repository.Agent {
	return &agentServiceRepo{
		URL: conn,
	}
}

func (r *agentServiceRepo) Ping(ctx context.Context) error {
	newURL, err := url.Parse(r.URL)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "error parsing url")
		return err
	}
	newURL.Path = "v1/agents/ping"

	resp, err := http.Get(newURL.String())
	logutil.LogObj.SetDebugLog(map[string]interface{}{
		"resp": resp,
	}, "Data response")

	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "ping agent failed")
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

	if resp.StatusCode != http.StatusOK {
		err = errors.New("ping agent failed")
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error":         err,
			"response_body": string(body),
			"url":           newURL.String(),
		}, "Error ping agent")
		return err
	}

	return nil
}

func (r *agentServiceRepo) CreateEnrollmentEvent(ctx context.Context, data *entity.CreateEnrollmentEventCoordinator) error {
	newURL, err := url.Parse(r.URL)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "error parsing url")
		return err
	}
	newURL.Path = "v1/agents/event-enrollments"

	postData, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", newURL.String(), bytes.NewBuffer(postData))
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"url":   newURL.String(),
			"error": err,
		}, "Failed create new request to create enrollment event at agent service")
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
