package services

import (
	"note1/internal/repositories"
)

type TagServices interface {
}

type tagServiceImpl struct {
	repo *repositories.GORMTagRepository
}

func NewTagService(repo *repositories.GORMTagRepository) TagServices {
	return &tagServiceImpl{repo: repo}
}
