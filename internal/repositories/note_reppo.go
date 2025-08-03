package repositories

import (
	"context"
	"gorm.io/gorm"
	"note1/internal/models"
)

type GORMNoteRepository struct {
	db *gorm.DB
}

func NewGORMNoteRepository(db *gorm.DB) *GORMNoteRepository {
	return &GORMNoteRepository{db: db}
}

func (r *GORMNoteRepository) Create(ctx context.Context, note *models.Note) error {
	return r.db.WithContext(ctx).Create(note).Error
}

func (r *GORMNoteRepository) FindByID(ctx context.Context, id uint) (*models.Note, error) {
	var note models.Note
	err := r.db.WithContext(ctx).
		Preload("Tags").
		Where("note_id = ?", id).
		First(&note).Error
	if err != nil {
		return nil, err
	}
	return &note, nil
}

func (r *GORMNoteRepository) Update(ctx context.Context, note *models.Note) error {
	return r.db.WithContext(ctx).Save(note).Error
}

func (r *GORMNoteRepository) Delete(ctx context.Context, id uint) error {
	var note models.Note
	if err := r.db.WithContext(ctx).Delete(&note, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *GORMNoteRepository) GetNotesFrom(ctx context.Context, pageInt, limit int) (*[]models.Note, error) {

	offset := (pageInt - 1) * limit

	var notes []models.Note

	err := r.db.WithContext(ctx).
		Preload("Tags").
		Limit(limit).
		Offset(offset).
		Find(&notes).Error

	if err != nil {
		return nil, err
	}
	return &notes, nil
}

func (r *GORMNoteRepository) FindAllNotesWithTags() ([]models.Note, error) {
	var tag []models.Tags

	var notes []models.Note
	err := r.db.Model(&tag).Association("Notes").Find(&notes)
	if err != nil {
		return nil, err
	}

	return notes, nil

}
