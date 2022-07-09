package database

import (
	"web_app/internal/domain"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type PicturesTable struct {
	db *gorm.DB
}

func NewPicturesTableIFNotExist(db *gorm.DB) PicturesTable {
	if !db.HasTable(&domain.Picture{}) {
		return PicturesTable{
			db: db.AutoMigrate(&domain.Picture{}),
		}
	} else {
		return PicturesTable{db: db}
	}

}

func (obj PicturesTable) AddPicture(pic domain.Picture) error {
	err := obj.db.Create(&pic).Error
	if err != nil {
		return errors.Wrap(err, "cant insert into database")
	}
	return nil
}

func (obj PicturesTable) GetPictureById(id uint) (domain.Picture, error) {
	picture := domain.Picture{}
	err := obj.db.Where(&domain.Picture{ID: id}).Find(&picture).Error
	if err != nil {
		return domain.Picture{}, errors.Wrap(err, "cant find picture by id")
	}
	return picture, nil

}
func (obj PicturesTable) GetPicturePathById(id uint) (string, error) {
	picture := domain.Picture{}
	err := obj.db.Where(&domain.Picture{ID: id}).Find(&picture).Error
	if err != nil {
		return "", errors.Wrap(err, "cant find picture path by id")
	}
	return picture.Picture_path, nil
}
