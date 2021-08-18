package presenter

import (
	"encoding/json"
	"time"
)

// CoordinatorRequest is handler to processs from upload
type CoordinatorRequest struct {
	Agent       string          `json:"agent"`
	EventAction string          `json:"event_action"`
	Images      []string        `json:"images"`
	Payload     json.RawMessage `json:"payload"`
}

type CoordinatorResponse struct {
	EventID     string                 `json:"event_id"`
	Agent       string                 `json:"agent"`
	EventAction string                 `json:"event_action"`
	Images      []string               `json:"images"`
	Payload     EventEnrollmentPayload `json:"payload"`
	CreatedAt   time.Time              `json:"created_at"`
}

// EventEnrollmentPayload is payload from enrollment event
type EventEnrollmentPayload struct {
	FaceID         uint64 `json:"face_id"`
	Name           string `json:"name"`
	IdentityNumber string `json:"identity_number"`
	Status         string `json:"status"`
}
