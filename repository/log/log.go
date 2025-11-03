package logrepo

import (
	"gym/model"

	"gorm.io/gorm"
)

type Repo interface {
	Create(l *model.ExerciseLog) error
	ListByUser(userID uint) ([]model.ExerciseLog, error)
}

type repo struct{ db *gorm.DB }

func New(db *gorm.DB) Repo { return &repo{db} }

func (r *repo) Create(l *model.ExerciseLog) error { return r.db.Create(l).Error }

func (r *repo) ListByUser(userID uint) ([]model.ExerciseLog, error) {
	var rows []model.ExerciseLog
	err := r.db.Preload("Exercise").Where("user_id = ?", userID).Order("logged_at DESC").Find(&rows).Error
	return rows, err
}
