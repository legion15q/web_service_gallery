package database

import (
	"web_app/internal/domain"

	"github.com/jinzhu/gorm"
)

type Tables struct {
	Pictures Pictures
	//todo. Добавить users, возможно admins
}

//todo. Добавить users, возможно admins
//Интерфейс взаимодейтсвия с таблицей картин
type Pictures interface {
	AddPicture(pic domain.Picture) error
	GetPictureById(id uint) (domain.Picture, error)
	GetPicturePathById(id uint) (string, error)
}

func NewTables(db *gorm.DB) *Tables {

	return &Tables{
		Pictures: NewPicturesTableIFNotExist(db),
		//todo. Добавить users, возможно admins
	}
}
