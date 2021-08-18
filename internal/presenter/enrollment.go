package presenter

import (
	"time"

	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/pkg/util"
)

// EnrollmentRequest is handler to processs from upload
type EnrollmentRequest struct {
	Images         []*ImageFile `json:"images"`
	Name           string       `json:"name"`
	IdentityNumber string       `json:"identity_number"`
	Status         string       `json:"status"`
	FaceID         string       `json:"face_id"`
}

// ImageFile is struct for handle request image from user
type ImageFile struct {
	Image []byte `json:"image"`
}

// EnrollmentResponse is handling response enrollment
type EnrollmentResponse struct {
	ID             uint64              `json:"id"`
	Name           string              `json:"name"`
	IdentityNumber string              `json:"identity_number"`
	Status         string              `json:"status"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      *time.Time          `json:"updated_at"`
	DeletedAt      *time.Time          `json:"deleted_at"`
	FaceID         uint64              `json:"face_id"`
	Faces          []*entity.FaceImage `json:"faces"`
}

// EnrollmentPaging is struct for handling list enrollment with pagination
type EnrollmentPaging struct {
	util.PaginationDetails
	Enrollments []*EnrollmentResponse `json:"enrollments"`
}
