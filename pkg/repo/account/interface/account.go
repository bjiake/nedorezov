package interfaces

import (
	"context"
	"nedorezov/pkg/domain/account"
)

type AccountRepository interface {
	Migrate(ctx context.Context) error
	Registration(ctx context.Context, newAccount account.Registration) (*account.Info, error)
	Login(ctx context.Context, acc account.Login) (int64, error)
	Put(ctx context.Context, id int64, updateAcc account.Registration) (*account.Info, error)
	Balance(ctx context.Context, id int64) (*account.Info, error)
	ChangeBalance(ctx context.Context, id int64, balance float64) (*float64, error)
}
