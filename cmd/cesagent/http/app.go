package cesagent

import (
	"context"
	"net/url"
	"os"
	"os/exec"

	logutil "gitlab.com/nodefluxio/goutils/pkg/log"
	httphandler "gitlab.com/nodefluxio/vanilla-dashboard/cmd/cesagent/http/handler"
	dbinfra "gitlab.com/nodefluxio/vanilla-dashboard/internal/infrastructure/db/psql"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/interactor"
	httpRepo "gitlab.com/nodefluxio/vanilla-dashboard/internal/repository/http"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/repository/psql"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/usecase/agent"
)

// Options is cli arguments to start the app
type Options struct {
	LogLevel             string `validate:"required"`
	AppPort              string
	DatabaseHost         string `validate:"required"`
	DatabaseUsername     string `validate:"required"`
	DatabasePassword     string `validate:"required"`
	DatabaseName         string `validate:"required"`
	SyncPeriod           string `validate:"required"`
	AgentName            string `validate:"required"`
	CoordinatorURL       string `validate:"required"`
	EnrollmentVanillaURL string `validate:"required"`
	TotalEventSync       string `validate:"required"`
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
	migrationDirName := "migrations_ces_agent"
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

	// Initialize infrastuctures
	// Psql database
	psqlRepo := dbinfra.NewPsqlRepository(dbURL.String(), opt.LogLevel)
	defer psqlRepo.Close()
	latestTimestampRepo := psql.NewLatestTimestampRepository(psqlRepo)

	// http repository
	coordinatorRepo := httpRepo.NewCoordinatorServiceRepo(opt.CoordinatorURL)
	enrollmentVanillaRepo := httpRepo.NewEnrollmentVanillaServiceRepo(opt.EnrollmentVanillaURL)

	agentSvc := &agent.ServiceImpl{
		LatestTimestampRepo:   latestTimestampRepo,
		CoordinatorRepo:       coordinatorRepo,
		EnrollmentVanillaRepo: enrollmentVanillaRepo,
		SyncPeriod:            opt.SyncPeriod,
		AgentName:             opt.AgentName,
		TotalEventSync:        opt.TotalEventSync,
	}

	agentSvc.Sync(context.Background())
	agentSvc.CronjobSyncEnrollment(context.Background())

	// Tidying up all services to an interactor
	interactor := &interactor.AppInteractor{
		AgentSvc: agentSvc,
	}

	httphandler.New(opt.AppPort, interactor)
}
