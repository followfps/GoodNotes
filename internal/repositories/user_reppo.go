package repositories

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"note1/internal/models"
	"strings"
)

type GORMUserRepository struct {
	db *gorm.DB
}

func NewGORMUserRepository(db *gorm.DB) *GORMUserRepository {
	return &GORMUserRepository{db: db}
}

func (u *GORMUserRepository) FindUserById(id uint) (*models.Users, error) {
	var user models.Users
	err := u.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *GORMUserRepository) FindUserByEmail(email string) (*models.Users, error) {
	var user models.Users
	err := u.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *GORMUserRepository) FindUserByUsername(username string) (*models.Users, error) {
	var user models.Users
	err := u.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *GORMUserRepository) CreateUser(user *models.Users) error {
	err := u.db.Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *GORMUserRepository) UpdateUser(user *models.Users) error {
	err := u.db.Save(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *GORMUserRepository) DeleteUser(user *models.Users) error {
	err := u.db.Delete(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *GORMUserRepository) FindAllUsers() ([]*models.Users, error) {
	var users []*models.Users
	err := u.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *GORMUserRepository) EmailExists(email string) (bool, error) {
	var count int64
	u.db.Model(&models.Users{}).Where("email = ?", email).Count(&count)
	return count > 0, nil
}

func (u *GORMUserRepository) ValidationEmailCheck(email string) bool {
	if !strings.Contains(email, "@") {
		return true
	}
	if !strings.Contains(email, ".") {
		return true
	}
	parts := strings.Split(email, "@")

	if len(parts) != 2 {
		return true
	}
	domain := parts[1]

	if !strings.Contains(domain, ".") {
		return true
	}
	return false

}

func (u *GORMUserRepository) FindUserByID(ID uuid.UUID) (*models.Users, error) {
	user := &models.Users{}
	err := u.db.Where("user_id = ?", ID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
