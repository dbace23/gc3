package workoutrepo

import (
	"gym/model"

	"gorm.io/gorm"
)

type Repo interface {
	ListByUser(userID uint) ([]model.Workout, error)
	FindByID(id uint) (*model.Workout, error)
	Create(w *model.Workout) error
	Update(w *model.Workout) error
	Delete(id uint) error
}

type repo struct{ db *gorm.DB }

func New(db *gorm.DB) Repo { return &repo{db} }

func (r *repo) ListByUser(userID uint) ([]model.Workout, error) {
	var rows []model.Workout
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&rows).Error
	return rows, err
}

func (r *repo) FindByID(id uint) (*model.Workout, error) {
	var w model.Workout
	if err := r.db.Preload("Exercises").First(&w, id).Error; err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *repo) Create(w *model.Workout) error { return r.db.Create(w).Error }
func (r *repo) Update(w *model.Workout) error { return r.db.Save(w).Error }
func (r *repo) Delete(id uint) error          { return r.db.Delete(&model.Workout{}, id).Error }
