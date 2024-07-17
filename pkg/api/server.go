package api

import (
	"github.com/gin-gonic/gin"
	"nedorezov/pkg/api/handler"
)

type ServerHTTP struct {
	engine *gin.Engine
}

func NewServerHTTP(userHandler *handler.Handler) *ServerHTTP {
	engine := gin.New()

	// Use logger from Gin
	engine.Use(gin.Logger())

	engine.POST("/registration", userHandler.Registration)
	engine.POST("/login", userHandler.Login)

	// Use middleware from Gin
	engine.Use(userHandler.AuthMiddleware())

	//Peoples
	engine.POST("/accounts/deposit", userHandler.Deposit)
	engine.POST("/accounts/withdraw", userHandler.WithDraw)
	engine.GET("/accounts/balance", userHandler.GetBalance)
	engine.PUT("/accounts", userHandler.PutAccount)

	return &ServerHTTP{engine: engine}
}

func (sh *ServerHTTP) Start() {
	sh.engine.Run("0.0.0.0:8001")
}
