package database

import (
	"database/sql"
	"web_app/internal/domain"
)

type Tables struct {
	Pictures Pictures
}

type Pictures interface {
	AddPicture(db *sql.DB, pic domain.Picture) error
	GetPictureById(db *sql.DB, id int) (domain.Picture, error)
	GetPicturePathById(db *sql.DB, id int) (string, error)
}

func NewTables(db *sql.DB) *Tables {
	return &Tables{
		Pictures: NewPicturesTable(db),
	}
}
