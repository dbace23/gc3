package workoutsrv

import (
	"errors"
	"gym/model"
	workoutrepo "gym/repository/workout"
)

type Service interface {
	ListMy(userID uint) ([]model.Workout, error)
	GetMy(userID, workoutID uint) (*model.Workout, error)
	Create(userID uint, name, desc string) (*model.Workout, error)
	Update(userID, workoutID uint, name, desc string) (*model.Workout, error)
	Delete(userID, workoutID uint) error
}

type service struct{ r workoutrepo.Repo }

func New(r workoutrepo.Repo) Service { return &service{r} }

func (s *service) ListMy(uid uint) ([]model.Workout, error) { return s.r.ListByUser(uid) }

func (s *service) GetMy(uid, wid uint) (*model.Workout, error) {
	w, err := s.r.FindByID(wid)
	if err != nil {
		return nil, err
	}
	if w.UserID != uid {
		return nil, errors.New("forbidden: not owner")
	}
	return w, nil
}

func (s *service) Create(uid uint, name, desc string) (*model.Workout, error) {
	w := &model.Workout{UserID: uid, Name: name, Description: desc}
	if err := s.r.Create(w); err != nil {
		return nil, err
	}
	return w, nil
}

func (s *service) Update(uid, wid uint, name, desc string) (*model.Workout, error) {
	w, err := s.r.FindByID(wid)
	if err != nil {
		return nil, err
	}
	if w.UserID != uid {
		return nil, errors.New("forbidden: not owner")
	}
	w.Name, w.Description = name, desc
	if err := s.r.Update(w); err != nil {
		return nil, err
	}
	return w, nil
}

func (s *service) Delete(uid, wid uint) error {
	w, err := s.r.FindByID(wid)
	if err != nil {
		return err
	}
	if w.UserID != uid {
		return errors.New("forbidden: not owner")
	}
	return s.r.Delete(wid)
}
