package models

import (
	"time"

	"gorm.io/gorm"
)

type LocationData struct {
	gorm.Model
	UIC           string `gorm:"column:uic"`
	NLC           string
	CRS           string
	FareGroup     string
	Description   string
	StartDate     *time.Time
	EndDate       *time.Time
	GroupedMember bool
}

func (LocationData) TableName() string {
	return "location"
}
