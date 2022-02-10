package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/pkg/util"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/presenter"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/usecase/user"
)

// UserHandler holds app interactor to be used by handler function
type UserHandler struct {
	UserSvc user.Service
}

// Create for insert new user
func (h *UserHandler) Create(c *gin.Context) {
	data := presenter.UserRequest{}
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
	dataUser, err := h.UserSvc.Create(c, &data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"ok":      false,
			"code":    "bad-request",
			"errors":  []string{err.Error()},
			"message": "Error when trying create user",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ok":      true,
		"message": "Successfully add new user",
		"site":    dataUser,
	})
}

// Update for Update existing user
func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"ok":      false,
			"code":    "bad-request",
			"errors":  []string{err.Error()},
			"message": "Error when trying to parse parameter",
		})
		return
	}
	data := presenter.UserRequest{
		ID: id,
	}
	err = c.ShouldBindJSON(&data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"ok":      false,
			"code":    "bad-request",
			"errors":  []string{err.Error()},
			"message": "Error when trying to parse to presenter",
		})
		return
	}
	err = h.UserSvc.Update(c, &data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"ok":      false,
			"code":    "bad-request",
			"errors":  []string{err.Error()},
			"message": "Error when trying update user",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ok":      true,
		"message": "Successfully update user data",
	})
}

// ChangePassword for update password
func (h *UserHandler) ChangePassword(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"ok":      false,
			"code":    "bad-request",
			"errors":  []string{err.Error()},
			"message": "Error when trying to parse parameter",
		})
		return
	}
	data := presenter.UserChangePassRequest{
		ID: id,
	}
	err = c.ShouldBindJSON(&data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"ok":      false,
			"code":    "bad-request",
			"errors":  []string{err.Error()},
			"message": "Error when trying to parse to presenter",
		})
		return
	}
	err = h.UserSvc.ChangePassword(c, &data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"ok":      false,
			"code":    "bad-request",
			"errors":  []string{err.Error()},
			"message": "Error when trying change password user",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ok":      true,
		"message": "Successfully change password user",
	})
}

// Delete for delete user data
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"ok":      false,
			"code":    "bad-request",
			"errors":  []string{err.Error()},
			"message": "Error when trying to parse parameter",
		})
		return
	}
	err = h.UserSvc.Delete(c, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"ok":      false,
			"code":    "bad-request",
			"errors":  []string{err.Error()},
			"message": "Error when trying delete user",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ok":      true,
		"message": "Successfully delete user data",
	})
}

// Detail for get detail user data
func (h *UserHandler) Detail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"ok":      false,
			"code":    "bad-request",
			"errors":  []string{err.Error()},
			"message": "Error when trying to parse parameter",
		})
		return
	}
	userData, err := h.UserSvc.GetDetail(c, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"ok":      false,
			"code":    "bad-request",
			"errors":  []string{err.Error()},
			"message": "Error when trying delete user",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ok":      true,
		"message": "Successfully get detail user",
		"user":    userData,
	})
}

//GetList for get all data user
func (h *UserHandler) GetList(c *gin.Context) {
	var paging util.Pagination
	paging.Limit, _ = strconv.Atoi(c.Query("limit"))
	paging.Sort = c.QueryMap("sort")
	paging.Page, _ = strconv.Atoi(c.Query("page"))
	paging.Filter = c.QueryMap("filter")
	paging.Search = c.Query("search")

	events, err := h.UserSvc.GetList(c, &paging)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"ok":      false,
			"code":    "bad-request",
			"errors":  []string{},
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ok":      true,
		"message": "Successfully get list of user",
		"results": events,
	})
}
