package exerciserepo

import (
	"gym/model"

	"gorm.io/gorm"
)

type Repo interface {
	FindByID(id uint) (*model.Exercise, error)
	Create(e *model.Exercise) error
	Delete(id uint) error
	FindWorkoutByExerciseID(id uint) (uint, error) // returns workout.user_id owner
}

type repo struct{ db *gorm.DB }

func New(db *gorm.DB) Repo { return &repo{db} }

func (r *repo) FindByID(id uint) (*model.Exercise, error) {
	var e model.Exercise
	if err := r.db.First(&e, id).Error; err != nil {
		return nil, err
	}
	return &e, nil
}
func (r *repo) Create(e *model.Exercise) error { return r.db.Create(e).Error }
func (r *repo) Delete(id uint) error           { return r.db.Delete(&model.Exercise{}, id).Error }

func (r *repo) FindWorkoutByExerciseID(id uint) (uint, error) {
	var e model.Exercise
	if err := r.db.First(&e, id).Error; err != nil {
		return 0, err
	}
	var w model.Workout
	if err := r.db.Select("user_id").First(&w, e.WorkoutID).Error; err != nil {
		return 0, err
	}
	return w.UserID, nil
}
