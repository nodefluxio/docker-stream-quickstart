package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	logutil "gitlab.com/nodefluxio/goutils/pkg/log"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/presenter"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/usecase/auth"
)

// EventHandler holds app interactor to be used by handler function
type AuthHandler struct {
	AuthSvc auth.Service
}

// Create for add enrollment
func (h *AuthHandler) Login(c *gin.Context) {
	data := presenter.LoginInput{}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"should bind json is failed")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"ok":      false,
			"errors":  []string{err.Error()},
			"message": "Invalid email or password",
		})
		return
	}
	result, err := h.AuthSvc.Login(c, &data)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"failed in return of usecase login")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"ok":      false,
			"errors":  []string{err.Error()},
			"message": "Invalid email or password",
		})
		return
	}
	logutil.LogObj.SetInfoLog(map[string]interface{}{
		"info": result,
	},
		"Success Login")
	c.JSON(http.StatusCreated, gin.H{
		"ok":            true,
		"message":       "Successfull Login",
		"access_token":  result.AccessToken,
		"refresh_token": result.RefreshToken,
		"name":          result.Name,
	})
	return
}

// will use at next iteration
// //RefreshToken handler user to refresh token if access token is expired
// func (h *AuthHandler) RefreshToken(c *gin.Context) {
// 	claims, err := parseToClaimsJWT(c.Request.Header.Get("Authorization"))

// 	resultClaims, err := h.AuthSvc.RefreshToken()
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
// 			"ok":      false,
// 			"code":    "fatal-error",
// 			"message": err})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{
// 		"ok":            true,
// 		"access_token":  resultClaims.AccessToken,
// 		"refresh_token": resultClaims.RefreshToken,
// 		"user":          resultClaims.Object})
// }
// func parseToClaimsJWT(authorizationHeader string) (*ClaimsStruct, error) {
// 	if strings.Contains(authorizationHeader, JWTSpecialKey) == true {
// 		tokenString := strings.Replace(authorizationHeader, JWTSpecialKey, "", -1)
// 		token, err := jwt.ParseWithClaims(tokenString, &ClaimsStruct{}, func(token *jwt.Token) (interface{}, error) {
// 			return verifyKey, nil
// 		})
// 		if err != nil {
// 			return nil, errors.New("Error when parse claims")

// 		} else if !token.Valid {
// 			return nil, errors.New("token not valid")

// 		}
// 		fields := make(map[string]interface{})
// 		fields["jwt_token"] = token
// 		logutil.LogObj.SetWarnLog(fields, "[Debugging] JWT Token")
// 		claims, _ := token.Claims.(*ClaimsStruct)
// 		return claims, nil
// 	}
// 	return nil, errors.New("No Contain JWT Special Keys")
// }
