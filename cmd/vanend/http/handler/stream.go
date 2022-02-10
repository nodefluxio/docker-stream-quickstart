package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	logutil "gitlab.com/nodefluxio/goutils/pkg/log"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/pkg/util"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/presenter"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/usecase/stream"
)

// StreamHandler holds app interactor to be used by handler function
type StreamHandler struct {
	StreamSvc stream.Service
}

// GetList for get all data stream
func (h *StreamHandler) GetList(c *gin.Context) {
	var paging util.Pagination
	accInfo := c.MustGet("account_info").(*presenter.AuthInfoResponse)
	paging.Filter = c.QueryMap("filter")

	streamList, err := h.StreamSvc.GetList(c, &paging, accInfo)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"ok":      false,
			"code":    "internal-server-error",
			"errors":  []string{err.Error()},
			"message": "failed get list of stream",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ok":      true,
		"message": "successfully get list stream",
		"data":    streamList,
	})
}

// GetDetail for get detail data stream
func (h *StreamHandler) GetDetail(c *gin.Context) {
	accInfo := c.MustGet("account_info").(*presenter.AuthInfoResponse)
	nodeID, err := strconv.ParseInt(c.Param("node_id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"ok":      false,
			"code":    "bad-request",
			"errors":  []string{err.Error()},
			"message": "failed parsing node id, value data must be number",
		})
		return
	}
	req := presenter.StreamRequest{
		NodeNumber: nodeID,
		StreamID:   c.Param("stream_id"),
	}

	stream, err := h.StreamSvc.GetDetail(c, &req, accInfo)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"request_data": req,
			"err":          err,
		},
			"error when get detail stream usecase")

		httpStatus := http.StatusInternalServerError
		code := "internal-server-error"

		if err.Error() == "stream not found" {
			httpStatus = http.StatusBadRequest
			code = "bad-request"
		}

		c.AbortWithStatusJSON(httpStatus, gin.H{
			"ok":      false,
			"code":    code,
			"errors":  []string{err.Error()},
			"message": "failed get detail of stream",
		})
		return

	}
	c.JSON(http.StatusOK, gin.H{
		"ok":      true,
		"message": "successfully get detail stream",
		"stream":  stream,
	})
}
