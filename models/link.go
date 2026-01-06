package models

import (
	"gorm.io/gorm"
)

type Link struct {
	gorm.Model
	URL        string `gorm:"uniqueIndex"`
	Type       string
	HTML       *string
	Word       *string `gorm:"index"`
	ParentID   *uint
	ParentURL  *string
	Vocabulary string `gorm:"default:sum.in.ua;not null"`
	Desc       *string
	Title      *string `gorm:"index"`
}
