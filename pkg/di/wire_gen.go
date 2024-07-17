package di

import (
	"context"
	http "nedorezov/pkg/api"
	"nedorezov/pkg/api/handler"
	"nedorezov/pkg/config"
	"nedorezov/pkg/db"
	"nedorezov/pkg/repo/account"
	"nedorezov/pkg/service"
)

func InitializeAPI(cfg config.Config) (*http.ServerHTTP, error) {
	bd, err := db.ConnectToBD(cfg)
	if err != nil {
		return nil, err
	}
	// Repository
	accountRepository := account.NewAccountDataBase(bd)

	//service - logic
	userService := service.NewService(accountRepository)

	// Init Migrate
	err = userService.Migrate(context.Background())
	if err != nil {
		return nil, err
	}

	userHandler := handler.NewHandler(userService)
	serverHTTP := http.NewServerHTTP(userHandler)

	return serverHTTP, nil
}
