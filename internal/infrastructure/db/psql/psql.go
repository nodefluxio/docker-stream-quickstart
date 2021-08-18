package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // define postgres dialect
)

// NewPsqlRepository initiate postgres database client connection
func NewPsqlRepository(dbURL string, logLevel string) *gorm.DB {
	client, err := gorm.Open("postgres", dbURL)
	if err != nil {
		panic("failed to connect database")
	}
	client.SingularTable(true)
	if logLevel != "info" {
		client.LogMode(true)
	}

	return client
}
