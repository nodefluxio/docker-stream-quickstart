package presenter

// SiteRequest is struct for handling data site request
type SiteRequest struct {
	Name string `json:"name"`
	ID   uint64
}

// AssignStreamRequest is struct for handling data assign stream to site request
type AssignStreamRequest struct {
	SiteID   uint64
	StreamID string `json:"stream_id"`
}
