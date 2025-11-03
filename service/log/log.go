package logsrv

import (
	"errors"
	"time"

	"gym/model"
	exerciserepo "gym/repository/exercise"
	logrepo "gym/repository/log"
)

type Service interface {
	Create(uid uint, exerciseID uint, setCount, repCount, weight int, loggedAt time.Time) error
	ListMy(uid uint) (any, error)
}

type service struct {
	lr logrepo.Repo
	ex exerciserepo.Repo
}

func New(lr logrepo.Repo, ex exerciserepo.Repo) Service { return &service{lr: lr, ex: ex} }

func (s *service) Create(uid, exerciseID uint, setCount, repCount, weight int, loggedAt time.Time) error {
	ownerID, err := s.ex.FindWorkoutByExerciseID(exerciseID)
	if err != nil {
		return err
	}
	if ownerID != uid {
		return errors.New("forbidden: exercise not owned by user")
	}
	return s.lr.Create(&model.ExerciseLog{
		ExerciseID: exerciseID, UserID: uid,
		SetCount: setCount, RepCount: repCount, Weight: weight,
		LoggedAt: loggedAt,
	})
}

func (s *service) ListMy(uid uint) (any, error) {
	return s.lr.ListByUser(uid)
}
