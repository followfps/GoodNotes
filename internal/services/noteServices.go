package services

import (
	"context"
	"errors"
	"note1/internal/models"
	"note1/internal/repositories"
)

type NoteService interface {
	CreateNote(ctx context.Context, note *models.Note) error
	GetNoteByID(ctx context.Context, id uint) (*models.Note, error)
	DeleteNote(ctx context.Context, id uint) error
	GetNotesFrom(ctx context.Context, pageInt, limit int) (*[]models.Note, error)
}

type noteServiceImpl struct {
	repo *repositories.GORMNoteRepository
}

func NewNoteService(repo *repositories.GORMNoteRepository) NoteService {
	return &noteServiceImpl{repo: repo}
}

func (s *noteServiceImpl) CreateNote(ctx context.Context, note *models.Note) error {
	if note.NoteName == "" || note.NoteBody == "" {
		return errors.New("note name and body are required")
	}
	return s.repo.Create(ctx, note)
}

func (s *noteServiceImpl) GetNoteByID(ctx context.Context, id uint) (*models.Note, error) {
	note, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return note, nil
}

func (s *noteServiceImpl) DeleteNote(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *noteServiceImpl) GetNotesFrom(ctx context.Context, pageInt, limit int) (*[]models.Note, error) {

	notes, err := s.repo.GetNotesFrom(ctx, pageInt, limit)
	if err != nil {
		return nil, err
	}
	return notes, nil
}
