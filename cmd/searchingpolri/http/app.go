package searchingpolri

import (
	"strconv"

	logutil "gitlab.com/nodefluxio/goutils/pkg/log"
	httphandler "gitlab.com/nodefluxio/vanilla-dashboard/cmd/searchingpolri/http/handler"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/interactor"
	httpRepo "gitlab.com/nodefluxio/vanilla-dashboard/internal/repository/http"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/usecase/polrisearching"
)

// Options is cli arguments to start the app
type Options struct {
	LogLevel             string `validate:"required"`
	AppPort              string
	PolriBaseURL         string `validate:"required"`
	PolriBasicUsername   string `validate:"required"`
	PolriBasicPassword   string `validate:"required"`
	PolriSearchPlatePath string `validate:"required"`
	PolriSearchNikPath   string `validate:"required"`
	SeagateBaseURL       string `validate:"required"`
	FremisnURL           string `validate:"required"`
	FremisKeyspace       string `validate:"required"`
	MaxImageSize         string `validate:"required"`
}

func Start(opt *Options) {

	// http repository
	polriRepo := httpRepo.NewPolriServiceRepo(
		opt.PolriBaseURL,
		opt.PolriBasicUsername,
		opt.PolriBasicPassword,
		opt.PolriSearchPlatePath,
		opt.PolriSearchNikPath,
	)
	seagateRepo := httpRepo.NewSeagateServiceRepo(
		opt.SeagateBaseURL,
	)
	fremisRepo := httpRepo.NewFremisRepository(
		opt.FremisnURL,
		opt.FremisKeyspace,
	)

	MaxImageSize, err := strconv.ParseInt(opt.MaxImageSize, 10, 64)
	if err != nil {
		logutil.LogObj.SetPanicLog(map[string]interface{}{
			"error": err,
		}, "failed convert MaxSizeImageEnrollment to int64")
	}
	kolantasSearchingSvc := &polrisearching.ServiceImpl{
		PolriRepo:          polriRepo,
		SeagateRepo:        seagateRepo,
		FRemisRepo:         fremisRepo,
		MaxSizeImageUpload: MaxImageSize,
	}

	// Tidying up all services to an interactor
	interactor := &interactor.AppInteractor{
		PolriSearchingSvc: kolantasSearchingSvc,
	}

	httphandler.New(opt.AppPort, interactor)
}
