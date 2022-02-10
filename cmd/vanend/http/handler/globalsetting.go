package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/presenter"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/usecase/globalsetting"
)

// GlobalSettingHandler holds app interactor to be used by handler function
type GlobalSettingHandler struct {
	GlobalSettingSvc globalsetting.Service
}

//GlobalSettingHandler for get detail data enrollment
func (h *GlobalSettingHandler) GetDetail(c *gin.Context) {

	detail, err := h.GlobalSettingSvc.GetDetail(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"ok":      false,
			"code":    "unprocessable-entity",
			"errors":  []string{err.Error()},
			"message": "Error when trying get detail setting",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ok":      true,
		"message": "succesfully get global setting",
		"setting": detail,
	})
}

// CreateOrUpdate for change globalsetting
func (h *GlobalSettingHandler) CreateOrUpdate(c *gin.Context) {
	data := presenter.GlobalSettingRequest{}
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
	result, errSvc := h.GlobalSettingSvc.CreateOrUpdate(c, &data)
	if errSvc != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"ok":      false,
			"code":    "unprocessable-entity",
			"errors":  []string{errSvc.Error()},
			"message": "Error when trying to create",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"ok":             true,
		"message":        "Successfully Create or Update result",
		"global_setting": result,
	})
	return
}
