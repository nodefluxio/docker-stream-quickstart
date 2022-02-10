package httphandler

import (
	"bytes"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/usecase/polrisearching"
)

// PolriSearchingHandler holds app interactor to be used by handler function
type PolriSearchingHandler struct {
	PolriSearchingSvc polrisearching.Service
}

// GetPlate for get detail data plate number at database polri
func (h *PolriSearchingHandler) GetPlate(c *gin.Context) {
	nopol, _ := c.GetQuery("nopol")
	if nopol == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"ok":      false,
			"code":    "bad-request",
			"errors":  []string{"parameter nopol is required"},
			"message": "error get parameter",
		})
		return
	}
	data, err := h.PolriSearchingSvc.SearchPlateNumber(c, nopol)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"ok":      false,
			"code":    "bad-request",
			"errors":  []string{err.Error()},
			"message": "Error when trying get data plate number",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ok":         true,
		"message":    "succesfully get data plate number",
		"plate_info": data,
	})
}

// GetNik for get detail data nik
func (h *PolriSearchingHandler) GetNik(c *gin.Context) {
	nik, _ := c.GetQuery("nik")
	if nik == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"ok":      false,
			"code":    "bad-request",
			"errors":  []string{"parameter nik is required"},
			"message": "error get parameter",
		})
		return
	}
	data, err := h.PolriSearchingSvc.SearchNik(c, nik)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"ok":      false,
			"code":    "bad-request",
			"errors":  []string{err.Error()},
			"message": "Error when trying get data nik",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ok":           true,
		"message":      "succesfully get data nik",
		"citizen_data": data,
	})
}

func (h *PolriSearchingHandler) GetFaceSearchToken(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"ok":      false,
			"code":    "bad-request",
			"errors":  []string{err.Error()},
			"message": "data is invalid",
		})
		return
	}
	files := form.File["image"]
	if len(files) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"ok":      false,
			"code":    "bad-request",
			"errors":  []string{err.Error()},
			"message": "field image is Required",
		})
		return
	}
	var image []byte
	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"ok":      false,
				"code":    "bad-request",
				"errors":  []string{err.Error()},
				"message": "Error when open file",
			})
			return
		}
		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, src); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"ok":      false,
				"code":    "bad-request",
				"errors":  []string{err.Error()},
				"message": "Error when trying to get face search token",
			})
			return
		}
		image = buf.Bytes()
		limit, _ := strconv.ParseUint(c.Request.PostFormValue("limit"), 10, 64)
		token, errSvc := h.PolriSearchingSvc.GetFaceSearchToken(c, image, limit)
		if errSvc != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"ok":      false,
				"code":    "bad-request",
				"errors":  []string{errSvc.Error()},
				"message": "Error when trying to get face search token",
			})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"ok":      true,
			"message": "Successfully get face search token",
			"token":   token,
		})
	}
}

// GetNik for get detail data nik
func (h *PolriSearchingHandler) GetFaceSearchResult(c *gin.Context) {
	token, _ := c.GetQuery("token")
	if token == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"ok":      false,
			"code":    "bad-request",
			"errors":  []string{"parameter token is requred"},
			"message": "error get parameter",
		})
		return
	}
	data, err := h.PolriSearchingSvc.GetFaceSearchResult(c, token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"ok":      false,
			"code":    "bad-request",
			"errors":  []string{err.Error()},
			"message": "Error when trying face result",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ok":      true,
		"message": "succesfully get data",
		"results": data,
	})
}
