package exercisesrv

import (
	"errors"
	"gym/model"
	exerciserepo "gym/repository/exercise"
	workoutrepo "gym/repository/workout"
)

type Service interface {
	Create(uid uint, workoutID uint, name, desc string) (*model.Exercise, error)
	Delete(uid, exerciseID uint) error
}

type service struct {
	ex exerciserepo.Repo
	wr workoutrepo.Repo
}

func New(ex exerciserepo.Repo, wr workoutrepo.Repo) Service { return &service{ex: ex, wr: wr} }

func (s *service) Create(uid, workoutID uint, name, desc string) (*model.Exercise, error) {
	w, err := s.wr.FindByID(workoutID)
	if err != nil {
		return nil, err
	}
	if w.UserID != uid {
		return nil, errors.New("forbidden: workout not owned by user")
	}
	e := &model.Exercise{WorkoutID: workoutID, Name: name, Description: desc}
	if err := s.ex.Create(e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *service) Delete(uid, exerciseID uint) error {
	ownerID, err := s.ex.FindWorkoutByExerciseID(exerciseID)
	if err != nil {
		return err
	}
	if ownerID != uid {
		return errors.New("forbidden: not owner")
	}
	return s.ex.Delete(exerciseID)
}
