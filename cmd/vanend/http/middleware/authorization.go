package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	logutil "gitlab.com/nodefluxio/goutils/pkg/log"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/presenter"
)

type AuthorizeMiddleware struct{}

// CheckSite middleware for checking this user have access to this site id or not
func (m *AuthorizeMiddleware) CheckSite(c *gin.Context) {
	accInfo := c.MustGet("account_info").(*presenter.AuthInfoResponse)
	paramSiteID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err":          err,
			"parameter_id": c.Param("id"),
		},
			"error when parsing data")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"ok":      false,
			"code":    "bad-request",
			"errors":  []string{err.Error()},
			"message": "failed parsing, value data must be number",
		})
		return
	}

	logutil.LogObj.SetInfoLog(map[string]interface{}{
		"account_info":      accInfo,
		"parameter_site_id": paramSiteID,
	},
		"error when get list site with stream")
	if accInfo.Role == string(entity.UserRoleOperator) {
		found := false
		for _, v := range accInfo.SiteID {
			if v == paramSiteID {
				found = true
			}
		}
		if !found {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"ok":      false,
				"code":    "forbidden",
				"errors":  []string{},
				"message": "not allowed for access this resource",
			})
		}
		return
	}
	c.Next()
}

// CheckAccess middleware for checking this user role have access to the resource
func (m *AuthorizeMiddleware) CheckAccess(c *gin.Context) {
	accInfo := c.MustGet("account_info").(*presenter.AuthInfoResponse)
	if accInfo.Role == string(entity.UserRoleOperator) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"ok":      false,
			"code":    "forbidden",
			"errors":  []string{},
			"message": "not allowed for access this resource",
		})
		return
	}
	c.Next()
}
