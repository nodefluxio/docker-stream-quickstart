package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	validatorHelper "gitlab.com/nodefluxio/vanilla-dashboard/internal/delivery/http/middleware/validator/helper"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/pkg/util"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/presenter"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/usecase/vehicle"
	"gopkg.in/go-playground/validator.v9"
)

// vehicleHandler holds app interactor to be used by handler function
type VehicleHandler struct {
	VehicleSvc vehicle.Service
	Validator  *validator.Validate
	Translator ut.Translator
}

//GetList for get all data vehicle
func (h *VehicleHandler) GetList(c *gin.Context) {
	var paging util.Pagination
	paging.Limit, _ = strconv.Atoi(c.Query("limit"))
	paging.Sort = c.QueryMap("sort")
	paging.Page, _ = strconv.Atoi(c.Query("page"))
	paging.Filter = c.QueryMap("filter")
	paging.Search = c.Query("search")

	list, err := h.VehicleSvc.GetList(c, &paging)
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
		"message": "succesfully get vehicles",
		"result":  list,
	})
}

//GetDetail for get detail data vehicle
func (h *VehicleHandler) GetDetail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	detail, err := h.VehicleSvc.GetDetail(c, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"ok":      false,
			"code":    "unprocessable-entity",
			"errors":  []string{err.Error()},
			"message": "Error when trying get detail vehicle",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ok":      true,
		"message": "succesfully get detail vehicles",
		"vehicle": detail,
	})
}

// Create for add vehicle
func (h *VehicleHandler) Create(c *gin.Context) {
	postData := presenter.VehicleEnrollmentRequest{}
	err := c.ShouldBindJSON(&postData)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"ok":      false,
			"code":    "unprocessable-entity",
			"errors":  []string{err.Error()},
			"message": "Error when trying to parse to presenter",
		})
		return
	}

	err = h.Validator.Struct(postData)
	fmt.Println(err)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"ok":      false,
			"code":    "unprocessable-entity",
			"errors":  validatorHelper.ErrorMessageTranslator(err, h.Translator),
			"message": "Error when trying to add new vehicle",
		})
		return
	}
	result, errSvc := h.VehicleSvc.Create(c, &postData)
	if errSvc != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"ok":      false,
			"code":    "unprocessable-entity",
			"errors":  []string{errSvc.Error()},
			"message": "Error when trying to add new vehicle",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"ok":      true,
		"message": "Successfully register new vehicle",
		"vehicle": result,
	})
	return
}

// Delete for delete tier
func (h *VehicleHandler) Delete(c *gin.Context) {
	// checking if id is as string

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {

	}
	err = h.VehicleSvc.Delete(c, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"ok":      false,
			"code":    "bad-request",
			"errors":  []string{err.Error()},
			"message": "Error when trying to delete vehicle",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok":      true,
		"message": "succesfully delete vehicle",
	})
}

// DeleteAll for delete all vehicle
func (h *VehicleHandler) DeleteAll(c *gin.Context) {
	// checking if id is as string

	err := h.VehicleSvc.DeleteAll(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"ok":      false,
			"code":    "bad-request",
			"errors":  []string{err.Error()},
			"message": "Error when trying to delete vehicle",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok":      true,
		"message": "succesfully delete vehicle",
	})
}

// Update for update vehicle
func (h *VehicleHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	data := presenter.VehicleEnrollmentRequest{}
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
	err = h.Validator.Struct(data)
	fmt.Println(err)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"ok":      false,
			"code":    "unprocessable-entity",
			"errors":  validatorHelper.ErrorMessageTranslator(err, h.Translator),
			"message": "Error when trying to add new vehicle",
		})
		return
	}
	errSvc := h.VehicleSvc.Update(c, id, &data)
	if errSvc != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"ok":      false,
			"code":    "unprocessable-entity",
			"errors":  []string{errSvc.Error()},
			"message": "Error when trying to update vehicle",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok":      true,
		"message": "succesfully update plate",
		"vehicle": map[string]interface{}{"plate": data.Plate},
	})
}
