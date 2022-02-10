package presenter

import "gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"

// StreamRequest is struct represent data reuqest for stream
type StreamRequest struct {
	NodeNumber int64  `json:"node_number"`
	StreamID   string `json:"stream_id"`
}

type StreamResponse struct {
	StreamNumber int                    `json:"stream_number"`
	Streams      []StreamDetailWithSite `json:"streams"`
}

type StreamDetailWithSite struct {
	entity.VisionaireStreamDetail
	StreamSiteID   uint64 `json:"stream_site_id"`
	StreamSiteName string `json:"stream_site_name"`
}
