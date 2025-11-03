package model

import "time"

type ExerciseLog struct {
	ID         uint      `gorm:"primaryKey"`
	ExerciseID uint      `gorm:"not null;index"`
	UserID     uint      `gorm:"not null;index"`
	SetCount   int       `gorm:"not null;check:set_count_gt_zero,set_count > 0"`
	RepCount   int       `gorm:"not null;check:rep_count_gt_zero,rep_count > 0"`
	Weight     int       `gorm:"not null;check:weight_gte_zero,weight >= 0"`
	LoggedAt   time.Time `gorm:"not null"`

	Exercise Exercise `gorm:"constraint:OnDelete:CASCADE"`
	User     User     `gorm:"constraint:OnDelete:CASCADE"`
}

func (ExerciseLog) TableName() string { return "exercise_logs" }
