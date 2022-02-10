package entity

import "time"

// Plate represent attributes on enrolled plate
type Vehicle struct {
	ID        uint64     `json:"id" db:"id"`
	Plate     string     `json:"plate_number" db:"plate_number" gorm:"column:plate_number"`
	Type      string     `json:"type" db:"type"`
	Brand     string     `json:"brand" db:"brand"`
	Color     string     `json:"color" db:"color"`
	Name      string     `json:"name" db:"name"`
	UniqueID  string     `json:"unique_id" db:"unique_id"`
	Status    string     `json:"status" db:"status"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}
