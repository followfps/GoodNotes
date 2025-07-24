package repositories

import (
	"gorm.io/gorm"
	"note1/internal/models"
)

type GORMNoteRepository struct {
	db *gorm.DB
}

func NewGORMNoteRepository(db *gorm.DB) *GORMNoteRepository {
	return &GORMNoteRepository{db: db}
}

func (r *GORMNoteRepository) Create(note *models.Note) error {
	return r.db.Create(note).Error
}

func (r *GORMNoteRepository) FindByID(id uint) (*models.Note, error) {
	var note models.Note
	if err := r.db.Where("id = ?", id).First(&note).Error; err != nil {
		return nil, err
	}
	return &note, nil
}

func (r *GORMNoteRepository) Update(note *models.Note) error {
	return r.db.Save(note).Error
}

func (r *GORMNoteRepository) Delete(id uint) error {
	var note models.Note
	if err := r.db.Delete(&note, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *GORMNoteRepository) GetNotesFrom(pageInt, limit int) (*[]models.Note, error) {

	offset := (pageInt - 1) * limit

	var notes []models.Note

	err := r.db.Limit(limit).Offset(offset).Find(&notes).Error
	if err != nil {
		return nil, err
	}

	return &notes, nil
}
