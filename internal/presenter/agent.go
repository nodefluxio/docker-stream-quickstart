package presenter

import "time"

type AgentStatus struct {
	Status            string    `json:"status"`
	LastSyncTimestamp time.Time `json:"last_sync_timestamp"`
}
