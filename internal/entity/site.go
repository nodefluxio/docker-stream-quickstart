package entity

import "time"

type Site struct {
	ID        uint64     `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

type MapSiteStream struct {
	ID        uint64    `json:"id" db:"id"`
	SiteID    uint64    `json:"site_id" db:"site_id"`
	StreamID  string    `json:"stream_id" db:"stream_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type SiteWithStream struct {
	ID        uint64     `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	StreamID  string     `json:"stream_id" db:"sream_id"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}
