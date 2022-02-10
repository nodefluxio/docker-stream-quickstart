package entity

import "time"

type EnrollmentEventCoordinator struct {
	Message string `json:"message"`
	Ok      bool   `json:"ok"`
	Results []struct {
		EventID     string                            `json:"event_id"`
		Agent       string                            `json:"agent"`
		EventAction string                            `json:"event_action"`
		Images      []string                          `json:"images"`
		Payload     PayloadEnrollmentEventCoordinator `json:"payload"`
		CreatedAt   string                            `json:"created_at"`
	} `json:"results"`
}

type PayloadEnrollmentEventCoordinator struct {
	FaceID         uint64    `json:"face_id"`
	Name           string    `json:"name"`
	IdentityNumber string    `json:"identity_number"`
	Gender         string    `json:"gender" db:"gender"`
	BirthPlace     string    `json:"birth_place" db:"birth_place"`
	BirthDate      time.Time `json:"birth_date" db:"birth_date"`
	Status         string    `json:"status"`
}

type CreateEnrollmentEventCoordinator struct {
	Agent       string                            `json:"agent"`
	EventAction string                            `json:"event_action"`
	Images      []string                          `json:"images"`
	Payload     PayloadEnrollmentEventCoordinator `json:"payload"`
}
