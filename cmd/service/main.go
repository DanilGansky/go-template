package main

import (
	"net/http"

	"github.com/littlefut/go-template/pkg/db"
	"github.com/littlefut/go-template/pkg/server"

	"github.com/littlefut/go-template/pkg/logger"

	"github.com/littlefut/go-template/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/littlefut/go-template/internal/api"
	"github.com/littlefut/go-template/internal/auth"
	"github.com/littlefut/go-template/internal/hash"
	"github.com/littlefut/go-template/internal/repository"
	"github.com/littlefut/go-template/internal/user"
)

func main() {
	cfg := config.Get()
	log := logger.Get(cfg.LogLevel())
	dbConn := db.Get(cfg.DSN)

	log.Debugf("PROVIDED CONFIG: %v", cfg)

	userRepo := repository.NewUserRepository(dbConn)
	hashSvc := hash.NewService(cfg.Cost)
	tokenSvc, err := hash.NewTokenService(cfg.Secret, cfg.Issuer)
	if err != nil {
		log.Panicf("failed to create token service: %s", err.Error())
	}

	userSvc := user.NewService(hashSvc, userRepo)
	authSvc := auth.NewService(hashSvc, tokenSvc, userSvc)

	router := gin.Default()
	authMiddleware := api.AuthorizationMiddleware(userSvc, tokenSvc)
	opts := &api.Options{
		Log:     log,
		Timeout: cfg.MaxTimeout(),
		Router:  router,
	}

	api.NewUserController(userSvc, opts, authMiddleware)
	api.NewAuthController(authSvc, userSvc, opts)

	httpServer := &http.Server{Addr: cfg.Addr, Handler: router}
	if err = server.Run(httpServer, cfg.MaxTimeout(), log); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}
}
