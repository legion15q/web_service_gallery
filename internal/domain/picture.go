package domain

type Picture struct {
	ID                  uint    `json:"-" gorm:"primary_key"`
	Picture_name        string  `json:"picture_name"`
	Picture_description string  `json:"picture_description"`
	Author              string  `json:"author"`
	Price               float32 `json:"price"`
	Is_purchased        bool    `json:"is_purchased"`
	Picture_path        string  `json:"picture_path"`
}

//todo. Добавить описание таблицы users, возможно admins. Либо в этот go файл, либо в отдельный
