package psql

import (
	"context"

	"github.com/jinzhu/gorm"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/repository"
)

type psqlTransactionRepo struct {
	Conn *gorm.DB
}

// NewPsqlTransactionRepository is method to initiate PsqlTransaction repo
func NewPsqlTransactionRepository(conn *gorm.DB) repository.PsqlTransaction {
	return &psqlTransactionRepo{
		Conn: conn,
	}
}

func (p *psqlTransactionRepo) BeginTransaction(ctx context.Context) *gorm.DB {
	return p.Conn.Begin()
}

func (p *psqlTransactionRepo) CommitTransaction(ctx context.Context, tx *gorm.DB) *gorm.DB {
	return tx.Commit()
}

func (p *psqlTransactionRepo) RollbackTransaction(ctx context.Context, tx *gorm.DB) *gorm.DB {
	return tx.Rollback()
}
