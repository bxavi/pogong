// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package db

import (
	"context"
)

type Querier interface {
	CreateAccounts(ctx context.Context, arg CreateAccountsParams) (*Account, error)
	DeleteAccounts(ctx context.Context, id int64) error
	GetAccounts(ctx context.Context, id int64) (*Account, error)
	ListAccounts(ctx context.Context) ([]*Account, error)
	UpdateAccounts(ctx context.Context, arg UpdateAccountsParams) (*Account, error)
}

var _ Querier = (*Queries)(nil)