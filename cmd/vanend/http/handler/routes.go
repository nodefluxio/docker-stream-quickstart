package httphandler

import (
	"context"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	httphandler "gitlab.com/nodefluxio/vanilla-dashboard/internal/delivery/http"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/interactor"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// NewHandler initialize http service for the app
func New(appPort string, interactor *interactor.AppInteractor) {
	router := httphandler.New()
	enrollmentHandler := &EnrollmentHandler{
		EnrollmentSvc: interactor.EnrollmentSvc,
	}
	eventHandler := &EventHandler{
		EventSvc: interactor.EventSvc,
	}
	globalSettingHandler := &GlobalSettingHandler{
		GlobalSettingSvc: interactor.GlobalSettingSvc,
	}
	authHandler := &AuthHandler{
		AuthSvc: interactor.AuthSvc,
	}

	apiGroup := router.Group("api")
	{
		authRoute := apiGroup.Group("auth")
		{
			authRoute.POST("token", authHandler.Login)
			// authRoute.PATCH("", authHandler.RefreshToken) will use next interation
		}

		enrollmentRoute := apiGroup.Group("enrollment")
		{
			enrollmentRoute.GET("", enrollmentHandler.GetList)
			enrollmentRoute.GET("/:id", enrollmentHandler.GetDetail)
			enrollmentRoute.DELETE("", enrollmentHandler.DeleteAll)
			enrollmentRoute.POST("", enrollmentHandler.Create)
			enrollmentRoute.PUT("/:id", enrollmentHandler.Update)
			enrollmentRoute.DELETE("/:id", enrollmentHandler.Delete)
		}

		apiGroup.GET("face/image/:id", enrollmentHandler.GetImage)

		apiGroup.GET("event_channel", eventHandler.GetEvent)

		eventRoute := apiGroup.Group("events")
		{
			eventRoute.GET("", eventHandler.History)
			eventRoute.GET("export", eventHandler.Export)
			eventRoute.GET("export/status", eventHandler.CheckExportStatus)
			eventRoute.GET("export/download", eventHandler.DownloadExport)
		}

		globalsettingRoute := apiGroup.Group("settings")
		{
			globalsettingRoute.GET("", globalSettingHandler.GetDetail)
			globalsettingRoute.POST("", globalSettingHandler.CreateOrUpdate)
		}
	}

	eventHandler.Dumping(context.Background())
	// Only in docker

	router.Use(static.Serve("/", static.LocalFile("./", true)))
	router.NoRoute(func(c *gin.Context) {
		c.File("./index.html")
	})
	router.Run(":" + appPort)
}
