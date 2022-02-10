package interactor

import (
	ut "github.com/go-playground/universal-translator"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/usecase/agent"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/usecase/auth"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/usecase/coordinator"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/usecase/enrollment"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/usecase/event"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/usecase/globalsetting"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/usecase/polrisearching"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/usecase/site"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/usecase/stream"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/usecase/user"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/usecase/vehicle"
	"gopkg.in/go-playground/validator.v9"
)

// AppInteractor represents app's interactor object
type AppInteractor struct {
	EnrollmentSvc     enrollment.Service
	EventSvc          event.Service
	GlobalSettingSvc  globalsetting.Service
	CoordinatorSvc    coordinator.Service
	AgentSvc          agent.Service
	AuthSvc           auth.Service
	VehicleSvc        vehicle.Service
	Validator         *validator.Validate
	Translator        ut.Translator
	PolriSearchingSvc polrisearching.Service
	SiteSvc           site.Service
	UserSvc           user.Service
	StreamScv         stream.Service
}
