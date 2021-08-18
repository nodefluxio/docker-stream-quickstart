package entity

import "time"

// EnrolledFace represent attributes on enrolled face
type EnrolledFace struct {
	ID             uint64     `json:"id" db:"id"`
	FaceID         uint64     `json:"face_id" db:"face_id"`
	Name           string     `json:"name" db:"name"`
	IdentityNumber string     `json:"identity_number" db:"identity_number"`
	Status         string     `json:"status" db:"status"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at" db:"deleted_at"`
}

// EnrolledFaceWithImage represent attributes on enrolled face with image
type EnrolledFaceWithImage struct {
	ID             uint64     `json:"id" db:"id"`
	FaceID         uint64     `json:"face_id" db:"face_id"`
	Name           string     `json:"name" db:"name"`
	IdentityNumber string     `json:"identity_number" db:"identity_number"`
	Status         string     `json:"status" db:"status"`
	Image          []byte     `json:"image" db:"image"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at" db:"deleted_at"`
}
