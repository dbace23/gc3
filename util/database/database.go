package database

import (
	"context"
	"fmt"
	"time"

	"github.com/labstack/gommon/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Opener interface {
	Open(ctx context.Context, dsn string) (*gorm.DB, error)
}

type opener struct{}

func New() Opener { return &opener{} }

func (o *opener) Open(ctx context.Context, dsn string) (*gorm.DB, error) {
	cfg := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	}
	db, err := gorm.Open(postgres.Open(dsn), cfg)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	sql, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("native db: %w", err)
	}
	sql.SetMaxOpenConns(3)
	sql.SetMaxIdleConns(2)
	sql.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := sql.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}
	log.Info("database connected")
	return db, nil
}

