package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/pkg/util"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/presenter"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/usecase/enrollment"
)

// EnrollmentHandler holds app interactor to be used by handler function
type EnrollmentHandler struct {
	EnrollmentSvc enrollment.Service
}

//GetList for get all data enrollment
func (h *EnrollmentHandler) GetList(c *gin.Context) {
	var paging util.Pagination
	paging.Limit, _ = strconv.Atoi(c.Query("limit"))
	paging.Sort = c.QueryMap("sort")
	paging.Page, _ = strconv.Atoi(c.Query("page"))
	paging.Filter = c.QueryMap("filter")
	paging.Search = c.Query("search")

	events, err := h.EnrollmentSvc.GetList(c, &paging)
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
		"message": "succesfully get enrolled person",
		"results": events,
	})
}

//GetDetail for get detail data enrollment
func (h *EnrollmentHandler) GetDetail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	detail, err := h.EnrollmentSvc.GetDetail(c, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"ok":      false,
			"code":    "unprocessable-entity",
			"errors":  []string{err.Error()},
			"message": "Error when trying get detail enrollment",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ok":         true,
		"message":    "succesfully get detail enrolled faces",
		"enrollment": detail,
	})
}

//Backup for get detail data enrollment
func (h *EnrollmentHandler) Backup(c *gin.Context) {
	backupFile, err := h.EnrollmentSvc.Backup(c)
	csvFilePath := fmt.Sprintf("%s/%s", filepath.Dir(backupFile.Name()), backupFile.Name())
	file, err := os.Open(csvFilePath) //Create a file
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"ok":      false,
			"code":    "not-found",
			"errors":  []string{err.Error()},
			"message": "Error when trying to backup enrollment data",
		})
		return
	}
	defer file.Close()
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.Writer.Header().Add("Content-Disposition", "attachment; filename="+backupFile.Name())
	_, err = io.Copy(c.Writer, file)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"ok":      false,
			"code":    "not-found",
			"errors":  []string{err.Error()},
			"message": "Error when trying to backup enrollment data",
		})
		return
	}
}

// Create for add enrollment
func (h *EnrollmentHandler) Create(c *gin.Context) {
	name := c.Request.PostFormValue("name")
	identityNumber := c.Request.PostFormValue("identity_number")
	status := c.Request.PostFormValue("status")
	faceID := c.Request.PostFormValue("face_id")
	gender := c.Request.PostFormValue("gender")
	birthPlace := c.Request.PostFormValue("birth_place")
	birthDate := c.Request.PostFormValue("birth_date")
	isAgent := c.GetHeader("Is-Agent")

	form, _ := c.MultipartForm()
	files := form.File["images"]
	if len(files) == 0 {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"ok":      false,
			"code":    "unprocessable-entity",
			"errors":  []string{},
			"message": "field images is Required",
		})
		return
	}
	var nImages []*presenter.ImageFile
	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"ok":      false,
				"code":    "unprocessable-entity",
				"errors":  []string{err.Error()},
				"message": "Error when open file",
			})
			return
		}
		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, src); err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"ok":      false,
				"code":    "unprocessable-entity",
				"errors":  []string{err.Error()},
				"message": "Error when trying to enrollment",
			})
			return
		}
		nImages = append(nImages, &presenter.ImageFile{
			Image: buf.Bytes(),
		})

	}
	data := presenter.EnrollmentRequest{
		Images:         nImages,
		Name:           name,
		IdentityNumber: identityNumber,
		Gender:         gender,
		BirthPlace:     birthPlace,
		BirthDate:      birthDate,
		Status:         status,
		FaceID:         faceID,
	}

	result, errSvc := h.EnrollmentSvc.Create(c, &data, isAgent)
	if errSvc != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"ok":      false,
			"code":    "unprocessable-entity",
			"errors":  []string{errSvc.Error()},
			"message": "Error when trying to enrollment",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"ok":         true,
		"message":    "Successfully Enroll New Face",
		"enrollment": result,
	})
}

// Delete for delete enrollment
func (h *EnrollmentHandler) Delete(c *gin.Context) {
	// checking if id is as string
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"ok":      false,
			"code":    "bad-request",
			"errors":  []string{err.Error()},
			"message": "Error when parsing id",
		})
		return
	}
	isAgent := c.GetHeader("Is-Agent")
	err = h.EnrollmentSvc.Delete(c, id, isAgent)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"ok":      false,
			"code":    "bad-request",
			"errors":  []string{err.Error()},
			"message": "Error when trying to delete enrollment",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok":      true,
		"message": "succesfully delete enrolled person",
	})
}

// DeleteAll for delete all enrollment
func (h *EnrollmentHandler) DeleteAll(c *gin.Context) {
	// checking if id is as string

	err := h.EnrollmentSvc.DeleteAll(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"ok":      false,
			"code":    "bad-request",
			"errors":  []string{err.Error()},
			"message": "Error when trying to delete enrollment",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok":      true,
		"message": "succesfully delete enrolled person",
	})
}

// Update for update enrollment
func (h *EnrollmentHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	name := c.Request.PostFormValue("name")
	identityNumber := c.Request.PostFormValue("identity_number")
	status := c.Request.PostFormValue("status")
	gender := c.Request.PostFormValue("gender")
	birthPlace := c.Request.PostFormValue("birth_place")
	birthDate := c.Request.PostFormValue("birth_date")
	deletedVariations, _ := c.GetPostFormArray("deleted_variations")
	form, _ := c.MultipartForm()
	files := form.File["images"]
	isAgent := c.GetHeader("Is-Agent")

	var nImages []*presenter.ImageFile
	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"ok":      false,
				"code":    "unprocessable-entity",
				"errors":  []string{err.Error()},
				"message": "Error when open file",
			})
			return
		}
		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, src); err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"ok":      false,
				"code":    "unprocessable-entity",
				"errors":  []string{err.Error()},
				"message": "Error when trying copy file for buffer",
			})
			return
		}
		nImages = append(nImages, &presenter.ImageFile{
			Image: buf.Bytes(),
		})

	}
	data := presenter.EnrollmentRequest{
		Images:            nImages,
		Name:              name,
		IdentityNumber:    identityNumber,
		Gender:            gender,
		BirthPlace:        birthPlace,
		BirthDate:         birthDate,
		Status:            status,
		DeletedVariations: deletedVariations,
	}
	errSvc := h.EnrollmentSvc.Update(c, id, &data, isAgent)
	if errSvc != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"ok":      false,
			"code":    "unprocessable-entity",
			"errors":  []string{errSvc.Error()},
			"message": "Error when trying to update enrollment",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok":         true,
		"message":    "succesfully update enrolled person",
		"enrollment": map[string]interface{}{"name": name},
	})
}

// GetImage is handler for show image face
func (h *EnrollmentHandler) GetImage(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	img, _ := h.EnrollmentSvc.GetFaceImage(c, id)
	c.Header("Content-Type", "image/jpeg")
	c.Data(http.StatusOK, "image/jpeg", img)
}
