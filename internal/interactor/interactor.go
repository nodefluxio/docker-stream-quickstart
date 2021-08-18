package interactor

import (
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/usecase/agent"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/usecase/auth"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/usecase/coordinator"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/usecase/enrollment"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/usecase/event"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/usecase/globalsetting"
)

// AppInteractor represents app's interactor object
type AppInteractor struct {
	EnrollmentSvc    enrollment.Service
	EventSvc         event.Service
	GlobalSettingSvc globalsetting.Service
	CoordinatorSvc   coordinator.Service
	AgentSvc         agent.Service
	AuthSvc          auth.Service
}
