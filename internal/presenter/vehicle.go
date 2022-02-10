package presenter

import (
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/pkg/util"
)

// VehicleEnrollmentRequest is handler to processs from upload
type VehicleEnrollmentRequest struct {
	Plate    string `json:"plate_number" validate:"required"`
	UniqueID string `json:"unique_id"  validate:"required"`
	Type     string `json:"type"`
	Brand    string `json:"brand"`
	Color    string `json:"color"`
	Status   string `json:"status"`
	Name     string `json:"name"`
}

// VehicleEnrollmentResponse is handling response enrollment
type VehicleEnrollmentResponse struct {
	ID       uint64 `json:"id"`
	Plate    string `json:"plate_number"`
	Type     string `json:"type"`
	Brand    string `json:"brand"`
	Color    string `json:"color"`
	Status   string `json:"status"`
	Name     string `json:"name"`
	UniqueID string `json:"unique_id"`
}

// EnrollmentPaging is struct for handling list enrollment with pagination
type VehiclePaging struct {
	util.PaginationDetails
	Vehicles []*VehicleEnrollmentResponse `json:"vehicles"`
}
