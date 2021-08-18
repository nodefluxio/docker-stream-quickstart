package entity

import (
	"time"
)

// Event represent attributes on event data
type GlobalSetting struct {
	ID             uint64    `json:"id" db:"id"`
	Similarity 	   float64 `json:"similarity" db:"similarity"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}