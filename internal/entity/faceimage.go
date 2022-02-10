package entity

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"time"
)

// FaceImage represent list of image use in enrolled face
type FaceImage struct {
	ID             uint64     `json:"id" db:"id"`
	EnrolledFaceID uint64     `json:"enrolled_face_id,omitempty" db:"enrolled_face_id"`
	Variation      string     `json:"variation" db:"variation"`
	Image          []byte     `json:"image,omitempty" db:"image"`
	ImageThumbnail []byte     `json:"image_thumbnail,omitempty" db:"image_thumbnail"`
	CreatedAt      time.Time  `json:"created_at,omitempty" db:"created_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

func (fi *FaceImage) WriteToJpeg(ctx context.Context, fileName string) (*os.File, error) {
	out, _ := os.Create(fmt.Sprintf("%s.jpg", fileName))
	defer out.Close()

	image, _, err := image.Decode(bytes.NewReader(fi.Image))
	if err != nil {
		return nil, err
	}

	err = jpeg.Encode(out, image, nil)
	if err != nil {
		return nil, err
	}

	return out, nil
}

type EventEnrollmentFaceImage struct {
	ID                uint64    `json:"id" db:"id"`
	EventEnrollmentID uint64    `json:"event_enrollment_id,omitempty" db:"event_enrollment_id"`
	Image             []byte    `json:"image,omitempty" db:"image"`
	CreatedAt         time.Time `json:"created_at,omitempty" db:"created_at"`
}

func (EventEnrollmentFaceImage) TableName() string {
	return "face_image"
}
