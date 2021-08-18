package entity

import "time"

type EventEnrollment struct {
	ID          uint64    `json:"id" db:"id"`
	EventID     string    `json:"event_id" db:"event_id"`
	Agent       string    `json:"agent" db:"agent"`
	EventAction string    `json:"event_action" db:"event_action"`
	Payload     []byte    `json:"payload" db:"payload"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// EventEnrollmentPayload is payload from enrollment event
type EventEnrollmentImage struct {
	Image []byte `json:"Image"`
}
