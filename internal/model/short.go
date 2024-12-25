package model

import "gorm.io/gorm"

type Short struct {
	gorm.Model
}

func (m *Short) TableName() string {
    return "short"
}
