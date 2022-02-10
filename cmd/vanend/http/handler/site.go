package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/pkg/util"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/presenter"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/usecase/site"
)

// SiteHandler holds app interactor to be used by handler function
type SiteHandler struct {
	SiteSvc site.Service
}

//GetList for get all data enrollment
func (h *SiteHandler) GetList(c *gin.Context) {
	var paging util.Pagination
	paging.Sort = c.QueryMap("sort")
	paging.Filter = c.QueryMap("filter")
	paging.Search = c.Query("search")
	accInfo := c.MustGet("account_info").(*presenter.AuthInfoResponse)

	sites, err := h.SiteSvc.GetList(c, &paging, accInfo)
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
		"message": "successfully get all site",
		"sites":   sites,
	})
}

// Create for create site data
func (h *SiteHandler) Create(c *gin.Context) {
	data := presenter.SiteRequest{}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"ok":      false,
			"code":    "bad-request",
			"errors":  []string{err.Error()},
			"message": "Error when trying to parse to presenter",
		})
		return
	}
	dataSite, err := h.SiteSvc.Create(c, &data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"ok":      false,
			"code":    "bad-request",
			"errors":  []string{err.Error()},
			"message": "Error when trying create site",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ok":      true,
		"message": "succesfully create site",
		"site":    dataSite,
	})
}

// Update is for update site data
func (h *SiteHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	data := presenter.SiteRequest{}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"ok":      false,
			"code":    "bad-request",
			"errors":  []string{err.Error()},
			"message": "Error when trying to parse to presenter",
		})
		return
	}
	data.ID = id
	err = h.SiteSvc.Update(c, &data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"ok":      false,
			"code":    "bad-request",
			"errors":  []string{err.Error()},
			"message": "Error when trying update site",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ok":      true,
		"message": "succesfully update site",
	})
}

// Delete is for delete site data
func (h *SiteHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	err := h.SiteSvc.Delete(c, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"ok":      false,
			"code":    "bad-request",
			"errors":  []string{err.Error()},
			"message": "Error when trying delete site",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ok":      true,
		"message": "succesfully delete site",
	})
}

// AssignStream for assign stream id to site
func (h *SiteHandler) AssignStream(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	data := presenter.AssignStreamRequest{}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"ok":      false,
			"code":    "bad-request",
			"errors":  []string{err.Error()},
			"message": "Error when trying to parse to presenter",
		})
		return
	}
	data.SiteID = id
	err = h.SiteSvc.AssignStreamToSite(c, &data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"ok":      false,
			"code":    "bad-request",
			"errors":  []string{err.Error()},
			"message": "Error when trying create site",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ok":      true,
		"message": "successfully assign stream site",
	})
}
