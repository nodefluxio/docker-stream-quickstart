package entity

// VisionaireStreamDetail represent attributes for in stream details api
type VisionaireStreamDetail struct {
	Pipelines []string `json:"pipelines"`
	Seats     []struct {
		AnalyticID   string `json:"analytic_id"`
		SerialNumber string `json:"serial_number"`
	} `json:"seats"`
	StreamAddress    string      `json:"stream_address"`
	StreamCustomData interface{} `json:"stream_custom_data"`
	StreamID         string      `json:"stream_id"`
	StreamLatitude   float64     `json:"stream_latitude"`
	StreamLongitude  float64     `json:"stream_longitude"`
	StreamName       string      `json:"stream_name"`
	StreamNodeNum    int         `json:"stream_node_num"`
	StreamStats      struct {
		Fps          int    `json:"fps"`
		FrameHeight  int    `json:"frame_height"`
		FrameWidth   int    `json:"frame_width"`
		LastErrorMsg string `json:"last_error_msg"`
		State        string `json:"state"`
	} `json:"stream_stats"`
}

// VisionaireStream attributes in streams list
type VisionaireStream struct {
	StreamNumber int                      `json:"stream_number"`
	Streams      []VisionaireStreamDetail `json:"streams"`
}
