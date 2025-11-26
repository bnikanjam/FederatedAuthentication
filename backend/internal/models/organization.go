package models

import (
	"gorm.io/gorm"
)

type Organization struct {
	gorm.Model
	Domain      string `gorm:"uniqueIndex;size:255;not null"`
	Auth0OrgID  string `gorm:"not null"`
	DisplayName string
}
