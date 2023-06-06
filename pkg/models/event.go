package models

import "gorm.io/gorm"

type RelayEvent struct {
	gorm.Model
	Data string `gorm:"type:varchar(512)"`
}
