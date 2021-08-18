package vanend

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strconv"

	logutil "gitlab.com/nodefluxio/goutils/pkg/log"
	httphandler "gitlab.com/nodefluxio/vanilla-dashboard/cmd/vanend/http/handler"
	dbinfra "gitlab.com/nodefluxio/vanilla-dashboard/internal/infrastructure/db/psql"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/interactor"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/repository/goroutine"
	httpRepo "gitlab.com/nodefluxio/vanilla-dashboard/internal/repository/http"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/repository/psql"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/usecase/auth"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/usecase/enrollment"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/usecase/event"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/usecase/globalsetting"
)

// Options is cli arguments to start the app
type Options struct {
	LogLevel               string `validate:"required"`
	FEAppHost              string
	AppPort                string
	NodeEnv                string
	DatabaseHost           string `validate:"required"`
	DatabaseUsername       string `validate:"required"`
	DatabasePassword       string `validate:"required"`
	DatabaseName           string `validate:"required"`
	FremisnURL             string `validate:"required"`
	FremisKeyspace         string `validate:"required"`
	VisionaireHost         string `validate:"required"`
	CronjobPartitionSpec   string `validate:"required"`
	UseCES                 string `validate:"required"`
	AgentURL               string `validate:"required"`
	MaxSizeImageEnrollment string `validate:"required"`
}

func initFE(opt *Options) {
	logutil.LogObj.SetInfoLog(map[string]interface{}{}, "init frontend")
	os.Setenv("ENV_DEST", "/")
	os.Setenv("NODE_ENV", opt.NodeEnv)
	os.Setenv("REACT_APP_VISIONAIRE_HOST", opt.VisionaireHost)
	os.Setenv("REACT_APP_HOST", opt.FEAppHost)
	cmd := exec.Command("env-gen")
	cmd.Stdout = os.Stdout
	err := cmd.Start()
	if err != nil {
		logutil.LogObj.SetPanicLog(map[string]interface{}{
			"error": err,
		}, "failed init frontend")
	}
}

func Start(opt *Options) {
	// construct database url
	dbURL := &url.URL{
		Scheme:   "postgres",
		Host:     opt.DatabaseHost,
		User:     url.UserPassword(opt.DatabaseUsername, opt.DatabasePassword),
		Path:     opt.DatabaseName,
		RawQuery: "sslmode=disable",
	}
	if opt.DatabasePassword == "" {
		dbURL.User = url.User(opt.DatabaseUsername)
	}

	// Run database migration
	logutil.LogObj.SetInfoLog(map[string]interface{}{}, "starting database migration...")
	migrationDirName := "migrations"
	cmd := exec.Command("./script/migration.sh", migrationDirName, dbURL.String(), "up")
	cmd.Stdout = os.Stdout
	err := cmd.Start()
	if err != nil {
		logutil.LogObj.SetPanicLog(map[string]interface{}{
			"error": err,
		}, "failed run database migration")
	}
	cmd.Wait()
	logutil.LogObj.SetInfoLog(map[string]interface{}{}, "finished run database migration")

	if opt.NodeEnv == "production" {
		initFE(opt)
	}

	// Initialize infrastuctures
	// Psql database
	psqlRepo := dbinfra.NewPsqlRepository(dbURL.String(), opt.LogLevel)
	defer psqlRepo.Close()
	enrolledFaceRepo := psql.NewEnrollFaceRepository(psqlRepo)
	eventRepo := psql.NewEventRepository(psqlRepo)
	faceImageRepo := psql.NewFaceImageRepository(psqlRepo)
	globalSettingRepo := psql.NewGlobalSettingRepository(psqlRepo)
	psqlTransactionRepo := psql.NewPsqlTransactionRepository(psqlRepo)

	// http repository
	streamRepo := httpRepo.NewStreamServiceRepo(fmt.Sprintf("http://%s", opt.VisionaireHost))
	fremisRepo := httpRepo.NewFremisRepository(opt.FremisnURL, opt.FremisKeyspace)
	agentRepo := httpRepo.NewAgentServiceRepo(opt.AgentURL)

	// wshub repository
	hubRepo := goroutine.NewHub()
	go hubRepo.Run()

	MaxSizeImageEnrollment, err := strconv.ParseInt(opt.MaxSizeImageEnrollment, 10, 64)
	if err != nil {
		logutil.LogObj.SetPanicLog(map[string]interface{}{
			"error": err,
		}, "failed convert MaxSizeImageEnrollment to int64")
	}
	enrollmentSvc := &enrollment.ServiceImpl{
		EnrolledFaceRepo:       enrolledFaceRepo,
		FRemisRepo:             fremisRepo,
		FaceImageRepo:          faceImageRepo,
		UseCES:                 opt.UseCES,
		AgentRepo:              agentRepo,
		PsqlTransactionRepo:    psqlTransactionRepo,
		MaxSizeImageEnrollment: MaxSizeImageEnrollment,
	}

	visionaireWebsocket := fmt.Sprintf("ws://%s/event_channel", opt.VisionaireHost)
	eventSvc := &event.ServiceImpl{
		EnrolledFaceRepo:     enrolledFaceRepo,
		FRemisRepo:           fremisRepo,
		WSHubRepo:            hubRepo,
		URLGridLiteWS:        visionaireWebsocket,
		CronjobPartitionSpec: opt.CronjobPartitionSpec,
		EventRepo:            eventRepo,
		StreamRepo:           streamRepo,
		GlobalSettingRepo:    globalSettingRepo,
	}
	authSvc := &auth.ServiceImpl{}
	globalSettingSvc := &globalsetting.ServiceImpl{
		GlobalSettingRepo: globalSettingRepo,
	}
	// runnting event partition
	eventSvc.Partition(context.Background())
	eventSvc.CronjobPartition(context.Background())

	// Tidying up all services to an interactor
	interactor := &interactor.AppInteractor{
		EnrollmentSvc:    enrollmentSvc,
		EventSvc:         eventSvc,
		GlobalSettingSvc: globalSettingSvc,
		AuthSvc:          authSvc,
	}

	httphandler.New(opt.AppPort, interactor)
}
