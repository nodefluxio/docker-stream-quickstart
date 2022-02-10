package httphandler

import (
	httphandler "gitlab.com/nodefluxio/vanilla-dashboard/internal/delivery/http"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/interactor"
)

// NewHandler initialize http service for the app
func New(appPort string, interactor *interactor.AppInteractor) {
	router := httphandler.New()
	polriSearchingHandler := &PolriSearchingHandler{
		PolriSearchingSvc: interactor.PolriSearchingSvc,
	}

	apiGroup := router.Group("api")
	{

		searcingRoute := apiGroup.Group("search")
		{
			polriGroup := searcingRoute.Group("polri")
			{
				polriGroup.GET("plate", polriSearchingHandler.GetPlate)
				polriGroup.GET("nik", polriSearchingHandler.GetNik)

				faceSimiGroup := polriGroup.Group("face-similarity")
				{
					faceSimiGroup.POST("token", polriSearchingHandler.GetFaceSearchToken)
					faceSimiGroup.GET("results", polriSearchingHandler.GetFaceSearchResult)
				}
			}
		}
	}

	router.Run(":" + appPort)
}
