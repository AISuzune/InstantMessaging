package model

import (
	"github.com/jinzhu/gorm"
)

type Community struct {
	gorm.Model
	Name    string
	OwnerId uint //群主ID
	Img     string
	Type    string
	Desc    string
}
