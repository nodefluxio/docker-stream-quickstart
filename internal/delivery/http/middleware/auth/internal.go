package httpauth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	logutil "gitlab.com/nodefluxio/goutils/pkg/log"
)

// Check make sure every request is authorized
func Check(c *gin.Context) {
	logutil.LogObj.SetDebugLog(map[string]interface{}{"header value": c.Request.Header}, "Header payment service")
	if c.GetHeader("x-nodeflux-can-access-internal") != "true" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"ok":      false,
			"code":    "unauthorized-request",
			"message": "You are not authorized to do this.",
			"errors":  []string{},
		})
		return
	}
	c.Next()
}
