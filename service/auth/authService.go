package authsvc

import (
	"errors"
	"gym/model"
	"time"

	userrepo "gym/repository/user"
	"gym/util/hash"
	jwtutil "gym/util/jwt"
)

type Service interface {
	Register(p RegisterParams) (*RegisterResult, error)
	Login(p LoginParams, jwtSecret string) (*LoginResult, error)
}

type service struct{ users userrepo.Repo }

func New(users userrepo.Repo) Service { return &service{users: users} }

type RegisterParams struct {
	FullName string
	Email    string
	Username string
	Password string
	WeightKG int
	HeightCM int
}

type RegisterResult struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
	WeightKG int    `json:"weight"`
	HeightCM int    `json:"height"`
}

func (s *service) Register(p RegisterParams) (*RegisterResult, error) {
	if dup, err := s.users.ExistsEmail(p.Email); err != nil {
		return nil, err
	} else if dup {
		return nil, errors.New("email already registered")
	}
	if dup, err := s.users.ExistsUsername(p.Username); err != nil {
		return nil, err
	} else if dup {
		return nil, errors.New("username already taken")
	}
	hashed, err := hash.Make(p.Password)
	if err != nil {
		return nil, err
	}
	u := &model.User{
		Email:        p.Email,
		Username:     p.Username,
		FullName:     p.FullName,
		PasswordHash: hashed,
		WeightKG:     p.WeightKG,
		HeightCM:     p.HeightCM,
	}
	if err := s.users.Create(u); err != nil {
		return nil, err
	}
	return &RegisterResult{
		ID:       u.ID,
		Email:    u.Email,
		Username: u.Username,
		FullName: u.FullName,
		WeightKG: u.WeightKG,
		HeightCM: u.HeightCM,
	}, nil
}

type LoginParams struct {
	Email    string
	Password string
}

type LoginResult struct {
	AccessToken string `json:"access_token"`
}

func (s *service) Login(p LoginParams, jwtSecret string) (*LoginResult, error) {
	u, err := s.users.FindByEmail(p.Email)
	if err != nil {
		if userrepo.IsNotFound(err) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}
	if !hash.Check(u.PasswordHash, p.Password) {
		return nil, errors.New("invalid email or password")
	}
	tok, err := jwtutil.SignHS256(jwtSecret, u.ID, 60*time.Minute)
	if err != nil {
		return nil, err
	}
	return &LoginResult{AccessToken: tok}, nil
}
