package interfaces

import (
	"context"
	"nedorezov/pkg/domain/account"
)

type ServiceUseCase interface {
	Migrate(ctx context.Context) error

	//Account
	Registration(ctx context.Context, userId string, newAccount account.Registration) (*account.Info, error)
	Login(ctx context.Context, acc account.Login) (int64, error)
	PutAccount(ctx context.Context, id string, updateAcc account.Registration) (*account.Info, error)

	GetBalance(ctx context.Context, id string) (*float64, error)
	Deposit(ctx context.Context, id string, amount float64) (*float64, error)
	WithDraw(ctx context.Context, id string, amount float64) (*float64, error)
}
