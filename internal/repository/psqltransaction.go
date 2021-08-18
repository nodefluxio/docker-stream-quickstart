package repository

import (
	"context"

	"github.com/jinzhu/gorm"
)

// PsqlTransaction is interface for transaction level at psql database
type PsqlTransaction interface {
	BeginTransaction(ctx context.Context) *gorm.DB
	CommitTransaction(ctx context.Context, tx *gorm.DB) *gorm.DB
	RollbackTransaction(ctx context.Context, tx *gorm.DB) *gorm.DB
}
