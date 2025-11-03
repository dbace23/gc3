package usersvc

import (
	userrepo "gym/repository/user"
)

type Service interface {
	FindByID(id uint) (*UserDTO, error)
}

type service struct {
	users userrepo.Repo
}

func New(users userrepo.Repo) Service {
	return &service{users: users}
}

type UserDTO struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
	WeightKG int    `json:"weight"`
	HeightCM int    `json:"height"`
}

func (s *service) FindByID(id uint) (*UserDTO, error) {
	u, err := s.users.FindByID(id)
	if err != nil {
		return nil, err
	}

	return &UserDTO{
		ID:       u.ID,
		Email:    u.Email,
		Username: u.Username,
		FullName: u.FullName,
		WeightKG: u.WeightKG,
		HeightCM: u.HeightCM,
	}, nil
}
