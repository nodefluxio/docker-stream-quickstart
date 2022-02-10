package entity

import (
	"bytes"
	"time"
)

type VanillaEnrollmentPayload struct {
	FaceID         uint64    `json:"face_id"`
	Name           string    `json:"name"`
	IdentityNumber string    `json:"identity_number"`
	Gender         string    `json:"gender"`
	BirthPlace     string    `json:"birth_place"`
	BirthDate      time.Time `json:"birth_date"`
	Status         string    `json:"status"`
}

type EnrollmentImage struct {
	Image *bytes.Buffer `json:"image"`
}

type VanillaEnrollmentData struct {
	Message string `json:"message"`
	Ok      bool   `json:"ok"`
	Results struct {
		Limit       int `json:"limit"`
		CurrentPage int `json:"current_page"`
		TotalData   int `json:"total_data"`
		TotalPage   int `json:"total_page"`
		Enrollments []struct {
			ID             uint64      `json:"id"`
			Name           string      `json:"name"`
			IdentityNumber string      `json:"identity_number"`
			Gender         string      `json:"gender"`
			BirthPlace     string      `json:"birth_place"`
			BirthDate      time.Time   `json:"birth_date"`
			Status         string      `json:"status"`
			CreatedAt      time.Time   `json:"created_at"`
			UpdatedAt      time.Time   `json:"updated_at"`
			DeletedAt      interface{} `json:"deleted_at"`
			FaceID         int64       `json:"face_id"`
			Faces          []struct {
				ID             int       `json:"id"`
				EnrolledFaceID int       `json:"enrolled_face_id"`
				Variation      string    `json:"variation"`
				CreatedAt      time.Time `json:"created_at"`
			} `json:"faces"`
		} `json:"enrollments"`
	} `json:"results"`
}
