package service

import (
	"os"
	"web_app/internal/database"
	"web_app/internal/domain"
	sender "web_app/pkg/tempaltes_parser"
)

type Pictures interface {
	CreatePictureRecord(pic domain.Picture) error
	GetPicturePathById(id uint) (string, error)
	GetPictureById(id uint) (domain.Picture, error)
	GetMaxIdValue() (uint, error)
	//todo. Добавить signIn, signUp...
}

type StorageManager interface {
	CreateUnicFile(file_extension string) (*os.File, string, error)
	GetFileStoragePath() string
}

type TemplateParser interface {
	GenerateBodyFromHTML(data interface{}) string
}

type Services struct {
	Pictures       Pictures
	StorageManager StorageManager
	TemplateParser TemplateParser
	//todo. Добавить, возможно Users, Files...
}

type Deps struct {
	Tables            *database.Tables
	File_storage_path string
	TemplateFileName  string
	//todo. Добавить TokenManager, PasswordHasher, Cacher...
	Environment string
	Domain      string
}

func NewServices(deps Deps) *Services {
	return &Services{
		Pictures:       NewPicturesService(deps.Tables.Pictures),
		StorageManager: NewStorageManagerService(deps.File_storage_path),
		TemplateParser: sender.NewTemplateParser(deps.TemplateFileName),
	}
}
