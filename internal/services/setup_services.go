package services

import (
	"note1/internal/config"
	"note1/internal/repositories"
)

// ServicesContainer содержит все сервисы приложения.
type ServicesContainer struct {
	NoteService NoteService
	UserService UserService
	TagService  TagServices
}

func NewServicesContainer() *ServicesContainer {
	noteRepo := repositories.NewGORMNoteRepository(config.DBNote)
	noteService := NewNoteService(noteRepo)

	userRepo := repositories.NewGORMUserRepository(config.DBUsers)
	userService := NewUserService(userRepo)

	tagRepo := repositories.NewTegReppo(config.DBTags)
	TagService := NewTagService(tagRepo)

	return &ServicesContainer{
		NoteService: noteService,
		UserService: userService,
		TagService:  TagService,
	}
}
