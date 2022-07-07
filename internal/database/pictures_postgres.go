package database

import (
	"database/sql"
	"fmt"

	_ "github.com/jinzhu/gorm"
)

type PicturesTable struct {
	db *sql.DB
}

func NewPicturesTable(db *sql.DB) PicturesTable {
	fmt.Println("asd")
	return &PicturesTable{
		db: sql.DB,
	}
}
