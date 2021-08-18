package httphandler

import (
	httphandler "gitlab.com/nodefluxio/vanilla-dashboard/internal/delivery/http"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/interactor"
)

// NewHandler initialize http service for the app
func New(appPort string, interactor *interactor.AppInteractor) {
	router := httphandler.New()
	coordinatorHandler := &CoordinatorHandler{
		CoordinatorSvc: interactor.CoordinatorSvc,
	}

	v1Route := router.Group("/v1")
	{
		enrollmentRoute := v1Route.Group("/coordinators")
		{
			enrollmentRoute.GET("ping", coordinatorHandler.Ping)
			enrollmentRoute.GET("", coordinatorHandler.GetEventEnrollment)
			enrollmentRoute.POST("", coordinatorHandler.Create)
		}
	}

	router.Run(":" + appPort)
}
