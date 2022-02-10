package router

import (
	"context"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"gitlab.com/nodefluxio/vanilla-dashboard/cmd/vanend/http/handler"
	"gitlab.com/nodefluxio/vanilla-dashboard/cmd/vanend/http/middleware"
	httphandler "gitlab.com/nodefluxio/vanilla-dashboard/internal/delivery/http"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/interactor"

	validatorHelper "gitlab.com/nodefluxio/vanilla-dashboard/internal/delivery/http/middleware/validator/helper"

	validation "gopkg.in/go-playground/validator.v9"
)

// NewHandler initialize http service for the app
func New(appPort string, interactor *interactor.AppInteractor) {

	validator := validation.New()
	i18n := validatorHelper.RequiredErrorMessage(validator)
	router := httphandler.New()

	// prepare handler
	enrollmentHandler := &handler.EnrollmentHandler{
		EnrollmentSvc: interactor.EnrollmentSvc,
	}
	eventHandler := &handler.EventHandler{
		EventSvc: interactor.EventSvc,
	}
	globalSettingHandler := &handler.GlobalSettingHandler{
		GlobalSettingSvc: interactor.GlobalSettingSvc,
	}
	authHandler := &handler.AuthHandler{
		AuthSvc: interactor.AuthSvc,
	}
	vehicleHandler := &handler.VehicleHandler{
		VehicleSvc: interactor.VehicleSvc,
		Validator:  validator,
		Translator: i18n,
	}
	SiteHandler := &handler.SiteHandler{
		SiteSvc: interactor.SiteSvc,
	}
	UserHandler := &handler.UserHandler{
		UserSvc: interactor.UserSvc,
	}
	StreamHandler := &handler.StreamHandler{
		StreamSvc: interactor.StreamScv,
	}

	// prepare middleware
	AuthMiddlware := &middleware.AuthMiddleware{
		AuthSvc: interactor.AuthSvc,
	}
	AuthorizeMiddleware := &middleware.AuthorizeMiddleware{}

	apiGroup := router.Group("api")
	{
		authRoute := apiGroup.Group("auth")
		{
			authRoute.POST("token", authHandler.Login)
			// authRoute.PATCH("", authHandler.RefreshToken) will use next interation
		}

		vehicleRoute := apiGroup.Group("vehicles", AuthMiddlware.IsLoggedIn)
		{
			vehicleRoute.GET("", vehicleHandler.GetList)
			vehicleRoute.GET("/:id", vehicleHandler.GetDetail)
			vehicleRoute.DELETE("", vehicleHandler.DeleteAll)
			vehicleRoute.POST("", vehicleHandler.Create)
			vehicleRoute.PUT("/:id", vehicleHandler.Update)
			vehicleRoute.DELETE("/:id", vehicleHandler.Delete)
		}

		filesRoute := apiGroup.Group("files", AuthMiddlware.IsLoggedIn)
		{
			filesRoute.GET("/enrollment", enrollmentHandler.Backup)
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
			eventRoute.GET("", AuthMiddlware.IsLoggedIn, eventHandler.History)
			eventRoute.GET("export", AuthMiddlware.IsLoggedIn, eventHandler.Export)
			eventRoute.GET("export/status", AuthMiddlware.IsLoggedIn, eventHandler.CheckExportStatus)
			eventRoute.GET("export/download", eventHandler.DownloadExport)
			eventRoute.GET("insight", AuthMiddlware.IsLoggedIn, eventHandler.Insight)
		}

		globalsettingRoute := apiGroup.Group("settings", AuthMiddlware.IsLoggedIn)
		{
			globalsettingRoute.GET("", globalSettingHandler.GetDetail)
			globalsettingRoute.POST("", globalSettingHandler.CreateOrUpdate)
		}

		siteRoute := apiGroup.Group("sites", AuthMiddlware.IsLoggedIn)
		{
			siteRoute.GET("", SiteHandler.GetList)
			siteRoute.POST("", AuthorizeMiddleware.CheckAccess, SiteHandler.Create)
			siteRoute.PUT("/:id", AuthorizeMiddleware.CheckSite, SiteHandler.Update)
			siteRoute.DELETE("/:id", AuthorizeMiddleware.CheckSite, SiteHandler.Delete)
			siteRoute.POST("/:id/assign-stream", AuthorizeMiddleware.CheckSite, SiteHandler.AssignStream)
		}

		userRoute := apiGroup.Group("manage-users", AuthMiddlware.IsLoggedIn)
		{
			userRoute.GET("", UserHandler.GetList)
			userRoute.GET("/:id", UserHandler.Detail)
			userRoute.POST("", UserHandler.Create)
			userRoute.PUT("/:id", UserHandler.Update)
			userRoute.PUT("/:id/change-password", UserHandler.ChangePassword)
			userRoute.DELETE("/:id", UserHandler.Delete)
		}

		streamRoute := apiGroup.Group("streams", AuthMiddlware.IsLoggedIn)
		{
			streamRoute.GET("", StreamHandler.GetList)
			streamRoute.GET("/:node_id/:stream_id", StreamHandler.GetDetail)
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
