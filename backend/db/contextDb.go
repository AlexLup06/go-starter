package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"alexlupatsiy.com/personal-website/backend/repository"
	"gorm.io/gorm"
)

type transactionState struct {
	session    *gorm.DB
	committed  bool
	rolledBack bool
}

type contextKey string

const (
	dbKey contextKey = "db"
	txKey contextKey = "txState"
)

type ContextDb struct {
	gorm *gorm.DB
}

func NewContextDb(gorm *gorm.DB) repository.ContextDb {
	return &ContextDb{gorm: gorm}
}

func (d ContextDb) WithContext(parentContext context.Context) (context.Context, error) {
	session := d.gorm.WithContext(parentContext).Begin(&sql.TxOptions{ReadOnly: false, Isolation: sql.LevelReadCommitted})
	if err := session.Error; err != nil {
		return nil, err
	}
	ctx := context.WithValue(parentContext, dbKey, session)
	ctx = context.WithValue(ctx, txKey, &transactionState{})
	go d.cleanup(ctx)
	return ctx, nil
}

func (d ContextDb) IsCommitted(ctx context.Context) bool {
	txState, ok := ctx.Value(txKey).(*transactionState)
	if !ok {
		return false
	}
	return txState.committed
}

func (d ContextDb) IsRolledBack(ctx context.Context) bool {
	txState, ok := ctx.Value(txKey).(*transactionState)
	if !ok {
		return false
	}
	return txState.rolledBack
}

func (d ContextDb) IsCommittedOrRolledBack(ctx context.Context) bool {
	return d.IsCommitted(ctx) || d.IsRolledBack(ctx)
}

func (d ContextDb) Commit(ctx context.Context) error {
	db, err := getContextDb(ctx)
	if err != nil {
		return err
	}
	err = db.Commit().Error
	if err == nil {
		committed := true
		if err = updateTransactionState(ctx, &committed, nil); err != nil {
			return err
		}
	}
	return err
}

func (d ContextDb) Rollback(ctx context.Context) error {
	if d.IsCommitted(ctx) {
		return errors.New("transaction is already committed, cannot rollback")
	}
	db, err := getContextDb(ctx)
	if err != nil {
		return err
	}
	err = db.Rollback().Error
	if err == nil {
		rolledBack := true
		if err = updateTransactionState(ctx, nil, &rolledBack); err != nil {
			return err
		}
	}
	return err
}

func (d ContextDb) cleanup(ctx context.Context) {
	<-ctx.Done()
	if !d.IsCommittedOrRolledBack(ctx) {
		_ = d.Rollback(ctx)
	}
}

func (d ContextDb) WithTransaction(ctx context.Context, op func(context.Context) error) error {
	dbCtx, cancelFn, err := d.WithCancel(ctx)
	defer cancelFn()
	if err != nil {
		return err
	}
	err = op(dbCtx)
	if err != nil {
		if !d.IsCommittedOrRolledBack(dbCtx) {
			if rbErr := d.Rollback(dbCtx); rbErr != nil {
				return rbErr
			}
			return err
		}
	}
	if !d.IsCommittedOrRolledBack(dbCtx) {
		if cmErr := d.Commit(dbCtx); cmErr != nil {
			return cmErr
		}
		return err
	}
	return nil
}

func (d ContextDb) WithCancel(parentContext context.Context) (context.Context, context.CancelFunc, error) {
	ctx, cancelFn := context.WithCancel(parentContext)
	dbCtx, err := d.WithContext(ctx)
	return dbCtx, cancelFn, err
}

// getContextDb gets gorm session from context
func getContextDb(ctx context.Context) (*gorm.DB, error) {
	session, ok := ctx.Value(dbKey).(*gorm.DB)
	if !ok {
		return nil, fmt.Errorf("context does not contain a gorm session")
	}
	return session, nil
}

// updateTransactionState update transaction state in context
func updateTransactionState(ctx context.Context, committed *bool, rolledBack *bool) error {
	txState, ok := ctx.Value(txKey).(*transactionState)
	if !ok {
		return errors.New("context does not contain a transaction State")
	}
	if committed != nil {
		txState.committed = *committed
	}
	if rolledBack != nil {
		txState.rolledBack = *rolledBack
	}
	return nil
}
