package model

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey"`
	Email        string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	Username     string    `gorm:"type:varchar(50);uniqueIndex;not null"`
	FullName     string    `gorm:"type:varchar(120);not null"`
	PasswordHash string    `gorm:"type:varchar(255);not null"`
	WeightKG     int       `gorm:"not null;check:weight_kg_gt_zero,weight_kg > 0"`
	HeightCM     int       `gorm:"not null;check:height_cm_gt_zero,height_cm > 0"`
	CreatedAt    time.Time `gorm:"not null;default:now()"`
	UpdatedAt    time.Time `gorm:"not null;default:now()"`
	Workouts     []Workout
}

func (User) TableName() string { return "users" }
