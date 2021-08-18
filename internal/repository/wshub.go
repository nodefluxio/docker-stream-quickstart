package repository

import "gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"

// WSHub interface abstracts the repository layer and should be implemented in repository
type WSHub interface {
	Run()
	BroadcastMessage(message []byte)
	RegisterClient(*entity.Client)
	UnregisterClient(*entity.Client)
}
