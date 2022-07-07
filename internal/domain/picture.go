package domain

import "github.com/jinzhu/gorm"

type Picture struct {
	gorm.Model
	Picture_name        string
	Picture_description string
	Author              string
	Price               float32
	Is_purchased        bool
	Picture_path        string
}
