package model

import "time"

type Exercise struct {
	ID          uint      `gorm:"primaryKey"`
	WorkoutID   uint      `gorm:"not null;index"`
	Name        string    `gorm:"type:varchar(120);not null"`
	Description string    `gorm:"type:text;not null"`
	CreatedAt   time.Time `gorm:"not null;default:now()"`
	UpdatedAt   time.Time `gorm:"not null;default:now()"`

	Workout Workout `gorm:"constraint:OnDelete:CASCADE"`
	Logs    []ExerciseLog
}

func (Exercise) TableName() string { return "exercises" }
