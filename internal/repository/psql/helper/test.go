package helper

import (
	"database/sql"
	"database/sql/driver"
	"log"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	logutil "gitlab.com/nodefluxio/goutils/pkg/log"
)

// AnyTime is helper for testing time
type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

// Suite is struct declarete variable use at setup
type Suite struct {
	Conn *sql.DB
	Db   *gorm.DB
	Mock sqlmock.Sqlmock
}

// Setup is helper for init psql connection interface in unit test
func Setup() *Suite {
	conn, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	db, err := gorm.Open("postgres", conn)
	if err != nil {
		panic("failed to connect database")
	}
	logutil.Init("info")

	return &Suite{
		Conn: conn,
		Db:   db,
		Mock: mock,
	}
}
