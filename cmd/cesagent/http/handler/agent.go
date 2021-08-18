package httphandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/presenter"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/usecase/agent"
)

// AgentHandler holds app interactor to be used by handler function
type AgentHandler struct {
	AgentSvc agent.Service
}

// Ping for add data event enrollment on CES coordinator
func (h *AgentHandler) Ping(c *gin.Context) {
	err := h.AgentSvc.PingCoordinator(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"ok":      false,
			"code":    "internal-server-error",
			"errors":  []string{err.Error()},
			"message": "Error when trying to ping CES coordinator",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ok":      true,
		"message": "pong",
	})
}

// Create for add data event enrollment on CES coordinator
func (h *AgentHandler) Create(c *gin.Context) {
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
	err = h.AgentSvc.Create(c, &data)
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

// GetStatus for get sync status at ces agent
func (h *AgentHandler) GetStatus(c *gin.Context) {
	data, err := h.AgentSvc.GetStatus(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"ok":      false,
			"code":    "internal-server-error",
			"errors":  []string{err.Error()},
			"message": "Error when trying to get agent status",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok":      data.Status != "error",
		"message": "successfully get agent status",
		"result":  data,
	})
}
