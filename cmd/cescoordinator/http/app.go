package cescoordinator

import (
	"context"
	"net/url"
	"os"
	"os/exec"

	logutil "gitlab.com/nodefluxio/goutils/pkg/log"
	httphandler "gitlab.com/nodefluxio/vanilla-dashboard/cmd/cescoordinator/http/handler"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
	dbinfra "gitlab.com/nodefluxio/vanilla-dashboard/internal/infrastructure/db/psql"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/interactor"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/repository/psql"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/usecase/coordinator"
)

// Options is cli arguments to start the app
type Options struct {
	LogLevel                    string `validate:"required"`
	AppPort                     string
	DatabaseHost                string `validate:"required"`
	DatabaseUsername            string `validate:"required"`
	DatabasePassword            string `validate:"required"`
	DatabaseName                string `validate:"required"`
	DatabaseMaxIdleConn         string `validate:"required"`
	DatabaseMaxOpenConn         string `validate:"required"`
	DatabaseMaxLifetimeInMinute string `validate:"required"`
	CronjobPartitionSpec        string `validate:"required"`
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
	migrationDirName := "migrations_ces_coordinator"
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
	dbOption := &entity.PsqlDBConnOption{
		URL:                 dbURL.String(),
		MaxIdleConn:         opt.DatabaseMaxIdleConn,
		MaxOpenConn:         opt.DatabaseMaxOpenConn,
		MaxLifetimeInMinute: opt.DatabaseMaxLifetimeInMinute,
	}
	psqlRepo := dbinfra.NewPsqlRepository(dbOption, opt.LogLevel)
	defer psqlRepo.Close()
	eventEnrollmentRepo := psql.NewEventEnrollmentRepository(psqlRepo)
	eventEnrollmentFaceImageRepo := psql.NewEventEnrollmentFaceImageRepository(psqlRepo)
	// http repository
	coordinatorSvc := &coordinator.ServiceImpl{
		EventEnrollmentRepo:          eventEnrollmentRepo,
		EventEnrollmentFaceImageRepo: eventEnrollmentFaceImageRepo,
	}
	// runnting event partition
	coordinatorSvc.Partition(context.Background())
	coordinatorSvc.CronjobPartition(context.Background())
	// Tidying up all services to an interactor
	interactor := &interactor.AppInteractor{
		CoordinatorSvc: coordinatorSvc,
	}

	httphandler.New(opt.AppPort, interactor)
}
