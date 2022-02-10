package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/usecase/auth"
)

type AuthMiddleware struct {
	AuthSvc auth.Service
}

// IsLoggedIn middleware for checking token user is valid or not
func (m *AuthMiddleware) IsLoggedIn(c *gin.Context) {
	requestHeaders := make(map[string]string)
	for name, headers := range c.Request.Header {
		for _, h := range headers {
			requestHeaders[strings.ToLower(name)] = h
		}
	}

	headerAuth := requestHeaders["authorization"]
	splitToken := strings.Split(headerAuth, "Bearer ")
	if len(splitToken) != 2 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"ok":      false,
			"code":    "not-authenticate",
			"errors":  []string{},
			"message": "Bearer token not in proper format",
		})
		return
	}

	token := splitToken[1]
	accInfo, err := m.AuthSvc.GetInfoAuthToken(c, token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"ok":      false,
			"code":    "not-authenticate",
			"errors":  []string{err.Error()},
			"message": "Failed verify token",
		})
		return
	}
	c.Set("account_info", accInfo)
	c.Next()
}
