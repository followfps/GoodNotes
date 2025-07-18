package services

import (
	"note1/internal/config"
	"note1/internal/repositories"
)

// ServicesContainer содержит все сервисы приложения.
type ServicesContainer struct {
	NoteService NoteService
	UserService UserService
}

func NewServicesContainer() *ServicesContainer {
	noteRepo := repositories.NewGORMNoteRepository(config.DBNote)
	noteService := NewNoteService(noteRepo)

	userRepo := repositories.NewGORMUserRepository(config.DBUsers)
	userService := NewUserService(userRepo)

	return &ServicesContainer{
		NoteService: noteService,
		UserService: userService,
	}
}
