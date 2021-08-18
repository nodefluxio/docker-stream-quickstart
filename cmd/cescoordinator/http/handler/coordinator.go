package httphandler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/pkg/util"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/presenter"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/usecase/coordinator"
)

// CoordinatorHandler holds app interactor to be used by handler function
type CoordinatorHandler struct {
	CoordinatorSvc coordinator.Service
}

// Ping for add data event enrollment on CES coordinator
func (h *CoordinatorHandler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"ok":      true,
		"message": "pong",
	})
}

// Create for add data event enrollment on CES coordinator
func (h *CoordinatorHandler) Create(c *gin.Context) {
	var data presenter.CoordinatorRequest
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"ok":      false,
			"code":    "unprocessable-entity",
			"errors":  []string{err.Error()},
			"message": "Error when trying to parse to presenter",
		})
		return
	}
	err = h.CoordinatorSvc.Create(c, &data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"ok":      false,
			"code":    "unprocessable-entity",
			"errors":  []string{err.Error()},
			"message": "Error when trying to submit new event enrollement",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"ok":      true,
		"message": "Successfully submit new event enrollement",
	})
}

// GetEventEnrollment is handler for event GetEventEnrollment router
func (h *CoordinatorHandler) GetEventEnrollment(c *gin.Context) {
	var paging util.Pagination
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"ok":      false,
			"code":    "unprocessable-entity",
			"errors":  []string{err.Error()},
			"message": "query parameter limit must be number",
		})
		return
	}
	paging.Limit = limit
	paging.Filter = c.QueryMap("filter")

	data, err := h.CoordinatorSvc.Get(c, &paging)
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
		"message": "successfully fetch data event enrollement",
		"results": data,
	})
}
