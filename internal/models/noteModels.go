package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// Note Структура для хранения заметки.
//
// swagger:model Note
type Note struct {
	// ID Уникальный идентификатор заметки.
	// required: true
	ID uint `json:"id" gorm:"primarykey"`
	// CreatedAt Дата и время создания заметки.
	// required: true
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt Дата и время последнего обновления заметки.
	// required: true
	UpdatedAt time.Time `json:"updated_at"`
	// DeletedAt Дата и время удаления заметки (если удалена). Используется для мягкого удаления.
	// format: date-time
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index" swaggertype:"string" format:"date-time"`
	// NoteName Название заметки.
	// example: Моя заметка
	NoteName string `json:"note_name"`

	// NoteBody Текст заметки.
	// example: Это содержимое моей заметки.
	NoteBody string `json:"note_body"`

	// FileBucket Имя бакета в MinIO, где хранятся файлы заметки.
	// example: user_bucket_abc123
	FileBucket string `json:"file_bucket"`

	// FilePrefix Префикс для файлов этой заметки в MinIO.
	// example: note_12345678
	FilePrefix string `json:"file_prefix"`

	// CreatedBy UUID пользователя, создавшего заметку.
	// example: 550e8400-e29b-41d4-a716-446655440000
	// format: uuid
	CreatedBy uuid.UUID `json:"created_by"`
}
