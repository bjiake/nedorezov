//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	http "nedorezov/pkg/api"
	"nedorezov/pkg/api/handler"
	"nedorezov/pkg/config"
	"nedorezov/pkg/db"
	"nedorezov/pkg/service"
)

func InitializeAPI(cfg config.Config) (*http.ServerHTTP, error) {
	wire.Build(db.ConnectToBD, service.NewService, handler.NewHandler, http.NewServerHTTP)

	return &http.ServerHTTP{}, nil
}
