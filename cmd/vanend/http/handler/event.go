package httphandler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	logutil "gitlab.com/nodefluxio/goutils/pkg/log"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/pkg/util"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/usecase/event"
)

// EventHandler holds app interactor to be used by handler function
type EventHandler struct {
	EventSvc event.Service
}

//GetEvent for get all data in default tier
func (h *EventHandler) GetEvent(c *gin.Context) {
	// step
	// handle connection to vanilla dashboard
	// handle check to database
	// handler transfer things
	wsupgrader.CheckOrigin = func(r *http.Request) bool { return true }
	streamID := c.Query("stream_id")

	conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"failed when update connection")
		return
	}
	h.EventSvc.InitiateDataStream(streamID, conn)
}

//Dumping for get all data in default tier
func (h *EventHandler) Dumping(c context.Context) {
	go h.EventSvc.Dumping(c)
}

// History is handler for event history router
func (h *EventHandler) History(c *gin.Context) {
	var paging util.Pagination
	paging.Limit, _ = strconv.Atoi(c.Query("limit"))
	paging.Sort = c.QueryMap("sort")
	paging.Page, _ = strconv.Atoi(c.Query("page"))
	paging.Filter = c.QueryMap("filter")
	paging.Search = c.Query("search")

	events, err := h.EventSvc.GetHistory(c, &paging)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"ok":      false,
			"code":    "internal-server-error",
			"errors":  []string{},
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ok":      true,
		"message": "successfully fetch event data",
		"results": events,
	})
}

// Export is handler for event Export router
func (h *EventHandler) Export(c *gin.Context) {
	var paging util.Pagination
	paging.Sort = c.QueryMap("sort")
	paging.Page, _ = strconv.Atoi(c.Query("page"))
	paging.Filter = c.QueryMap("filter")
	paging.Search = c.Query("search")

	err := h.EventSvc.ExportEvent(c, &paging)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"ok":      false,
			"code":    "bad-request",
			"errors":  []string{},
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ok":      true,
		"message": "successfully export event data",
	})
}

// CheckExportStatus is handler for check event Export router
func (h *EventHandler) CheckExportStatus(c *gin.Context) {
	status, err := h.EventSvc.CheckExportedEvent(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"ok":      false,
			"code":    "internal-server-error",
			"errors":  []string{},
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ok":      true,
		"message": "successfully get status export event",
		"status":  status,
	})
}

// DownloadExport is handler for download event Export router
func (h *EventHandler) DownloadExport(c *gin.Context) {
	status, err := h.EventSvc.CheckExportedEvent(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"ok":      false,
			"code":    "internal-server-error",
			"errors":  []string{},
			"message": err.Error(),
		})
		return
	}
	if status != "ready" {
		c.JSON(http.StatusOK, gin.H{
			"ok":      false,
			"message": "file export not yet ready to download",
		})

		return
	}
	filename := "exported_event.zip"
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.File(fmt.Sprintf("./tmp/%s", filename))
	h.EventSvc.UpdateStatusExportDownload(c)
}
