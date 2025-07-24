package services

import (
	"errors"
	"note1/internal/models"
	"note1/internal/repositories"
)

type NoteService interface {
	CreateNote(note *models.Note) error
	GetNoteByID(id uint) (*models.Note, error)
	DeleteNote(id uint) error
	GetNotesFrom(pageInt, limit int) (*[]models.Note, error)
}

type noteServiceImpl struct {
	repo *repositories.GORMNoteRepository
}

func NewNoteService(repo *repositories.GORMNoteRepository) NoteService {
	return &noteServiceImpl{repo: repo}
}

func (s *noteServiceImpl) CreateNote(note *models.Note) error {
	if note.NoteName == "" || note.NoteBody == "" {
		return errors.New("note name and body are required")
	}
	return s.repo.Create(note)
}

func (s *noteServiceImpl) GetNoteByID(id uint) (*models.Note, error) {
	note, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return note, nil
}

func (s *noteServiceImpl) DeleteNote(id uint) error {
	return s.repo.Delete(id)
}

func (s *noteServiceImpl) GetNotesFrom(pageInt, limit int) (*[]models.Note, error) {

	notes, err := s.repo.GetNotesFrom(pageInt, limit)
	if err != nil {
		return nil, err
	}
	return notes, nil
}
