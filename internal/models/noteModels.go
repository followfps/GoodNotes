package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Note struct {
	gorm.Model
	NoteName   string
	NoteBody   string
	FileBucket string
	CreatedBy  uuid.UUID
}
