package httphandler

import (
	httphandler "gitlab.com/nodefluxio/vanilla-dashboard/internal/delivery/http"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/interactor"
)

// NewHandler initialize http service for the app
func New(appPort string, interactor *interactor.AppInteractor) {
	router := httphandler.New()
	AgentHandler := &AgentHandler{
		AgentSvc: interactor.AgentSvc,
	}
	v1Route := router.Group("/v1")
	{
		enrollmentRoute := v1Route.Group("/agents")
		{
			enrollmentRoute.GET("ping", AgentHandler.Ping)
			enrollmentRoute.POST("event-enrollments", AgentHandler.Create)
			enrollmentRoute.GET("status", AgentHandler.GetStatus)
		}
	}

	router.Run(":" + appPort)
}
