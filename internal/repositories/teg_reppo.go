package repositories

import "gorm.io/gorm"

type GORMTagRepository struct {
	db *gorm.DB
}

func NewTegReppo(db *gorm.DB) *GORMTagRepository {
	return &GORMTagRepository{db: db}
}
