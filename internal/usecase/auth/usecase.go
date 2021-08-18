package auth

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	logutil "gitlab.com/nodefluxio/goutils/pkg/log"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/pkg/util"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/presenter"
	"golang.org/x/crypto/bcrypt"
)

const (
	accessTokenExpireDuration  = time.Minute * 30
	refreshTokenExpireDuration = time.Minute * 60
)

type ServiceImpl struct {
}

func (s *ServiceImpl) Login(ctx context.Context, postData *presenter.LoginInput) (*presenter.LoginResponse, error) {
	// 1. get username and password default
	// 2. checking
	// 3. create claim
	// 4. return the data
	const emailorusername = "admin@admin.com"
	const passwordHash = "$2y$12$RWzFqxS7dCLluJuvxzMXleuiWQd7aehac4SIhwZtYp8.0z3XsbUCOe"
	status := checkPasswordHash(passwordHash, postData.Password)
	if status == false || postData.UserAccess != emailorusername {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err":    "status is false or user access is not same",
			"status": status,
			"email":  postData.UserAccess},
			"checking status and user access")
		return nil, errors.New("invalid email or password")
	}

	accessToken := jwt.New(jwt.SigningMethodHS256)
	accessClaims := accessToken.Claims.(jwt.MapClaims)
	name := "admin"
	accessClaims["grant_type"] = "access_token"
	accessClaims["email"] = "admin@admin.com"
	accessClaims["name"] = name
	accessClaims["exp"] = time.Now().Add(accessTokenExpireDuration).Unix()

	signedAccessToken, err := accessToken.SignedString([]byte(os.Getenv("SECRET_TOKEN")))
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"failed to signed access token")
		return nil, err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshClaims["grant_type"] = "refresh_token"
	refreshClaims["exp"] = time.Now().Add(refreshTokenExpireDuration).Unix()

	signedRefreshToken, err := refreshToken.SignedString([]byte(util.ReverseString(os.Getenv("SECRET_TOKEN"))))
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"failed to signedRefreshToken")
		return nil, err
	}

	loginResponse := &presenter.LoginResponse{
		AccessToken:  signedAccessToken,
		RefreshToken: signedRefreshToken,
		Name:         name,
	}

	return loginResponse, nil
}

func (s *ServiceImpl) RefreshToken(ctx context.Context, postData *presenter.LoginResponse) (*presenter.LoginResponse, error) {
	var Refresh *presenter.LoginResponse

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(postData.RefreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_TOKEN")), nil
	})
	if err != nil && len(claims) == 0 {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"failed to parse to claims")
		return nil, errors.New("Error while parsing token, invalid token")
	}

	accessToken := jwt.New(jwt.SigningMethodHS256)
	acessClaims := accessToken.Claims.(jwt.MapClaims)
	acessClaims["grant_type"] = "access_token"
	acessClaims["id"] = claims["id"]
	acessClaims["exp"] = time.Now().Add(accessTokenExpireDuration).Unix()
	acessClaims["key"] = claims["key"]

	signedaccessToken, err := accessToken.SignedString([]byte(os.Getenv("SECRET_TOKEN")))
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"failed to create access token")
		return nil, err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshClaims["grant_type"] = "refresh_token"
	refreshClaims["exp"] = time.Now().Add(refreshTokenExpireDuration).Unix()

	signedrefreshToken, err := refreshToken.SignedString([]byte(util.ReverseString(os.Getenv("SECRET_TOKEN"))))
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"failed to create refresh token")

		return nil, err
	}

	Refresh.AccessToken = signedaccessToken
	Refresh.RefreshToken = signedrefreshToken

	return Refresh, nil
}

// CheckPasswordHash check hashed password with unhashed data
func checkPasswordHash(hashPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	return err == nil
}
