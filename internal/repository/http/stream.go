package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	logutil "gitlab.com/nodefluxio/goutils/pkg/log"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/repository"
)

// streamServiceRepo connection visionaire v4 stream service
type streamServiceRepo struct {
	URL string
}

// NewStreamServiceRepo will create an object that represent the analyticsetting.Repository interface
func NewStreamServiceRepo(conn string) repository.Stream {
	return &streamServiceRepo{
		URL: conn,
	}
}

func (r *streamServiceRepo) GetDetail(ctx context.Context, nodeNumber int64, streamID string) (*entity.VisionaireStreamDetail, error) {
	var result entity.VisionaireStreamDetail
	url := fmt.Sprintf("%s/streams/%d/%s", r.URL, nodeNumber, streamID)
	resp, err := http.Get(url)
	logutil.LogObj.SetDebugLog(map[string]interface{}{
		"resp": resp,
	}, "Data response")

	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error":       err,
			"node_number": nodeNumber,
			"stream_id":   streamID,
		}, "Request get stream detail failed")
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

	// http code error handling
	if resp.StatusCode >= http.StatusBadRequest {
		err = fmt.Errorf("error %d response from visionaire streams API, with body : %s", resp.StatusCode, string(body))
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"url":         url,
			"status_code": resp.StatusCode,
			"error":       err,
		}, "Error when trying hit streams api")

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
		"node_number": nodeNumber,
		"stream_id":   streamID,
		"response":    &result,
	}, "Info get detail stream")

	return &result, nil
}

func (r *streamServiceRepo) GetList(ctx context.Context) (*entity.VisionaireStream, error) {
	var result entity.VisionaireStream
	newURL, err := url.Parse(r.URL)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
			"url":   r.URL,
		}, "error parsing url")
		return nil, err
	}
	newURL.Path = "streams"
	q := newURL.Query()
	newURL.RawQuery = q.Encode()

	resp, err := http.Get(newURL.String())
	logutil.LogObj.SetDebugLog(map[string]interface{}{
		"resp": resp,
	}, "Data response")

	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
			"url":   newURL.String(),
		}, "Request get list streams failed")
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	logutil.LogObj.SetDebugLog(map[string]interface{}{
		"body": string(body),
	}, "Data body")

	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"error": err,
		}, "Error when trying to read body")
		return nil, err
	}

	// http code error handling
	if resp.StatusCode >= http.StatusBadRequest {
		err = fmt.Errorf("error %d response from visionaire streams API, with body : %s", resp.StatusCode, string(body))
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"url":         newURL.String(),
			"status_code": resp.StatusCode,
			"error":       err,
		}, "Error when trying hit streams api")

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
	}, "Info get list streams")

	return &result, nil
}
