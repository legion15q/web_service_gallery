package domain

type Picture struct {
	ID                  uint    `json:"-" gorm:"primary_key;autoIncrement:false" schema:"-"`
	Picture_name        string  `json:"picture_name" schema:"picture_name"`
	Picture_description string  `json:"picture_description" schema:"picture_description"`
	Author              string  `json:"author" schema:"author"`
	Price               float32 `json:"price" schema:"price"`
	Is_purchased        bool    `json:"is_purchased" schema:"is_purchased"`
	Picture_path        string  `json:"picture_path" schema:"-"`
}

//todo. Добавить описание таблицы users, возможно admins. Либо в этот go файл, либо в отдельный
