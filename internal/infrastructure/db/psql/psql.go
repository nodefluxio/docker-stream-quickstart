package db

import (
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // define postgres dialect
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
)

// NewPsqlRepository initiate postgres database client connection
func NewPsqlRepository(option *entity.PsqlDBConnOption, logLevel string) *gorm.DB {
	client, err := gorm.Open("postgres", option.URL)
	if err != nil {
		panic("failed to connect database")
	}
	client.SingularTable(true)
	// Reference: https://www.alexedwards.net/blog/configuring-sqldb
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	maxIdleConn, _ := strconv.Atoi(option.MaxIdleConn)
	client.DB().SetMaxIdleConns(maxIdleConn)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	maxOpenConn, _ := strconv.Atoi(option.MaxOpenConn)
	client.DB().SetMaxOpenConns(maxOpenConn)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	connMaxLifeTimeInMinutes, _ := strconv.Atoi(option.MaxLifetimeInMinute)
	client.DB().SetConnMaxLifetime(time.Duration(connMaxLifeTimeInMinutes) * time.Minute)
	if logLevel != "info" {
		client.LogMode(true)
	}

	return client
}
