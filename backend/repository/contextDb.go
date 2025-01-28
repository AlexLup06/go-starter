package repository

import (
	"context"
)

type ContextDb interface {
	WithContext(parentContext context.Context) (context.Context, error)
	WithCancel(parentContext context.Context) (context.Context, context.CancelFunc, error)
	IsCommitted(ctx context.Context) bool
	IsRolledBack(ctx context.Context) bool
	IsCommittedOrRolledBack(ctx context.Context) bool
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	WithTransaction(ctx context.Context, op func(context.Context) error) error
}
