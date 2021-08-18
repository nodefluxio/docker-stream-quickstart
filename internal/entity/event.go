package entity

import (
	"time"
)

// Event represent attributes on event data
type Event struct {
	ID             uint64    `json:"id" db:"id"`
	EventType      string    `json:"event_type" db:"type" gorm:"column:type"`
	StreamID       string    `json:"stream_id" db:"stream_id"`
	Detection      Message   `json:"detection" db:"detection"`
	PrimaryImage   []byte    `json:"primary_image" db:"primary_image"`
	SecondaryImage []byte    `json:"secondary_image" db:"secondary_image"`
	Result         []byte    `json:"result" db:"result"`
	Status         string    `json:"status" db:"status"`
	EventTime      time.Time `json:"event_time" db:"event_time"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// EventWithoutImage represent attributes on event data
type EventWithoutImage struct {
	ID        uint64    `json:"id" db:"id"`
	EventType string    `json:"event_type" db:"type" gorm:"column:type"`
	StreamID  string    `json:"stream_id" db:"stream_id"`
	Detection []byte    `json:"detection" db:"detection"`
	Result    []byte    `json:"result" db:"result"`
	Status    string    `json:"status" db:"status"`
	EventTime time.Time `json:"event_time" db:"event_time"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// EventResult represent static attributes for event history
type EventResult struct {
	Label     string    `json:"label"`
	Result    string    `json:"result"`
	Location  string    `json:"location"`
	Timestamp time.Time `json:"timestamp"`
}

// EventWebSocket represent attributes for websocket event data
type EventWebSocket struct {
	AnalyticID     string    `json:"analytic_id"`
	StreamID       string    `json:"stream_id"`
	Label          string    `json:"label"`
	Location       string    `json:"location"`
	PrimaryImage   []byte    `json:"primary_image"`
	SecondaryImage []byte    `json:"secondary_image"`
	Result         string    `json:"result"`
	Timestamp      time.Time `json:"timestamp"`
}

// ExportEventStatus enum
type ExportEventStatus string

const (
	// ExportEventStatusRunning running status
	ExportEventStatusRunning ExportEventStatus = "running"
	// ExportEventStatusDone done status
	ExportEventStatusDone ExportEventStatus = "done"
	// ExportEventStatusError error status
	ExportEventStatusError ExportEventStatus = "error"
	// ExportEventStatusDownloaded error status
	ExportEventStatusDownloaded ExportEventStatus = "downloaded"
)

type ExportEventTemplate struct {
	Items []*ExportEventItem
}

type ExportEventItem struct {
	ID             uint64
	EventType      string
	StreamID       string
	Detection      string
	PrimaryImage   string
	SecondaryImage string
	Result         string
	Status         string
	EventTime      string
	CreatedAt      string
}
