package model

import "time"

type Workout struct {
	ID          uint      `gorm:"primaryKey"`
	UserID      uint      `gorm:"not null;index"`
	Name        string    `gorm:"type:varchar(120);not null"`
	Description string    `gorm:"type:text;not null"`
	CreatedAt   time.Time `gorm:"not null;default:now()"`
	UpdatedAt   time.Time `gorm:"not null;default:now()"`

	Exercises []Exercise
	User      User `gorm:"constraint:OnDelete:CASCADE"`
}

func (Workout) TableName() string { return "workouts" }
