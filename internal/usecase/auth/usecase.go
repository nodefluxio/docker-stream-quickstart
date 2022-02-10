package auth

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	logutil "gitlab.com/nodefluxio/goutils/pkg/log"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/pkg/util"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/presenter"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

const (
	accessTokenExpireDuration  = time.Minute * 30
	refreshTokenExpireDuration = time.Minute * 60
)

type ServiceImpl struct {
	UserRepo repository.User
}

func (s *ServiceImpl) Login(ctx context.Context, postData *presenter.LoginInput) (*presenter.LoginResponse, error) {
	// 1. get username and password default
	// 2. checking
	// 3. create claim
	// 4. return the data
	// const emailorusername = "admin@admin.com"
	// const passwordHash = "$2y$12$RWzFqxS7dCLluJuvxzMXleuiWQd7aehac4SIhwZtYp8.0z3XsbUCOe"
	var user *entity.User
	var usedUserAccess string

	isEmail := util.EmailValidator(postData.UserAccess)
	if isEmail {
		userData, err := s.UserRepo.GetByEmail(ctx, postData.UserAccess)
		if err != nil {
			return nil, err
		}
		user = userData
		usedUserAccess = userData.Email
	} else {
		userData, err := s.UserRepo.GetByUsername(ctx, postData.UserAccess)
		if err != nil {
			return nil, err
		}
		user = userData
		usedUserAccess = userData.Username
	}
	status := checkPasswordHash(user.Password, postData.Password)
	if status == false || postData.UserAccess != usedUserAccess {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err":              "status is false or user access is not same",
			"status":           status,
			"user_access":      postData.UserAccess,
			"email":            user.Email,
			"username":         user.Username,
			"used_user_access": usedUserAccess,
		},
			"checking status and user access")
		return nil, errors.New("invalid email or password")
	}

	accessToken := jwt.New(jwt.SigningMethodHS256)
	accessClaims := accessToken.Claims.(jwt.MapClaims)
	accessClaims["grant_type"] = "access_token"
	accessClaims["id"] = user.ID
	accessClaims["email"] = user.Email
	accessClaims["username"] = user.Username
	accessClaims["name"] = user.Fullname
	accessClaims["role"] = user.Role
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
		Name:         user.Fullname,
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

func (s *ServiceImpl) GetInfoAuthToken(ctx context.Context, token string) (*presenter.AuthInfoResponse, error) {
	var claims entity.ClaimsStruct
	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_TOKEN")), nil
	})
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"Error while parsing token, invalid token")
		return nil, errors.New("Error while parsing token, invalid token")
	}
	UserInfo, err := s.UserRepo.GetByEmail(ctx, claims.Email)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err":   err,
			"email": claims.Email,
		},
			"Error when validating email")
		return nil, errors.New("Invalid email or password")
	}
	if UserInfo.Role != claims.Role {
		err := errors.New("Invalid role")
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err":        err,
			"email":      claims.Email,
			"user_role":  UserInfo.Role,
			"claim_role": UserInfo.Role,
		},
			"Error when validating role, user role and claims role is different")
		return nil, err
	}
	response := &presenter.AuthInfoResponse{
		Name:   UserInfo.Fullname,
		Email:  claims.Email,
		UserID: claims.ID,
		SiteID: UserInfo.SiteID,
		Role:   UserInfo.Role,
	}
	return response, nil
}

// CheckPasswordHash check hashed password with unhashed data
func checkPasswordHash(hashPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	return err == nil
}
