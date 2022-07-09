package service

import (
	"web_app/internal/database"
	"web_app/internal/domain"
)

type Pictures interface {
	CreatePictureRecord(pic domain.Picture) error
	GetPicturePathById(id uint) (string, error)
	GetPictureById(id uint) (domain.Picture, error)
	//todo. Добавить signIn, signUp...
}

type Services struct {
	Pictures Pictures
	//todo. Добавить, возможно Users, Files...
}

type Deps struct {
	Tables *database.Tables
	//todo. Добавить TokenManager, PasswordHasher, Cacher...
	Environment string
	Domain      string
}

func NewServices(deps Deps) *Services {
	return &Services{
		Pictures: NewPicturesService(deps.Tables.Pictures),
	}
}
