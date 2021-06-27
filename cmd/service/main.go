package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
	"github.com/littlefut/go-template/config"
	"github.com/littlefut/go-template/internal/api"
	"github.com/littlefut/go-template/internal/auth"
	"github.com/littlefut/go-template/internal/hash"
	"github.com/littlefut/go-template/internal/repository"
	"github.com/littlefut/go-template/internal/user"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Get()
	fmt.Printf("PROVIDED CONFIG: %v", cfg)

	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{})
	log.SetLevel(logrus.DebugLevel)

	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		log.Panicf("failed to connect to db: %s", err.Error())
	}

	userRepo := repository.NewUserRepository(db)
	hashSvc := hash.NewService(cfg.Cost)
	tokenSvc, err := hash.NewTokenService(cfg.Secret)
	if err != nil {
		log.Panicf("failed to create token service: %s", err.Error())
	}

	userSvc := user.NewService(hashSvc, userRepo)
	authSvc := auth.NewService(hashSvc, tokenSvc, userSvc)

	router := gin.Default()

	authMiddleware := api.AuthorizationMiddleware(userSvc, tokenSvc)
	api.NewUserController(userSvc, log, cfg.MaxTimeout(), router, authMiddleware)
	api.NewAuthController(authSvc, userSvc, log, cfg.MaxTimeout(), router)

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	go func() {
		log.Infof("server started at: %s", cfg.Addr)
		if err = server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatalf("server error: %s", err.Error())
			}
		}
	}()

	<-signals
	log.Info("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.MaxTimeout())
	defer cancel()

	if err = server.Shutdown(ctx); err != nil {
		log.Warnf("shutdown error: %s", err.Error())
	}
}
