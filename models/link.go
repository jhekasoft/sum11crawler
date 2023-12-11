package models

import (
	"gorm.io/gorm"
)

type Link struct {
	gorm.Model
	URL       string `gorm:"uniqueIndex"`
	Type      string
	HTML      *string
	Word      *string
	ParentID  *uint
	ParentURL *string
}
