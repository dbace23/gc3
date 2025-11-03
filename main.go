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
	"log"
	"log/slog"
	"os"
	"runtime"
	"time"

	"github.com/labstack/echo/v4"

	echoServer "gym/app/echoServer"
	"gym/config"
	userrepo "gym/repository/user"
	authsvc "gym/service/auth"
	usersvc "gym/service/user"
	"gym/util/database"
	"gym/util/rapid"
)

 
func logMem(tag string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	log.Printf("[%s] Alloc=%.2fMB Sys=%.2fMB NumGC=%v\n",
		tag, float64(m.Alloc)/1024/1024, float64(m.Sys)/1024/1024, m.NumGC)
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}))
	slog.SetDefault(logger)

	ctx := context.Background()
	logger.Info("starting gym-api")

	// 1️⃣ Load config
	cfg, err := config.Load(ctx)
	if err != nil {
		logger.Error("failed to load config", slog.Any("err", err))
		os.Exit(1)
	}

	// 2️⃣ Connect DB  
	db, err := database.New().Open(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("database connection failed", slog.Any("err", err))
		os.Exit(1)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(5)
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	logMem("after DB open")

	// ⚠️ Disable migrations on Heroku
	/*
	if err := db.AutoMigrate(&model.User{}, &model.Workout{}, &model.Exercise{}, &model.ExerciseLog{}); err != nil {
		logger.Error("migration failed", slog.Any("err", err))
		os.Exit(1)
	}
	*/

	// 3️⃣ Dependency injection
	userRepo := userrepo.New(db)
	authSvc := authsvc.New(userRepo)
	userSvc := usersvc.New(userRepo)
	bmiClient := rapid.NewBMI(cfg.RapidKey, cfg.RapidHost)

	// 4️⃣ HTTP server
	e := echo.New()
	logger.Info("mounting routes")
	echoServer.Mount(e, cfg, echoServer.Controllers{
		Auth: echoServer.NewAuthController(authSvc, cfg.JWTSecret),
		User: echoServer.NewUserController(userSvc, bmiClient),
	})

	// Use Heroku's $PORT or fallback to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	logger.Info("server running", slog.String("port", port))

	logMem("before start")
	if err := e.Start(":" + port); err != nil {
		logger.Error("server stopped", slog.Any("err", err))
	}
}
