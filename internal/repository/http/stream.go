package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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

func (r *streamServiceRepo) GetDetail(ctx context.Context, nodeNumber int64, streamID string) (*entity.StreamDetail, error) {
	var result entity.StreamDetail
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
