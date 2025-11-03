// @title Gym REST API
// @version 1.0
// @description RESTful API for Gym App (Hacktiv8 Graded Challenge 3)
// @termsOfService http://swagger.io/terms/
// @contact.name Halim Iskandar
// @contact.email halim@example.com
// @license.name MIT
// @host localhost:8080
// @BasePath /
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"

	echoServer "gym/app/echoServer"
	"gym/config"
	//"gym/model"
	userrepo "gym/repository/user"
	authsvc "gym/service/auth"
	usersvc "gym/service/user"
	"gym/util/database"
	"gym/util/rapid"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}))
	slog.SetDefault(logger)
	ctx := context.Background()
	logger.Info("starting gym-api")

	//load config
	cfg, err := config.Load(ctx)
	if err != nil {
		logger.Error("failed to load config", slog.Any("err", err))
		os.Exit(1)
	}

	//open db
	db, err := database.New().Open(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("database connection failed", slog.Any("err", err))
		os.Exit(1)
	}

	// auto migrate
	//if err := db.AutoMigrate(&model.User{}, &model.Workout{}, &model.Exercise{}, &model.ExerciseLog{}); err != nil {
		//logger.Error("migration failed", slog.Any("err", err))
		//os.Exit(1)
	//}
	//logger.Info("migration success")

	//dependency inject
	userRepo := userrepo.New(db)
	authSvc := authsvc.New(userRepo)
	userSvc := usersvc.New(userRepo)
	bmiClient := rapid.NewBMI(cfg.RapidKey, cfg.RapidHost)

	// HTTP server
	e := echo.New()
	logger.Info("mounting routes")
	echoServer.Mount(e, cfg, echoServer.Controllers{
		Auth: echoServer.NewAuthController(authSvc, cfg.JWTSecret),
		User: echoServer.NewUserController(userSvc, bmiClient),
	})

	logger.Info("server running", slog.String("port", cfg.AppPort))
	if err := e.Start(":" + cfg.AppPort); err != nil {
		logger.Error("server stopped", slog.Any("err", err))
	}
}


