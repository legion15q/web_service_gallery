package service

import (
	"web_app/internal/database"
	"web_app/internal/domain"
)

type PicturesService struct {
	//todo. Add hasher       hash.PasswordHasher
	table database.Pictures
}

func NewPicturesService(table database.Pictures) *PicturesService {
	return &PicturesService{
		table: table,
	}
}

func (obj *PicturesService) CreatePictureRecord(pic domain.Picture) error {
	return obj.table.AddPicture(pic)
}

func (obj *PicturesService) GetPicturePathById(id uint) (string, error) {
	return obj.table.GetPicturePathById(id)
}

func (obj *PicturesService) GetPictureById(id uint) (domain.Picture, error) {
	return obj.table.GetPictureById(id)
}
