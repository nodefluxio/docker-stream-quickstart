package presenter

import (
	"os"
	"time"

	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/pkg/util"
)

// EnrollmentRequest is handler to processs from upload
type EnrollmentRequest struct {
	Images            []*ImageFile `json:"images"`
	Name              string       `json:"name"`
	IdentityNumber    string       `json:"identity_number"`
	Gender            string       `json:"gender" db:"gender"`
	BirthPlace        string       `json:"birth_place" db:"birth_place"`
	BirthDate         string       `json:"birth_date" db:"birth_date"`
	Status            string       `json:"status"`
	FaceID            string       `json:"face_id"`
	DeletedVariations []string     `json:"deleted_variations"`
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
	Gender         string              `json:"gender" db:"gender"`
	BirthPlace     string              `json:"birth_place" db:"birth_place"`
	BirthDate      string              `json:"birth_date" db:"birth_date"`
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

// EnrollmentBackupFiles is a response for backup enrollment
type EnrollmentBackupFiles struct {
	CSVFile        *os.File
	FaceImagesFile []*os.File
}
