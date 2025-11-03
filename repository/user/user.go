package userrepo

import (
	"errors"

	"gym/model"

	"gorm.io/gorm"
)

type Repo interface {
	Create(u *model.User) error
	FindByEmail(email string) (*model.User, error)
	FindByID(id uint) (*model.User, error)
	ExistsEmail(email string) (bool, error)
	ExistsUsername(username string) (bool, error)
}

type repo struct{ db *gorm.DB }

func New(db *gorm.DB) Repo { return &repo{db} }

func (r *repo) Create(u *model.User) error { return r.db.Create(u).Error }

func (r *repo) FindByEmail(email string) (*model.User, error) {
	var u model.User
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repo) FindByID(id uint) (*model.User, error) {
	var u model.User
	if err := r.db.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repo) ExistsEmail(email string) (bool, error) {
	var n int64
	if err := r.db.Model(&model.User{}).Where("email = ?", email).Count(&n).Error; err != nil {
		return false, err
	}
	return n > 0, nil
}

func (r *repo) ExistsUsername(username string) (bool, error) {
	var n int64
	if err := r.db.Model(&model.User{}).Where("username = ?", username).Count(&n).Error; err != nil {
		return false, err
	}
	return n > 0, nil
}

func IsNotFound(err error) bool { return errors.Is(err, gorm.ErrRecordNotFound) }
