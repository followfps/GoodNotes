package services

import (
	"errors"
	"github.com/google/uuid"
	"note1/internal/models"
	"note1/internal/repositories"
)

type UserService interface {
	CreateUser(user *models.Users) error
	EmailExists(email string) (bool, error)
	FindUserByEmail(email string) (*models.Users, error)
	ValidationEmailCheck(email string) error
	FindUserById(id *uuid.UUID) (*models.Users, error)
}

type userServiceImpl struct {
	repo *repositories.GORMUserRepository
}

func NewUserService(repo *repositories.GORMUserRepository) UserService {
	return &userServiceImpl{repo: repo}
}

func (u *userServiceImpl) CreateUser(user *models.Users) error {
	return u.repo.CreateUser(user)
}

func (u *userServiceImpl) EmailExists(email string) (bool, error) {
	return u.repo.EmailExists(email)
}

func (u *userServiceImpl) FindUserByEmail(email string) (*models.Users, error) {
	return u.repo.FindUserByEmail(email)
}

func (u *userServiceImpl) ValidationEmailCheck(email string) error {
	if u.repo.ValidationEmailCheck(email) {
		return errors.New("email dont validate")
	}
	return nil
}

func (u *userServiceImpl) FindUserById(id *uuid.UUID) (*models.Users, error) {
	return u.repo.FindUserByID(*id)
}
