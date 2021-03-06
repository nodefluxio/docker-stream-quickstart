package presenter

import (
	"time"

	"gitlab.com/nodefluxio/vanilla-dashboard/internal/pkg/util"
)

// EventHistoryPaging is struct for handling events pagination
type EventHistoryPaging struct {
	util.PaginationDetails
	Events []*EventGroup `json:"events"`
}

// EventGroup is struct for handling event grouping
type EventGroup struct {
	Timestamp string       `json:"timestamp"`
	Data      []*EventData `json:"data"`
}

// EventData is struct for handling detail event
type EventData struct {
	ID             uint64    `json:"id"`
	AnalyticID     string    `json:"analytic_id"`
	PrimaryImage   []byte    `json:"primary_image"`
	SecondaryImage []byte    `json:"secondary_image"`
	Label          string    `json:"label"`
	Result         string    `json:"result"`
	Location       string    `json:"location"`
	Timestamp      time.Time `json:"timestamp"`
}

// EventInsightData is struct for handling data in insight event
type EventInsightData struct {
	TotalToday     uint64 `json:"total_today"`
	TotalYesterday uint64 `json:"total_yesterday"`
	TotalWeek      uint64 `json:"total_week"`
	TotalMonth     uint64 `json:"total_month"`
}
